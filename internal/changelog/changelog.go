package changelog

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	ErrNoUnreleasedChanges         = errors.New("automatic tag detection does not work without unreleased changelogs")
	ErrEntriesWithoutSeverityFound = errors.New("all entries must have a severity for the automatic tag detection to work")
	ErrLastTagMalformed            = errors.New("last tag is malformed")
)

type Changelog struct {
	Config  Config
	Catalog Catalog
}

func New() *Changelog {
	return &Changelog{
		Config: DefaultConfig,
	}
}

func (c *Changelog) GetCatalogFilePath() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir, c.Config.CatalogFileName)
}

func (c *Changelog) GetChangelogDir() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir)
}

func (c *Changelog) GetUnreleasedDir() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir, c.Config.UnreleasedDir)
}

func (c *Changelog) GetUnreleasedKeepFilePath() string {
	return filepath.Join(c.GetUnreleasedDir(), ".keep")
}

func (c *Changelog) getNewUnreleasedEntryFilePath() string {
	return filepath.Join(c.GetUnreleasedDir(), time.Now().Format(c.Config.UnreleasedFilePattern)+".yaml")
}

func (c *Changelog) getUnreleasedChangelogFilePaths() ([]string, error) {
	entries, err := os.ReadDir(c.GetUnreleasedDir())
	if err != nil {
		return []string{}, fmt.Errorf("failed to get unreleased changelogs; %w", err)
	}

	var filePaths []string
	for _, entry := range entries {
		if !entry.IsDir() &&
			filepath.Ext(strings.ToLower(entry.Name())) == ".yaml" {
			filePaths = append(filePaths, filepath.Join(c.GetUnreleasedDir(), entry.Name()))
		}
	}
	return filePaths, nil
}

func (c *Changelog) GetChangelogFilePath() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogFileName)
}

func (c *Changelog) loadConfig() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory; %w", err)
	}

	if err := c.Config.LoadFromDir(workingDir); err != nil {
		return errors.New("failed to load config file")
	}

	return nil
}

func (c *Changelog) loadCatalog() error {
	return c.Catalog.LoadFromFile(c.GetCatalogFilePath())
}

func (c *Changelog) loadUnreleased() error {
	c.Catalog.Unreleased = []Entry{}
	filePaths, err := c.getUnreleasedChangelogFilePaths()
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		var entry Entry
		if err := entry.LoadFromFile(filePath); err != nil {
			return fmt.Errorf("failed to load unreleased changelog; %w", err)
		}
		c.Catalog.Unreleased = append(c.Catalog.Unreleased, entry)
	}
	return nil
}

func (c *Changelog) deleteUnreleasedChangelogs() error {
	filePaths, err := c.getUnreleasedChangelogFilePaths()
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to clean up unreleased changelog entries; %w", err)
		}
	}
	return nil
}

func (c *Changelog) Load() error {
	if err := c.loadConfig(); err != nil {
		return err
	}
	if err := c.loadCatalog(); err != nil {
		return err
	}
	if err := c.loadUnreleased(); err != nil {
		return err
	}
	return nil
}

func (c *Changelog) SaveCatalog() error {
	if err := c.Catalog.SaveToFile(c.GetCatalogFilePath()); err != nil {
		return fmt.Errorf("failed to save catalog to file; %w", err)
	}
	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (c *Changelog) Init() error {
	if fileExists(ConfigFilename) {
		fmt.Printf("WARNING: Skipping creation of %q because it already exists.\n", ConfigFilename)
	} else {
		if err := c.Config.SaveToFile(ConfigFilename); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(c.GetChangelogDir(), 0o775); err != nil { //nolint:gosec // Those permissions are fine
		return err
	}

	if err := os.MkdirAll(c.GetUnreleasedDir(), 0o775); err != nil { //nolint:gosec // Those permissions are fine
		return err
	}

	if fileExists(c.GetUnreleasedKeepFilePath()) {
		fmt.Printf("WARNING: Skipping creation of %q because it already exists.\n", c.GetUnreleasedKeepFilePath())
	} else {
		if err := os.WriteFile(c.GetUnreleasedKeepFilePath(), []byte{}, 0o664); err != nil { //nolint:gosec // Those permissions are fine
			return err
		}
	}

	if fileExists(c.GetCatalogFilePath()) {
		fmt.Printf("WARNING: Skipping creation of %q because it already exists.\n", c.GetCatalogFilePath())
	} else {
		if err := c.SaveCatalog(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Changelog) Add(description string, severity Severity) error {
	entry := Entry{
		Description: description,
		Severity:    severity,
	}
	c.Catalog.Unreleased = append(c.Catalog.Unreleased, entry)

	if err := entry.SaveToFile(c.getNewUnreleasedEntryFilePath()); err != nil {
		return fmt.Errorf("failed to add changelog entry; %w", err)
	}

	return nil
}

func (c *Changelog) RenderChangelog() error {
	templ, err := template.New("changelog").Funcs(template.FuncMap{"trim": strings.TrimSpace}).Parse(c.Config.Template)
	if err != nil {
		return fmt.Errorf("failed to parse changelog template; %w", err)
	}

	file, err := os.Create(c.GetChangelogFilePath())
	if err != nil {
		return fmt.Errorf("failed to create changelog file; %w", err)
	}
	defer file.Close() //nolint:errcheck

	if err := templ.Execute(file, c); err != nil {
		return fmt.Errorf("failed to render changelog template; %w", err)
	}
	return nil
}

func (c *Changelog) CreateRelease(version string) error {
	release := Release{
		Version: version,
		Date:    time.Now().Truncate(24 * time.Hour),
		Entries: c.Catalog.Unreleased,
	}
	c.Catalog.Releases = append([]Release{release}, c.Catalog.Releases...)
	c.Catalog.Unreleased = []Entry{}
	if err := c.Catalog.SaveToFile(c.GetCatalogFilePath()); err != nil {
		return fmt.Errorf("failed to save changelog catalog; %w", err)
	}

	return c.deleteUnreleasedChangelogs()
}

func (c *Changelog) NextTag() (string, error) {
	if len(c.Catalog.Unreleased) == 0 {
		return "", ErrNoUnreleasedChanges
	}

	maxSeverity, err := c.getMaximumSeverityFromUnreleased()
	if err != nil {
		return "", err
	}

	lastTag := c.getLastTag()

	pureTag := strings.Trim(lastTag, "v")
	tagParts := strings.Split(pureTag, ".")

	if len(tagParts) != 3 {
		return "", ErrLastTagMalformed
	}

	for severityPositionInTag, tagPart := range tagParts {
		// Since severity 0 is reserved for unknown due to backwards compatibility the enum starts at 1.
		// To iterate through the different parts of the tag, corresponding to the different severities, and we need to subtract this offset here again.
		severityPositionToBump := int(maxSeverity) - 1

		if severityPositionInTag < severityPositionToBump {
			continue
		}
		if severityPositionInTag > severityPositionToBump {
			tagParts[severityPositionInTag] = "0"
			continue
		}

		parsedTagPart, err := strconv.Atoi(tagPart)
		if err != nil {
			return "", ErrLastTagMalformed
		}
		tagParts[severityPositionInTag] = strconv.Itoa(parsedTagPart + 1)
	}

	newTag := "v" + strings.Join(tagParts, ".")
	return newTag, nil
}

func (c *Changelog) getMaximumSeverityFromUnreleased() (Severity, error) {
	maxSeverity := Severity(math.MaxInt)
	for _, unreleased := range c.Catalog.Unreleased {
		if unreleased.Severity == SeverityUnknown {
			return SeverityUnknown, ErrEntriesWithoutSeverityFound
		}
		maxSeverity = min(maxSeverity, unreleased.Severity)
	}
	return maxSeverity, nil
}

func (c *Changelog) getLastTag() string {
	lastTag := "v0.0.0"
	if len(c.Catalog.Releases) > 0 {
		lastTag = c.Catalog.Releases[0].Version
	}
	return lastTag
}

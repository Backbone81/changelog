package changelog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Changelog struct {
	Config  Config
	Catalog Catalog
}

func New() *Changelog {
	return &Changelog{}
}

func (c *Changelog) getCatalogFilePath() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir, c.Config.CatalogFileName)
}

func (c *Changelog) getChangelogDir() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir)
}

func (c *Changelog) getUnreleasedDir() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogDir, c.Config.UnreleasedDir)
}

func (c *Changelog) getUnreleasedKeepFilePath() string {
	return filepath.Join(c.getUnreleasedDir(), ".keep")
}

func (c *Changelog) getUnreleasedEntryFilePath() string {
	return filepath.Join(c.getUnreleasedDir(), time.Now().Format(c.Config.UnreleasedFilePattern)+".yaml")
}

func (c *Changelog) LoadConfig() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory; %w", err)
	}

	if err := c.Config.loadFromDir(workingDir); err != nil {
		return fmt.Errorf("failed to load config file")
	}

	return nil
}

func (c *Changelog) SaveCatalog() error {
	if err := c.Catalog.SaveToFile(c.getCatalogFilePath()); err != nil {
		return fmt.Errorf("failed to save catalog to file; %w", err)
	}
	return nil
}

func (c *Changelog) Init() error {
	c.Config = DefaultConfig

	if err := c.Config.saveToFile(ConfigFilename); err != nil {
		return err
	}

	if err := os.MkdirAll(c.getChangelogDir(), 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(c.getUnreleasedDir(), 0777); err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.getUnreleasedKeepFilePath(), []byte{}, 0666); err != nil {
		return err
	}

	if err := c.SaveCatalog(); err != nil {
		return err
	}

	return nil
}

func (c *Changelog) Add(description string) error {
	entry := Entry{
		Description: description,
	}

	data, err := yaml.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry; %w", err)
	}

	if err := ioutil.WriteFile(c.getUnreleasedEntryFilePath(), data, 0666); err != nil {
		return fmt.Errorf("failed to write unreleased changelog file; %w", err)
	}

	return nil
}

func (c *Changelog) LoadCatalog() error {
	return c.Catalog.LoadFromFile(c.getCatalogFilePath())
}

func (c *Changelog) getUnreleasedChangelogFilePaths() ([]string, error) {
	entries, err := os.ReadDir(c.getUnreleasedDir())
	if err != nil {
		return []string{}, fmt.Errorf("failed to get unreleased changelogs; %w", err)
	}

	var filePaths []string
	for _, entry := range entries {
		if !entry.IsDir() &&
			filepath.Ext(strings.ToLower(entry.Name())) == ".yaml" {
			filePaths = append(filePaths, filepath.Join(c.getUnreleasedDir(), entry.Name()))
		}
	}
	return filePaths, nil
}

func (c *Changelog) LoadUnreleased() error {
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

func (c *Changelog) getChangelogFilePath() string {
	return filepath.Join(c.Config.BaseDir, c.Config.ChangelogFileName)
}

func (c *Changelog) RenderChangelog() error {
	templ, err := template.New("changelog").Funcs(template.FuncMap{"trim": strings.TrimSpace}).Parse(c.Config.Template)
	if err != nil {
		return fmt.Errorf("failed to parse changelog template; %w", err)
	}

	file, err := os.Create(c.getChangelogFilePath())
	if err != nil {
		return fmt.Errorf("failed to create changelog file; %w", err)
	}
	defer file.Close()

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
	if err := c.Catalog.SaveToFile(c.getCatalogFilePath()); err != nil {
		return fmt.Errorf("failed to save changelog catalog; %w", err)
	}

	return c.deleteUnreleasedChangelogs()
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

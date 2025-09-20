package changelog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Catalog is a structure which holds all releases and all unreleased changelog entries.
type Catalog struct {
	Unreleased []Entry   // Unreleased holds a list of all unreleased changelog entries
	Releases   []Release // Releases holds a list of all releases with all their changelog entries
}

// SaveToFile writes the catalog as yaml file to disk.
func (c *Catalog) SaveToFile(filePath string) error {
	data, err := yaml.Marshal(*c)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog; %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o664); err != nil { //nolint:gosec // Those permissions are fine
		return fmt.Errorf("failed to write catalog file; %w", err)
	}

	return nil
}

// LoadFromFile reads the catalog as yaml file from disk.
func (c *Catalog) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to read catalog file; %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal catalog; %w", err)
	}

	return nil
}

package changelog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Catalog is a structure which holds all releases and all unreleased changelog entries
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

	if err := ioutil.WriteFile(filePath, data, 0666); err != nil {
		return fmt.Errorf("failed to write catalog file; %w", err)
	}

	return nil
}

// LoadFromFile reads the catalog as yaml file from disk.
func (c *Catalog) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read catalog file; %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal catalog; %w", err)
	}

	return nil
}

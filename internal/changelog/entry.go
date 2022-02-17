package changelog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Entry is a struct which holds all relevant information about a single changelog entry
type Entry struct {
	Description string // Description provides some plain text describing the change
}

// SaveToFile writes the entry as yaml file to disk.
func (e *Entry) SaveToFile(filePath string) error {
	data, err := yaml.Marshal(*e)
	if err != nil {
		return fmt.Errorf("failed to marshal entry; %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0666); err != nil {
		return fmt.Errorf("failed to write unreleased changelog file; %w", err)
	}

	return nil
}

// LoadFromFile reads the entry as yaml file from disk.
func (e *Entry) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read changelog entry; %w", err)
	}

	if err := yaml.Unmarshal(data, e); err != nil {
		return fmt.Errorf("failed to unmarshal changelog entry; %w", err)
	}

	return nil
}

package changelog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Entry is a struct which holds all relevant information about a single changelog entry.
type Entry struct {
	// Description provides some plain text describing the change
	Description string

	// Severity describes the level of a change. Possible values are major, minor and patch
	Severity Severity `yaml:"severity,omitempty"`
}

// SaveToFile writes the entry as yaml file to disk.
func (e *Entry) SaveToFile(filePath string) error {
	data, err := yaml.Marshal(*e)
	if err != nil {
		return fmt.Errorf("failed to marshal entry; %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o664); err != nil { //nolint:gosec // Those permissions are fine
		return fmt.Errorf("failed to write unreleased changelog file; %w", err)
	}

	return nil
}

// LoadFromFile reads the entry as yaml file from disk.
func (e *Entry) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to read changelog entry; %w", err)
	}

	if err := yaml.Unmarshal(data, e); err != nil {
		return fmt.Errorf("failed to unmarshal changelog entry; %w", err)
	}

	return nil
}

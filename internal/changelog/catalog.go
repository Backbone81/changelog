package changelog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Catalog struct {
	Unreleased []Entry
	Releases   []Release
}

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

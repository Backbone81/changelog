package changelog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Entry struct {
	Description string
}

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

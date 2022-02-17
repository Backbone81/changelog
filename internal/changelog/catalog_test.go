package changelog

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var testCatalog = Catalog{
	Unreleased: []Entry{
		{
			Description: "Some feature",
		},
		{
			Description: "Some bugfix",
		},
	},
	Releases: []Release{
		{
			Version: "v1.5.3",
			Date:    time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC),
			Entries: []Entry{
				{
					Description: "Great feature",
				},
				{
					Description: "Small feature",
				},
			},
		},
	},
}

const testCatalogYaml = `unreleased:
- description: Some feature
- description: Some bugfix
releases:
- version: v1.5.3
  date: 2020-04-02T00:00:00Z
  entries:
  - description: Great feature
  - description: Small feature
`

func TestCatalog_LoadFromFile(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	filePath := "test-catalog.yaml"
	if err := ioutil.WriteFile(filePath, []byte(testCatalogYaml), 0666); err != nil {
		t.Errorf("Expected ioutil.WriteFile() to succeed; got %v", err)
	}

	var catalog Catalog
	if err := catalog.LoadFromFile(filePath); err != nil {
		t.Errorf("Expected catalog.LoadFromFile(%q) to succeed; got %v", filePath, err)
	}

	if !reflect.DeepEqual(catalog, testCatalog) {
		t.Errorf("Expected catalog to be:\n%v\ngot:\n%v", testCatalog, catalog)
	}
}

func TestCatalog_SaveToFile(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	filePath := "test-catalog.yaml"
	if err := testCatalog.SaveToFile(filePath); err != nil {
		t.Errorf("Expected testCatalog.SaveToFile() to succeed; got %v", err)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Expected to read catalog yaml; got %v", err)
	}
	if string(data) != testCatalogYaml {
		t.Errorf("Catalog yaml does not match expectations.\n--- want: ---\n%s------\n--- got: ---\n%s------", testCatalogYaml, string(data))
	}
}

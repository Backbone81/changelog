package changelog

import (
	"io/ioutil"
	"reflect"
	"testing"
)

var testEntry = Entry{
	Description: "Some feature",
}

const testEntryYaml = `description: Some feature
`

func TestEntry_LoadFromFile(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	filePath := "test-entry.yaml"
	if err := ioutil.WriteFile(filePath, []byte(testEntryYaml), 0666); err != nil {
		t.Errorf("Expected ioutil.WriteFile() to succeed; got %v", err)
	}

	var entry Entry
	if err := entry.LoadFromFile(filePath); err != nil {
		t.Errorf("Expected entry.LoadFromFile() to succeed; got %v", err)
	}

	if !reflect.DeepEqual(entry, testEntry) {
		t.Errorf("Expected entry to be:\n%v\ngot:\n%v", testEntry, entry)
	}
}

func TestEntry_SaveToFile(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	filePath := "test-entry.yaml"
	if err := testEntry.SaveToFile(filePath); err != nil {
		t.Errorf("Expected testEntry.SaveToFile() to succeed; got %v", err)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Expected to read entry yaml; got %v", err)
	}
	if string(data) != testEntryYaml {
		t.Errorf("Entry yaml does not match expectations.\n--- want: ---\n%s------\n--- got: ---\n%s------", testEntryYaml, string(data))
	}
}

package changelog_test

import (
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/changelog/internal/changelog"
)

var testCatalog = changelog.Catalog{
	Unreleased: []changelog.Entry{
		{
			Description: "Some feature",
		},
		{
			Description: "Some bugfix",
		},
	},
	Releases: []changelog.Release{
		{
			Version: "v1.5.3",
			Date:    time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC),
			Entries: []changelog.Entry{
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

var _ = Describe("changelog.Catalog", func() {
	var ctx TestContext
	const filePath = "test-catalog.yaml"

	BeforeEach(func() {
		ctx.Init()
	})

	AfterEach(func() {
		ctx.Shutdown()
	})

	It("should load from file", func() {
		Expect(os.WriteFile(filePath, []byte(testCatalogYaml), 0o664)).Error().ToNot(HaveOccurred(), //nolint:gosec // Those permissions are fine
			"File path %q", filePath)

		var catalog changelog.Catalog
		Expect(catalog.LoadFromFile(filePath)).Error().ToNot(HaveOccurred(),
			"File path %q", filePath)
		Expect(reflect.DeepEqual(catalog, testCatalog)).To(BeTrue())
	})

	It("should save to file", func() {
		Expect(testCatalog.SaveToFile(filePath)).Error().ToNot(HaveOccurred(),
			"File path %q", filePath)

		data, err := os.ReadFile(filePath)
		Expect(err).ToNot(HaveOccurred(),
			"File path %q", filePath)
		Expect(string(data)).To(Equal(testCatalogYaml))
	})
})

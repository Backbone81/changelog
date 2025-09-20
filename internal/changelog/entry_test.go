package changelog_test

import (
	"os"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/changelog/internal/changelog"
)

var testEntry = changelog.Entry{
	Description: "Some feature",
}

const testEntryYaml = `description: Some feature
`

var _ = Describe("changelog.Entry", func() {
	var ctx TestContext
	const filePath = "test-entry.yaml"

	BeforeEach(func() {
		ctx.Init()
	})

	AfterEach(func() {
		ctx.Shutdown()
	})

	It("should load from file", func() {
		Expect(os.WriteFile(filePath, []byte(testEntryYaml), 0o664)).NotTo(HaveOccurred(), //nolint:gosec // Those permissions are fine
			"File path %q", filePath)

		var entry changelog.Entry
		Expect(entry.LoadFromFile(filePath)).NotTo(HaveOccurred(),
			"File path %q", filePath)
		Expect(reflect.DeepEqual(entry, testEntry)).To(BeTrue())
	})

	It("should save to file", func() {
		Expect(testEntry.SaveToFile(filePath)).NotTo(HaveOccurred(),
			"File path %q", filePath)

		data, err := os.ReadFile(filePath)
		Expect(err).NotTo(HaveOccurred(),
			"File path %q", filePath)
		Expect(string(data)).To(Equal(testEntryYaml))
	})
})

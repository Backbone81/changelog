package changelog_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gitlab.com/backbone/changelog/internal/changelog"
	"io/ioutil"
	"reflect"
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
		Expect(ioutil.WriteFile(filePath, []byte(testEntryYaml), 0666)).NotTo(HaveOccurred(),
			"File path %q", filePath)

		var entry changelog.Entry
		Expect(entry.LoadFromFile(filePath)).NotTo(HaveOccurred(),
			"File path %q", filePath)
		Expect(reflect.DeepEqual(entry, testEntry)).To(BeTrue())
	})

	It("should save to file", func() {
		Expect(testEntry.SaveToFile(filePath)).NotTo(HaveOccurred(),
			"File path %q", filePath)

		data, err := ioutil.ReadFile(filePath)
		Expect(err).NotTo(HaveOccurred(),
			"File path %q", filePath)
		Expect(string(data)).To(Equal(testEntryYaml))
	})
})

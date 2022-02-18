package changelog_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"
)

var _ = Describe("changelog.Changelog", func() {
	var ctx TestContext

	BeforeEach(func() {
		ctx.Init()
	})

	AfterEach(func() {
		ctx.Shutdown()
	})

	It("should init", func() {
		clog := changelog.New()
		Expect(clog.Init()).Error().ToNot(HaveOccurred())

		Expect(os.Stat(changelog.ConfigFilename)).Error().ToNot(HaveOccurred())
		Expect(os.Stat(clog.GetChangelogDir())).Error().ToNot(HaveOccurred())
		Expect(os.Stat(clog.GetUnreleasedDir())).Error().ToNot(HaveOccurred())
		Expect(os.Stat(clog.GetUnreleasedKeepFilePath())).Error().ToNot(HaveOccurred())
		Expect(os.Stat(clog.GetCatalogFilePath())).Error().ToNot(HaveOccurred())
	})

	Describe("initialized changelog", func() {
		var clog *changelog.Changelog

		BeforeEach(func() {
			clog = ctx.PrepareChangelog()
		})

		It("should increment the catalog size on add", func() {
			beforeUnreleased := len(clog.Catalog.Unreleased)
			Expect(clog.Add("some new unreleased changelog")).Error().ToNot(HaveOccurred())
			Expect(len(clog.Catalog.Unreleased)).To(Equal(beforeUnreleased + 1))
		})

		It("should increment the file count on add", func() {
			beforeUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(clog.Add("some new unreleased changelog")).Error().ToNot(HaveOccurred())

			afterUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(len(afterUnreleasedFiles)).To(Equal(len(beforeUnreleasedFiles) + 1))
		})

		It("should decrement the catalog size on release", func() {
			Expect(len(clog.Catalog.Unreleased)).To(BeNumerically(">", 0))
			Expect(clog.CreateRelease("v1.5.2")).Error().ToNot(HaveOccurred())
			Expect(len(clog.Catalog.Unreleased)).To(BeNumerically("==", 0))
		})

		It("should decrement the file count on release", func() {
			beforeUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())
			wantFiles := len(beforeUnreleasedFiles) - len(clog.Catalog.Unreleased)

			Expect(clog.CreateRelease("v1.5.2")).Error().ToNot(HaveOccurred())

			afterUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(len(afterUnreleasedFiles)).To(BeNumerically("==", wantFiles))
		})

		It("should render the changelog file", func() {
			Expect(os.Remove(clog.GetChangelogFilePath())).Error().ToNot(HaveOccurred())
			Expect(os.Stat(clog.GetChangelogFilePath())).Error().To(HaveOccurred())
			Expect(clog.RenderChangelog()).Error().ToNot(HaveOccurred())
			Expect(os.Stat(clog.GetChangelogFilePath())).Error().ToNot(HaveOccurred())
		})
	})
})

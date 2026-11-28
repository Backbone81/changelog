package changelog_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/changelog/internal/changelog"
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

	It("should not overwrite existing files on init", func() {
		const testContent = `some file content not relevant to the project`
		clog := changelog.New()

		Expect(os.WriteFile(changelog.ConfigFilename, []byte(testContent), 0o664)).Error().ToNot(HaveOccurred())         //nolint:gosec // Those permissions are fine
		Expect(os.MkdirAll(filepath.Dir(clog.GetUnreleasedKeepFilePath()), 0o777)).Error().ToNot(HaveOccurred())         //nolint:gosec // This permission is fine
		Expect(os.WriteFile(clog.GetUnreleasedKeepFilePath(), []byte(testContent), 0o664)).Error().ToNot(HaveOccurred()) //nolint:gosec // Those permissions are fine
		Expect(os.MkdirAll(filepath.Dir(clog.GetCatalogFilePath()), 0o777)).Error().ToNot(HaveOccurred())                //nolint:gosec // This permission is fine
		Expect(os.WriteFile(clog.GetCatalogFilePath(), []byte(testContent), 0o664)).Error().ToNot(HaveOccurred())        //nolint:gosec // Those permissions are fine

		Expect(clog.Init()).Error().ToNot(HaveOccurred())

		Expect(os.ReadFile(changelog.ConfigFilename)).To(Equal([]byte(testContent)))
		Expect(os.ReadFile(clog.GetUnreleasedKeepFilePath())).To(Equal([]byte(testContent)))
		Expect(os.ReadFile(clog.GetCatalogFilePath())).To(Equal([]byte(testContent)))
	})

	Describe("initialized changelog", func() {
		var clog *changelog.Changelog

		BeforeEach(func() {
			clog = ctx.PrepareChangelog()
		})

		It("should increment the catalog size on add", func() {
			beforeUnreleased := len(clog.Catalog.Unreleased)
			Expect(clog.Add("some new unreleased changelog", 0)).Error().ToNot(HaveOccurred())
			Expect(clog.Catalog.Unreleased).To(HaveLen(beforeUnreleased + 1))
		})

		It("should increment the file count on add", func() {
			beforeUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(clog.Add("some new unreleased changelog", 0)).Error().ToNot(HaveOccurred())

			afterUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(afterUnreleasedFiles).To(HaveLen(len(beforeUnreleasedFiles) + 1))
		})

		It("should decrement the catalog size on release", func() {
			Expect(clog.Catalog.Unreleased).ToNot(BeEmpty())
			Expect(clog.CreateRelease("v1.5.2")).Error().ToNot(HaveOccurred())
			Expect(clog.Catalog.Unreleased).To(BeEmpty())
		})

		It("should decrement the file count on release", func() {
			beforeUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())
			wantFiles := len(beforeUnreleasedFiles) - len(clog.Catalog.Unreleased)

			Expect(clog.CreateRelease("v1.5.2")).Error().ToNot(HaveOccurred())

			afterUnreleasedFiles, err := os.ReadDir(clog.GetUnreleasedDir())
			Expect(err).ToNot(HaveOccurred())

			Expect(afterUnreleasedFiles).To(HaveLen(wantFiles))
		})

		It("should render the changelog file", func() {
			Expect(os.Remove(clog.GetChangelogFilePath())).Error().ToNot(HaveOccurred())
			Expect(os.Stat(clog.GetChangelogFilePath())).Error().To(HaveOccurred())
			Expect(clog.RenderChangelog()).Error().ToNot(HaveOccurred())
			Expect(os.Stat(clog.GetChangelogFilePath())).Error().ToNot(HaveOccurred())
		})
	})

	Describe("NextTag", func() {
		var clog *changelog.Changelog

		BeforeEach(func() {
			clog = changelog.New()
			Expect(clog.Init()).NotTo(HaveOccurred())
			Expect(clog.Load()).NotTo(HaveOccurred())
		})

		It("should fail without unreleased changes", func() {
			Expect(clog.NextTag()).Error().To(BeEquivalentTo(changelog.ErrNoUnreleasedChanges))
		})

		It("should fail with changes without severity", func() {
			Expect(clog.Add("another big feature", 0)).To(Succeed())
			Expect(clog.NextTag()).Error().To(BeEquivalentTo(changelog.ErrEntriesWithoutSeverityFound))
		})

		It("should fail with malformed last tag (invalid numbers)", func() {
			Expect(clog.CreateRelease("vx.x.x")).To(Succeed())
			Expect(clog.Add("another big feature", changelog.SeverityPatch)).To(Succeed())
			Expect(clog.NextTag()).Error().To(BeEquivalentTo(changelog.ErrLastTagMalformed))
		})

		It("should fail with malformed last tag (invalid format)", func() {
			Expect(clog.CreateRelease("v0.0")).To(Succeed())
			Expect(clog.Add("another big feature", changelog.SeverityPatch)).To(Succeed())
			Expect(clog.NextTag()).Error().To(BeEquivalentTo(changelog.ErrLastTagMalformed))
		})

		It("should succeed from scratch with patch", func() {
			Expect(clog.Add("another big feature", changelog.SeverityPatch)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v0.0.1"))
		})

		It("should succeed from scratch with minor", func() {
			Expect(clog.Add("another big feature", changelog.SeverityMinor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v0.1.0"))
		})

		It("should succeed from scratch with major", func() {
			Expect(clog.Add("another big feature", changelog.SeverityMajor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v1.0.0"))
		})

		It("should succeed with preceding patch tag", func() {
			Expect(clog.Add("another big feature", changelog.SeverityPatch)).To(Succeed())
			Expect(clog.CreateRelease("v0.0.1")).Error().ToNot(HaveOccurred())
			Expect(clog.Add("another big feature", changelog.SeverityPatch)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v0.0.2"))
		})

		It("should succeed with preceding minor tag", func() {
			Expect(clog.Add("another big feature", changelog.SeverityMinor)).To(Succeed())
			Expect(clog.CreateRelease("v0.1.0")).Error().ToNot(HaveOccurred())
			Expect(clog.Add("another big feature", changelog.SeverityMinor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v0.2.0"))
		})

		It("should succeed with preceding major tag", func() {
			Expect(clog.Add("another big feature", changelog.SeverityMajor)).To(Succeed())
			Expect(clog.CreateRelease("v1.0.0")).Error().ToNot(HaveOccurred())
			Expect(clog.Add("another big feature", changelog.SeverityMajor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v2.0.0"))
		})

		It("should succeed with minor while resetting patch to 0", func() {
			Expect(clog.Add("another big feature", 0)).To(Succeed())
			Expect(clog.CreateRelease("v1.5.3")).Error().ToNot(HaveOccurred())
			Expect(clog.Add("another big feature", changelog.SeverityMinor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v1.6.0"))
		})

		It("should succeed with major while resetting minor and patch to 0", func() {
			Expect(clog.Add("another big feature", 0)).To(Succeed())
			Expect(clog.CreateRelease("v1.5.3")).Error().ToNot(HaveOccurred())
			Expect(clog.Add("another big feature", changelog.SeverityMajor)).To(Succeed())
			tag, err := clog.NextTag()
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v2.0.0"))
		})
	})
})

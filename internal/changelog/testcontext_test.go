package changelog_test

import (
	"os"

	. "github.com/onsi/gomega"

	"github.com/backbone81/changelog/internal/changelog"
)

func init() {
	// change the default file pattern to include fractional seconds
	// otherwise we would overwrite the same file if we create multiple unreleased changelog entries
	// in quick succession
	changelog.DefaultConfig.UnreleasedFilePattern += ".000000"
}

type TestContext struct {
	TempDir       string
	WorkdirBackup string
}

func (tc *TestContext) Init() {
	tempDir, err := os.MkdirTemp("", "changelog")
	Expect(err).NotTo(HaveOccurred())
	tc.TempDir = tempDir

	currDir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	tc.WorkdirBackup = currDir

	Expect(os.Chdir(tc.TempDir)).NotTo(HaveOccurred())
}

func (tc *TestContext) Shutdown() {
	Expect(os.Chdir(tc.WorkdirBackup)).NotTo(HaveOccurred())
	Expect(os.RemoveAll(tc.TempDir)).NotTo(HaveOccurred())
}

func (tc *TestContext) PrepareChangelog() *changelog.Changelog {
	clog := changelog.New()
	Expect(clog.Init()).NotTo(HaveOccurred())
	Expect(clog.Load()).NotTo(HaveOccurred())

	Expect(clog.Add("some new feature", 0)).NotTo(HaveOccurred())
	Expect(clog.Add("some bugfix", 0)).NotTo(HaveOccurred())
	Expect(clog.CreateRelease("v1.4.3")).NotTo(HaveOccurred())

	Expect(clog.Add("some other feature", 0)).NotTo(HaveOccurred())
	Expect(clog.Add("some other bugfix", 0)).NotTo(HaveOccurred())
	Expect(clog.RenderChangelog()).NotTo(HaveOccurred())

	return clog
}

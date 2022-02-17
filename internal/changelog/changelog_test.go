package changelog

import (
	"fmt"
	"os"
	"testing"
)

type TestContext struct {
	TempDir       string
	WorkdirBackup string
}

func init() {
	// change the default file pattern to include fractional seconds
	// otherwise we would overwrite the same file if we create multiple unreleased changelog entries
	// in quick succession
	DefaultConfig.UnreleasedFilePattern = DefaultConfig.UnreleasedFilePattern + ".000000"
}

func (tc *TestContext) Init() {
	tempDir, err := os.MkdirTemp("", "changelog")
	if err != nil {
		fmt.Println("ERROR: Failed to create temporary directory for running tests.")
		fmt.Println(err)
		os.Exit(1)
	}
	tc.TempDir = tempDir

	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: Failed to get current working directory.")
		fmt.Println(err)
		os.Exit(1)
	}
	tc.WorkdirBackup = currDir

	if err := os.Chdir(tc.TempDir); err != nil {
		fmt.Println("ERROR: Failed to change working directory to temporary directory.")
		fmt.Println(err)
		os.Exit(1)
	}
}

func (tc *TestContext) Shutdown() {
	if err := os.Chdir(tc.WorkdirBackup); err != nil {
		fmt.Println("ERROR: Failed to change working directory back to original.")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.RemoveAll(tc.TempDir); err != nil {
		fmt.Println("ERROR: Failed to cleanup temporary directory.")
		fmt.Println(err)
		os.Exit(1)
	}
}

func (tc *TestContext) PrepareChangelog(t *testing.T) *Changelog {
	clog := New()
	if err := clog.Init(); err != nil {
		t.Errorf("Expected clog.Init() to succeed; got %v", err)
	}
	if err := clog.Load(); err != nil {
		t.Errorf("Expected clog.Load() to succeed; got %v", err)
	}

	if err := clog.Add("some new feature"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}
	if err := clog.Add("some bugfix"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}

	if err := clog.CreateRelease("v1.4.3"); err != nil {
		t.Errorf("Expected clog.CreateRelease() to succeed; got %v", err)
	}

	if err := clog.Add("some other feature"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}
	if err := clog.Add("some other bugfix"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}

	if err := clog.RenderChangelog(); err != nil {
		t.Errorf("Expected clog.RenderChangelog() to succeed; got %v", err)
	}

	return clog
}

func TestChangelog_Init(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := New()
	err := clog.Init()
	if err != nil {
		t.Errorf("Expected changelog.Init() to succeed")
	}

	if _, err := os.Stat(ConfigFilename); err != nil {
		t.Errorf("Expected %s to be present in working directory", ConfigFilename)
	}

	if _, err := os.Stat(clog.getChangelogDir()); err != nil {
		t.Errorf("Expected %s to be present in working directory", clog.getChangelogDir())
	}

	if _, err := os.Stat(clog.getUnreleasedDir()); err != nil {
		t.Errorf("Expected %s to be present in working directory", clog.getUnreleasedDir())
	}

	if _, err := os.Stat(clog.getUnreleasedKeepFilePath()); err != nil {
		t.Errorf("Expected %s to be present in working directory", clog.getUnreleasedKeepFilePath())
	}

	if _, err := os.Stat(clog.getCatalogFilePath()); err != nil {
		t.Errorf("Expected %s to be present in working directory", clog.getCatalogFilePath())
	}
}

func TestChangelog_Add_IncreaseCatalog(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := ctx.PrepareChangelog(t)

	// get current state for catalog count
	beforeUnreleased := len(clog.Catalog.Unreleased)

	// add new entry
	if err := clog.Add("some new unreleased changelog"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}

	// get current state for catalog count
	gotUnreleased := len(clog.Catalog.Unreleased)
	wantUnreleased := beforeUnreleased + 1

	// we want the catalog count to be one more than before the add
	if gotUnreleased != wantUnreleased {
		t.Errorf("Expected unreleased catalog entries to be %d; got %d", wantUnreleased, gotUnreleased)
	}
}

func TestChangelog_Add_IncreaseFiles(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := ctx.PrepareChangelog(t)

	// get current state for file count
	beforeUnreleasedFiles, err := os.ReadDir(clog.getUnreleasedDir())
	if err != nil {
		t.Errorf("Failed to get file entries from unreleased dir; %v", err)
	}

	// add new entry
	if err := clog.Add("some new unreleased changelog"); err != nil {
		t.Errorf("Expected clog.Add() to succeed; got %v", err)
	}

	// get current state for file count
	afterUnreleasedFiles, err := os.ReadDir(clog.getUnreleasedDir())
	if err != nil {
		t.Errorf("Failed to get file entries from unreleased dir; %v", err)
	}
	wantUnreleasedFiles := len(beforeUnreleasedFiles) + 1

	// we want the file count to be one more than before the add
	if len(afterUnreleasedFiles) != wantUnreleasedFiles {
		t.Errorf("Expected unreleased file entries to be %d; got %d", wantUnreleasedFiles, len(afterUnreleasedFiles))
	}
}

func TestChangelog_CreateRelease_DecreaseCatalog(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := ctx.PrepareChangelog(t)

	beforeUnreleased := len(clog.Catalog.Unreleased)
	if beforeUnreleased <= 0 {
		t.Errorf("Expected unreleased changelogs to be present; got %d", beforeUnreleased)
	}

	if err := clog.CreateRelease("v1.5.2"); err != nil {
		t.Errorf("Expected clog.CreateRelease() to succeed; got %v", err)
	}

	afterUnreleased := len(clog.Catalog.Unreleased)
	if afterUnreleased != 0 {
		t.Errorf("Expected no unreleased changelogs to be present; got %d", afterUnreleased)
	}
}

func TestChangelog_CreateRelease_DecreaseFiles(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := ctx.PrepareChangelog(t)

	beforeUnreleasedFiles, err := os.ReadDir(clog.getUnreleasedDir())
	if err != nil {
		t.Errorf("Failed to get file entries from unreleased dir; %v", err)
	}
	wantFiles := len(beforeUnreleasedFiles) - len(clog.Catalog.Unreleased)

	if err := clog.CreateRelease("v1.5.2"); err != nil {
		t.Errorf("Expected clog.CreateRelease() to succeed; got %v", err)
	}

	afterUnreleasedFiles, err := os.ReadDir(clog.getUnreleasedDir())
	if err != nil {
		t.Errorf("Failed to get file entries from unreleased dir; %v", err)
	}

	if len(afterUnreleasedFiles) != wantFiles {
		t.Errorf("Expected unreleased changelog files to be %d; got %d", wantFiles, len(afterUnreleasedFiles))
	}
}

func TestChangelog_RenderChangelog(t *testing.T) {
	ctx := TestContext{}
	ctx.Init()
	defer ctx.Shutdown()

	clog := ctx.PrepareChangelog(t)

	if err := os.Remove(clog.getChangelogFilePath()); err != nil {
		t.Errorf("Failed to remove %s: got %v", clog.getChangelogFilePath(), err)
	}
	if _, err := os.Stat(clog.getChangelogFilePath()); err == nil {
		t.Errorf("Expected %s to be missing in working directory", clog.getChangelogFilePath())
	}
	if err := clog.RenderChangelog(); err != nil {
		t.Errorf("Expected clog.RenderChangelog() to succeed; got %v", err)
	}
	if _, err := os.Stat(clog.getChangelogFilePath()); err != nil {
		t.Errorf("Expected %s to be present in working directory", clog.getChangelogFilePath())
	}
}

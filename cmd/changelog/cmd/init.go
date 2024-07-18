package cmd

import (
	"fmt"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new changelog environment",
	Long: fmt.Sprintf(`Initializes a new changelog environment by creating the config
file %q and creating the subdirectory %q to hold the changelogs.`, changelog.ConfigFilename, changelog.DefaultConfig.ChangelogDir),
	Run: func(cmd *cobra.Command, args []string) {
		clogInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func clogInit() {
	clog := changelog.New()
	err := clog.Init()
	if err != nil {
		fmt.Println("ERROR: Failed to initialize changelog.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully initialized changelog.")
}

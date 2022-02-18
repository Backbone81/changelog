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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

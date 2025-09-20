package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/backbone81/changelog/internal/changelog"
)

// updateCmd represents the update command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the end user changelog file with the current state",
	Long: fmt.Sprintf(`Updates the %q file with the current state of all releases and unreleased changelogs.`,
		changelog.DefaultConfig.ChangelogFileName),
	Run: func(cmd *cobra.Command, args []string) {
		clogUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func clogUpdate() {
	clog := changelog.New()
	if err := clog.Load(); err != nil {
		fmt.Println("ERROR: Failed to load changelog. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}
	if err := clog.SaveCatalog(); err != nil {
		fmt.Println("ERROR: Failed to save changelog catalog.")
		fmt.Println(err)
		os.Exit(1)
	}
	if err := clog.RenderChangelog(); err != nil {
		fmt.Println("ERROR: Failed to render changelog.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated changelog.")
}

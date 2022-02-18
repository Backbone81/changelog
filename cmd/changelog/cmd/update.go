/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

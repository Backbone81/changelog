package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of this tool.",
	Long:  `Displays the version of this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Changelog v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

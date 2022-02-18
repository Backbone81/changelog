package cmd

import (
	"bufio"
	"fmt"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new changelog entry",
	Long:  `Adds a new changelog entry as unreleased changelog.`,
	Run: func(cmd *cobra.Command, args []string) {
		clogAdd()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clogAdd() {
	clog := changelog.New()
	if err := clog.Load(); err != nil {
		fmt.Println("ERROR: Failed to load changelog. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print("New changelog description: ")
	reader := bufio.NewReader(os.Stdin)
	description, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERROR: Failed to add new changelog entry.")
		fmt.Println(err)
		os.Exit(1)
	}

	err = clog.Add(description)
	if err != nil {
		fmt.Println("ERROR: Failed to add new changelog entry.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully added new changelog entry.")
}

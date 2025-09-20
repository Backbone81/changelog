package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/backbone81/changelog/internal/changelog"
)

// addCmd represents the add command.
var (
	description string
	addCmd      = &cobra.Command{
		Use:   "add",
		Short: "Adds a new changelog entry",
		Long:  `Adds a new changelog entry as unreleased changelog.`,
		Run: func(cmd *cobra.Command, args []string) {
			clogAdd()
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"A string describing the change for the end user.")
}

func clogAdd() {
	clog := changelog.New()
	if err := clog.Load(); err != nil {
		fmt.Println("ERROR: Failed to load changelog. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}

	if description == "" {
		fmt.Print("New changelog description: ")
		reader := bufio.NewReader(os.Stdin)
		stdinDescription, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERROR: Failed to add new changelog entry.")
			fmt.Println(err)
			os.Exit(1)
		}
		description = stdinDescription
	}

	err := clog.Add(description)
	if err != nil {
		fmt.Println("ERROR: Failed to add new changelog entry.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully added new changelog entry.")
}

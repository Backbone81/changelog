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
	severity    string
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
	addCmd.Flags().StringVarP(
		&severity,
		"severity",
		"s",
		"",
		"An enum describing the severity of the change. (major|minor|patch)")
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

	addingSeverity := changelog.SeverityUnknown
	if severity != "" {
		var err error
		addingSeverity, err = changelog.ParseSeverity(severity)
		if err != nil {
			fmt.Println("ERROR: Failed to add new changelog entry.")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	err := clog.Add(description, addingSeverity)
	if err != nil {
		fmt.Println("ERROR: Failed to add new changelog entry.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully added new changelog entry.")
}

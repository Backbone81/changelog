package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/backbone81/changelog/internal/changelog"
)

// releaseCmd represents the release command.
var (
	releaseTag string
	autoTag    bool
	releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Creates a new release",
		Long:  `Creates a new release and consolidates all unreleased changelogs into it.`,
		Run: func(cmd *cobra.Command, args []string) {
			clogRelease()
		},
	}
)

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVarP(
		&releaseTag,
		"tag",
		"t",
		"",
		"The version string to use for the release (i.e. 'v1.5.2')")
	releaseCmd.Flags().BoolVarP(
		&autoTag,
		"auto-tag",
		"a",
		false,
		"When set to true, tag will be ignored and changelog will generate the tag by using the severities from the changelogs.")
}

func clogRelease() {
	clog := changelog.New()
	if err := clog.Load(); err != nil {
		fmt.Println("ERROR: Failed to load changelog. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}

	if autoTag {
		tag, err := clog.NextTag()
		if err != nil {
			fmt.Println("ERROR: Failed to create a changelog release.")
			fmt.Println(err)
			os.Exit(1)
		}
		releaseTag = tag
	} else if releaseTag == "" {
		fmt.Print("Release version (i.e. 'v1.5.3'): ")
		reader := bufio.NewReader(os.Stdin)
		stdioReleaseTag, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERROR: Failed to create a changelog release.")
			fmt.Println(err)
			os.Exit(1)
		}
		releaseTag = stdioReleaseTag
	}

	if err := clog.CreateRelease(releaseTag); err != nil {
		fmt.Println("ERROR: Failed to create a changelog release.")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := clog.RenderChangelog(); err != nil {
		fmt.Println("ERROR: Failed to render changelog.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully created release.")
}

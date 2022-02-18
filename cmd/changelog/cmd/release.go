package cmd

import (
	"bufio"
	"fmt"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseTag string
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Creates a new release",
	Long:  `Creates a new release and consolidates all unreleased changelogs into it.`,
	Run: func(cmd *cobra.Command, args []string) {
		clogRelease()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseCmd.Flags().StringVarP(
		&releaseTag,
		"tag",
		"t",
		"",
		"The version string to use for the release (i.e. 'v1.5.2')")
}

func clogRelease() {
	clog := changelog.New()
	if err := clog.Load(); err != nil {
		fmt.Println("ERROR: Failed to load changelog. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}

	if releaseTag == "" {
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

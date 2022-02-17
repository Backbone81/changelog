package main

import (
	"bufio"
	"fmt"
	"gitlab.com/backbone/changelog/internal/changelog"
	"os"
)

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

func clogAdd() {
	clog := changelog.New()
	if err := clog.LoadConfig(); err != nil {
		fmt.Println("ERROR: Failed to load changelog config. Did you run init?")
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

func clogUpdate() {
	clog := changelog.New()
	if err := clog.LoadConfig(); err != nil {
		fmt.Println("ERROR: Failed to load changelog config. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}
	if err := clog.LoadCatalog(); err != nil {
		fmt.Println("ERROR: Failed to load changelog catalog.")
		fmt.Println(err)
		os.Exit(1)
	}
	if err := clog.LoadUnreleased(); err != nil {
		fmt.Println("ERROR: Failed to load unreleased changelog entries.")
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

func clogRelease() {
	clog := changelog.New()
	if err := clog.LoadConfig(); err != nil {
		fmt.Println("ERROR: Failed to load changelog config. Did you run init?")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := clog.LoadCatalog(); err != nil {
		fmt.Println("ERROR: Failed to load changelog catalog.")
		fmt.Println(err)
		os.Exit(1)
	}
	if err := clog.LoadUnreleased(); err != nil {
		fmt.Println("ERROR: Failed to load unreleased changelog entries.")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print("Release version (i.e. 'v1.5.3'): ")
	reader := bufio.NewReader(os.Stdin)
	version, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERROR: Failed to create a changelog release.")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := clog.CreateRelease(version); err != nil {
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Expected 'init', 'new', 'update' or 'release' subcommands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		clogInit()
	case "add":
		clogAdd()
	case "update":
		clogUpdate()
	case "release":
		clogRelease()
	}
}

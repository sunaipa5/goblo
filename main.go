package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	SiteUrl          string            `json:"site-url"`
	SiteTitle        string            `json:"site-title"`
	TemplateSettings map[string]string `json:"template-settings"`
}

var projectPath string

var config Config

func main() {
	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println(`Goblo v1.0.3

+----------+-----------------------------------------+
| Commands |                  Usage                  |
+----------+-----------------------------------------+
| create   | ./project-dir                           |
| build    | ./project-dir                           |
| live     |  ./project-dir                          |
+----------+-----------------------------------------+

Live command flags:
+--------------+----------------------------+---------+
|    Flags     |            Usage           | Default |
+--------------+----------------------------+---------+
| --port=1111  | live port                  | 5151    |
| --open=false | open browser automatically | true    |
+--------------+----------------------------+---------+`)
		fmt.Fprintln(os.Stderr, "Usage: goblo <command> [path]")
		fmt.Fprintln(os.Stderr, "Commands: build, create")
		os.Exit(1)
	}

	cmd := os.Args[1]

	path := "."
	if len(os.Args) >= 3 {
		path = os.Args[2]
	}

	switch cmd {
	case "build":
		buildFlags := flag.NewFlagSet("build", flag.ExitOnError)

		args := []string{}
		if len(os.Args) > 3 {
			args = os.Args[3:]
		}
		buildFlags.Parse(args)

		projectPath = path

		if err := build(); err != nil {
			werror("Failed to build", err)
			return
		}

	case "create":
		createFlags := flag.NewFlagSet("create", flag.ExitOnError)

		args := []string{}
		if len(os.Args) > 3 {
			args = os.Args[3:]
		}
		createFlags.Parse(args)

		projectPath = path
		if err := create(path); err != nil {
			werror("Failed to create project", err)
			return
		}
	case "live":
		liveFlags := flag.NewFlagSet("live", flag.ExitOnError)

		port := liveFlags.Int("port", 5151, "Port to run live server")
		openb := liveFlags.Bool("open", true, "Open browser automatically")

		err := liveFlags.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
			os.Exit(1)
		}

		args := liveFlags.Args()
		if len(args) > 0 {
			projectPath = args[0]
		} else {
			projectPath = "."
		}

		if err := build(); err != nil {
			werror("Failed to build", err)
			return
		}

		if err := live(projectPath, *port, *openb); err != nil {
			werror("Failed to start live server", err)
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Done %s\n", elapsed)
}

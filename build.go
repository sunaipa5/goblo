package main

import (
	"fmt"
	"os"
	"path/filepath"

	json "github.com/json-iterator/go"
)

func build() {
	checkPaths := []string{
		"posts",
		"static",
		"goblo.json",
		"index.html",
	}

	for _, p := range checkPaths {
		if typ, err := checkPath(filepath.Join(projectPath, p)); typ < 3 {
			if err != nil {
				werror("Failed to check '"+p+"'", err)
			} else {
				fmt.Println("Please create a project! You can use create command.")
			}
			return
		}
	}

	//Init config
	configByte, err := os.ReadFile(filepath.Join(projectPath, "goblo.json"))
	if err != nil {
		werror("Failed to read config file", err)
		return
	}

	if err := json.Unmarshal(configByte, &config); err != nil {
		werror("Failed to parse config file", err)
		return
	}

	//Parse templates
	if err := parse_templates(); err != nil {
		werror("Failed to parse templates", err)
		return
	}

	//Parse posts
	if err := parse_posts(); err != nil {
		werror("Failed to parse posts", err)
		return
	}

	//Copy static files
	err = copyDir(filepath.Join(projectPath, "static"), filepath.Join(projectPath, "dist", "static"), nil)
	if err != nil {
		werror("Failed to copy static dir", err)
		return
	}

	exclude := map[string]bool{
		"goblo.json": true,
		"dist":       true,
		"posts":      true,
		"static":     true,
		"templates":  true,
	}

	err = copyDir(filepath.Join(projectPath), filepath.Join(projectPath, "dist"), exclude)
	if err != nil {
		werror("Failed to copy static dir", err)
		return
	}
}

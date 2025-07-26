package main

import (
	"fmt"
	"os"
	"path/filepath"

	json "github.com/json-iterator/go"
)

func build() error {

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
				return fmt.Errorf("Failed to check %s %v", p, err)
			} else {
				return fmt.Errorf("Please create a project! You can use create command.")
			}
		}
	}

	if err := resetDir(filepath.Join(projectPath, "dist")); err != nil {
		return fmt.Errorf("Failed to clear dist dir %v", err)
	}

	//Init config
	configByte, err := os.ReadFile(filepath.Join(projectPath, "goblo.json"))
	if err != nil {
		return fmt.Errorf("Failed to read config file %v", err)
	}

	if err := json.Unmarshal(configByte, &config); err != nil {
		return fmt.Errorf("Failed to parse config file %v", err)
	}

	//Parse templates
	if err := parse_templates(); err != nil {
		return fmt.Errorf("Failed to parse templates %v", err)
	}

	//Parse posts
	if err := parse_posts(); err != nil {
		return fmt.Errorf("Failed to parse posts %v", err)
	}

	//Copy static files
	err = copyDir(filepath.Join(projectPath, "static"), filepath.Join(projectPath, "dist", "static"), nil)
	if err != nil {
		return fmt.Errorf("Failed to copy static dir %v", err)
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
		return fmt.Errorf("Failed to copy static dir %v", err)
	}

	if err := create_sitemap(); err != nil {
		return fmt.Errorf("Failed to create sitemap %v", err)
	}

	if err := create_feed(); err != nil {
		return fmt.Errorf("Failed to create feed %V", err)
	}
	return nil
}

func live_build() error {
	if err := resetDir(filepath.Join(projectPath, "dist")); err != nil {
		return err
	}

	//Init config
	configByte, err := os.ReadFile(filepath.Join(projectPath, "goblo.json"))
	if err != nil {
		return fmt.Errorf("Failed to read config file %v", err)
	}

	if err := json.Unmarshal(configByte, &config); err != nil {
		return fmt.Errorf("Failed to parse config file %v", err)
	}

	//Parse templates
	if err := parse_templates(); err != nil {
		return fmt.Errorf("Failed to parse templates %v", err)
	}

	//Parse posts
	if err := parse_posts(); err != nil {
		return fmt.Errorf("Failed to parse posts %v", err)
	}

	//Copy static files
	err = copyDir(filepath.Join(projectPath, "static"), filepath.Join(projectPath, "dist", "static"), nil)
	if err != nil {
		return fmt.Errorf("Failed to copy static dir %v", err)
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
		return fmt.Errorf("Failed to copy static dir %v", err)
	}

	return nil
}

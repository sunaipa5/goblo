package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	Before  string
	Content string
	After   string
}

type Posts struct {
	Meta    string
	Content string
}

var templates = make(map[string]Template)

func parse_templates() error {
	tpath := filepath.Join(projectPath, "templates")

	entries, err := os.ReadDir(tpath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			ipath := filepath.Join(tpath, entry.Name(), "index.tmpl")

			fileByte, err := os.ReadFile(ipath)
			if err != nil {
				return err
			}

			fileContent := string(fileByte)

			parts := strings.Split(fileContent, "{{content}}")
			if len(parts) != 2 {
				return fmt.Errorf("Failed to find {{content}} %s please read doc.", ipath)
			}

			tpath := filepath.Join(tpath, entry.Name())

			before := readIncludes(tpath, parts[0])
			after := readIncludes(tpath, parts[1])

			newTemplate := Template{
				Before: before,
				After:  after,
			}

			templates[entry.Name()] = newTemplate
		}
	}

	return nil
}

func parse_posts() error {
	ppath := filepath.Join(projectPath, "posts")

	entries, err := os.ReadDir(ppath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".html") {
			htmlByte, err := os.ReadFile(filepath.Join(ppath, name))
			if err != nil {
				return err
			}

			headTags, content := extractHeadAndContent(string(htmlByte))

			var parsedContent string
			if v, ok := config[name]; ok {
				if vt, okt := templates[v.Template]; okt {
					pstr, err := inject_post(vt.Before, headTags)
					if err != nil {
						return fmt.Errorf("%s %v", v.Template, err)
					}
					vt.Before = pstr
					vt.Content = content
					parsedContent = vt.Before + vt.Content + vt.After
				}
			} else if v, ok := config["*"]; ok {
				if vt, okt := templates[v.Template]; okt {
					pstr, err := inject_post(vt.Before, headTags)
					if err != nil {
						return fmt.Errorf("%s template %v", v.Template, err)
					}
					vt.Before = pstr
					vt.Content = content
					parsedContent = vt.Before + vt.Content + vt.After
				}
			} else {
				return fmt.Errorf("Could not find any template %s", name)
			}

			lastContent := removeEmptyLinesFast(parsedContent)

			ppath := filepath.Join(projectPath, "dist", "posts")

			if err := os.MkdirAll(ppath, os.ModePerm); err != nil {
				return err
			}

			if err := os.WriteFile(filepath.Join(ppath, name), []byte(lastContent), 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

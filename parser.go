package main

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Template struct {
	Before  string
	Content string
	After   string
}

var templates = make(map[string]Template)
var postsFeed []Entry

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

	popath := filepath.Join(projectPath, "dist", "posts")
	if err := os.MkdirAll(popath, os.ModePerm); err != nil {
		return err
	}

	sem := make(chan struct{}, 8)
	errCh := make(chan error, len(entries))
	var wg sync.WaitGroup

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		sem <- struct{}{}
		wg.Add(1)
		go func(e fs.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			name := entry.Name()
			if strings.HasSuffix(name, ".html") {
				htmlByte, err := os.ReadFile(filepath.Join(ppath, name))
				if err != nil {
					errCh <- err
					return
				}

				headTags, content := extractHeadAndContent(string(htmlByte))

				t, dat, des, img := findMeta(headTags)

				var imgPath string
				if img != "" {
					imgPath, err = url.JoinPath(config.SiteUrl, img)
					if err != nil {
						errCh <- err
						return
					}
				}

				link, err := url.JoinPath(config.SiteUrl, "posts", name)
				if err != nil {
					errCh <- err
					return
				}

				entry := Entry{
					Title: t,
					Link: Link{
						Href: link,
						Rel:  "alternate",
					},
					ID:      name,
					Updated: dat,
					Summary: des,
					Content: Content{
						Type: "html",
						Body: content,
					},
				}

				if imgPath != "" {
					entry.Enclosure = &Enclosure{
						URL:    imgPath,
						Type:   getImageMimeType(img),
						Length: 1,
					}
				}

				postsFeed = append(postsFeed, entry)

				var parsedContent string
				if v, ok := config.TemplateSettings[name]; ok {
					if vt, okt := templates[v]; okt {
						pstr, err := inject_post(vt.Before, headTags)
						if err != nil {
							errCh <- fmt.Errorf("%s %v", v, err)
							return
						}

						parsedContent = pstr + content + vt.After
					}
				} else if v, ok := config.TemplateSettings["*"]; ok {
					if vt, okt := templates[v]; okt {
						pstr, err := inject_post(vt.Before, headTags)
						if err != nil {
							errCh <- fmt.Errorf("%s template %v", v, err)
							return
						}
						parsedContent = pstr + content + vt.After
					}
				} else {
					errCh <- fmt.Errorf("Could not find any template %s", name)
					return
				}

				lastContent := removeEmptyLinesFast(parsedContent)

				if err := os.WriteFile(filepath.Join(popath, name), []byte(lastContent), 0755); err != nil {
					errCh <- err
					return
				}
			}
		}(entry)
	}
	wg.Wait()

	return nil
}

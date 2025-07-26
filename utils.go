package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func werror(msg string, err error) {
	fmt.Printf("\033[31m[ ERROR ][ %s ][ %v ] \033[0m\n", msg, err)
}

func checkPath(path string) (int, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 1, nil
	} else if err != nil {
		return 2, err
	}
	return 3, nil
}

func readIncludes(rootPath, s string) string {
	re := regexp.MustCompile(`{{\s*([^{}]+\.html)\s*}}`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		filename := strings.Trim(match, "{} \t\n")
		fpath := filepath.Join(rootPath, filename)
		data, err := os.ReadFile(fpath)
		if err != nil {
			return fmt.Sprintf("<!-- %s not found -->", fpath)
		}
		return string(data)
	})
}

func inject_post(htmlStr, metaContent string) (string, error) {
	closingHead := "</head>"

	index := strings.Index(strings.ToLower(htmlStr), closingHead)
	if index == -1 {
		return "", fmt.Errorf("Could not find </head> tag")
	}

	return htmlStr[:index] + "\n" + metaContent + "\n" + htmlStr[index:], nil
}

func removeEmptyLinesFast(input string) string {
	var b strings.Builder
	start := 0
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' {
			line := strings.TrimSpace(input[start:i])
			if line != "" {
				b.WriteString(line)
				b.WriteByte('\n')
			}
			start = i + 1
		}
	}

	if start < len(input) {
		line := strings.TrimSpace(input[start:])
		if line != "" {
			b.WriteString(line)
		}
	}
	return b.String()
}

func copyDir(src string, dst string, exclude map[string]bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		if exclude != nil && exclude[name] {
			continue
		}

		srcPath := filepath.Join(src, name)
		dstPath := filepath.Join(dst, name)

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath, exclude)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return os.Chmod(dstFile, srcInfo.Mode())
}

func extractHeadAndContent(htmlStr string) (head string, content string) {
	headStart := strings.Index(htmlStr, "<head>")
	headEnd := strings.Index(htmlStr, "</head>")
	contentStart := strings.Index(htmlStr, "<content>")
	contentEnd := strings.Index(htmlStr, "</content>")

	if headStart == -1 || headEnd == -1 || contentStart == -1 || contentEnd == -1 {
		return "", ""
	}

	head = htmlStr[headStart+len("<head>") : headEnd]
	content = htmlStr[contentStart+len("<content>") : contentEnd]

	return strings.TrimSpace(head), strings.TrimSpace(content)
}

// title, date, description
func findMeta(html string) (string, string, string, string) {
	re := regexp.MustCompile(`<meta\s+(?:name|itemprop)="(title|date|description|image)"[^>]*content="([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	result := map[string]string{}
	for _, m := range matches {
		if len(m) == 3 {
			result[m[1]] = m[2]
		}
	}

	return result["title"], result["date"], result["description"], result["image"]
}

func getImageMimeType(imageURL string) string {
	ext := strings.ToLower(path.Ext(imageURL))
	if ext == "" {
		return "image/*"
	}

	ext = ext[1:]

	return "image/" + ext
}

func resetDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

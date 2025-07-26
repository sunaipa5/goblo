package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func create(path string) error {
	paths := []string{"posts", "static", filepath.Join("templates", "default")}

	for _, p := range paths {
		if err := os.MkdirAll(filepath.Join(path, p), 0755); err != nil {
			return err
		}
	}

	goblo_json := `{"*": {"template": "default"}}`
	if err := os.WriteFile(filepath.Join(path, "goblo.json"), []byte(goblo_json), 0755); err != nil {
		return err
	}

	index_html := `<!doctype html>
<html lang="en">
<head>
<title>Goblo</title>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</head>
<body>
<h1>Hello goblo - index.html</h1>
</body>
</html>`

	if err := os.WriteFile(filepath.Join(path, "index.html"), []byte(index_html), 0755); err != nil {
		return err
	}

	index_tmpl := `<!doctype html>
<html lang="en">
{{header.html}}
<body>
{{content}}
{{footer.html}}
</body>
</html>
`

	if err := os.WriteFile(filepath.Join(path, "templates", "default", "index.tmpl"), []byte(index_tmpl), 0755); err != nil {
		return err
	}

	header_html := `<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<link rel="stylesheet" href="/static/style.css" />
</head>
`

	if err := os.WriteFile(filepath.Join(path, "templates", "default", "header.html"), []byte(header_html), 0755); err != nil {
		return err
	}

	footer_html := `<footer>
<p>My Blog</p>
</footer>
`

	fmt.Println(filepath.Join(path, "templates", "default", "footer.html"))

	if err := os.WriteFile(filepath.Join(path, "templates", "default", "footer.html"), []byte(footer_html), 0755); err != nil {
		return err
	}

	style_css := `* {
    margin: 0;
    padding: 0;
}

body {
    background-color: #ccc;
    color: bisque;
}`

	if err := os.WriteFile(filepath.Join(path, "static", "style.css"), []byte(style_css), 0755); err != nil {
		return err
	}

	first_post_html := `<head>
    <meta name="title" content="An awesome post" />
    <meta name="date" content="2025-07-25" />
    <meta name="description" content="New post!" />
    <title>Title</title>
</head>

<content>
    <article>
        <h1>This a post</h1>
        <p>Welcome to the post</p>
    </article>
</content>
`

	if err := os.WriteFile(filepath.Join(path, "posts", "first_post.html"), []byte(first_post_html), 0755); err != nil {
		return err
	}

	return nil
}

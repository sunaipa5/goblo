# Goblo

Fast and easy static site generator written in Go

## Features

- Easy to use
- Easy to learn
- HTML based

## Quick start

```bash
go install github.com/sunaipa5/goblo@latest

goblo create
goblo live
```

## Commands

```bash
$ goblo
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
+--------------+----------------------------+---------+
Usage: goblo <command> [path]
Commands: build, create
```

## Project Tree

```bash
example
├── dist
│   ├── index.html
│   ├── posts
│   │   └── first_post.html
│   └── static
│       └── style.css
├── goblo.json
├── index.html
├── posts
│   └── first_post.html
├── static
│   └── style.css
└── templates
    └── default
        ├── footer.html
        ├── header.html
        └── index.tmpl
```

## Usage

### Posts

> The contents of the `<head>` tag are pasted into the `<head>` tag in the template.

> The contents of the `<content>` tag are pasted into the `{{content}}` tag in the template.

```html
<head>
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
```

## Templating

```bash
default
├── footer.html
├── header.html
└── index.tmpl
```

First, create a folder in the templates folder. The name of this folder will be the template name, which you can use with the config.

### index.tmpl

```html
<!doctype html>
<html lang="en">
  {{header.html}}
  <body>
    {{content}} {{footer.html}}
  </body>
</html>
```

`{{content}}` equal to `<content></content>` in posts

`{{*.html}}` equal to `*.html` files in template folder

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func scaffold(dir, tmplType string, defaults fs.FS) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	dirs := []string{
		"content",
		"content/posts",
		"templates",
		"static",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(dir, d), 0755); err != nil {
			return err
		}
	}

	configContent := fmt.Sprintf(`[site]
title = "My Site"
description = "A site built with sitemaker"
base_url = "/"
author = "Your Name"
language = "en"

[[site.nav]]
label = "Home"
url = "/"

[[site.nav]]
label = "Posts"
url = "/posts/"

[[site.nav]]
label = "About"
url = "/about.html"

[build]
content_dir = "content"
template_dir = "templates"
static_dir = "static"
output_dir = "dist"

[blog]
posts_dir = "posts"
date_format = "2006-01-02"
posts_per_page = 10
enable_rss = true

[urls]
style = "flat"
strip_date = true

template = "%s"
`, tmplType)

	if err := writeScaffoldFile(filepath.Join(dir, "sitemaker.toml"), configContent); err != nil {
		return err
	}

	indexContent := `[TITLE]: # (Welcome)
[INCLUDES]: # (H, F)

# Welcome to My Site

This is your new site, built with **sitemaker**.

Edit this file at ` + "`content/index.md`" + ` to get started.

## Features

- Obsidian-flavored markdown
- Syntax highlighting
- Blog with RSS
- Tailwind CSS styling
- Live reload dev server
`
	if err := writeScaffoldFile(filepath.Join(dir, "content", "index.md"), indexContent); err != nil {
		return err
	}

	aboutContent := `[TITLE]: # (About)
[INCLUDES]: # (H, F)

# About

This is the about page. Edit it at ` + "`content/about.md`" + `.
`
	if err := writeScaffoldFile(filepath.Join(dir, "content", "about.md"), aboutContent); err != nil {
		return err
	}

	postContent := `[TITLE]: # (Hello World)
[DATE]: # (2025-01-15)
[TAGS]: # (welcome, blog)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# Hello World

Welcome to my first blog post!

## Getting Started

This is a sample post to get you started. You can find it at ` + "`content/posts/2025-01-15-hello-world.md`" + `.

## Code Blocks

Here's a syntax-highlighted code block:

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, sitemaker!")
}
` + "```" + `

## Callouts

> [!tip] Pro Tip
> You can use Obsidian-style callouts for notes, warnings, tips, and more.

> [!note] Note
> Sitemaker supports most Obsidian markdown extensions.

## Links

You can use [[about|wikilinks]] to link between pages.

---

That's it for now. Happy writing!
`
	if err := writeScaffoldFile(filepath.Join(dir, "content", "posts", "2025-01-15-hello-world.md"), postContent); err != nil {
		return err
	}

	gitignore := "dist/\n"
	if err := writeScaffoldFile(filepath.Join(dir, ".gitignore"), gitignore); err != nil {
		return err
	}

	return nil
}

func writeScaffoldFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

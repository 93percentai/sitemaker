import fs from "node:fs";
import path from "node:path";

export function scaffold(dir: string, tmplType: string): void {
  fs.mkdirSync(dir, { recursive: true });

  for (const d of ["content", "content/posts", "templates", "static"]) {
    fs.mkdirSync(path.join(dir, d), { recursive: true });
  }

  writeFile(
    path.join(dir, "sitemaker.toml"),
    `[site]
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
date_format = "YYYY-MM-DD"
posts_per_page = 10
enable_rss = true

[urls]
style = "flat"
strip_date = true

template = "${tmplType}"
`
  );

  writeFile(
    path.join(dir, "content", "index.md"),
    `[TITLE]: # (Welcome)
[INCLUDES]: # (H, F)

# Welcome to My Site

This is your new site, built with **sitemaker**.

Edit this file at \`content/index.md\` to get started.

## Features

- Obsidian-flavored markdown
- Syntax highlighting
- Blog with RSS
- Tailwind CSS styling
- Live reload dev server
`
  );

  writeFile(
    path.join(dir, "content", "about.md"),
    `[TITLE]: # (About)
[INCLUDES]: # (H, F)

# About

This is the about page. Edit it at \`content/about.md\`.
`
  );

  writeFile(
    path.join(dir, "content", "posts", "2025-01-15-hello-world.md"),
    `[TITLE]: # (Hello World)
[DATE]: # (2025-01-15)
[TAGS]: # (welcome, blog)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# Hello World

Welcome to my first blog post!

## Code Blocks

\`\`\`javascript
function greet(name) {
    console.log(\`Hello, \${name}!\`);
}

greet("World");
\`\`\`

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
  );

  writeFile(path.join(dir, ".gitignore"), "dist/\nnode_modules/\n");
}

function writeFile(filePath: string, content: string): void {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });
  fs.writeFileSync(filePath, content);
}

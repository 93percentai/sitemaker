# sitemaker

A fast, opinionated static site generator written in TypeScript -- built for personal sites, blogs, and small org websites with Obsidian-flavored markdown (via [markdown-it](https://github.com/markdown-it/markdown-it)), [Tailwind CSS v4](https://github.com/tailwindlabs/tailwindcss), and zero configuration needed to get started.

## Features

- **Obsidian-flavored Markdown** -- wikilinks, callouts, footnotes, strikethrough, task lists, tables, and raw HTML
- **Syntax highlighting** -- code blocks with language-aware highlighting via [highlight.js](https://highlightjs.org/) (GitHub theme)
- **Blog engine** -- date-based posts, tags, tag index pages, post listing pages, prev/next navigation, reading time estimates
- **RSS feed** -- auto-generated `rss.xml` and `feed.xml` for blog posts
- **Two built-in templates** -- `personal` (clean, minimal) and `org` (sidebar layout with card-style content area)
- **Directive system** -- per-page control over layout components, template inheritance, metadata, and drafts via markdown-compatible directives
- **Table of contents** -- auto-generated from headings, rendered as a sticky sidebar
- **Tailwind CSS v4.3** -- pre-built CSS bundled with the project, scanned from templates
- **Live reload dev server** -- file watcher with SSE-based instant browser refresh
- **Flat or directory URL styles** -- `/page.html` or `/page/index.html`
- **Static file passthrough** -- anything in `static/` is copied to the output verbatim
- **Cloudflare Pages ready** -- auto-generated `_headers` file with security and caching headers
- **Nunjucks templates** -- powerful template engine with inheritance, includes, and filters
- **MIT License**

## Quick Start

```bash
# 1. Install
git clone https://github.com/93percentai/sitemaker.git
cd sitemaker && npm install && npm run build

# 2. Initialize a new project
node dist/index.js init mysite
cd mysite

# 3. Build the site
node ../dist/index.js build

# 4. Start the dev server with live reload
node ../dist/index.js serve
```

Open [http://localhost:8080](http://localhost:8080) to see your site.

## Installation

### From source

```bash
git clone https://github.com/93percentai/sitemaker.git
cd sitemaker
npm install
npm run build
```

### As a global CLI

```bash
npm link   # from the sitemaker directory
sitemaker build
```

Requires Node.js 18 or later.

## CLI Reference

### `sitemaker build`

Build the site to the output directory.

```
sitemaker build [options]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `sitemaker.toml` | Path to the config file |

### `sitemaker serve`

Start a development server with live reload. Watches `content/`, `templates/`, and `static/` for changes and rebuilds automatically.

```
sitemaker serve [options]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `sitemaker.toml` | Path to the config file |
| `-port` | `8080` | Port to serve on |

The dev server injects a live-reload script into every HTML page. When any source file changes, the browser refreshes automatically via Server-Sent Events. File system events are debounced (200ms) to batch rapid saves.

### `sitemaker init`

Scaffold a new sitemaker project with a config file, example content, and directory structure.

```
sitemaker init [options] [directory]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-template` | `personal` | Template type: `personal` or `org` |

If no directory is given, the current directory is used. Creates:

```
directory/
  sitemaker.toml
  content/
    index.md
    about.md
    posts/
      2025-01-15-hello-world.md
  templates/
  static/
  .gitignore
```

## Configuration

All configuration lives in `sitemaker.toml`. Every field has a sensible default -- the file is optional for basic sites.

### Full reference

```toml
# --- Site metadata ---
[site]
title = "My Site"                              # Site title, used in <title> and templates
description = "A site built with sitemaker"    # Meta description
base_url = "/"                                 # Base URL for RSS links and canonical URLs
author = "Your Name"                           # Author name for meta tags and footer
language = "en"                                # HTML lang attribute

# Navigation links (rendered in header, footer, sidebar)
[[site.nav]]
label = "Home"
url = "/"

[[site.nav]]
label = "Posts"
url = "/posts/"

[[site.nav]]
label = "About"
url = "/about.html"

# --- Build paths ---
[build]
content_dir = "content"      # Where markdown files live
template_dir = "templates"   # Where user template overrides live
static_dir = "static"        # Static files copied verbatim to output
output_dir = "dist"          # Build output directory (cleaned on each build)

# --- Blog settings ---
[blog]
posts_dir = "posts"          # Subdirectory of content_dir for blog posts
date_format = "2006-01-02"   # Go time format for date parsing and display
posts_per_page = 10          # Posts per page (used in list template data)
enable_rss = true            # Generate rss.xml and feed.xml

# --- URL generation ---
[urls]
style = "flat"               # "flat" = /page.html, "directory" = /page/index.html
strip_date = true            # Remove YYYY-MM-DD prefix from post slugs in URLs

# --- Template selection ---
template = "personal"        # Built-in template: "personal" or "org"
```

### `[site]` section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `title` | string | `"My Site"` | Site title |
| `description` | string | `"A site built with sitemaker"` | Meta description |
| `base_url` | string | `"/"` | Base URL for absolute links (RSS, canonical) |
| `author` | string | `""` | Author name |
| `language` | string | `"en"` | HTML `lang` attribute |
| `nav` | array of `{label, url}` | `[]` | Navigation links |

### `[build]` section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `content_dir` | string | `"content"` | Markdown content directory |
| `template_dir` | string | `"templates"` | User template overrides directory |
| `static_dir` | string | `"static"` | Static assets directory |
| `output_dir` | string | `"dist"` | Build output directory |

### `[blog]` section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `posts_dir` | string | `"posts"` | Subdirectory within `content_dir` for posts |
| `date_format` | string | `"2006-01-02"` | Go-style date format string |
| `posts_per_page` | int | `10` | Number of posts per list page |
| `enable_rss` | bool | `true` | Generate RSS feed files |

### `[urls]` section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `style` | string | `"flat"` | `"flat"` produces `/slug.html`, `"directory"` produces `/slug/index.html` |
| `strip_date` | bool | `true` | Strip `YYYY-MM-DD-` prefix from blog post filenames when generating URLs |

### `template`

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `template` | string | `"personal"` | Built-in template set: `"personal"` or `"org"` |

## Directives

Directives are metadata annotations placed at the top of any markdown file. They use a markdown-compatible syntax that renders as invisible link references, so your files remain valid markdown in any editor.

**Syntax:** `[DIRECTIVE]: # (value)`

Directives must appear at the top of the file before any body content. Blank lines between directives are allowed.

### Complete reference

| Directive | Value | Description |
|-----------|-------|-------------|
| `[TITLE]: # (...)` | Text | Page title. Overrides the first `# Heading` and the filename-based fallback. |
| `[DATE]: # (...)` | Date string | Publication date. Parsed using the `date_format` from config. Falls back to the `YYYY-MM-DD` prefix in the filename. |
| `[TAGS]: # (...)` | Comma-separated | Tags for the page. Each tag generates a tag index page at `/tags/<tag>/`. |
| `[SLUG]: # (...)` | Text | Custom URL slug. Overrides the filename-derived slug. |
| `[DRAFT]: # (true)` | `true` | Marks the page as a draft. Draft pages are excluded from the build output. |
| `[INCLUDES]: # (...)` | Comma-separated | Layout components and asset files to include (see below). |
| `[INHERITS]: # (...)` | Template filename | Template to use for rendering. Defaults to `post.html` for blog posts, `base.html` for pages. |
| `[OVERRIDE]: # (...)` | Template filename | Replaces the entire page with a custom HTML template from `templates/`. No base template or partials are used. |

### INCLUDES directive in detail

The `INCLUDES` directive accepts a comma-separated list of **layout component codes** and **asset file paths**.

**Layout components:**

| Code | Component | Description |
|------|-----------|-------------|
| `H` | Header | Site header with navigation and mobile menu |
| `F` | Footer | Site footer with copyright, nav links, and RSS link |
| `TOC` | Table of Contents | Auto-generated from headings, rendered as a sticky sidebar |
| `SB` | Sidebar | Navigation sidebar (always visible in `org` template, toggled in `personal`) |

**Asset files:**

Any value containing a `.` is treated as an asset file path. CSS files are added as `<link>` tags in the `<head>`, and JS files are added as `<script>` tags before `</body>`.

```markdown
[INCLUDES]: # (H, F, TOC, /assets/custom.css, /js/charts.js)
```

Components not listed in `INCLUDES` are rendered as empty blocks, effectively hiding them.

### Examples

**Basic page with header and footer:**

```markdown
[TITLE]: # (About Us)
[INCLUDES]: # (H, F)

# About Us

Content goes here...
```

**Blog post with all components:**

```markdown
[TITLE]: # (My First Post)
[DATE]: # (2025-03-15)
[TAGS]: # (go, webdev, tutorial)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# My First Post

This post has a table of contents sidebar, navigation header, and footer.
```

**Draft post:**

```markdown
[TITLE]: # (Work in Progress)
[DATE]: # (2025-06-01)
[DRAFT]: # (true)
[INCLUDES]: # (H, F)

# Not ready yet

This page will not appear in the build output.
```

**Custom slug:**

```markdown
[TITLE]: # (Frequently Asked Questions)
[SLUG]: # (faq)
[INCLUDES]: # (H, F)

# FAQ

This page will be served at /faq.html (flat) or /faq/ (directory).
```

**Full HTML override (landing page, custom layout):**

```markdown
[TITLE]: # (Welcome)
[OVERRIDE]: # (landing.html)

Content here is available as .Page.HTMLContent in the template,
but the template has full control over the HTML document.
```

## Template System

### Built-in templates

Sitemaker ships with two template sets, selectable via the `template` key in config:

**`personal`** -- A clean, minimal layout centered on readability. Content area is capped at `max-w-4xl` (or `max-w-3xl` for posts). White background, no persistent sidebar. Best for personal blogs and portfolios.

**`org`** -- A wider layout with a persistent left sidebar for navigation. Content renders inside a bordered card on a light gray background (`bg-gray-50`). Content area uses `max-w-7xl`. Best for documentation sites, team blogs, and organizational pages.

Each template set includes three files:

| File | Purpose |
|------|---------|
| `base.html` | Default layout for non-post pages |
| `post.html` | Layout for blog posts (adds date, tags, reading time, prev/next nav) |
| `list.html` | Layout for post listing and tag index pages |

### Shared partials

Both template sets share the same partial templates:

| Partial | Defined name | Description |
|---------|-------------|-------------|
| `header.html` | `header` | Responsive nav bar with site title, nav links, mobile hamburger menu, RSS icon |
| `footer.html` | `footer` | Copyright line, nav links, RSS link |
| `sidebar.html` | `sidebar` | Sticky navigation sidebar using `site.nav` links |
| `toc.html` | `toc` | Sticky "On this page" table of contents generated from headings |

### Template resolution order

When rendering a page, sitemaker looks for templates in this priority order:

1. **User templates** -- `templates/<name>` (your project's `template_dir`)
2. **Built-in templates** -- `defaults/templates/<template-type>/<name>` (embedded in the binary)
3. **Built-in partials** -- `defaults/templates/partials/<name>`

This means you can override any template by placing a file with the same name in your `templates/` directory. For partials, place them in `templates/partials/`.

### INHERITS -- choosing a template

The `INHERITS` directive selects which template file wraps the page. It follows the same resolution order above.

```markdown
[INHERITS]: # (post.html)
```

If omitted, blog posts default to `post.html` and all other pages default to `base.html`.

### OVERRIDE -- raw HTML pages

The `OVERRIDE` directive bypasses the normal template system entirely. The specified template file (loaded only from user `templates/`) receives the full template data context but is responsible for the entire HTML document.

```markdown
[OVERRIDE]: # (landing.html)
```

This is useful for landing pages, custom layouts, or pages that need complete control over the markup.

### Template data context

All templates receive a `TemplateData` struct:

```
.Site
  .Title          string          -- Site title
  .Description    string          -- Site description
  .BaseURL        string          -- Base URL
  .Author         string          -- Author name
  .Language        string          -- Language code
  .Nav            []NavItem       -- Navigation items (.Label, .URL)
  .BuildTime      time.Time       -- Build timestamp

.Page
  .Title          string          -- Page title
  .HTMLContent    template.HTML   -- Rendered markdown HTML
  .TOCContent     template.HTML   -- Rendered table of contents HTML
  .URL            string          -- Page URL
  .Date           time.Time       -- Page date
  .DateStr        string          -- Formatted date string
  .Tags           []string        -- Tag list
  .ReadingTime    int             -- Estimated reading time in minutes
  .HasHeader      bool            -- Whether header component is included
  .HasFooter      bool            -- Whether footer component is included
  .HasTOC         bool            -- Whether TOC component is included
  .HasSidebar     bool            -- Whether sidebar component is included
  .CSSAssets      []string        -- Extra CSS file paths from INCLUDES
  .JSAssets       []string        -- Extra JS file paths from INCLUDES
  .IsBlogPost     bool            -- Whether this is a blog post
  .Prev           *PageRef        -- Previous post (.Title, .URL) or nil
  .Next           *PageRef        -- Next post (.Title, .URL) or nil

.Posts            []PageData      -- All posts (available on list pages)
.Tags             map[string][]PageData -- Tag index (available on list pages)
.CurrentTag       string          -- Current tag being rendered (tag pages)
```

### Template functions

The following functions are available in all templates:

| Function | Description | Example |
|----------|-------------|---------|
| `lower` | Lowercase string | `{{ lower .Page.Title }}` |
| `upper` | Uppercase string | `{{ upper "hello" }}` |
| `title` | Title-case string | `{{ title "hello world" }}` |
| `join` | Join string slice | `{{ join .Page.Tags ", " }}` |
| `split` | Split string | `{{ split "a,b,c" "," }}` |
| `contains` | String contains | `{{ if contains .Page.URL "/posts" }}` |
| `hasPrefix` | String starts with | `{{ hasPrefix .Page.URL "/" }}` |
| `hasSuffix` | String ends with | `{{ hasSuffix .Page.URL ".html" }}` |
| `trimPrefix` | Remove prefix | `{{ trimPrefix .Page.URL "/" }}` |
| `trimSuffix` | Remove suffix | `{{ trimSuffix .Page.URL ".html" }}` |
| `replace` | Replace all occurrences | `{{ replace .Page.Title " " "-" }}` |
| `formatDate` | Format a time.Time | `{{ formatDate .Page.Date "Jan 2, 2006" }}` |
| `now` | Current time | `{{ now.Year }}` |
| `safeHTML` | Mark string as safe HTML | `{{ safeHTML "<b>bold</b>" }}` |
| `seq` | Generate integer sequence | `{{ range seq 5 }}...{{ end }}` |
| `add` | Integer addition | `{{ add 1 2 }}` |
| `sub` | Integer subtraction | `{{ sub 10 3 }}` |
| `slice` | Create a slice from args | `{{ slice "a" "b" "c" }}` |

## Markdown Features

Sitemaker uses [Goldmark](https://github.com/yuin/goldmark) with GitHub Flavored Markdown (GFM) extensions and additional Obsidian-compatible features.

### Standard markdown

All standard markdown syntax is supported: headings, paragraphs, bold, italic, links, images, blockquotes, ordered/unordered lists, horizontal rules, and inline code.

### GitHub Flavored Markdown (GFM)

- **Tables**
- **Strikethrough** (`~~deleted~~`)
- **Task lists** (`- [x] done`, `- [ ] todo`)
- **Autolinks** (URLs are automatically linked)

### Wikilinks

Obsidian-style `[[wikilinks]]` are supported. They generate links based on the target page's slug.

```markdown
Link to [[another page]].
Link with custom text: [[another-page|click here]].
```

Renders as:

```html
<a href="/another-page.html" class="wikilink">another page</a>
<a href="/another-page.html" class="wikilink">click here</a>
```

Wikilinks are styled with a dotted underline that becomes solid on hover.

### Callouts

Obsidian-style callouts are rendered from blockquote syntax:

```markdown
> [!note] Title here
> Callout body content.

> [!tip] Pro Tip
> This is a tip callout.

> [!warning]
> If no title is provided, the callout type is used as the title.
```

**Supported callout types:**

| Type | Icon | Color |
|------|------|-------|
| `note` | Pencil | Blue |
| `tip` | Light bulb | Green |
| `info` | Info | Cyan |
| `warning` | Warning sign | Amber |
| `danger` | Hazard | Red |
| `bug` | Bug | Red |
| `example` | Clipboard | Purple |
| `quote` | Quote mark | Gray |
| `todo` | Checkbox | Blue |
| `success` | Checkmark | Green |
| `question` | Question mark | Amber |
| `failure` | X mark | Red |
| `abstract` | Document | Cyan |

Each callout renders as a styled `<div>` with a colored left border, background tint, icon, and title.

### Code blocks with syntax highlighting

Fenced code blocks with a language identifier get syntax highlighting via [Chroma](https://github.com/alecthomas/chroma) using CSS classes (GitHub theme). All languages supported by Chroma are available.

````markdown
```go
func main() {
    fmt.Println("Hello!")
}
```
````

### Footnotes

Standard footnote syntax is supported:

```markdown
This has a footnote[^1].

[^1]: The footnote content appears at the bottom of the page.
```

### Raw HTML and script blocks

HTML is rendered unsafe by default -- you can embed arbitrary HTML, `<script>` tags, `<style>` blocks, and custom elements directly in markdown:

```markdown
<div style="padding: 1rem; background: #f0f9ff;">
  <strong>Custom HTML</strong> rendered inline.
</div>

<script>
console.log("Scripts work too.");
</script>
```

### Auto heading IDs

All headings automatically get an `id` attribute derived from their text, enabling anchor links and the table of contents feature.

## Content Organization

### Directory structure

```
content/
  index.md              --> /              (site homepage)
  about.md              --> /about.html    (flat) or /about/ (directory)
  projects.md           --> /projects.html
  posts/
    2025-01-15-hello.md --> /posts/hello.html
    2025-02-20-update.md --> /posts/update.html
  docs/
    getting-started.md  --> /docs/getting-started.html
```

### Pages vs blog posts

Any `.md` file inside the blog `posts_dir` (default: `content/posts/`) is treated as a **blog post**. Everything else is a **page**.

Blog posts get additional behavior:

- Sorted by date (newest first)
- Included in the post listing page at `/posts/`
- Included in the tag index pages at `/tags/<tag>/`
- Linked with prev/next navigation
- Included in the RSS feed
- Automatically use the `post.html` template (unless overridden)

### URL generation

URLs are derived from the file path relative to `content_dir`:

| File | Flat style | Directory style |
|------|------------|-----------------|
| `content/about.md` | `/about.html` | `/about/` |
| `content/posts/2025-01-15-hello.md` | `/posts/hello.html` | `/posts/hello/` |
| `content/docs/setup.md` | `/docs/setup.html` | `/docs/setup/` |
| `content/index.md` | `/` | `/` |
| `content/docs/index.md` | `/docs/` | `/docs/` |

**Date stripping:** When `urls.strip_date = true` (the default), the `YYYY-MM-DD-` prefix is removed from blog post filenames when generating the slug and URL.

**Custom slugs:** The `[SLUG]` directive overrides the filename-derived slug entirely.

**Index files:** Files named `index.md` produce `index.html` in their parent directory and get a trailing-slash URL.

### Title resolution

Page titles are resolved in this order:

1. `[TITLE]` directive
2. First `# Heading` in the markdown body
3. Filename (without extension)

## Blog Features

### Dates

Post dates are resolved from (in priority order):

1. The `[DATE]` directive value, parsed with the configured `date_format`
2. The `YYYY-MM-DD` prefix in the filename (e.g., `2025-01-15-hello-world.md`)

### Tags

Tags are set via the `[TAGS]` directive and are comma-separated:

```markdown
[TAGS]: # (go, webdev, tutorial)
```

Each tag automatically generates an index page at `/tags/<tag>/` listing all posts with that tag.

### Post listing

A listing page is automatically generated at `/<posts_dir>/` (default: `/posts/`) showing all non-draft posts sorted by date, newest first.

### Prev/next navigation

Blog posts include prev/next links automatically. Posts are ordered by date (newest first), so "previous" goes to the older post and "next" goes to the newer post.

### Reading time

An estimated reading time is calculated for every page based on character count (assumes ~200 words per minute, with 5 characters per word). Displayed in post templates as "X min read".

### RSS feed

When `blog.enable_rss` is `true` (the default), an RSS 2.0 feed is generated at both `/rss.xml` and `/feed.xml`. The feed includes up to 20 of the most recent posts with their full HTML content, publication dates, tags as categories, and proper `<guid>` elements.

The RSS link is automatically included in the HTML `<head>` of all template pages:

```html
<link rel="alternate" type="application/rss+xml" title="Site Title" href="/rss.xml">
```

## Templates

### `personal` template

A minimal, content-focused layout ideal for personal blogs and portfolios.

- White background (`bg-white`)
- Centered content, narrow max width (`max-w-4xl` for pages, `max-w-3xl` for posts)
- Optional sidebar (hidden by default, shown when `SB` is in INCLUDES)
- Clean typography with Tailwind's prose classes

### `org` template

A documentation-style layout with a persistent sidebar, suited for teams and organizations.

- Light gray background (`bg-gray-50`) with white content card
- Persistent left sidebar (64px wide on large screens)
- Bordered, rounded content area with shadow
- Wider layout (`max-w-7xl`)
- Always shows sidebar navigation regardless of `SB` in INCLUDES

### Customizing templates

To customize any template, copy it from the built-in defaults to your project's `templates/` directory and edit it:

```bash
# Override the base layout
cp defaults/templates/personal/base.html templates/base.html

# Override a partial
mkdir -p templates/partials
cp defaults/templates/partials/header.html templates/partials/header.html
```

User templates in `templates/` always take priority over built-in templates.

### Creating custom templates

You can create entirely new templates and reference them with `INHERITS`:

```html
<!-- templates/docs.html -->
<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <title>{{ .Page.Title }} - {{ .Site.Title }}</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body>
    {{ template "header" . }}
    <main class="max-w-prose mx-auto py-12">
        {{ .Page.HTMLContent }}
    </main>
    {{ template "footer" . }}
</body>
</html>
```

```markdown
[INHERITS]: # (docs.html)
[INCLUDES]: # (H, F)

# Documentation Page

Content here uses the custom docs.html template.
```

## Tailwind CSS

Sitemaker includes [Tailwind CSS](https://tailwindcss.com) via the CDN play script. The Tailwind configuration is embedded directly in the base templates:

```html
<script src="https://cdn.tailwindcss.com"></script>
<script>
tailwind.config = {
    theme: {
        extend: {
            colors: {
                accent: '#2563eb',
            },
        },
    },
}
</script>
```

### Customizing

To change the accent color or extend the theme, override the base template and modify the inline `tailwind.config` object. The `org` template also defines an `accent-dark` color (`#1d4ed8`).

### Production optimization

The CDN version is convenient for development but includes the full Tailwind runtime. For production sites where bundle size matters, you can override the base template to use a locally built Tailwind CSS file instead:

1. Install Tailwind CLI and generate a production CSS file
2. Place it in `static/css/tailwind.css`
3. Replace the CDN `<script>` tag in your template override with a `<link>` tag

## Deployment

### Cloudflare Pages

Sitemaker generates a `_headers` file in the output directory that is compatible with [Cloudflare Pages](https://pages.cloudflare.com/).

**Setup:**

| Setting | Value |
|---------|-------|
| Build command | `sitemaker build` |
| Build output directory | `dist` |

The generated `_headers` file includes:

```
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY

/assets/*
  Cache-Control: public, max-age=31536000, immutable

/*.html
  Cache-Control: public, max-age=0, must-revalidate
```

This configures security headers globally, immutable caching for assets, and revalidation for HTML pages.

### Other platforms

Sitemaker outputs a standard static site to the `dist/` directory (configurable via `output_dir`). It works with any static hosting provider:

- **Netlify** -- set build command to `sitemaker build`, publish directory to `dist`
- **GitHub Pages** -- build in CI and deploy the `dist/` directory
- **Vercel** -- set build command and output directory
- **S3 + CloudFront** -- upload the `dist/` directory

## Project Structure

```
sitemaker/
  main.go                           # CLI entry point (build, serve, init commands)
  scaffold.go                       # Project scaffolding for `sitemaker init`
  go.mod
  defaults/                         # Embedded default assets and templates
    css/
      syntax.css                    # Chroma syntax highlighting styles
    js/
      main.js                       # Smooth scroll for anchor links
    templates/
      personal/                     # Personal template set
        base.html                   # Default page layout
        post.html                   # Blog post layout
        list.html                   # Post listing layout
      org/                          # Org template set
        base.html
        post.html
        list.html
      partials/                     # Shared partial templates
        header.html
        footer.html
        sidebar.html
        toc.html
  internal/
    build/
      build.go                      # Build pipeline orchestration
    config/
      config.go                     # TOML config parsing with defaults
    content/
      content.go                    # Page loading, URL/slug/title/date resolution
      blog.go                       # Post classification, tag index, prev/next linking
    directive/
      directive.go                  # Directive parsing (TITLE, DATE, TAGS, etc.)
    markdown/
      markdown.go                   # Goldmark renderer setup (GFM, footnotes, highlighting)
      wikilink.go                   # Obsidian [[wikilink]] parser and renderer
      callout.go                    # Obsidian callout parser and renderer
    rss/
      rss.go                        # RSS 2.0 feed generation
    server/
      server.go                     # Dev server with file watching and live reload
    template/
      template.go                   # Template engine with inheritance and partials
      funcs.go                      # Custom template functions
    toc/
      toc.go                        # Table of contents extraction and rendering
  example/
    sitemaker.toml                  # Example config
    content/
      posts/
        2025-01-15-hello-world.md   # Example blog post
```

## Development

### Prerequisites

- Go 1.26 or later

### Building

```bash
go build -o sitemaker .
```

### Running tests

```bash
go test ./...
```

### Project architecture

The build pipeline follows this sequence:

1. **Load config** -- parse `sitemaker.toml`, fill in defaults
2. **Load content** -- walk `content_dir`, parse directives, resolve titles/slugs/dates/URLs
3. **Classify pages** -- separate posts from pages, sort posts by date
4. **Build tag index** -- group posts by tag
5. **Link prev/next** -- wire up navigation between adjacent posts
6. **Render markdown** -- convert markdown to HTML with all extensions (wikilinks, callouts, highlighting, GFM)
7. **Extract TOC** -- parse headings for pages that include the TOC component
8. **Render templates** -- execute Go templates with page data, partials, and template functions
9. **Generate listing pages** -- post index and per-tag index pages
10. **Generate RSS** -- RSS 2.0 feed with up to 20 recent posts
11. **Copy static files** -- from `static/` to output
12. **Copy default assets** -- embedded CSS/JS to `assets/` in output
13. **Write headers** -- generate `_headers` for Cloudflare Pages

## License

[The Unlicense](LICENSE) -- public domain. Use it however you want.

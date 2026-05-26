[TITLE]: # (Markdown Features Deep Dive)
[DATE]: # (2025-02-20)
[TAGS]: # (markdown, tutorial)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# Markdown Features Deep Dive

A comprehensive look at all the markdown features sitemaker supports.

## Obsidian Wikilinks

You can link to other pages using [[hello-world|wikilink syntax]]. This creates a link to the hello-world page.

## Strikethrough

~~This text is struck through.~~

## Emphasis

*Italic*, **bold**, ***bold italic***, and `inline code`.

## Blockquotes

> This is a blockquote.
> It can span multiple lines.
>
> And have multiple paragraphs.

## Images

Images work with standard markdown syntax:

```markdown
![Alt text](/path/to/image.png)
```

## Horizontal Rules

---

## Nested Lists

1. First item
   - Nested bullet
   - Another nested bullet
     1. Deep nested numbered
     2. Another deep nested
2. Second item
3. Third item

## Definition-style Content

Using callouts for definitions:

> [!info] What is a Static Site Generator?
> A static site generator (SSG) is a tool that generates a full static HTML website based on raw data and templates. Unlike dynamic websites, static sites are pre-built and served as-is.

> [!example] Example Usage
> Run `sitemaker build` to generate your site, then deploy the `dist/` folder to any static hosting provider.

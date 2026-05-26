[TITLE]: # (Every Markdown Feature Sitemaker Supports)
[DATE]: # (2026-05-15)
[TAGS]: # (markdown, reference)
[INCLUDES]: # (H, F, TOC, SB)
[INHERITS]: # (post.html)

# Every Markdown Feature Sitemaker Supports

A comprehensive reference of every markdown formatting option available in sitemaker, powered by markdown-it.

## Inline Formatting

**Bold text** with `**double asterisks**`

*Italic text* with `*single asterisks*`

***Bold and italic*** with `***triple asterisks***`

~~Strikethrough~~ with `~~double tildes~~`

`Inline code` with `` `backticks` ``

[Links](https://example.com) with `[text](url)`

[[about|Wikilinks]] with `[[target|display text]]`

## Headings

Headings from `# H1` through `###### H6` — each gets an auto-generated ID for linking:

### H3 Heading
#### H4 Heading
##### H5 Heading
###### H6 Heading

## Paragraphs and Line Breaks

Regular paragraphs are separated by blank lines.

This is a second paragraph. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

## Blockquotes

> This is a blockquote.
>
> It can span multiple paragraphs.
>
> > And can be nested.
> >
> > Like this.

## Unordered Lists

- First item
- Second item
  - Nested item A
  - Nested item B
    - Deeply nested
- Third item

Alternate syntax:

* Asterisk item
* Another item

## Ordered Lists

1. First step
2. Second step
   1. Sub-step A
   2. Sub-step B
3. Third step

Starting at a different number:

5. Item five
6. Item six
7. Item seven

## Task Lists

- [x] Write the markdown parser
- [x] Add syntax highlighting
- [x] Implement callouts
- [ ] Add search functionality
- [ ] Write more documentation

## Code Blocks

Fenced code blocks with language annotation:

```javascript
// JavaScript
const greet = (name) => `Hello, ${name}!`;
console.log(greet("World"));
```

```python
# Python
def factorial(n: int) -> int:
    if n <= 1:
        return 1
    return n * factorial(n - 1)

print(factorial(10))
```

```rust
// Rust
fn main() {
    let numbers: Vec<i32> = (1..=10).collect();
    let sum: i32 = numbers.iter().sum();
    println!("Sum: {}", sum);
}
```

```sql
-- SQL
SELECT
    u.name,
    COUNT(p.id) AS post_count,
    MAX(p.created_at) AS latest_post
FROM users u
LEFT JOIN posts p ON p.author_id = u.id
WHERE u.active = true
GROUP BY u.name
HAVING COUNT(p.id) > 5
ORDER BY post_count DESC;
```

```bash
# Shell
#!/bin/bash
for file in *.md; do
    echo "Processing: $file"
    wc -w "$file"
done
```

```yaml
# YAML
site:
  title: My Blog
  description: A sitemaker blog
  nav:
    - label: Home
      url: /
    - label: Posts
      url: /posts/
```

```diff
- removed line
+ added line
  unchanged line
- another removed line
+ another added line
```

Inline code: `const x = 42;`

Indented code (4 spaces):

    function oldSchool() {
        return "This also works";
    }

## Tables

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:--------------:|--------------:|
| Left | Center | Right |
| Data | Data | Data |
| Longer content here | Short | 42 |

## Horizontal Rules

Three or more hyphens:

---

Three or more asterisks:

***

Three or more underscores:

___

## Links

[Inline link](https://example.com)

[Link with title](https://example.com "Example Title")

Autolinks: https://example.com

Email autolinks: user@example.com

## Images

```markdown
![Alt text](https://via.placeholder.com/600x200 "Image title")
```

## Footnotes

Here's a sentence with a footnote[^1], and another[^longnote].

[^1]: This is a simple footnote.

[^longnote]: This is a longer footnote with multiple paragraphs.

    Subsequent paragraphs are indented to show they belong to the previous footnote.

## Callouts (Obsidian-flavored)

All supported callout types:

> [!note] Note
> Default informational callout.

> [!tip] Tip
> Helpful advice for the reader.

> [!warning] Warning
> Something to be cautious about.

> [!danger] Danger
> Critical warning — something could break.

> [!info] Info
> Additional context or background.

> [!example] Example
> A concrete example or demonstration.

> [!quote] Quote
> A notable quotation.

> [!bug] Bug
> Known issue or bug report.

> [!todo] Todo
> Action item or task.

> [!success] Success
> Something that worked or passed.

> [!question] Question
> An open question or inquiry.

> [!failure] Failure
> Something that failed or didn't work.

> [!abstract] Abstract
> Summary or overview.

## Wikilinks

Link to another page: [[about|About Page]]

Link by filename: [[building-a-rate-limiter-in-go|Rate Limiter Post]]

## Escaping

\*Not italic\*

\`Not code\`

\# Not a heading

## Typographer Replacements

Quotes become "smart quotes" and 'single quotes'.

Dashes: -- becomes an en-dash, --- becomes an em-dash.

Ellipsis: ... becomes an ellipsis.

## HTML Entities

&copy; &reg; &trade; &mdash; &ndash; &hellip; &frac12; &frac14;

---

That covers every markdown feature sitemaker supports. Use this page as a formatting reference when writing your content.

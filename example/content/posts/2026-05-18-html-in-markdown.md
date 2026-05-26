[TITLE]: # (HTML Inside Markdown: A Complete Guide)
[DATE]: # (2026-05-18)
[TAGS]: # (html, markdown, tutorial)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# HTML Inside Markdown: A Complete Guide

Sitemaker renders markdown with `html: true`, meaning you can freely mix HTML and markdown. Here's every pattern you might need.

## Styled Containers

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin: 1.5rem 0;">
  <div style="padding: 1.5rem; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 0.75rem; color: white;">
    <h3 style="margin: 0 0 0.5rem 0; color: white;">Card One</h3>
    <p style="margin: 0; opacity: 0.9;">Gradient backgrounds with custom grid layouts — all inside markdown.</p>
  </div>
  <div style="padding: 1.5rem; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); border-radius: 0.75rem; color: white;">
    <h3 style="margin: 0 0 0.5rem 0; color: white;">Card Two</h3>
    <p style="margin: 0; opacity: 0.9;">CSS Grid works perfectly inside HTML blocks embedded in markdown.</p>
  </div>
</div>

## Embedded Media

### YouTube Video (iframe)

<div style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; margin: 1.5rem 0; border-radius: 0.5rem; background: #000;">
  <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: #666; text-align: center;">
    <p style="font-size: 3rem; margin: 0;">▶</p>
    <p>Video placeholder — replace with an iframe</p>
  </div>
</div>

Usage:
```html
<iframe src="https://www.youtube.com/embed/VIDEO_ID"
  width="100%" height="400"
  frameborder="0" allowfullscreen>
</iframe>
```

### SVG Inline

<svg viewBox="0 0 200 60" style="width: 200px; margin: 1rem 0;">
  <rect x="0" y="0" width="200" height="60" rx="8" fill="#2563eb"/>
  <text x="100" y="35" text-anchor="middle" fill="white" font-family="sans-serif" font-size="14" font-weight="bold">Inline SVG</text>
</svg>

## Tables with HTML Styling

Standard markdown tables work:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | List all users |
| POST | `/api/users` | Create a user |
| PUT | `/api/users/:id` | Update a user |
| DELETE | `/api/users/:id` | Delete a user |

But you can also use HTML tables for more control:

<table style="width: 100%; border-collapse: collapse; margin: 1rem 0;">
<thead>
<tr style="background: #1a1a2e; color: white;">
<th style="padding: 0.75rem; text-align: left;">Status</th>
<th style="padding: 0.75rem; text-align: left;">Code</th>
<th style="padding: 0.75rem; text-align: left;">Meaning</th>
</tr>
</thead>
<tbody>
<tr style="background: #ecfdf5;">
<td style="padding: 0.75rem;">🟢 Success</td>
<td style="padding: 0.75rem;"><code>200</code></td>
<td style="padding: 0.75rem;">Request succeeded</td>
</tr>
<tr style="background: #fffbeb;">
<td style="padding: 0.75rem;">🟡 Redirect</td>
<td style="padding: 0.75rem;"><code>301</code></td>
<td style="padding: 0.75rem;">Moved permanently</td>
</tr>
<tr style="background: #fef2f2;">
<td style="padding: 0.75rem;">🔴 Error</td>
<td style="padding: 0.75rem;"><code>500</code></td>
<td style="padding: 0.75rem;">Internal server error</td>
</tr>
</tbody>
</table>

## Keyboard Shortcuts Display

Use `<kbd>` tags for keyboard shortcuts:

Press <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>P</kbd> to open the command palette.

Use <kbd>⌘</kbd> + <kbd>K</kbd> on macOS.

## Mixing Markdown Inside HTML

> [!warning] Important
> Markdown inside HTML block elements is **not parsed** unless there's a blank line separating them. Inline HTML like `<kbd>`, `<mark>`, `<abbr>` works anywhere.

This text has <mark>highlighted words</mark> and <abbr title="Static Site Generator">SSG</abbr> abbreviations mixed in naturally.

## Definition List (HTML)

<dl>
<dt><strong>SSG</strong></dt>
<dd>Static Site Generator — a tool that produces HTML files from templates and content at build time.</dd>
<dt><strong>CDN</strong></dt>
<dd>Content Delivery Network — a distributed network of servers that cache and serve static content.</dd>
<dt><strong>SSE</strong></dt>
<dd>Server-Sent Events — a web API for receiving push notifications from a server via HTTP.</dd>
</dl>

[TITLE]: # (Pico CSS Demo)
[DATE]: # (2026-05-25)
[TAGS]: # (css, tutorial)
[INCLUDES]: # (H, F, /assets/pico/pico.min.css)
[INHERITS]: # (post.html)

# Using Pico CSS Alongside Tailwind

This page demonstrates loading an external CSS framework — [Pico CSS](https://picocss.com) — via the `INCLUDES` directive. The directive `[INCLUDES]: # (H, F, /assets/pico/pico.min.css)` adds a `<link>` tag for Pico's stylesheet.

## How It Works

In your markdown file, add any `.css` or `.js` file to the `INCLUDES` directive:

```markdown
[INCLUDES]: # (H, F, /assets/pico/pico.min.css)
```

Sitemaker detects the file extension and injects it as a `<link rel="stylesheet">` in the `<head>`. JS files get `<script>` tags before `</body>`.

## Pico-Styled Elements

Pico CSS styles semantic HTML automatically. Here are some elements that Pico enhances:

### Forms

<form>
  <label for="name">Full Name</label>
  <input type="text" id="name" placeholder="Dana Kim">
  
  <label for="email">Email</label>
  <input type="email" id="email" placeholder="dana@example.com">
  
  <label for="bio">Bio</label>
  <textarea id="bio" placeholder="Tell us about yourself..."></textarea>
  
  <fieldset>
    <legend>Notifications</legend>
    <label><input type="checkbox" checked> Email notifications</label>
    <label><input type="checkbox"> SMS notifications</label>
  </fieldset>
  
  <button type="submit">Submit</button>
</form>

### Progress Bars

<progress value="75" max="100"></progress>

### Details/Summary

<details>
  <summary>Click to expand</summary>
  <p>This is hidden content revealed by clicking the summary. Pico CSS styles this natively with smooth transitions.</p>
</details>

<details>
  <summary>Another expandable section</summary>
  <p>You can nest any markdown or HTML content inside these elements.</p>
</details>

## Mixing Pico + Tailwind

Since Tailwind uses utility classes and Pico targets semantic HTML, they coexist well. Tailwind handles layout, Pico handles element styling.

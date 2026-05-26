[TITLE]: # (Interactive JavaScript in Markdown)
[DATE]: # (2026-05-22)
[TAGS]: # (javascript, tutorial)
[INCLUDES]: # (H, F, TOC)
[INHERITS]: # (post.html)

# Interactive JavaScript in Markdown

Sitemaker allows raw HTML and `<script>` blocks directly in your markdown files. This page demonstrates several patterns for adding interactivity.

## Inline Script Block

You can drop a `<script>` tag anywhere in your markdown. Here's a live counter:

<div id="counter-demo" style="padding: 1.5rem; background: #f0f9ff; border-radius: 0.5rem; margin: 1rem 0; text-align: center;">
  <p style="font-size: 2rem; font-weight: bold; margin: 0;" id="count">0</p>
  <div style="display: flex; gap: 0.5rem; justify-content: center; margin-top: 0.5rem;">
    <button onclick="decrement()" style="padding: 0.5rem 1rem; background: #e5e7eb; border: none; border-radius: 0.25rem; cursor: pointer;">−</button>
    <button onclick="increment()" style="padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 0.25rem; cursor: pointer;">+</button>
    <button onclick="resetCounter()" style="padding: 0.5rem 1rem; background: #f3f4f6; border: 1px solid #d1d5db; border-radius: 0.25rem; cursor: pointer;">Reset</button>
  </div>
</div>

<script>
let count = 0;
function increment() { document.getElementById('count').textContent = ++count; }
function decrement() { document.getElementById('count').textContent = --count; }
function resetCounter() { count = 0; document.getElementById('count').textContent = '0'; }
</script>

The code for that counter:

```html
<div id="counter-demo">
  <p id="count">0</p>
  <button onclick="increment()">+</button>
</div>

<script>
let count = 0;
function increment() {
  document.getElementById('count').textContent = ++count;
}
</script>
```

## Including External JS via INCLUDES

For larger scripts, put the file in `static/` and reference it:

```markdown
[INCLUDES]: # (H, F, /js/charts.js)
```

This injects `<script src="/js/charts.js"></script>` before `</body>`.

## Dynamic Content Example

Here's a clock that updates every second:

<div id="clock" style="font-family: monospace; font-size: 1.5rem; padding: 1rem; background: #1a1a2e; color: #10b981; border-radius: 0.5rem; text-align: center; margin: 1rem 0;"></div>

<script>
function updateClock() {
  var now = new Date();
  var time = now.toLocaleTimeString('en-US', { hour12: false });
  var date = now.toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });
  document.getElementById('clock').innerHTML = time + '<br><span style="font-size: 0.75rem; color: #6b7280;">' + date + '</span>';
}
updateClock();
setInterval(updateClock, 1000);
</script>

## Collapsible Code Viewer

<div style="border: 1px solid #e5e7eb; border-radius: 0.5rem; margin: 1rem 0;">
  <button onclick="this.nextElementSibling.style.display = this.nextElementSibling.style.display === 'none' ? 'block' : 'none'" style="width: 100%; padding: 0.75rem 1rem; background: #f9fafb; border: none; text-align: left; cursor: pointer; font-weight: 600;">
    Toggle Source Code
  </button>
  <div style="display: none; padding: 1rem; border-top: 1px solid #e5e7eb;">

```javascript
// This code block is inside a toggleable HTML container
function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

for (let i = 0; i < 10; i++) {
    console.log(`fib(${i}) = ${fibonacci(i)}`);
}
```

  </div>
</div>

> [!tip] Script Safety
> Since sitemaker generates static sites from your own content, raw HTML and scripts are safe. You are the author, not an untrusted user.

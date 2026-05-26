import MarkdownIt from "markdown-it";
import footnotePlugin from "markdown-it-footnote";
import taskListPlugin from "markdown-it-task-lists";
import hljs from "highlight.js";

function wikilinkPlugin(md: MarkdownIt): void {
  md.inline.ruler.after("emphasis", "wikilink", (state, silent) => {
    const src = state.src;
    const pos = state.pos;

    if (src[pos] !== "[" || src[pos + 1] !== "[") return false;

    const end = src.indexOf("]]", pos + 2);
    if (end < 0) return false;

    if (!silent) {
      const content = src.slice(pos + 2, end);
      let target = content;
      let display = content;

      const pipeIdx = content.indexOf("|");
      if (pipeIdx >= 0) {
        target = content.slice(0, pipeIdx).trim();
        display = content.slice(pipeIdx + 1).trim();
      }

      const slug = target
        .toLowerCase()
        .replace(/\s+/g, "-")
        .replace(/[^a-z0-9\-/]/g, "");

      const tokenOpen = state.push("wikilink_open", "a", 1);
      tokenOpen.attrSet("href", `/${slug}.html`);
      tokenOpen.attrSet("class", "wikilink");

      const tokenText = state.push("text", "", 0);
      tokenText.content = display;

      state.push("wikilink_close", "a", -1);
    }

    state.pos = end + 2;
    return true;
  });
}

function calloutPlugin(md: MarkdownIt): void {
  const defaultRender =
    md.renderer.rules.blockquote_open ||
    function (tokens, idx, options, _env, self) {
      return self.renderToken(tokens, idx, options);
    };

  const defaultCloseRender =
    md.renderer.rules.blockquote_close ||
    function (tokens, idx, options, _env, self) {
      return self.renderToken(tokens, idx, options);
    };

  const CALLOUT_RE = /^\[!(\w+)\]\s*(.*)?$/;
  const ICONS: Record<string, string> = {
    note: "&#9998;",
    tip: "&#128161;",
    warning: "&#9888;",
    danger: "&#9762;",
    info: "&#8505;",
    example: "&#128203;",
    quote: "&#10077;",
    bug: "&#128027;",
    todo: "&#9745;",
    success: "&#10004;",
    question: "&#10067;",
    failure: "&#10008;",
    abstract: "&#128196;",
  };

  md.renderer.rules.blockquote_open = function (
    tokens,
    idx,
    options,
    env,
    self
  ) {
    const contentToken = findBlockquoteContent(tokens, idx);
    if (!contentToken) return defaultRender(tokens, idx, options, env, self);

    const match = contentToken.content.match(CALLOUT_RE);
    if (!match) return defaultRender(tokens, idx, options, env, self);

    const calloutType = match[1].toLowerCase();
    const title =
      match[2] || calloutType.charAt(0).toUpperCase() + calloutType.slice(1);
    const icon = ICONS[calloutType] || "";

    contentToken.content = contentToken.content.replace(
      CALLOUT_RE,
      ""
    ).trim();
    if (contentToken.content === "") {
      contentToken.hidden = true;
    }

    (tokens[idx] as any).__callout = true;

    const iconHtml = icon
      ? `<span class="callout-icon">${icon}</span> `
      : "";
    return `<div class="callout callout-${calloutType}">\n<div class="callout-title">${iconHtml}${md.utils.escapeHtml(title)}</div>\n<div class="callout-content">\n`;
  };

  md.renderer.rules.blockquote_close = function (
    tokens,
    idx,
    options,
    env,
    self
  ) {
    const openIdx = findMatchingOpen(tokens, idx);
    if (openIdx >= 0 && (tokens[openIdx] as any).__callout) {
      return `</div>\n</div>\n`;
    }
    return defaultCloseRender(tokens, idx, options, env, self);
  };
}

function findBlockquoteContent(
  tokens: any[],
  openIdx: number
): any | null {
  let depth = 0;
  for (let i = openIdx + 1; i < tokens.length; i++) {
    if (tokens[i].type === "blockquote_open") depth++;
    if (tokens[i].type === "blockquote_close") {
      if (depth === 0) return null;
      depth--;
    }
    if (depth === 0 && tokens[i].type === "inline" && tokens[i].content) {
      return tokens[i];
    }
  }
  return null;
}

function findMatchingOpen(tokens: any[], closeIdx: number): number {
  let depth = 0;
  for (let i = closeIdx - 1; i >= 0; i--) {
    if (tokens[i].type === "blockquote_close") depth++;
    if (tokens[i].type === "blockquote_open") {
      if (depth === 0) return i;
      depth--;
    }
  }
  return -1;
}

export function createRenderer(): MarkdownIt {
  const md = new MarkdownIt({
    html: true,
    linkify: true,
    typographer: true,
    highlight(str: string, lang: string): string {
      if (lang && hljs.getLanguage(lang)) {
        try {
          const result = hljs.highlight(str, { language: lang });
          return `<pre class="hljs"><code class="language-${lang}">${result.value}</code></pre>`;
        } catch {
          // fall through
        }
      }
      return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`;
    },
  });

  md.use(footnotePlugin);
  md.use(taskListPlugin, { enabled: true, label: true });
  wikilinkPlugin(md);
  calloutPlugin(md);

  return md;
}

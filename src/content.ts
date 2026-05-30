import fs from "node:fs";
import path from "node:path";
import { Config } from "./config.js";
import { Directives, parseDirectives, hasComponent, assetFiles } from "./directive.js";

export interface Page {
  sourcePath: string;
  rawContent: string;
  body: string;
  directives: Directives;
  htmlContent: string;
  tocContent: string;
  title: string;
  slug: string;
  url: string;
  date: Date | null;
  dateStr: string;
  tags: string[];
  readingTime: number;
  isBlogPost: boolean;
  draft: boolean;
  hasHeader: boolean;
  hasFooter: boolean;
  hasTOC: boolean;
  hasSidebar: boolean;
  cssAssets: string[];
  jsAssets: string[];
  prev: { title: string; url: string } | null;
  next: { title: string; url: string } | null;
}

export interface Site {
  config: Config;
  pages: Page[];
  posts: Page[];
  tagIndex: Record<string, Page[]>;
}

const DATE_PREFIX_RE = /^(\d{4}-\d{2}-\d{2})-(.+)$/;

export function loadSite(cfg: Config): Site {
  const site: Site = {
    config: cfg,
    pages: [],
    posts: [],
    tagIndex: {},
  };

  walkDir(cfg.build.content_dir, cfg.build.content_dir, site, cfg);
  return site;
}

function walkDir(dir: string, contentDir: string, site: Site, cfg: Config): void {
  if (!fs.existsSync(dir)) return;

  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walkDir(fullPath, contentDir, site, cfg);
    } else if (entry.name.endsWith(".md")) {
      const page = loadPage(fullPath, contentDir, cfg);
      site.pages.push(page);
    }
  }
}

function loadPage(filePath: string, contentDir: string, cfg: Config): Page {
  const raw = fs.readFileSync(filePath, "utf-8");
  const relPath = path.relative(contentDir, filePath).replace(/\\/g, "/");
  const { directives, body } = parseDirectives(raw);

  const isBlogPost = relPath.startsWith(cfg.blog.posts_dir + "/");
  const title = resolveTitle(directives, body, relPath);
  const slug = resolveSlug(directives, relPath, isBlogPost, cfg);
  const date = resolveDate(directives, relPath);
  const url = resolveURL(relPath, slug, cfg);

  const assets = assetFiles(directives);

  return {
    sourcePath: relPath,
    rawContent: raw,
    body,
    directives,
    htmlContent: "",
    tocContent: "",
    title,
    slug,
    url,
    date,
    dateStr: date ? formatDate(date) : "",
    tags: directives.tags,
    readingTime: estimateReadingTime(body),
    isBlogPost,
    draft: directives.draft,
    hasHeader: hasComponent(directives, "H"),
    hasFooter: hasComponent(directives, "F"),
    hasTOC: hasComponent(directives, "TOC"),
    hasSidebar: hasComponent(directives, "SB"),
    cssAssets: assets.filter((a) => a.endsWith(".css")),
    jsAssets: assets.filter((a) => a.endsWith(".js")),
    prev: null,
    next: null,
  };
}

function resolveTitle(d: Directives, body: string, relPath: string): string {
  if (d.title) return d.title;
  const match = body.match(/^#\s+(.+)$/m);
  if (match) return match[1];
  const base = path.basename(relPath, ".md");
  return base;
}

function resolveSlug(
  d: Directives,
  relPath: string,
  isBlogPost: boolean,
  cfg: Config
): string {
  if (d.slug) return d.slug;
  let base = path.basename(relPath, ".md");
  if (isBlogPost && cfg.urls.strip_date) {
    const m = base.match(DATE_PREFIX_RE);
    if (m) base = m[2];
  }
  return base;
}

function resolveDate(d: Directives, relPath: string): Date | null {
  if (d.date) {
    const t = new Date(d.date + "T00:00:00Z");
    if (!isNaN(t.getTime())) return t;
  }
  const base = path.basename(relPath, ".md");
  const m = base.match(DATE_PREFIX_RE);
  if (m) {
    const t = new Date(m[1] + "T00:00:00Z");
    if (!isNaN(t.getTime())) return t;
  }
  return null;
}

function resolveURL(relPath: string, slug: string, cfg: Config): string {
  const dir = path.dirname(relPath);
  const dirPrefix = dir === "." ? "" : dir;

  if (slug === "index") {
    return dirPrefix === "" ? "/" : `/${dirPrefix}/`;
  }

  const urlPath = dirPrefix ? `${dirPrefix}/${slug}` : slug;

  if (cfg.urls.style === "directory") {
    return `/${urlPath}/`;
  }
  return `/${urlPath}.html`;
}

export function outputPath(page: Page, cfg: Config): string {
  if (page.slug === "index") {
    const dir = path.dirname(page.sourcePath);
    if (dir === ".") return "index.html";
    return path.join(dir, "index.html");
  }

  const dir = path.dirname(page.sourcePath);
  const dirPrefix = dir === "." ? "" : dir;

  if (cfg.urls.style === "directory") {
    return path.join(dirPrefix, page.slug, "index.html");
  }
  return path.join(dirPrefix, page.slug + ".html");
}

export function classifyPages(site: Site): void {
  site.posts = site.pages
    .filter((p) => p.isBlogPost && !p.draft)
    .sort((a, b) => {
      const da = a.date?.getTime() ?? 0;
      const db = b.date?.getTime() ?? 0;
      return db - da;
    });
}

export function buildTagIndex(site: Site): void {
  site.tagIndex = {};
  for (const p of site.posts) {
    for (const tag of p.tags) {
      if (!site.tagIndex[tag]) site.tagIndex[tag] = [];
      site.tagIndex[tag].push(p);
    }
  }
}

export function linkPrevNext(site: Site): void {
  for (let i = 0; i < site.posts.length; i++) {
    if (i > 0) {
      site.posts[i].next = {
        title: site.posts[i - 1].title,
        url: site.posts[i - 1].url,
      };
    }
    if (i < site.posts.length - 1) {
      site.posts[i].prev = {
        title: site.posts[i + 1].title,
        url: site.posts[i + 1].url,
      };
    }
  }
}

function estimateReadingTime(content: string): number {
  const words = content.split(/\s+/).length;
  return Math.max(1, Math.ceil(words / 200));
}

function formatDate(d: Date): string {
  return d.toISOString().slice(0, 10);
}

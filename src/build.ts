import fs from "node:fs";
import path from "node:path";
import { Config } from "./config.js";
import {
  loadSite,
  classifyPages,
  buildTagIndex,
  linkPrevNext,
  outputPath,
} from "./content.js";
import { createRenderer } from "./markdown.js";
import { TemplateEngine } from "./template.js";
import { extractHeadings, renderTOC } from "./toc.js";
import { generateRSS } from "./rss.js";

export function build(cfg: Config, defaultsDir: string): void {
  const site = loadSite(cfg);
  classifyPages(site);
  buildTagIndex(site);
  linkPrevNext(site);

  const md = createRenderer();

  for (const page of site.pages) {
    if (page.directives.override) continue;

    page.htmlContent = md.render(page.body);

    if (page.hasTOC) {
      const headings = extractHeadings(page.htmlContent);
      page.tocContent = renderTOC(headings);
    }
  }

  const engine = new TemplateEngine(cfg, defaultsDir);

  // Clean output
  if (fs.existsSync(cfg.build.output_dir)) {
    fs.rmSync(cfg.build.output_dir, { recursive: true, force: true });
  }

  for (const page of site.pages) {
    if (page.draft) continue;

    const html = engine.render(page, site);
    const outPath = path.join(cfg.build.output_dir, outputPath(page, cfg));
    writeFile(outPath, html);
  }

  // Post listing
  if (site.posts.length > 0) {
    const listHTML = engine.renderList(site.posts, site, "Posts");
    const listPath = path.join(
      cfg.build.output_dir,
      cfg.blog.posts_dir,
      "index.html"
    );
    writeFile(listPath, listHTML);
  }

  // Tag pages
  for (const [tag, posts] of Object.entries(site.tagIndex)) {
    const tagHTML = engine.renderList(posts, site, `Tag: ${tag}`);
    const tagPath = path.join(cfg.build.output_dir, "tags", tag, "index.html");
    writeFile(tagPath, tagHTML);
  }

  // RSS
  if (cfg.blog.enable_rss && site.posts.length > 0) {
    generateRSS(site, cfg);
  }

  // Static files
  copyStatic(cfg);
  copyDefaultAssets(cfg, defaultsDir);
  writeHeaders(cfg);
}

function writeFile(filePath: string, content: string): void {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });
  fs.writeFileSync(filePath, content);
}

function copyStatic(cfg: Config): void {
  const staticDir = cfg.build.static_dir;
  if (!fs.existsSync(staticDir)) return;

  copyDirRecursive(staticDir, cfg.build.output_dir);
}

function copyDefaultAssets(cfg: Config, defaultsDir: string): void {
  for (const dir of ["css", "js"]) {
    const srcDir = path.join(defaultsDir, dir);
    if (!fs.existsSync(srcDir)) continue;

    const destDir = path.join(cfg.build.output_dir, "assets", dir);
    copyDirRecursive(srcDir, destDir);
  }
}

function copyDirRecursive(src: string, dest: string): void {
  if (!fs.existsSync(src)) return;

  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);

    if (entry.isDirectory()) {
      copyDirRecursive(srcPath, destPath);
    } else {
      fs.mkdirSync(path.dirname(destPath), { recursive: true });
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

function writeHeaders(cfg: Config): void {
  const headers = `/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY

/assets/*
  Cache-Control: public, max-age=31536000, immutable

/*.html
  Cache-Control: public, max-age=0, must-revalidate
`;
  writeFile(path.join(cfg.build.output_dir, "_headers"), headers);
}

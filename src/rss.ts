import fs from "node:fs";
import path from "node:path";
import { Config } from "./config.js";
import { Site } from "./content.js";

export function generateRSS(site: Site, cfg: Config): void {
  const posts = site.posts.slice(0, 20);

  let xml = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
  <title>${escXml(cfg.site.title)}</title>
  <link>${escXml(cfg.site.base_url)}</link>
  <description>${escXml(cfg.site.description)}</description>
  <language>${escXml(cfg.site.language)}</language>
  <lastBuildDate>${new Date().toUTCString()}</lastBuildDate>
  <atom:link href="${escXml(cfg.site.base_url)}rss.xml" rel="self" type="application/rss+xml"/>
`;

  for (const post of posts) {
    const fullURL = cfg.site.base_url.replace(/\/$/, "") + post.url;
    xml += `  <item>
    <title>${escXml(post.title)}</title>
    <link>${escXml(fullURL)}</link>
    <description><![CDATA[${post.htmlContent}]]></description>
    <pubDate>${post.date?.toUTCString() ?? ""}</pubDate>
    <guid>${escXml(fullURL)}</guid>
`;
    for (const tag of post.tags) {
      xml += `    <category>${escXml(tag)}</category>\n`;
    }
    xml += `  </item>\n`;
  }

  xml += `</channel>\n</rss>\n`;

  const outDir = cfg.build.output_dir;
  fs.mkdirSync(outDir, { recursive: true });
  fs.writeFileSync(path.join(outDir, "rss.xml"), xml);
  fs.writeFileSync(path.join(outDir, "feed.xml"), xml);
}

function escXml(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

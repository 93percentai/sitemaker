import fs from "node:fs";
import http from "node:http";
import path from "node:path";
import { Config } from "./config.js";
import { build } from "./build.js";

const MIME_TYPES: Record<string, string> = {
  ".html": "text/html; charset=utf-8",
  ".css": "text/css",
  ".js": "application/javascript",
  ".json": "application/json",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".gif": "image/gif",
  ".svg": "image/svg+xml",
  ".xml": "application/xml",
  ".ico": "image/x-icon",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
};

export async function serve(
  cfg: Config,
  defaultsDir: string,
  port: number
): Promise<void> {
  console.log("Building site...");
  build(cfg, defaultsDir);

  const clients = new Set<http.ServerResponse>();

  // File watcher
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  const watchDirs = [cfg.build.content_dir, cfg.build.template_dir, cfg.build.static_dir];

  for (const dir of watchDirs) {
    if (fs.existsSync(dir)) {
      fs.watch(dir, { recursive: true }, () => {
        if (debounceTimer) clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
          console.log("Rebuilding...");
          try {
            build(cfg, defaultsDir);
            console.log("Build complete.");
            for (const res of clients) {
              res.write("data: reload\n\n");
            }
          } catch (err) {
            console.error("Build error:", err);
          }
        }, 200);
      });
    }
  }

  const server = http.createServer((req, res) => {
    const urlPath = req.url || "/";

    // SSE endpoint for live reload
    if (urlPath === "/__reload") {
      res.writeHead(200, {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
        "Access-Control-Allow-Origin": "*",
      });
      clients.add(res);
      req.on("close", () => clients.delete(res));
      return;
    }

    let filePath = path.join(cfg.build.output_dir, urlPath);

    // Directory → index.html
    if (
      fs.existsSync(filePath) &&
      fs.statSync(filePath).isDirectory()
    ) {
      filePath = path.join(filePath, "index.html");
    }

    // Try without trailing slash
    if (!fs.existsSync(filePath) && !path.extname(filePath)) {
      const withHtml = filePath + ".html";
      if (fs.existsSync(withHtml)) filePath = withHtml;
    }

    if (!fs.existsSync(filePath)) {
      res.writeHead(404, { "Content-Type": "text/plain" });
      res.end("404 Not Found");
      return;
    }

    const ext = path.extname(filePath);
    const contentType = MIME_TYPES[ext] || "application/octet-stream";

    let content = fs.readFileSync(filePath);

    // Inject live reload script into HTML
    if (ext === ".html") {
      const html = content.toString();
      const reloadScript = `<script>
(function() {
  var es = new EventSource('/__reload');
  es.onmessage = function() { location.reload(); };
  es.onerror = function() { setTimeout(function() { location.reload(); }, 1000); };
})();
</script>`;
      content = Buffer.from(html.replace("</body>", reloadScript + "\n</body>"));
    }

    res.writeHead(200, { "Content-Type": contentType });
    res.end(content);
  });

  server.listen(port, () => {
    console.log(`Serving at http://localhost:${port}`);
    console.log("Live reload enabled. Watching for changes...");
  });
}

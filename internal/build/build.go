package build

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/93percentai/sitemaker/internal/config"
	"github.com/93percentai/sitemaker/internal/content"
	"github.com/93percentai/sitemaker/internal/markdown"
	"github.com/93percentai/sitemaker/internal/rss"
	tmpl "github.com/93percentai/sitemaker/internal/template"
	"github.com/93percentai/sitemaker/internal/toc"

	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func Build(cfg config.Config, defaults fs.FS) error {
	site, err := content.LoadSite(cfg)
	if err != nil {
		return fmt.Errorf("loading site: %w", err)
	}

	site.ClassifyPages()
	site.BuildTagIndex()
	site.LinkPrevNext()

	md := markdown.NewRenderer()

	for _, page := range site.Pages {
		if page.Directives.Override != "" {
			continue
		}

		html, err := md.Render(page.Body)
		if err != nil {
			return fmt.Errorf("rendering markdown for %s: %w", page.SourcePath, err)
		}
		page.HTMLContent = template.HTML(html)

		if page.HasTOC() {
			reader := text.NewReader(page.Body)
			doc := md.Parser().Parse(reader, parser.WithContext(parser.NewContext()))
			headings := toc.Extract(doc, page.Body)
			page.TOCContent = toc.RenderTOC(headings)
		}
	}

	engine := tmpl.NewEngine(cfg, defaults)

	if err := os.RemoveAll(cfg.Build.OutputDir); err != nil {
		return fmt.Errorf("cleaning output dir: %w", err)
	}

	for _, page := range site.Pages {
		if page.Draft {
			continue
		}

		output, err := engine.Render(page, site)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", page.SourcePath, err)
		}

		outPath := filepath.Join(cfg.Build.OutputDir, page.OutputPath(cfg))
		if err := writeFile(outPath, output); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
	}

	if len(site.Posts) > 0 {
		listHTML, err := engine.RenderList(site.Posts, site, "Posts")
		if err != nil {
			return fmt.Errorf("rendering post list: %w", err)
		}
		listPath := filepath.Join(cfg.Build.OutputDir, cfg.Blog.PostsDir, "index.html")
		if err := writeFile(listPath, listHTML); err != nil {
			return fmt.Errorf("writing post list: %w", err)
		}
	}

	for tag, posts := range site.TagIndex {
		tagHTML, err := engine.RenderList(posts, site, "Tag: "+tag)
		if err != nil {
			return fmt.Errorf("rendering tag page %s: %w", tag, err)
		}
		tagPath := filepath.Join(cfg.Build.OutputDir, "tags", tag, "index.html")
		if err := writeFile(tagPath, tagHTML); err != nil {
			return fmt.Errorf("writing tag page: %w", err)
		}
	}

	if cfg.Blog.EnableRSS && len(site.Posts) > 0 {
		if err := rss.Generate(site, cfg); err != nil {
			return fmt.Errorf("generating RSS: %w", err)
		}
	}

	if err := copyStatic(cfg); err != nil {
		return fmt.Errorf("copying static files: %w", err)
	}

	if err := copyDefaultAssets(cfg, defaults); err != nil {
		return fmt.Errorf("copying default assets: %w", err)
	}

	if err := writeHeaders(cfg); err != nil {
		return fmt.Errorf("writing headers: %w", err)
	}

	return nil
}

func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func copyStatic(cfg config.Config) error {
	staticDir := cfg.Build.StaticDir
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		relPath, _ := filepath.Rel(staticDir, path)
		destPath := filepath.Join(cfg.Build.OutputDir, relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return writeFile(destPath, data)
	})
}

func copyDefaultAssets(cfg config.Config, defaults fs.FS) error {
	for _, dir := range []string{"css", "js"} {
		err := fs.WalkDir(defaults, dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			data, err := fs.ReadFile(defaults, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(cfg.Build.OutputDir, "assets", path)
			return writeFile(destPath, data)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func writeHeaders(cfg config.Config) error {
	headers := `/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY

/assets/*
  Cache-Control: public, max-age=31536000, immutable

/*.html
  Cache-Control: public, max-age=0, must-revalidate
`
	return writeFile(filepath.Join(cfg.Build.OutputDir, "_headers"), []byte(headers))
}

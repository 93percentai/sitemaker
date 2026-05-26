package content

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/93percentai/sitemaker/internal/config"
	"github.com/93percentai/sitemaker/internal/directive"
)

type Page struct {
	SourcePath  string
	RawContent  []byte
	Body        []byte
	Directives  directive.Directives
	HTMLContent template.HTML
	TOCContent  template.HTML

	Title       string
	Slug        string
	URL         string
	Date        time.Time
	Tags        []string
	ReadingTime int
	IsBlogPost  bool
	Draft       bool

	Prev *Page
	Next *Page
}

type Site struct {
	Config   config.Config
	Pages    []*Page
	Posts    []*Page
	TagIndex map[string][]*Page
}

var datePrefix = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.+)$`)

func LoadSite(cfg config.Config) (*Site, error) {
	site := &Site{
		Config:   cfg,
		TagIndex: make(map[string][]*Page),
	}

	contentDir := cfg.Build.ContentDir
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		page, err := loadPage(path, contentDir, cfg)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}

		site.Pages = append(site.Pages, page)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return site, nil
}

func loadPage(path, contentDir string, cfg config.Config) (*Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(contentDir, path)
	relPath = filepath.ToSlash(relPath)

	directives, body := directive.Parse(data)

	page := &Page{
		SourcePath: relPath,
		RawContent: data,
		Body:       body,
		Directives: directives,
		Draft:      directives.Draft,
		Tags:       directives.Tags,
	}

	page.IsBlogPost = isPost(relPath, cfg)
	page.resolveTitle()
	page.resolveSlug(cfg)
	page.resolveDate(cfg)
	page.resolveURL(cfg)
	page.ReadingTime = estimateReadingTime(body)

	return page, nil
}

func isPost(relPath string, cfg config.Config) bool {
	return strings.HasPrefix(relPath, cfg.Blog.PostsDir+"/")
}

func (p *Page) resolveTitle() {
	if p.Directives.Title != "" {
		p.Title = p.Directives.Title
		return
	}
	for _, line := range strings.Split(string(p.Body), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			p.Title = strings.TrimPrefix(line, "# ")
			return
		}
	}
	base := filepath.Base(p.SourcePath)
	p.Title = strings.TrimSuffix(base, filepath.Ext(base))
}

func (p *Page) resolveSlug(cfg config.Config) {
	if p.Directives.Slug != "" {
		p.Slug = p.Directives.Slug
		return
	}

	base := strings.TrimSuffix(filepath.Base(p.SourcePath), ".md")

	if p.IsBlogPost && cfg.URLs.StripDate {
		if m := datePrefix.FindStringSubmatch(base); m != nil {
			base = m[2]
		}
	}

	p.Slug = base
}

func (p *Page) resolveDate(cfg config.Config) {
	if p.Directives.Date != "" {
		t, err := time.Parse(cfg.Blog.DateFormat, p.Directives.Date)
		if err == nil {
			p.Date = t
			return
		}
	}

	base := strings.TrimSuffix(filepath.Base(p.SourcePath), ".md")
	if m := datePrefix.FindStringSubmatch(base); m != nil {
		t, err := time.Parse("2006-01-02", m[1])
		if err == nil {
			p.Date = t
		}
	}
}

func (p *Page) resolveURL(cfg config.Config) {
	dir := filepath.Dir(p.SourcePath)
	if dir == "." {
		dir = ""
	}

	slug := p.Slug

	if slug == "index" {
		if dir == "" {
			p.URL = "/"
		} else {
			p.URL = "/" + dir + "/"
		}
		return
	}

	var urlPath string
	if dir != "" {
		urlPath = dir + "/" + slug
	} else {
		urlPath = slug
	}

	switch cfg.URLs.Style {
	case "directory":
		p.URL = "/" + urlPath + "/"
	default:
		p.URL = "/" + urlPath + ".html"
	}
}

func (p *Page) OutputPath(cfg config.Config) string {
	if p.Slug == "index" {
		dir := filepath.Dir(p.SourcePath)
		if dir == "." {
			return "index.html"
		}
		return filepath.Join(dir, "index.html")
	}

	dir := filepath.Dir(p.SourcePath)
	if dir == "." {
		dir = ""
	}

	switch cfg.URLs.Style {
	case "directory":
		return filepath.Join(dir, p.Slug, "index.html")
	default:
		return filepath.Join(dir, p.Slug+".html")
	}
}

func (p *Page) HasHeader() bool   { return p.Directives.HasComponent("H") }
func (p *Page) HasFooter() bool   { return p.Directives.HasComponent("F") }
func (p *Page) HasTOC() bool      { return p.Directives.HasComponent("TOC") }
func (p *Page) HasSidebar() bool  { return p.Directives.HasComponent("SB") }

func (p *Page) CSSAssets() []string {
	var css []string
	for _, a := range p.Directives.AssetFiles() {
		if strings.HasSuffix(a, ".css") {
			css = append(css, a)
		}
	}
	return css
}

func (p *Page) JSAssets() []string {
	var js []string
	for _, a := range p.Directives.AssetFiles() {
		if strings.HasSuffix(a, ".js") {
			js = append(js, a)
		}
	}
	return js
}

func estimateReadingTime(content []byte) int {
	words := utf8.RuneCount(content) / 5
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

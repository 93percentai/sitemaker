package tmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/93percentai/sitemaker/internal/config"
	"github.com/93percentai/sitemaker/internal/content"
)

type Engine struct {
	cfg        config.Config
	defaults   fs.FS
	funcMap    template.FuncMap
}

type TemplateData struct {
	Site       SiteData
	Page       PageData
	Posts      []PageData
	Tags       map[string][]PageData
	CurrentTag string
}

type SiteData struct {
	Title       string
	Description string
	BaseURL     string
	Author      string
	Language    string
	Nav         []config.NavItem
	BuildTime   time.Time
}

type PageData struct {
	Title       string
	HTMLContent template.HTML
	TOCContent  template.HTML
	URL         string
	Date        time.Time
	DateStr     string
	Tags        []string
	ReadingTime int
	HasHeader   bool
	HasFooter   bool
	HasTOC      bool
	HasSidebar  bool
	CSSAssets   []string
	JSAssets    []string
	IsBlogPost  bool
	Prev        *PageRef
	Next        *PageRef
}

type PageRef struct {
	Title string
	URL   string
}

func NewEngine(cfg config.Config, defaults fs.FS) *Engine {
	return &Engine{
		cfg:      cfg,
		defaults: defaults,
		funcMap:  FuncMap(),
	}
}

func (e *Engine) Render(page *content.Page, site *content.Site) ([]byte, error) {
	if page.Directives.Override != "" {
		return e.renderOverride(page, site)
	}

	tmplName := page.Directives.Inherits
	if tmplName == "" {
		if page.IsBlogPost {
			tmplName = "post.html"
		} else {
			tmplName = "base.html"
		}
	}

	tmpl := template.New("page").Funcs(e.funcMap)

	baseContent, err := e.loadTemplate(tmplName)
	if err != nil {
		return nil, fmt.Errorf("loading template %s: %w", tmplName, err)
	}
	tmpl, err = tmpl.Parse(baseContent)
	if err != nil {
		return nil, fmt.Errorf("parsing template %s: %w", tmplName, err)
	}

	partials := map[string]string{
		"H":   "header.html",
		"F":   "footer.html",
		"TOC": "toc.html",
		"SB":  "sidebar.html",
	}

	for code, file := range partials {
		if page.Directives.HasComponent(code) {
			partialContent, err := e.loadPartial(file)
			if err != nil {
				return nil, fmt.Errorf("loading partial %s: %w", file, err)
			}
			tmpl, err = tmpl.Parse(partialContent)
			if err != nil {
				return nil, fmt.Errorf("parsing partial %s: %w", file, err)
			}
		}
	}

	emptyPartials := map[string]string{
		"H":   `{{define "header"}}{{end}}`,
		"F":   `{{define "footer"}}{{end}}`,
		"TOC": `{{define "toc"}}{{end}}`,
		"SB":  `{{define "sidebar"}}{{end}}`,
	}
	for code, fallback := range emptyPartials {
		if !page.Directives.HasComponent(code) {
			tmpl, _ = tmpl.Parse(fallback)
		}
	}

	data := e.buildData(page, site)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}

func (e *Engine) RenderList(pages []*content.Page, site *content.Site, title string) ([]byte, error) {
	tmplContent, err := e.loadTemplate("list.html")
	if err != nil {
		return nil, fmt.Errorf("loading list template: %w", err)
	}

	tmpl := template.New("page").Funcs(e.funcMap)
	tmpl, err = tmpl.Parse(tmplContent)
	if err != nil {
		return nil, fmt.Errorf("parsing list template: %w", err)
	}

	for _, partial := range []string{"header.html", "footer.html", "sidebar.html", "toc.html"} {
		content, perr := e.loadPartial(partial)
		if perr == nil {
			tmpl, _ = tmpl.Parse(content)
		}
	}

	emptyPartials := []string{
		`{{define "header"}}{{end}}`,
		`{{define "footer"}}{{end}}`,
		`{{define "toc"}}{{end}}`,
		`{{define "sidebar"}}{{end}}`,
	}
	for _, ep := range emptyPartials {
		if tmpl.Lookup(templateNameFromDefine(ep)) == nil {
			tmpl, _ = tmpl.Parse(ep)
		}
	}

	var posts []PageData
	for _, p := range pages {
		posts = append(posts, pageToData(p, e.cfg))
	}

	data := TemplateData{
		Site:  e.siteData(site),
		Posts: posts,
		Page: PageData{
			Title:     title,
			HasHeader: true,
			HasFooter: true,
		},
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing list template: %w", err)
	}

	return buf.Bytes(), nil
}

func templateNameFromDefine(s string) string {
	return ""
}

func (e *Engine) renderOverride(page *content.Page, site *content.Site) ([]byte, error) {
	overrideContent, err := e.loadUserTemplate(page.Directives.Override)
	if err != nil {
		return nil, fmt.Errorf("loading override %s: %w", page.Directives.Override, err)
	}

	tmpl := template.New("override").Funcs(e.funcMap)
	tmpl, err = tmpl.Parse(overrideContent)
	if err != nil {
		return nil, fmt.Errorf("parsing override: %w", err)
	}

	data := e.buildData(page, site)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing override: %w", err)
	}

	return buf.Bytes(), nil
}

func (e *Engine) buildData(page *content.Page, site *content.Site) TemplateData {
	return TemplateData{
		Site: e.siteData(site),
		Page: pageToData(page, e.cfg),
	}
}

func (e *Engine) siteData(site *content.Site) SiteData {
	return SiteData{
		Title:       site.Config.Site.Title,
		Description: site.Config.Site.Description,
		BaseURL:     site.Config.Site.BaseURL,
		Author:      site.Config.Site.Author,
		Language:    site.Config.Site.Language,
		Nav:         site.Config.Site.Nav,
		BuildTime:   time.Now(),
	}
}

func pageToData(page *content.Page, cfg config.Config) PageData {
	pd := PageData{
		Title:       page.Title,
		HTMLContent: page.HTMLContent,
		TOCContent:  page.TOCContent,
		URL:         page.URL,
		Date:        page.Date,
		Tags:        page.Tags,
		ReadingTime: page.ReadingTime,
		HasHeader:   page.HasHeader(),
		HasFooter:   page.HasFooter(),
		HasTOC:      page.HasTOC(),
		HasSidebar:  page.HasSidebar(),
		CSSAssets:   page.CSSAssets(),
		JSAssets:    page.JSAssets(),
		IsBlogPost:  page.IsBlogPost,
	}

	if !page.Date.IsZero() {
		pd.DateStr = page.Date.Format(cfg.Blog.DateFormat)
	}

	if page.Prev != nil {
		pd.Prev = &PageRef{Title: page.Prev.Title, URL: page.Prev.URL}
	}
	if page.Next != nil {
		pd.Next = &PageRef{Title: page.Next.Title, URL: page.Next.URL}
	}

	return pd
}

func (e *Engine) loadTemplate(name string) (string, error) {
	userPath := filepath.Join(e.cfg.Build.TemplateDir, name)
	if data, err := os.ReadFile(userPath); err == nil {
		return string(data), nil
	}

	tmplType := e.cfg.Template
	if tmplType == "" {
		tmplType = "personal"
	}

	data, err := fs.ReadFile(e.defaults, filepath.Join("templates", tmplType, name))
	if err == nil {
		return string(data), nil
	}

	data, err = fs.ReadFile(e.defaults, filepath.Join("templates", "partials", name))
	if err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("template %s not found", name)
}

func (e *Engine) loadPartial(name string) (string, error) {
	userPath := filepath.Join(e.cfg.Build.TemplateDir, "partials", name)
	if data, err := os.ReadFile(userPath); err == nil {
		return string(data), nil
	}

	data, err := fs.ReadFile(e.defaults, filepath.Join("templates", "partials", name))
	if err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("partial %s not found", name)
}

func (e *Engine) loadUserTemplate(name string) (string, error) {
	userPath := filepath.Join(e.cfg.Build.TemplateDir, name)
	data, err := os.ReadFile(userPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Site     SiteConfig     `toml:"site"`
	Build    BuildConfig    `toml:"build"`
	Blog     BlogConfig     `toml:"blog"`
	URLs     URLConfig      `toml:"urls"`
	Template string         `toml:"template"`
}

type SiteConfig struct {
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	BaseURL     string    `toml:"base_url"`
	Author      string    `toml:"author"`
	Language    string    `toml:"language"`
	Nav         []NavItem `toml:"nav"`
}

type NavItem struct {
	Label string `toml:"label"`
	URL   string `toml:"url"`
}

type BuildConfig struct {
	ContentDir  string `toml:"content_dir"`
	TemplateDir string `toml:"template_dir"`
	StaticDir   string `toml:"static_dir"`
	OutputDir   string `toml:"output_dir"`
}

type BlogConfig struct {
	PostsDir     string `toml:"posts_dir"`
	DateFormat   string `toml:"date_format"`
	PostsPerPage int    `toml:"posts_per_page"`
	EnableRSS    bool   `toml:"enable_rss"`
}

type URLConfig struct {
	Style     string `toml:"style"`
	StripDate bool   `toml:"strip_date"`
}

func Default() Config {
	return Config{
		Site: SiteConfig{
			Title:       "My Site",
			Description: "A site built with sitemaker",
			BaseURL:     "/",
			Language:    "en",
		},
		Build: BuildConfig{
			ContentDir:  "content",
			TemplateDir: "templates",
			StaticDir:   "static",
			OutputDir:   "dist",
		},
		Blog: BlogConfig{
			PostsDir:     "posts",
			DateFormat:   "2006-01-02",
			PostsPerPage: 10,
			EnableRSS:    true,
		},
		URLs: URLConfig{
			Style:     "flat",
			StripDate: true,
		},
		Template: "personal",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

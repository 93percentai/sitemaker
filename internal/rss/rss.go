package rss

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"time"

	"github.com/93percentai/sitemaker/internal/config"
	"github.com/93percentai/sitemaker/internal/content"
)

type rssRoot struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	LastBuild   string `xml:"lastBuildDate"`
	Items       []item `xml:"item"`
}

type item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
	Categories  []string `xml:"category"`
}

func Generate(site *content.Site, cfg config.Config) error {
	feed := rssRoot{
		Version: "2.0",
		Channel: channel{
			Title:       cfg.Site.Title,
			Link:        cfg.Site.BaseURL,
			Description: cfg.Site.Description,
			Language:    cfg.Site.Language,
			LastBuild:   time.Now().Format(time.RFC1123Z),
		},
	}

	limit := 20
	if len(site.Posts) < limit {
		limit = len(site.Posts)
	}

	for _, post := range site.Posts[:limit] {
		fullURL := cfg.Site.BaseURL + post.URL
		feed.Channel.Items = append(feed.Channel.Items, item{
			Title:       post.Title,
			Link:        fullURL,
			Description: string(post.HTMLContent),
			PubDate:     post.Date.Format(time.RFC1123Z),
			GUID:        fullURL,
			Categories:  post.Tags,
		})
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	output := []byte(xml.Header)
	output = append(output, data...)

	outDir := cfg.Build.OutputDir
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(outDir, "rss.xml"), output, 0644); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outDir, "feed.xml"), output, 0644)
}

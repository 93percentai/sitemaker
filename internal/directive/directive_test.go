package directive

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := []byte(`[TITLE]: # (My Page)
[DATE]: # (2025-01-15)
[TAGS]: # (go, webdev, tutorial)
[INCLUDES]: # (H, F, TOC)
[INCLUDES]: # (custom.css, analytics.js)
[INHERITS]: # (post.html)
[DRAFT]: # (true)
[SLUG]: # (my-custom-slug)

# Hello World

This is the body.
`)

	d, body := Parse(input)

	if d.Title != "My Page" {
		t.Errorf("Title = %q, want %q", d.Title, "My Page")
	}
	if d.Date != "2025-01-15" {
		t.Errorf("Date = %q, want %q", d.Date, "2025-01-15")
	}
	if len(d.Tags) != 3 || d.Tags[0] != "go" || d.Tags[1] != "webdev" || d.Tags[2] != "tutorial" {
		t.Errorf("Tags = %v, want [go webdev tutorial]", d.Tags)
	}
	if len(d.Includes) != 5 {
		t.Errorf("Includes = %v, want 5 items", d.Includes)
	}
	if d.Inherits != "post.html" {
		t.Errorf("Inherits = %q, want %q", d.Inherits, "post.html")
	}
	if !d.Draft {
		t.Error("Draft = false, want true")
	}
	if d.Slug != "my-custom-slug" {
		t.Errorf("Slug = %q, want %q", d.Slug, "my-custom-slug")
	}

	if !contains(string(body), "# Hello World") {
		t.Errorf("body should contain '# Hello World', got %q", string(body))
	}
	if contains(string(body), "[TITLE]") {
		t.Error("body should not contain directives")
	}
}

func TestParseLayoutComponents(t *testing.T) {
	d := Directives{Includes: []string{"H", "F", "TOC", "custom.css", "app.js"}}

	comps := d.LayoutComponents()
	if len(comps) != 3 {
		t.Errorf("LayoutComponents = %v, want 3 items", comps)
	}

	assets := d.AssetFiles()
	if len(assets) != 2 {
		t.Errorf("AssetFiles = %v, want 2 items", assets)
	}
}

func TestParseOverride(t *testing.T) {
	input := []byte(`[OVERRIDE]: # (landing.html)

Some content.
`)
	d, _ := Parse(input)
	if d.Override != "landing.html" {
		t.Errorf("Override = %q, want %q", d.Override, "landing.html")
	}
}

func TestParseNoDirectives(t *testing.T) {
	input := []byte(`# Just a heading

Some content.
`)
	d, body := Parse(input)
	if d.Title != "" {
		t.Errorf("Title = %q, want empty", d.Title)
	}
	if !contains(string(body), "# Just a heading") {
		t.Errorf("body should contain heading, got %q", string(body))
	}
}

func TestHasComponent(t *testing.T) {
	d := Directives{Includes: []string{"H", "F", "TOC"}}

	if !d.HasComponent("H") {
		t.Error("HasComponent(H) = false, want true")
	}
	if !d.HasComponent("h") {
		t.Error("HasComponent(h) = false, want true (case insensitive)")
	}
	if d.HasComponent("SB") {
		t.Error("HasComponent(SB) = true, want false")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

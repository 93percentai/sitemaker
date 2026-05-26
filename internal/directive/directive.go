package directive

import (
	"bytes"
	"regexp"
	"strings"
)

type Directives struct {
	Includes []string
	Inherits string
	Override string
	Date     string
	Title    string
	Tags     []string
	Draft    bool
	Slug     string
}

var directiveRe = regexp.MustCompile(`(?i)^\[(INCLUDES|INHERITS|OVERRIDE|TITLE|DATE|TAGS|DRAFT|SLUG)\]:\s*#\s*\(([^)]*)\)\s*$`)

func Parse(content []byte) (Directives, []byte) {
	var d Directives
	lines := bytes.Split(content, []byte("\n"))
	bodyStart := 0

	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}

		m := directiveRe.FindSubmatch(trimmed)
		if m == nil {
			bodyStart = i
			break
		}

		bodyStart = i + 1
		key := strings.ToUpper(string(m[1]))
		val := strings.TrimSpace(string(m[2]))

		switch key {
		case "INCLUDES":
			d.Includes = append(d.Includes, splitCSV(val)...)
		case "INHERITS":
			d.Inherits = val
		case "OVERRIDE":
			d.Override = val
		case "TITLE":
			d.Title = val
		case "DATE":
			d.Date = val
		case "TAGS":
			d.Tags = append(d.Tags, splitCSV(val)...)
		case "DRAFT":
			d.Draft = strings.EqualFold(val, "true")
		case "SLUG":
			d.Slug = val
		}
	}

	body := bytes.Join(lines[bodyStart:], []byte("\n"))
	return d, body
}

func (d Directives) LayoutComponents() []string {
	var comps []string
	for _, inc := range d.Includes {
		if !strings.Contains(inc, ".") {
			comps = append(comps, strings.ToUpper(inc))
		}
	}
	return comps
}

func (d Directives) AssetFiles() []string {
	var assets []string
	for _, inc := range d.Includes {
		if strings.Contains(inc, ".") {
			assets = append(assets, inc)
		}
	}
	return assets
}

func (d Directives) HasComponent(code string) bool {
	code = strings.ToUpper(code)
	for _, c := range d.LayoutComponents() {
		if c == code {
			return true
		}
	}
	return false
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

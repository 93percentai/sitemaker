package tmpl

import (
	"html/template"
	"strings"
	"time"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      strings.Title,
		"join":       strings.Join,
		"split":      strings.Split,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"replace":    strings.ReplaceAll,
		"formatDate": formatDate,
		"now":        func() time.Time { return time.Now() },
		"safeHTML":   func(s string) template.HTML { return template.HTML(s) },
		"seq":        seq,
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"slice":      sliceHelper,
	}
}

func formatDate(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(layout)
}

func seq(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func sliceHelper(args ...interface{}) []interface{} {
	return args
}

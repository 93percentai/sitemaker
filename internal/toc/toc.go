package toc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark/ast"
)

type Heading struct {
	Level int
	ID    string
	Text  string
}

func Extract(doc ast.Node, source []byte) []Heading {
	var headings []Heading

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		id := ""
		if idAttr, found := h.AttributeString("id"); found {
			id = string(idAttr.([]byte))
		}

		txt := extractText(h, source)

		headings = append(headings, Heading{
			Level: h.Level,
			ID:    id,
			Text:  txt,
		})

		return ast.WalkContinue, nil
	})

	return headings
}

func extractText(n ast.Node, source []byte) string {
	var buf strings.Builder
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}

func RenderTOC(headings []Heading) template.HTML {
	if len(headings) == 0 {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(`<nav class="toc">`)
	buf.WriteString("\n<ul>\n")

	prevLevel := 0
	for _, h := range headings {
		if prevLevel == 0 {
			prevLevel = h.Level
		}

		if h.Level > prevLevel {
			for i := prevLevel; i < h.Level; i++ {
				buf.WriteString("<ul>\n")
			}
		} else if h.Level < prevLevel {
			for i := h.Level; i < prevLevel; i++ {
				buf.WriteString("</li>\n</ul>\n")
			}
			buf.WriteString("</li>\n")
		} else if prevLevel > 0 {
			buf.WriteString("</li>\n")
		}

		buf.WriteString(fmt.Sprintf(`<li><a href="#%s">%s</a>`, h.ID, h.Text))
		prevLevel = h.Level
	}

	for i := headings[0].Level; i <= prevLevel; i++ {
		buf.WriteString("</li>\n</ul>\n")
	}

	buf.WriteString("</nav>\n")
	return template.HTML(buf.String())
}

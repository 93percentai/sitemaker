package markdown

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type CalloutExtension struct{}

func (e *CalloutExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&calloutTransformer{}, 199),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&calloutRenderer{}, 199),
		),
	)
}

var KindCallout = ast.NewNodeKind("Callout")

type CalloutNode struct {
	ast.BaseBlock
	CalloutType string
	Title       string
}

func (n *CalloutNode) Kind() ast.NodeKind { return KindCallout }
func (n *CalloutNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Type":  n.CalloutType,
		"Title": n.Title,
	}, nil)
}

func (n *CalloutNode) IsRaw() bool { return false }

type calloutTransformer struct{}

func (t *calloutTransformer) Transform(node *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	var toReplace []struct {
		bq   *ast.Blockquote
		ctype string
		title string
	}

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		bq, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue, nil
		}

		firstChild := bq.FirstChild()
		if firstChild == nil {
			return ast.WalkContinue, nil
		}

		para, ok := firstChild.(*ast.Paragraph)
		if !ok {
			return ast.WalkContinue, nil
		}

		lines := para.Lines()
		if lines.Len() == 0 {
			return ast.WalkContinue, nil
		}

		firstLine := lines.At(0)
		lineText := strings.TrimSpace(string(source[firstLine.Start:firstLine.Stop]))

		if !strings.HasPrefix(lineText, "[!") {
			return ast.WalkContinue, nil
		}

		end := strings.Index(lineText, "]")
		if end < 0 {
			return ast.WalkContinue, nil
		}

		calloutType := strings.ToLower(lineText[2:end])
		title := strings.TrimSpace(lineText[end+1:])
		if title == "" {
			title = strings.Title(calloutType)
		}

		toReplace = append(toReplace, struct {
			bq    *ast.Blockquote
			ctype string
			title string
		}{bq, calloutType, title})

		return ast.WalkSkipChildren, nil
	})

	for _, item := range toReplace {
		callout := &CalloutNode{
			CalloutType: item.ctype,
			Title:       item.title,
		}

		firstChild := item.bq.FirstChild()
		if para, ok := firstChild.(*ast.Paragraph); ok {
			foundBreak := false
			var contentChildren []ast.Node
			for child := para.FirstChild(); child != nil; child = child.NextSibling() {
				if !foundBreak {
					if child.Kind() == ast.KindText {
						t := child.(*ast.Text)
						if t.SoftLineBreak() {
							foundBreak = true
						}
					}
					continue
				}
				contentChildren = append(contentChildren, child)
			}

			if len(contentChildren) > 0 {
				newPara := ast.NewParagraph()
				lines := para.Lines()
				for i := 1; i < lines.Len(); i++ {
					newPara.Lines().Append(lines.At(i))
				}
				for _, c := range contentChildren {
					newPara.AppendChild(newPara, c)
				}
				callout.AppendChild(callout, newPara)
			}

			child := para.NextSibling()
			for child != nil {
				next := child.NextSibling()
				callout.AppendChild(callout, child)
				child = next
			}
		}

		item.bq.Parent().ReplaceChild(item.bq.Parent(), item.bq, callout)
	}
}

type calloutRenderer struct{}

func (r *calloutRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCallout, r.renderCallout)
}

func (r *calloutRenderer) renderCallout(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*CalloutNode)
	if entering {
		w.WriteString(`<div class="callout callout-`)
		w.WriteString(n.CalloutType)
		w.WriteString("\">\n")
		w.WriteString(`<div class="callout-title">`)
		icon := calloutIcon(n.CalloutType)
		if icon != "" {
			w.WriteString(`<span class="callout-icon">`)
			w.WriteString(icon)
			w.WriteString(`</span> `)
		}
		w.WriteString(n.Title)
		w.WriteString("</div>\n")
		w.WriteString(`<div class="callout-content">`)
		w.WriteString("\n")
	} else {
		w.WriteString("</div>\n</div>\n")
	}
	return ast.WalkContinue, nil
}

func calloutIcon(t string) string {
	switch t {
	case "note":
		return "&#9998;"
	case "tip":
		return "&#128161;"
	case "warning":
		return "&#9888;"
	case "danger":
		return "&#9762;"
	case "info":
		return "&#8505;"
	case "example":
		return "&#128203;"
	case "quote":
		return "&#10077;"
	case "bug":
		return "&#128027;"
	case "todo":
		return "&#9745;"
	case "success":
		return "&#10004;"
	case "question":
		return "&#10067;"
	case "failure":
		return "&#10008;"
	case "abstract":
		return "&#128196;"
	default:
		return ""
	}
}

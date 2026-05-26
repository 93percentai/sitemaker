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

type WikilinkExtension struct{}

func (e *WikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&wikilinkParser{}, 199),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&wikilinkRenderer{}, 199),
		),
	)
}

var KindWikilink = ast.NewNodeKind("Wikilink")

type WikilinkNode struct {
	ast.BaseInline
	Target  string
	Display string
}

func (n *WikilinkNode) Kind() ast.NodeKind { return KindWikilink }
func (n *WikilinkNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Target":  n.Target,
		"Display": n.Display,
	}, nil)
}

type wikilinkParser struct{}

func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikilinkParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	if len(line) < 4 || line[0] != '[' || line[1] != '[' {
		return nil
	}

	end := strings.Index(string(line[2:]), "]]")
	if end < 0 {
		return nil
	}

	content := string(line[2 : 2+end])
	if content == "" {
		return nil
	}

	target := content
	display := content
	if idx := strings.Index(content, "|"); idx >= 0 {
		target = content[:idx]
		display = content[idx+1:]
	}

	node := &WikilinkNode{
		Target:  strings.TrimSpace(target),
		Display: strings.TrimSpace(display),
	}

	block.Advance(seg.Start + 2 + end + 2 - seg.Start)
	return node
}

type wikilinkRenderer struct{}

func (r *wikilinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.renderWikilink)
}

func (r *wikilinkRenderer) renderWikilink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*WikilinkNode)
	slug := slugify(n.Target)

	w.WriteString(`<a href="/`)
	w.WriteString(slug)
	w.WriteString(`.html" class="wikilink">`)
	w.WriteString(n.Display)
	w.WriteString(`</a>`)

	return ast.WalkContinue, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var result []byte
	for _, c := range []byte(s) {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '/' {
			result = append(result, c)
		}
	}
	return string(result)
}

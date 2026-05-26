package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)

type Renderer struct {
	md goldmark.Markdown
}

func NewRenderer() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			&WikilinkExtension{},
			&CalloutExtension{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return &Renderer{md: md}
}

func (r *Renderer) Render(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *Renderer) Parser() parser.Parser {
	return r.md.Parser()
}

package markdown

import (
	"gostatic/pkg/markup"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	pikchr "github.com/jchenry/goldmark-pikchr"
	mathjax "github.com/litao91/goldmark-mathjax"
	fences "github.com/stefanfritsch/goldmark-fences"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	frontmatter "go.abhg.dev/goldmark/frontmatter"
)

func NewMarkdownConverter() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
			&fences.Extender{},
			&pikchr.Extender{},
			mathjax.MathJax,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithWriter(HtmlWriter{}),
		),
	)
}

func Convert(fileSource []byte, doc *markup.Document, node *markup.Node) error {
	writer := NewTreeWriter(doc, node)
	defer writer.Free()

	md := NewMarkdownConverter()
	if err := md.Convert(fileSource, &writer); err != nil {
		return err
	}

	writer.Terminate()

	return nil
}

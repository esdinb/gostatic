package markdown

import (
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"gostatic/pkg/markup"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"

	d2 "github.com/FurqanSoftware/goldmark-d2"
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
			&d2.Extender{
				// Defaults when omitted
				//TODO: where are the imports?
				//Layout:  d2dagrelayout.Layout,
				//ThemeID: d2themescatalog.CoolClassics.ID,
			},
			mathjax.MathJax,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithWriter(HtmlWriter{}),
		),
	)
}

func ConvertFile(filePath string, doc *markup.Document, node *markup.Node) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	_ = info

	fileSource, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return ConvertMemory(fileSource, doc, node)
}

func ConvertMemory(fileSource []byte, doc *markup.Document, node *markup.Node) error {
	writer := NewTreeWriter(doc, node)
	defer writer.Free()

	writer.Write([]byte(fmt.Sprintf(`<div data-character-count="%s">`, strconv.Itoa(utf8.RuneCount(fileSource)))))

	md := NewMarkdownConverter()
	if err := md.Convert(fileSource, &writer); err != nil {
		return err
	}

	writer.Write([]byte("</div>"))
	writer.Terminate()
	return nil
}

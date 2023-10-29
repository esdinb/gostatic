package markdown

import (
	"gostatic/pkg/markup"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

type Converter struct {
	goldmark goldmark.Markdown
}

func rendererOptions() goldmark.Option {
	return goldmark.WithRendererOptions(
		html.WithWriter(HtmlWriter{}),
	)
}

func New(options ...goldmark.Option) *Converter {
	return &Converter{goldmark.New(append(options, rendererOptions())...)}
}

func (c *Converter) Convert(fileSource []byte, doc *markup.Document, node *markup.Node) error {
	writer := NewTreeWriter(doc, node)
	defer writer.Free()

	if err := c.goldmark.Convert(fileSource, &writer); err != nil {
		return err
	}

	writer.Terminate()

	return nil
}

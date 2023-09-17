package markdown

import (
    "os"
    "strconv"
    "unicode/utf8"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/renderer/html"

    frontmatter "go.abhg.dev/goldmark/frontmatter"
    pikchr "github.com/jchenry/goldmark-pikchr"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
    mathjax "github.com/litao91/goldmark-mathjax"
    d2 "github.com/FurqanSoftware/goldmark-d2"
    fences "github.com/stefanfritsch/goldmark-fences"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"

    "github.com/jbussdieker/golibxml"
)

const markdownHtmlParserOptions =
    golibxml.XML_PARSE_RECOVER &
    golibxml.XML_PARSE_NOENT &
    golibxml.XML_PARSE_PEDANTIC &
    golibxml.XML_PARSE_NONET

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

func ConvertFile(filePath string) (*golibxml.Document, error) {
    info, err := os.Stat(filePath)
    if err != nil {
        return nil, err
    }

    _ = info

    fileSource, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    return ConvertMemory(fileSource, filePath)
}

func ConvertMemory(fileSource []byte, filePath string) (*golibxml.Document, error) {
    writer := NewTreeWriter("<div>", filePath)
    defer writer.Free()

    md := NewMarkdownConverter()
    if err := md.Convert(fileSource, &writer); err != nil {
        return nil, err
    }

    writer.Terminate("</div>")
    doc := writer.Document()
    doc.Root().SetAttribute("data-character-count", strconv.Itoa(utf8.RuneCount(fileSource)))
    return doc, nil
}


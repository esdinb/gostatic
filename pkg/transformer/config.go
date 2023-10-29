package transformer

import (
	"gostatic/pkg/markdown"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/evanw/esbuild/pkg/api"
	pikchr "github.com/jchenry/goldmark-pikchr"
	mathjax "github.com/litao91/goldmark-mathjax"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/parser"
	frontmatter "go.abhg.dev/goldmark/frontmatter"
)

func configBundleOptions(rootPath string, loaders map[string]api.Loader) api.BuildOptions {
	return api.BuildOptions{
		Color:             api.ColorIfTerminal,
		LogLevel:          api.LogLevelDebug,
		Sourcemap:         api.SourceMapNone,
		Target:            api.ESNext,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		LineLimit:         80,
		Charset:           api.CharsetUTF8,
		TreeShaking:       api.TreeShakingTrue,
		LegalComments:     api.LegalCommentsInline,
		Bundle:            true,
		AbsWorkingDir:     rootPath,
		Platform:          api.PlatformBrowser,
		Format:            api.FormatESModule,
		Loader:            loaders,
	}
}

func configMarkdownConverter() *markdown.Converter {
	return markdown.New(
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
	)
}

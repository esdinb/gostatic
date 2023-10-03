package transformer

import (
	"errors"
	"gostatic/pkg/markup"
	"path/filepath"
	"strings"
)

func customLoader(context *Context) markup.DocLoaderFunc {
	return func(
		uri string,
		dict *markup.Dict,
		options markup.ParserOption,
		ctx *markup.DocLoaderContext,
		loadType markup.LoadType,
	) *markup.Document {
		templatePath := uri
		if !strings.HasPrefix(uri, context.RootPath) {
			templatePath = filepath.Join(context.RootPath, uri)
		}
		return markup.DefaultLoader(templatePath, dict, options, ctx, loadType)
	}
}

func TransformTransform(context *Context, args []string) (*Context, Status, error) {

	if len(args) < 1 {
		return context, Continue, errors.New("missing argument for transform")
	}

	markup.SetLoaderFunc(customLoader(context))

	var filename string
	var style *markup.Stylesheet

	filename = args[0]
	if filename == "inline" {
		style = markup.LoadStylesheetPI(context.Document)
		if style == nil {
			return context, Continue, errors.New("missing inline stylesheet")
		}
	} else {
		style = markup.ParseStylesheetFile(filename)
	}

	defer style.Free()

	params := []string{}

	result := markup.ApplyStylesheet(style, context.Document, params)
	if result == nil {
		return context, Continue, errors.New("error applying stylesheet")
	} else {
		context.Document.Free()
		context.Document = result
	}

	return context, Continue, nil
}

func init() {
	Registry.Register("transform", TransformTransform)
}

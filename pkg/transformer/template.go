package transformer

import (
	"context"
	"errors"
	"gostatic/pkg/markup"
	"path/filepath"
	"strings"
)

func customLoader(ctx context.Context) markup.DocLoaderFunc {
	return func(
		uri string,
		dict *markup.Dict,
		options markup.ParserOption,
		loaderCtx *markup.DocLoaderContext,
		loadType markup.LoadType,
	) *markup.Document {
		templatePath := uri
		rootPath := ctx.Value(RootPathContextKey).(string)
		if !strings.HasPrefix(uri, rootPath) {
			templatePath = filepath.Join(rootPath, uri)
		}
		return markup.DefaultLoader(templatePath, dict, options, loaderCtx, loadType)
	}
}

func TemplateTransform(ctx context.Context, args []string) (context.Context, Status, error) {

	if len(args) < 1 {
		return ctx, Continue, errors.New("missing argument for transform")
	}

	markup.SetLoaderFunc(customLoader(ctx))

	var filename string
	var style *markup.Stylesheet
	var document *markup.Document

	document = ctx.Value(DocumentContextKey).(*markup.Document)
	filename = args[0]
	if filename == "inline" {
		style = markup.LoadStylesheetPI(document)
		if style == nil {
			return ctx, Continue, errors.New("missing inline stylesheet")
		}
	} else {
		style = markup.ParseStylesheetFile(filename)
	}

	defer style.Free()

	params := ctx.Value(ParamsContextKey).([]string)
	strparams := ctx.Value(StringParamsContextKey).([]string)

	transformation := markup.ApplyStylesheet(style, document, params, strparams)
	if transformation == nil {
		return ctx, Continue, errors.New("error applying stylesheet")
	} else {
		document.Free()
		ctx = context.WithValue(ctx, DocumentContextKey, transformation)
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("template", TemplateTransform)
}

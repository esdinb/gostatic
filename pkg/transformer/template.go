package transformer

import (
	"context"
	"errors"
	"gostatic/pkg/config"
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
		rootPath := ctx.Value(config.RootPathContextKey).(string)
		if !strings.HasPrefix(uri, rootPath) {
			templatePath = filepath.Join(rootPath, uri)
		}
		return markup.DefaultLoader(templatePath, dict, options, loaderCtx, loadType)
	}
}

func TransformTemplate(ctx context.Context, args []string) (context.Context, Status, error) {

	if len(args) < 1 {
		return ctx, Continue, errors.New("missing argument for transform")
	}

	markup.SetLoaderFunc(customLoader(ctx))

	var filename string
	var style *markup.Stylesheet
	var document *markup.Document

	document = ctx.Value(config.DocumentContextKey).(*markup.Document)
	filename = args[0]
	if filename == "inline" {
		style = markup.LoadStylesheetPI(document)
		if style == nil {
			return ctx, Continue, errors.New("missing inline stylesheet")
		}
	} else {
		style = markup.ParseStylesheetFile(filename)
	}
	if style == nil {
		return ctx, Continue, errors.New("unable to parse stylesheet")
	}

	defer style.Free()

	params, ok := ctx.Value(config.ParamsContextKey).([]string)
	if !ok {
		return ctx, Continue, errors.New("missing params array")
	}
	strparams, ok := ctx.Value(config.StringParamsContextKey).([]string)
	if !ok {
		return ctx, Continue, errors.New("missing strparams array")
	}

	transformation := markup.ApplyStylesheetUser(style, document, params, strparams)
	if transformation == nil {
		return ctx, Continue, errors.New("error applying stylesheet")
	} else {
		document.Free()
		ctx = context.WithValue(ctx, config.DocumentContextKey, transformation)
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("template", TransformTemplate)
}

package transformer

import (
	"context"
	"errors"
	builder_context "gostatic/pkg/builder/context"
	"gostatic/pkg/markup"
	"log"
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
		rootPath := ctx.Value(builder_context.RootPathContextKey).(string)
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

	var logger *log.Logger
	var filename string
	var style *markup.Stylesheet
	var document *markup.Document

	logger = ctx.Value(builder_context.LoggerContextKey).(*log.Logger)
	document, ok := ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if !ok {
		return ctx, Continue, errors.New("missing input document to template transform")
	}
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

	params, ok := ctx.Value(builder_context.ParamsContextKey).([]string)
	if !ok {
		return ctx, Continue, errors.New("missing params array")
	}
	strparams, ok := ctx.Value(builder_context.StringParamsContextKey).([]string)
	if !ok {
		return ctx, Continue, errors.New("missing strparams array")
	}

	transformCtx := markup.NewTransformContext(style, document, logger)
	defer transformCtx.Free()

	transformation := transformCtx.ApplyStylesheet(style, document, params, strparams)
	if transformation == nil {
		return ctx, Continue, errors.New("error applying stylesheet")
	} else {
		document.Free()
		ctx = context.WithValue(ctx, builder_context.DocumentContextKey, transformation)
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("template", TransformTemplate)
}

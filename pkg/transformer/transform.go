package transformer

import (
    "errors"
    "path/filepath"

    "github.com/jbussdieker/golibxml"
)

func customLoader(context *Context) golibxml.DocLoaderFunc {
    return func(
        uri string,
        dict *golibxml.Dict,
        options golibxml.ParserOption,
        ctx *golibxml.DocLoaderContext,
        loadType golibxml.LoadType,
    ) *golibxml.Document {
        templatePath := filepath.Join(context.RootPath, uri)
        return golibxml.DefaultLoader(templatePath, dict, options, ctx, loadType)
    }
}

func TransformTransform(context *Context, args []string) (*Context, Status, error) {

    if len(args) < 1 {
        return context, Continue, errors.New("missing argument for transform")
    }

    golibxml.SetLoaderFunc(customLoader(context))

    var filename string
    var style *golibxml.Stylesheet

    filename = args[0]
    if filename == "inline" {
        if style = golibxml.LoadStylesheetPI(context.Document); style == nil {
            return context, Continue, errors.New("missing inline stylesheet")
        }
    } else {
        style = golibxml.ParseStylesheetFile(filename)
    }

    defer style.Free()

    params := []string{}

    if result := golibxml.ApplyStylesheet(style, context.Document, params); result == nil {
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


package transformer

import (
	"context"
	builder_context "gostatic/pkg/builder/context"
	"gostatic/pkg/markup"
)

func normalizeWhitespace(doc *markup.Document) {
	xpath := markup.NewXPathContext(doc)
	defer xpath.Free()

	// normalize whitespace
	whiteSpaceElements := xpath.Eval("////text()[not(normalize-space())]")
	if whiteSpaceElements != nil {
		defer whiteSpaceElements.Free()
		for _, node := range whiteSpaceElements.Results() {
			node.Unlink()
		}
	}
}

func TransformWhitespace(ctx context.Context, args []string) (context.Context, Status, error) {
	var subcommand string
	if len(args) > 0 {
		subcommand = args[0]
	} else {
		subcommand = "normalize"
	}
	document := ctx.Value(builder_context.DocumentContextKey).(*markup.Document)

	switch subcommand {
	case "normalize":
		normalizeWhitespace(document)
	default:
		break
	}
	return ctx, Continue, nil
}

func init() {
	Registry.Register("whitespace", TransformWhitespace)
}

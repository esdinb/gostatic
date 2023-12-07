package transformer

import (
	"context"
	"errors"
	builder_context "gostatic/pkg/builder/context"
	"gostatic/pkg/markup"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

func TransformBanner(ctx context.Context, args []string) (context.Context, Status, error) {
	text := "gostatic"
	font := "isometric1"
	switch l := len(args); l {
	case 0:
	case 1:
		text = args[0]
	case 2:
		font, text = args[0], args[1]
	default:
		font = args[0]
		text = strings.Join(args[1:], "")
	}
	banner := figure.NewFigure(text, font, true).String()
	document, ok := ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if !ok {
		return ctx, Continue, errors.New("missing input document")
	}
	comment := document.NewComment("\n" + banner + "\n\n")
	var node *markup.Node
	if node = document.FirstChild(); node != nil {
		node.AddPrevSibling(comment)
	} else if node = document.Root(); node != nil {
		node.AddChild(comment)
	} else {
		document.SetRoot(comment)
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("banner", TransformBanner)
}

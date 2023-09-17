package transformer

import (
    "strings"

    "github.com/common-nighthawk/go-figure"
)

func TransformBanner(context *Context, args []string) (*Context, Status, error) {
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
    comment := context.Document.NewComment("\n" + banner + "\n\n")
    context.Document.FirstChild().AddPrevSibling(*comment)
    return context, Continue, nil
}

func init() {
    Registry.Register("banner", TransformBanner)
}

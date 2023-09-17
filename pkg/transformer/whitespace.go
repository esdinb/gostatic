package transformer

import (
    "github.com/jbussdieker/golibxml"
)

func normalizeWhitespace(doc *golibxml.Document) {
    xpath := golibxml.NewXPathContext(doc)
    defer xpath.Free()

    // normalize whitespace
    whiteSpaceElements := xpath.Eval("////text()[not(normalize-space())]")
    if whiteSpaceElements != nil {
        defer whiteSpaceElements.Free()
        for node := range whiteSpaceElements.Results() {
            node.Unlink()
        }
    }
}

func TransformWhitespace(context *Context, args []string) (*Context, Status, error) {
    var subcommand string
    if len(args) > 0 {
        subcommand = args[0]
    } else {
        subcommand = "normalize"
    }
    switch subcommand {
    case "normalize":
        normalizeWhitespace(context.Document)
    default:
        break
    }
    return context, Continue, nil
}

func init() {
    Registry.Register("whitespace", TransformWhitespace)
}


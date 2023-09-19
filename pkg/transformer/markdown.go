package transformer

import (
    "os"
    "fmt"
    "path/filepath"
    "strings"
    "gostatic/pkg/markdown"
    "github.com/jbussdieker/golibxml"
)

func TransformMarkdown(context *Context, args []string) (*Context, Status, error) {

    xpath := golibxml.NewXPathContext(context.Document)
    defer xpath.Free()

    markdownFiles := xpath.Eval("//script[@type='text/markdown']")
    if markdownFiles == nil {
        return context, Continue, nil
    }

    defer markdownFiles.Free()

    for node := range markdownFiles.Results() {
        var mdDoc *golibxml.Document
        var path string
        var absPath string
        var err error

        content := strings.TrimSpace(node.GetContent())

        absPath = ""
        for attr := node.Attributes(); attr != nil; attr = attr.Next() {
            if (attr.Name() == "src") {
                path = attr.Children().String()
                if strings.HasPrefix(path, "/") {
                    path = filepath.Join(context.RootPath, path)
                } else {
                    fileInfo, err := os.Stat(context.InPath)
                    if err != nil {
                        return context, Continue, err
                    }
                    if fileInfo.IsDir() {
                        path = filepath.Join(context.InPath, path)
                    } else {
                        path = filepath.Join(filepath.Dir(context.InPath), path)
                    }
                }
                absPath, err = filepath.Abs(path)
                if err != nil {
                    fmt.Println("error getting abs path of source file", path)
                    return context, Continue, err
                }
                break
            }
        }

        if len(absPath) != 0 {
            mdDoc, err = markdown.ConvertFile(absPath)
        } else {
            mdDoc, err = markdown.ConvertMemory([]byte(content), "/")
        }
        if err != nil {
            return context, Continue, err
        }
        node.Replace(mdDoc.Root())
    }

    return context, Continue, nil
}

func init() {
    Registry.Register("markdown", TransformMarkdown)
}


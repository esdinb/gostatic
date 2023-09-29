package transformer

import (
	"fmt"
	"gostatic/pkg/markdown"
	"gostatic/pkg/markup"
	"os"
	"path/filepath"
	"strings"
)

func TransformMarkdown(context *Context, args []string) (*Context, Status, error) {

	xpath := markup.NewXPathContext(context.Document)
	defer xpath.Free()

	markdownFiles := xpath.Eval("//script[@type='text/markdown']")
	if markdownFiles == nil {
		return context, Continue, nil
	}

	defer markdownFiles.Free()

	for node := range markdownFiles.Results() {
		var (
			replacement *markup.Node
			path        string
			absPath     string
			err         error
		)

		content := strings.TrimSpace(node.GetContent())

		absPath = ""
		for attr := node.Attributes(); attr != nil; attr = attr.Next() {
			if attr.Name() == "src" {
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

		replacement = context.Document.NewFragment()
		node.AddNextSibling(*replacement)
		node.Unlink()
		if len(absPath) != 0 {
			err = markdown.ConvertFile(absPath, context.Document, replacement)
		} else {
			err = markdown.ConvertMemory([]byte(content), context.Document, replacement)
		}
		if err != nil {
			return context, Continue, err
		}
	}

	return context, Continue, nil
}

func init() {
	Registry.Register("markdown", TransformMarkdown)
}

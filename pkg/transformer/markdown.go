package transformer

import (
	"context"
	"fmt"
	"gostatic/pkg/markdown"
	"gostatic/pkg/markup"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func TransformMarkdown(ctx context.Context, args []string) (context.Context, Status, error) {

	document := ctx.Value(DocumentContextKey).(*markup.Document)
	xpath := markup.NewXPathContext(document)
	defer xpath.Free()

	markdownFiles := xpath.Eval("//script[@type='text/markdown']")
	if markdownFiles == nil {
		return ctx, Continue, nil
	}

	defer markdownFiles.Free()

	for node := range markdownFiles.Results() {
		var (
			replacement *markup.Node
			path        string
			absPath     string
			sourcePath  string
			err         error
		)

		content := strings.TrimSpace(node.GetContent())

		rootPath := ctx.Value(RootPathContextKey).(string)
		inPath := ctx.Value(InPathContextKey).(string)
		absPath = ""
		for attr := node.Attributes(); attr != nil; attr = attr.Next() {
			if attr.Name() == "src" {
				path = attr.Children().String()
				if strings.HasPrefix(path, "/") {
					path = filepath.Join(rootPath, path)
				} else {
					fileInfo, err := os.Stat(inPath)
					if err != nil {
						return ctx, Continue, err
					}
					if fileInfo.IsDir() {
						path = filepath.Join(inPath, path)
					} else {
						path = filepath.Join(filepath.Dir(inPath), path)
					}
				}
				absPath, err = filepath.Abs(path)
				if err != nil {
					fmt.Println("error getting abs path of source file", path)
					return ctx, Continue, err
				}
				break
			}
		}

		replacement = document.NewFragment()

		node.AddNextSibling(*replacement)
		node.Unlink()

		if len(absPath) != 0 {
			sourcePath = absPath
			err = markdown.ConvertFile(absPath, document, replacement)
		} else {
			sourcePath = inPath
			err = markdown.ConvertMemory([]byte(content), document, replacement)
		}
		if err != nil {
			return ctx, Continue, err
		}

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return ctx, Continue, err
		}
		strparams := ctx.Value(StringParamsContextKey).([]string)
		strparams = append(strparams, "modtime", fileInfo.ModTime().Format(time.DateOnly))
		ctx = context.WithValue(ctx, StringParamsContextKey, strparams)
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("markdown", TransformMarkdown)
}

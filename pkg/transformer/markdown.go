package transformer

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gostatic/pkg/markdown"
	"gostatic/pkg/markup"
)

func TransformMarkdown(ctx context.Context, args []string) (context.Context, Status, error) {

	document := ctx.Value(DocumentContextKey).(*markup.Document)
	xpath := markup.NewXPathContext(document)
	defer xpath.Free()

	markdownFiles := xpath.Eval("//*[@is='markdown-element']")
	if markdownFiles == nil {
		return ctx, Continue, nil
	}

	defer markdownFiles.Free()

	for node := range markdownFiles.Results() {
		var (
			path       string
			absPath    string
			sourcePath string
			content    string
			bytes      []byte
			err        error
		)

		content = strings.TrimSpace(node.GetContent())

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
					return ctx, Continue, err
				}
				break
			}
		}

		if len(absPath) != 0 {
			bytes, err = os.ReadFile(absPath)
			if err != nil {
				return ctx, Continue, err
			}
			sourcePath = absPath
		} else {
			bytes = []byte(content)
			sourcePath = inPath
		}
		node.SetContent("")
		err = markdown.Convert(bytes, document, node)
		if err != nil {
			return ctx, Continue, err
		}

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return ctx, Continue, err
		}

		node.SetAttribute("data-last-modified", fileInfo.ModTime().Format(time.RFC3339))
		node.SetAttribute("data-character-count", strconv.Itoa(utf8.RuneCount(bytes)))
		if attr := node.HasAttribute("is"); attr != nil {
			markup.RemoveAttribute(attr)
		}
		if attr := node.HasAttribute("src"); attr != nil {
			markup.RemoveAttribute(attr)
		}
	}

	return ctx, Continue, nil
}

func init() {
	Registry.Register("markdown", TransformMarkdown)
}

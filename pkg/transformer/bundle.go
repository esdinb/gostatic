package transformer

import (
	"context"
	"errors"
	"fmt"
	builder_context "gostatic/pkg/builder/context"
	"gostatic/pkg/markup"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

var (
	cssLoaders map[string]api.Loader
	jsLoaders  map[string]api.Loader
)

func buildBundle(buildOptions api.BuildOptions) string {
	var builder strings.Builder

	result := api.Build(buildOptions)
	for i := range result.OutputFiles {
		builder.Write(result.OutputFiles[i].Contents)
	}

	return strings.TrimSpace(builder.String())
}

func bundleInline(content string, rootPath string, stdinOptions api.StdinOptions, loaders map[string]api.Loader) string {
	buildOptions := configBundleOptions(rootPath, loaders)
	buildOptions.Stdin = &stdinOptions

	return buildBundle(buildOptions)
}

func bundle(filePaths []string, rootPath string, loaders map[string]api.Loader) string {
	buildOptions := configBundleOptions(rootPath, loaders)
	buildOptions.EntryPoints = filePaths

	return buildBundle(buildOptions)
}

func TransformBundle(ctx context.Context, args []string) (context.Context, Status, error) {

	document := ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if document == nil {
		return ctx, Continue, errors.New("missing input document")
	}

	rootPath, ok := ctx.Value(builder_context.RootPathContextKey).(string)
	if !ok {
		return ctx, Continue, errors.New("missing root path")
	}
	inPath, ok := ctx.Value(builder_context.InPathContextKey).(string)
	if !ok {
		return ctx, Continue, errors.New("missing input path")
	}

	xpath := markup.NewXPathContext(document)
	defer xpath.Free()

	styleElements := xpath.Eval("//style")
	for _, node := range styleElements.Results() {
		stdinOptions := api.StdinOptions{
			Contents:   node.GetContent(),
			ResolveDir: rootPath,
			Sourcefile: inPath,
			Loader:     api.LoaderCSS,
		}
		result := bundleInline(node.GetContent(), rootPath, stdinOptions, cssLoaders)
		node.SetContent(result)
	}

	var linkElements *markup.XPathObject

	linkElements = xpath.Eval("/html/head/link[@rel='stylesheet' and @href and string-length(@href) != 0]")
	linkPaths := []string{}
	var linkNode *markup.Node
	for _, node := range linkElements.Results() {
		linkPath := node.GetAttribute("href")
		linkPaths = append(linkPaths, linkPath)
		if linkNode != nil {
			node.Unlink()
		} else {
			linkNode = node
		}
		preloadLinkElements := xpath.Eval(fmt.Sprintf("/html/head/link[@rel='preload' and @href='%s']", linkPath))
		for _, preloadNode := range preloadLinkElements.Results() {
			preloadNode.Unlink()
		}
		if linkNode != nil && len(linkPaths) > 0 {
			result := bundle(linkPaths, rootPath, cssLoaders)
			newNode := document.NewNode(nil, "style", result)
			linkNode.Replace(newNode)
		}
	}

	linkElements = xpath.Eval("/html/body/link[@rel='stylesheet' and @href and string-length(@href) != 0]")
	for _, node := range linkElements.Results() {
		linkPath := node.GetAttribute("href")
		result := bundle([]string{linkPath}, rootPath, cssLoaders)
		newNode := document.NewNode(nil, "style", result)
		node.Replace(newNode)
	}

	var scriptElements *markup.XPathObject

	scriptElements = xpath.Eval("/html/head/script[@src and string-length(@src) != 0 and @type='module']")
	scriptPaths := []string{}
	var scriptNode *markup.Node
	for _, node := range scriptElements.Results() {
		srcAttr := node.HasAttribute("src")
		scriptPath := srcAttr.Children().String()
		scriptPaths = append(scriptPaths, scriptPath)
		if scriptNode != nil {
			node.Unlink()
		} else {
			scriptNode = node
		}
		preloadScriptElements := xpath.Eval(fmt.Sprintf("/html/head/link[@rel='preload' and @href='%s']", scriptPath))
		for _, preloadNode := range preloadScriptElements.Results() {
			preloadNode.Unlink()
		}
		if scriptNode != nil && len(scriptPaths) > 0 {
			result := bundle(scriptPaths, rootPath, jsLoaders)
			scriptNode.SetContent(result)
			markup.RemoveAttribute(srcAttr)
		}
	}

	scriptElements = xpath.Eval("/html/body/script[@src and string-length(@src) != 0 and @type='module']")
	for _, node := range scriptElements.Results() {
		srcAttr := node.HasAttribute("src")
		scriptPath := srcAttr.Children().String()
		result := bundle([]string{scriptPath}, rootPath, jsLoaders)
		node.SetContent(result)
		markup.RemoveAttribute(srcAttr)
	}

	scriptElements = xpath.Eval("//script[@src and string-length(@src) != 0 and not(@type)]")
	for _, node := range scriptElements.Results() {
		srcAttr := node.HasAttribute("src")
		scriptPath := srcAttr.Children().String()
		if strings.HasPrefix(scriptPath, "/") {
			scriptPath = filepath.Join(rootPath, scriptPath)
		} else {
			scriptPath = filepath.Join(filepath.Dir(inPath), scriptPath)
		}
		contents, err := os.ReadFile(scriptPath)
		if err != nil {
			return ctx, Continue, errors.New("cannot read script path")
		}
		node.SetContent(strings.TrimSpace(string(contents)))
		markup.RemoveAttribute(srcAttr)
	}

	return ctx, Continue, nil
}

func init() {
	cssLoaders = make(map[string]api.Loader)
	cssLoaders[""] = api.LoaderCSS

	jsLoaders := make(map[string]api.Loader)
	jsLoaders[""] = api.LoaderJS
	jsLoaders[".ts"] = api.LoaderTS
	jsLoaders[".d.ts"] = api.LoaderTS

	Registry.Register("bundle", TransformBundle)
}

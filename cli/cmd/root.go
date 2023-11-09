package cmd

import (
	_ "embed"
	builder_context "gostatic/pkg/builder/context"
	"os"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var VersionString string

var rootCmd = &cobra.Command{
	Use:   "gostatic",
	Short: "Make static 🥚s with 🥚 tech.",
	Long: `This tool applies XSLT transformations to HTML and XML input files.

Supported character encodings: utf-8.

Supported input document formats: .html, .xml.

Using stdin or stdout assumes XML document encoding.

One input file for each output file.

The program reads and writes relative to the document root (the location of a
build.yaml file or current working directory).
`,
	Version: VersionString,
}

func Execute() {
	ctx := builder_context.NewBuildContext()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "🤓 v%s" .Version}}
`)
}

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
	Long: `Make static websites with XSLT.

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

package cmd

import (
	"context"
	_ "embed"
	"gostatic/pkg/config"
	"log"
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
	logger := log.New(os.Stderr, "🐙 ", 0)
	ctx := context.WithValue(context.Background(), config.LoggerContextKey, logger)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "🤓 v%s" .Version}}
`)
}

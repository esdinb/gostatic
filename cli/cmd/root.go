package cmd

import (
    "os"
    "log"
    "context"

    "github.com/spf13/cobra"
)

type contextKey struct {
    name string
}

func (k *contextKey) String() string { return "gostatic cli context value " + k.name }

var LoggerContextKey = contextKey{"logger"}
var RootPathContextKey = contextKey{"rootpath"}
var BuildPathContextKey = contextKey{"buildpath"}
var ServerRootContextKey = contextKey{"serverroot"}
var ServerAddressContextKey = contextKey{"serveraddress"}
var ServerPortContextKey = contextKey{"serverport"}

var rootCmd = &cobra.Command{
    Use:   "gostatic",
    Short: "Make static 🥚s with 🥚 tech.",
    Long: `Make static websites with XSLT.

`,
    Version: "0.1",
}

func Execute() {
    logger := log.New(os.Stdout, "🐙 ", 0)
    ctx := context.WithValue(context.Background(), LoggerContextKey, logger)
    err := rootCmd.ExecuteContext(ctx)
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "🤓 v. %s" .Version}}
`)
}

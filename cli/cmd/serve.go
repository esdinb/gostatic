package cmd

import (
    "os"
	"fmt"
    "log"
    "errors"
    "io/fs"
    "net/http"

	"github.com/spf13/cobra"
)

var (
    address string
    port int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a http dev server",
	Long: `Start development http server.

    Listens on :8080 by default.

    Serves files from ./build or the directory named on the command line.
    `,
    Example: "server --address 0.0.0.0 --port 8080",
    Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        var path string
        var info fs.FileInfo
        var err error
        if len(args) == 1 {
            path = args[0]
        } else {
            path = "./build"
        }

        logger := log.New(os.Stdout, "🏄 ", 0)

        info, err = os.Lstat(path)
        if err != nil {
            logger.Panicf("%v 🤮\n", err)
        }

        if !info.IsDir() {
            logger.Println("root path is not a directory 🤥")
            os.Exit(1)
        }

        fs := http.FileServer(http.Dir(path))
        http.Handle("/", fs)

        logger.Printf("running http server on %s:%d\n", address, port)

        err = http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)

        if errors.Is(err, http.ErrServerClosed) {
            logger.Printf("closing http server\n")
        } else if err != nil {
            logger.Printf("error starting server: %v\n", err)
        }
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

    serveCmd.Flags().StringVarP(&address, "address", "a", "", "Server listen address")
    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Server listen port")
}

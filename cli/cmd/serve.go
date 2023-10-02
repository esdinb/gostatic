package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	serverRoot           string = "./build"
	serverAddress        string
	serverPort           int
	defaultServerAddress string = ""
	defaultServerPort    int    = 8080
)

func addServeFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&serverAddress, "address", "a", defaultServerAddress, "Server listen address")
	cmd.Flags().IntVarP(&serverPort, "port", "p", defaultServerPort, "Server listen port")
}

func loggingHandler(logger *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		logger.Println(r.Method, r.URL.Path, r.RemoteAddr)
	})
}

func runServer(ctx context.Context, wg *sync.WaitGroup, address string, port int, rootPath string) {

	logger := ctx.Value(LoggerContextKey).(*log.Logger)

	fsHandler := http.FileServer(http.Dir(rootPath))

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Handler:  loggingHandler(log.New(os.Stdout, "🌊  ", 0), fsHandler),
		ErrorLog: logger,
	}

	_ = context.AfterFunc(ctx, func() {
		ctx, _ := context.WithTimeout(ctx, 30*time.Second)
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Println(err)
		}
		wg.Done()
	})

	logger.Printf("starting server on %s:%d\n", address, port)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Printf("closing server\n")
	} else if err != nil {
		logger.Printf("error starting server: %v\n", err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a http dev server",
	Long: fmt.Sprintf(`Start development http server.

    Listens on %s:%d by default.

    Serves files from %s or the directory named on the command line.
    `, defaultServerAddress, defaultServerPort, serverRoot),
	Example: "serve --address 0.0.0.0 --port 8080",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		logger := ctx.Value(LoggerContextKey).(*log.Logger)
		logger.SetPrefix("🏄 ")

		serverRoot := args[0]

		info, err := os.Stat(serverRoot)
		if err != nil {
			logger.Fatalf("%v 🤮\n", err)
		}

		if !info.IsDir() {
			logger.Println("root path is not a directory 🤥")
			os.Exit(1)
		}

		wg := new(sync.WaitGroup)
		serverCtx, _ := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		wg.Add(1)
		go runServer(serverCtx, wg, serverAddress, serverPort, serverRoot)

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	addServeFlags(serveCmd)
}

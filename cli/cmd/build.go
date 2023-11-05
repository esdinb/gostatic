package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"gostatic/pkg/builder"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

const (
	configName string = "./build.yaml"
)

var (
	watchFiles bool
)

func getConfigurationPath(argPath string, configName string) (string, error) {
	if argPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return filepath.Join(cwd, configName), nil
	} else {
		fileInfo, err := os.Stat(argPath)
		if err != nil {
			return "", err
		}

		if fileInfo.IsDir() {
			argPath = filepath.Join(argPath, configName)
			fileInfo, err = os.Stat(argPath)
			if err != nil {
				return "", err
			}
		}

		return argPath, nil
	}
}

func readBuildConfiguration(path string) ([]builder.BuildSection, error) {
	var config []builder.BuildSection
	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func runBuild(ctx context.Context, wg *sync.WaitGroup, buildPath string) {

	logger := ctx.Value(LoggerContextKey).(*log.Logger)

	config, err := readBuildConfiguration(buildPath)
	if err != nil {
		logger.Println(err)
	}

	rootPath := filepath.Dir(buildPath)

	logger.Println("rebuilding...")
	for i := range config {
		section := config[i]
		err := section.Build(ctx, rootPath)
		if err != nil {
			logger.Println("build error:", err)
		}
	}

	wg.Done()
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Example: "gostatic build -s build/ .",
	Short:   "Build a site from configuration file",
	Long: `The build command reads build.yaml configuration file and applies a series of transformations
(a pipeline) to each input file to produce output files.

A transformation is a name and some arguments separated by ':'.

Transformations: template, bundle.

'template': Applies a XSL stylesheet to input. First argument is a path to a stylesheet.

'bundle': Bundles script and css using esbuild into a single HTML (or XML) file.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		ctx := cmd.Context()

		logger := ctx.Value(LoggerContextKey).(*log.Logger)
		logger.SetPrefix("🧱  ")

		argPath := ""
		if len(args) > 0 {
			argPath = args[0]
		}

		buildPath, err := getConfigurationPath(argPath, configName)
		if err != nil {
			logger.Fatal(err)
		}

		wg := new(sync.WaitGroup)

		serveFiles := cmd.Flags().Lookup("serve").Changed

		if watchFiles || serveFiles {
			watchCtx, _ := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

			runner := func(ctx context.Context) {
				buildCtx, _ := context.WithCancel(ctx)

				wg.Add(1)
				runBuild(buildCtx, wg, buildPath)
			}
			runner(watchCtx)

			inputPaths, err := collectInputPaths(buildPath)
			if err != nil {
				logger.Fatal(err)
			}

			wg.Add(1)
			go runWatcher(watchCtx, wg, inputPaths, []string{}, runner)

			if serveFiles {
				serverCtx, _ := context.WithCancel(watchCtx)

				wg.Add(1)
				go runServer(serverCtx, wg, serverAddress, serverPort, serverRoot)
			}
		} else {
			wg.Add(1)
			go runBuild(ctx, wg, buildPath)
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&watchFiles, "watch", "w", false, "Watch files")
	buildCmd.Flags().StringVarP(&serverRoot, "serve", "s", "./build", "Server root directory")
	addServeFlags(buildCmd)
}

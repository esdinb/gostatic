package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func collectOutputPaths(buildPath string) ([]string, error) {
	config, err := readBuildConfiguration(buildPath)
	if err != nil {
		return []string{}, err
	}

	rootPath := filepath.Dir(buildPath)
	paths := make([]string, len(config))
	for i := range config {
		path := filepath.Join(rootPath, config[i].Out)
		paths[i] = path
	}

	sort.Strings(paths)
	shortestPaths := []string{}
	lastPath := ""
	for i := range paths {
		nextPath := paths[i]
		if lastPath == "" || !strings.HasPrefix(nextPath, lastPath) {
			lastPath = nextPath
			shortestPaths = append(shortestPaths, lastPath)
		}
	}

	return shortestPaths, nil
}

func excludePathPrefix(prefixes []string, path string) bool {
	for i := range prefixes {
		if strings.HasPrefix(path, prefixes[i]) {
			return true
		}
	}

	return false
}

func includePathMatch(patterns []string, path string) bool {
	if len(patterns) == 0 {
		return true
	}

	for i := range patterns {
		pattern := patterns[i]
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
	}

	return false
}

func collectInputPaths(buildPath string) ([]string, error) {
	var err error

	outputPaths, err := collectOutputPaths(buildPath)
	if err != nil {
		return nil, err
	}

	rootPath := filepath.Dir(buildPath)
	inputPaths := []string{}

	err = filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, buildPath) && excludePathPrefix(outputPaths, path) {
			return nil
		}

		if info.IsDir() {
			if len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
				return fs.SkipDir
			} else {
				inputPaths = append(inputPaths, path)
			}
		}

		return nil
	})

	return inputPaths, nil
}

func runWatcher(ctx context.Context, wg *sync.WaitGroup, filePaths []string, matchPatterns []string, buildFunc func()) {
	var (
		watcher *fsnotify.Watcher
		rebuild bool
		err     error
	)

	logger := ctx.Value(LoggerContextKey).(*log.Logger)

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Println(err)
		goto WatchDone
	}

	for i := range filePaths {
		if err = watcher.Add(filePaths[i]); err != nil {
			logger.Println(err)
			goto WatchDone
		}
	}

	logger.Println("starting watch")
	for {
		select {
		case <-ctx.Done():
			watcher.Close()
			logger.Println("stopping watch")
			goto WatchDone
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			rebuild = false

			baseName := filepath.Base(event.Name)
			if strings.HasPrefix(baseName, ".") && strings.HasSuffix(baseName, ".swp") {
				continue
			}

			if strings.HasSuffix(baseName, "~") {
				continue
			}

			if baseName == "4913" { // Vim creates this file to test for write permission
				continue
			}

			fmt.Println("event", event)

			if !includePathMatch(matchPatterns, event.Name) {
				continue
			}

			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err != nil {
					continue
				} else {
					if info.IsDir() {
						_ = watcher.Add(event.Name)
					} else {
						rebuild = true
					}
				}
			}

			if event.Has(fsnotify.Write) {
				rebuild = true
			}

			if rebuild {
				logger.Println("change", event.Name)
				buildFunc()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Println(err)
		}
	}
WatchDone:
	wg.Done()
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch a project directory",
	Long: `Watch a project directory and rebuild files when they change.

A project directory is a directory with a build.yaml file.
    `,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		logger := ctx.Value(LoggerContextKey).(*log.Logger)
		logger.SetPrefix("👀 ")

		buildPath, err := getConfigurationPath(args[0], configName)
		if err != nil {
			logger.Fatal(err)
		}

		wg := new(sync.WaitGroup)
		watchCtx, _ := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		inputPaths, err := collectInputPaths(buildPath)
		if err != nil {
			logger.Fatal(err)
		}

		wg.Add(1)
		go runWatcher(watchCtx, wg, inputPaths, []string{}, func() {
			buildCtx, _ := context.WithCancel(watchCtx)
			wg.Add(1)
			runBuild(buildCtx, wg, buildPath)
		})

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&serverRoot, "serve", "s", "./build", "Server root directory")
	addServeFlags(watchCmd)
}

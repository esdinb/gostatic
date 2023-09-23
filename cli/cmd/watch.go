package cmd

import (
    "os"
    "os/signal"
    "syscall"
    "log"
    "context"
    "sync"
    "strings"
    "sort"
    "path/filepath"
    "io/fs"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/cobra"
)

func getOutputPaths(buildPath string) ([]string, error) {
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

func excludePath(prefixes []string, path string) bool {
    for i := range prefixes {
        if strings.HasPrefix(path, prefixes[i]) {
            return true
        }
    }

    return false
}

func runWatcher(ctx context.Context, wg *sync.WaitGroup, buildPath string) {

    var (
        watcher *fsnotify.Watcher
        outputPaths []string
        rebuild bool
        err error
    )

    logger := ctx.Value(LoggerContextKey).(*log.Logger)
    rootPath := filepath.Dir(buildPath)

    watcher, err = fsnotify.NewWatcher()
    if err != nil {
        logger.Println(err)
        goto WatchDone
    }

    outputPaths, err = getOutputPaths(buildPath)
    if err != nil {
        logger.Println(err)
        goto WatchDone
    }

    err = filepath.Walk(rootPath, func (path string, info fs.FileInfo, err error) error {
        if !strings.HasSuffix(path, buildPath) && excludePath(outputPaths, path) {
            return nil
        }

        if info.IsDir() {
            if len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
                return fs.SkipDir
            } else {
                err = watcher.Add(path)
                if err != nil {
                    return err
                }
            }
        }

        return nil
    })
    if err != nil {
        logger.Println(err)
        goto WatchDone
    }

    logger.Println("starting watch", rootPath)
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

            if event.Has(fsnotify.Create) {
                _ = watcher.Add(filepath.Dir(event.Name))
                rebuild = true
            }

            if event.Has(fsnotify.Write) {
                rebuild = true
            }

            if rebuild {
                logger.Println("change", event.Name)
                buildCtx, _ := context.WithCancel(ctx)
                wg.Add(1)
                go runBuild(buildCtx, wg, buildPath)
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

        wg.Add(1)
        go runWatcher(watchCtx, wg, buildPath)

        wg.Wait()
    },
}

func init() {
    rootCmd.AddCommand(watchCmd)

    watchCmd.Flags().StringVarP(&serverAddress, "address", "a", "", "Server listen address")
    watchCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Server listen port")
}

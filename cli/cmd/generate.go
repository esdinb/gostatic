package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gostatic/pkg/builder"

	"github.com/spf13/cobra"
)

func runGenerate(ctx context.Context, wg *sync.WaitGroup, section builder.BuildSection, rootPath string) {
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Example: "generate banner:example ./index.html .",
	Short:   "One-shot build from a template",
	Long: `Build a file from a template named on the command line.

Arguments to the generate command are an optional pipeline (a list of named transformations),
an input file path pattern and an output file path in this order. 

The arguments to the generate command matches keys of build.yaml configuration file sections.
    `,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		logger := log.New(os.Stdout, "🧪 ", 0)

		n := len(args)
		inPath, outPath, transformerNames := args[n-2], args[n-1], args[:n-2]
		pipeline := make([]builder.BuildTransformation, len(transformerNames))
		for i := range pipeline {
			pipeline[i] = builder.BuildTransformation(transformerNames[i])
		}

		rootPath := "."

		section := builder.BuildSection{inPath, outPath, ((builder.Pipeline)(pipeline))}

		wg := new(sync.WaitGroup)

		runner := func(ctx context.Context) {
			logger.Println("rebuilding")
			if err := section.Build(ctx, rootPath); err != nil {
				logger.Fatal(err)
			}
		}

		runner(ctx)

		if watchFiles {
			watchCtx, _ := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

			wg.Add(1)
			go runWatcher(watchCtx, wg, []string{"."}, []string{inPath}, runner)
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolVarP(&watchFiles, "watch", "w", false, "Watch files")
}

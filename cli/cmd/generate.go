package cmd

import (
    "os"
    "log"

    "gostatic/pkg/builder"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
    Aliases: []string{"gen"},
    Example: "generate banner:example ./index.html .",
	Short: "Generate an XML file from a template",
	Long: `Generate an XML file from an XML template.

This command takes as arguments a number of named transformations, a source path and a destination path. 

This command reads the source, applies the transformations to the source and writes the transformation result to destination.

    `,
    Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

        logger := log.New(os.Stdout, "🧪 ", 0)

        n := len(args)
        inPath, outPath, transformerNames := args[n - 2], args[n - 1], args[:n - 2]
        pipeline := make([]builder.BuildTransformation, len(transformerNames))
        for i := range pipeline {
            pipeline[i] = builder.BuildTransformation(transformerNames[i])
        }

        section := builder.BuildSection{inPath, outPath, ((builder.Pipeline)(pipeline))}
        if err := section.Build("."); err != nil {
            logger.Fatal(err)
        }
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

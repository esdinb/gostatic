package cmd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	configPerms    fs.FileMode = 0644
	configTemplate string      = `# build configuration
- in: /index.html
  out: /build/index.html
  pipeline:
    - transform:inline
    - banner:gostatic
`
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a build configuration template",
	Long: `Generate a template build.yaml file in the current directory or in the directory named on the command line.

    ⚡️ + 🥚 = 🐣.
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		var info fs.FileInfo
		var err error
		if len(args) == 1 {
			path = args[0]
		} else {
			path = "."
		}

		logger := log.New(os.Stdout, "🧩 ", 0)

		info, err = os.Lstat(path)
		if err != nil {
			logger.Panicf("🤮 %v\n", err)
		}

		if !info.IsDir() {
			logger.Println("output path is not a directory 🤥")
			os.Exit(1)
		}

		buffer := []byte(configTemplate)
		if err = os.WriteFile(filepath.Join(path, configName), buffer, configPerms); err != nil {
			logger.Panicln(err)
		}
		logger.Println("Write", configName, "to", path, "🐣")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

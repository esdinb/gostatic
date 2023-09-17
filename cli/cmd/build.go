package cmd

import (
    "os"
    "log"
    "path/filepath"

    "gostatic/pkg/builder"

	"github.com/spf13/cobra"
    yaml "gopkg.in/yaml.v3"
)

const (
    configName string = "./build.yaml"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a site from configuration file",
	Long: `

    The build command reads build.yaml configuration file and generates site accordingly.

            🛠️ 🪜 🔩 🪚 🧱 🧪 🛁

`,
    Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        var configPath string
        var err error

        logger := log.New(os.Stdout, "🛠️ ", 0)

        if len(args) == 0 {
            cwd, err := os.Getwd()
            if err != nil {
                logger.Fatal(err)
            }
            configPath = filepath.Join(cwd, configName)
        } else {
            fileInfo, err := os.Stat(args[0])
            if err != nil {
                logger.Fatal(err)
            }
            if fileInfo.IsDir() {
                configPath = filepath.Join(args[0], configName)
            }
        }

        bytes, err := os.ReadFile(configPath)
        if err != nil {
            logger.Fatal(err)
        }

        var config []builder.BuildSection
        err = yaml.Unmarshal(bytes, &config)
        if err != nil {
            logger.Fatal(err)
        }

        rootPath := filepath.Dir(configPath)

        for i := range config {
            section := config[i]
            err := section.Build(rootPath)
            if err != nil {
                logger.Println(err)
            }
        }
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

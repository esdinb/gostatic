package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version string",
	Run: func(cmd *cobra.Command, args []string) {
		root := cmd.Root()
		root.SetArgs([]string{"--version"})
		root.Execute()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

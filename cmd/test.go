package cmd

import (
	"fmt"
	"os"

	"minder-test/internal/runner"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [path]",
	Short: "Run Minder rule tests",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("Invalid path:", path)
			os.Exit(1)
		}

		if info.IsDir() {
			runner.RunDirectory(path)
		} else {
			fmt.Println("Provide a directory path")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

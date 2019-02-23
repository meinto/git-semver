package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmdOptions struct {
	NextVersionType string
	RepoPath        string
}

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.Flags().StringVarP(&rootCmdOptions.NextVersionType, "version", "v", "", "version of this program")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

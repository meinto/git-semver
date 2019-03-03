package cmd

//go:generate ./.circleci/generate-assets.sh

import (
	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion := cmdUtil.GetVersion()
	},
}

func Execute() {
	err := rootCmd.Execute()
	cmdUtil.LogFatalOnErr(err)
}

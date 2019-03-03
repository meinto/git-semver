package cmd

import (
	"fmt"

	"github.com/gobuffalo/packr"
	"github.com/meinto/git-semver/cmd/internal/flags"
	"github.com/meinto/git-semver/cmd/internal"
	"github.com/spf13/cobra"
)

func init() {
	flags.RootCmdFlags.Init(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	Run: func(cmd *cobra.Command, args []string) {
		box := packr.NewBox("../buildAssets")
		version, err := box.FindString("VERSION")
		internal.LogFatalOnErr(err)
		fmt.Printf("Version of git-semver: %s\n", version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	internal.LogFatalOnErr(err)
}

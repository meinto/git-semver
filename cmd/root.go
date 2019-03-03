package cmd

import (
	"fmt"

	"github.com/gobuffalo/packr"
	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	Run: func(cmd *cobra.Command, args []string) {
		box := packr.NewBox("../buildAssets")
		version, err := box.FindString("VERSION")
		cmdUtil.LogFatalOnErr(err)
		fmt.Printf("Version of git-semver: %s\n", version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	cmdUtil.LogFatalOnErr(err)
}

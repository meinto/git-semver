package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/meinto/git-semver/cmd/internal/flags"

	"github.com/meinto/git-semver/cmd/internal"
	"github.com/meinto/git-semver/util"
	"github.com/pkg/errors" 
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(getCmd) 
	flags.GetCmdFlags.Init(getCmd)
}

var getCmd = &cobra.Command{     
	Use:   "get",
	Short: "get version number",
	Run: func(cmd *cobra.Command, args []string) {
		gitRepoPath, err := filepath.Abs(flags.GetCmdFlags.RepoPath())
		internal.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

		pathToVersionFile := internal.VersionFilePath(gitRepoPath, viper.GetString("versionFileName"))

		_, err = os.Stat(pathToVersionFile) 
		internal.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

		versionFile, err := os.Open(pathToVersionFile)
		internal.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", viper.GetString("versionFileName"))))
		defer versionFile.Close()

		byteValue, err := ioutil.ReadAll(versionFile)
		internal.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
		currentVersion := internal.GetVersion(viper.GetString("versionFileType"), byteValue)

		if len(args) > 0 {
			nextVersionType := args[0]
			internal.ValidateNextVersionType(nextVersionType)

			nextVersion, err := util.NextVersion(currentVersion, nextVersionType)
			internal.LogFatalOnErr(err)

			internal.PrintNextVersion(nextVersionType, nextVersion, flags.GetCmdFlags.PrintRaw())
		} else {
			internal.PrintCurrentVersion(currentVersion, flags.GetCmdFlags.PrintRaw())
		}

	},
}

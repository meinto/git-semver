package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	semverUtil "github.com/meinto/git-semver/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmdOptions struct {
	RepoPath string
	PrintRaw bool
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
	getCmd.Flags().BoolVarP(&getCmdOptions.PrintRaw, "raw", "r", false, "print only the plain version number")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get version number",
	Run: func(cmd *cobra.Command, args []string) {
		gitRepoPath, err := filepath.Abs(getCmdOptions.RepoPath)
		cmdUtil.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

		pathToVersionFile := cmdUtil.VersionFilePath(gitRepoPath, viper.GetString("versionFileName"))

		_, err = os.Stat(pathToVersionFile)
		cmdUtil.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

		versionFile, err := os.Open(pathToVersionFile)
		cmdUtil.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", viper.GetString("versionFileName"))))
		defer versionFile.Close()

		byteValue, err := ioutil.ReadAll(versionFile)
		cmdUtil.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
		currentVersion := cmdUtil.GetVersion(viper.GetString("versionFileType"), byteValue)

		if len(args) > 0 {
			nextVersionType := args[0]
			cmdUtil.ValidateNextVersionType(nextVersionType)

			nextVersion, err := semverUtil.NextVersion(currentVersion, nextVersionType)
			cmdUtil.LogFatalOnErr(err)

			cmdUtil.PrintNextVersion(nextVersionType, nextVersion, getCmdOptions.PrintRaw)
		} else {
			cmdUtil.PrintCurrentVersion(currentVersion, getCmdOptions.PrintRaw)
		}

	},
}

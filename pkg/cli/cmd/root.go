package cmd

import (
	"log"
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"github.com/meinto/git-semver/util" 

	"github.com/gobuffalo/packr"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal/flags"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal"
	"github.com/spf13/cobra" 
	"github.com/pkg/errors" 
)

func init() {
	flags.RootCmdFlags.Init(rootCmd) 
}

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	PreRun: func(cmd *cobra.Command, args []string) {
		flags.RootCmdFlags.PreRun(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		if flags.RootCmdFlags.CreateTag() {
			gitRepoPath, err := filepath.Abs(flags.RootCmdFlags.RepoPath()) 
			internal.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

			pathToVersionFile := internal.VersionFilePath(gitRepoPath, flags.RootCmdFlags.VersionFile())

			_, err = os.Stat(pathToVersionFile) 
			internal.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

			versionFile, err := os.Open(pathToVersionFile)
			internal.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", flags.RootCmdFlags.VersionFile())))
			defer versionFile.Close()

			byteValue, err := ioutil.ReadAll(versionFile)
			internal.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
			currentVersion := internal.GetVersion(flags.RootCmdFlags.VersionFileFormat(), byteValue)
 
			util.MakeGitTag(gitRepoPath, currentVersion)
		}

		if flags.RootCmdFlags.Push() {
			if err := util.Push(flags.RootCmdFlags.RepoPath(), flags.RootCmdFlags.SSHFilePath()); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}

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

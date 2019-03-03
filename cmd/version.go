package cmd

import (
	"github.com/meinto/git-semver/cmd/internal/flags"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath" 
   
	"github.com/meinto/git-semver/cmd/internal"
	"github.com/meinto/git-semver/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra" 
)

func init() { 
	rootCmd.AddCommand(versionCmd)
	flags.VersionCmdFlags.Init(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "create new version for repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nextVersionType := args[0]
		internal.ValidateNextVersionType(nextVersionType)

		gitRepoPath, err := filepath.Abs(flags.VersionCmdFlags.RepoPath())
		internal.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

		pathToVersionFile := internal.VersionFilePath(
			gitRepoPath, 
			flags.VersionCmdFlags.VersionFile(),
		)

		_, err = os.Stat(pathToVersionFile)
		internal.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

		versionFile, err := os.Open(pathToVersionFile)
		internal.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", flags.VersionCmdFlags.VersionFile())))
		defer versionFile.Close()

		fileContent, err := ioutil.ReadAll(versionFile)
		internal.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
		currentVersion := internal.GetVersion(flags.VersionCmdFlags.VersionFileType(), fileContent)

		nextVersion, err := util.NextVersion(currentVersion, nextVersionType)
		internal.LogFatalOnErr(err)

		log.Println("new version: ", nextVersion)
		if flags.VersionCmdFlags.DryRun() {
			log.Println("dry run finished...")
			os.Exit(1)
		}

		internal.ValidateReadyForPushingChanges(
			flags.VersionCmdFlags.RepoPath(),
			flags.VersionCmdFlags.SSHFilePath(),
			flags.VersionCmdFlags.Push(),
		)

		err = internal.WriteVersion(
			flags.VersionCmdFlags.VersionFileType(),
			flags.VersionCmdFlags.VersionFile(),
			nextVersion,
			fileContent,
		)
		internal.LogFatalOnErr(err)

		err = util.AddVersionChanges(
			flags.VersionCmdFlags.RepoPath(),
			flags.VersionCmdFlags.VersionFile(),
			nextVersion,
			flags.VersionCmdFlags.Author(),
			flags.VersionCmdFlags.Email(),
		)
		internal.LogFatalOnErr(err)

		var createGitTagError error
		if flags.VersionCmdFlags.CreateTag() {
			createGitTagError = util.MakeGitTag(flags.VersionCmdFlags.RepoPath(), nextVersion)
		}

		if flags.VersionCmdFlags.Push() && createGitTagError == nil {
			if err = util.Push(flags.VersionCmdFlags.RepoPath(), flags.VersionCmdFlags.SSHFilePath()); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}
	},
}

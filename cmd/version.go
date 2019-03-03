package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath" 
   
	"github.com/meinto/git-semver/cmd/internal"
	"github.com/meinto/git-semver/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra" 
	"github.com/spf13/viper"
)

var versionCmdOptions struct {
	RepoPath          string 
	VersionFile       string
	VersionFileFormat string  
	DryRun            bool
	CreateTag         bool
	Push              bool
	Author            string
	Email             string
	SSHFilePath       string
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Author, "author", "a", "semver", "name of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Email, "email", "e", "semver@no-reply.git", "email of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.VersionFile, "outfile", "o", "semver.json", "name of version file")
	versionCmd.Flags().StringVarP(&versionCmdOptions.VersionFileFormat, "outfileFormat", "f", "json", "format of outfile (json, raw)")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.DryRun, "dryrun", "d", false, "only log how version number would change")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.CreateTag, "tag", "t", false, "create a git tag")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.Push, "push", "P", false, "push git tags and version changes")

	viper.BindPFlag("versionFileName", versionCmd.Flags().Lookup("outfile"))
	viper.BindPFlag("versionFileType", versionCmd.Flags().Lookup("outfileFormat"))
	viper.BindPFlag("tagVersions", versionCmd.Flags().Lookup("tag"))
	viper.BindPFlag("pushChanges", versionCmd.Flags().Lookup("push"))
	viper.BindPFlag("author", versionCmd.Flags().Lookup("author"))
	viper.BindPFlag("email", versionCmd.Flags().Lookup("email"))

	defaultSSHFilePath, err := util.GetDefaultSSHFilePath()
	internal.LogOnError(err)
	versionCmd.Flags().StringVar(&versionCmdOptions.SSHFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "create new version for repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nextVersionType := args[0]
		internal.ValidateNextVersionType(nextVersionType)

		gitRepoPath, err := filepath.Abs(versionCmdOptions.RepoPath)
		internal.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

		pathToVersionFile := internal.VersionFilePath(gitRepoPath, viper.GetString("versionFileName"))

		_, err = os.Stat(pathToVersionFile)
		internal.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

		versionFile, err := os.Open(pathToVersionFile)
		internal.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", viper.GetString("versionFileName"))))
		defer versionFile.Close()

		fileContent, err := ioutil.ReadAll(versionFile)
		internal.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
		currentVersion := internal.GetVersion(viper.GetString("versionFileType"), fileContent)

		nextVersion, err := util.NextVersion(currentVersion, nextVersionType)
		internal.LogFatalOnErr(err)

		log.Println("new version: ", nextVersion)
		if versionCmdOptions.DryRun {
			log.Println("dry run finished...")
			os.Exit(1)
		}

		internal.ValidateReadyForPushingChanges(
			versionCmdOptions.RepoPath,
			versionCmdOptions.SSHFilePath,
			viper.GetBool("pushChanges"),
		)

		err = internal.WriteVersion(
			viper.GetString("versionFileType"),
			viper.GetString("versionFileName"),
			nextVersion,
			fileContent,
		)
		internal.LogFatalOnErr(err)

		err = util.AddVersionChanges(
			versionCmdOptions.RepoPath,
			viper.GetString("versionFileName"),
			nextVersion,
			viper.GetString("author"),
			viper.GetString("email"),
		)
		internal.LogFatalOnErr(err)

		var createGitTagError error
		if viper.GetBool("tagVersions") {
			createGitTagError = util.MakeGitTag(versionCmdOptions.RepoPath, nextVersion)
		}

		if viper.GetBool("pushChanges") && createGitTagError == nil {
			if err = util.Push(versionCmdOptions.RepoPath, versionCmdOptions.SSHFilePath); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}
	},
}

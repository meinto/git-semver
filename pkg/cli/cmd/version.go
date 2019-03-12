package cmd

import (
	"log"
	"os"

	semver "github.com/meinto/git-semver"
	"github.com/meinto/git-semver/file"
	"github.com/meinto/git-semver/git"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmdFlags struct {
	dryRun bool
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionCmdFlags.dryRun, "dryrun", "d", false, "only log how version number would change")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "create new version for repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		l := internal.NewLogger(rootCmdFlags.verbose)

		gs := git.NewGitService(viper.GetString("gitPath"))
		repoPath, err := gs.GitRepoPath()
		l.LogFatalOnError(err)

		versionFilepath := repoPath + "/" + viper.GetString("versionFile")
		fs := file.NewVersionFileService(versionFilepath)

		currentVersion, err := fs.ReadVersionFromFile(viper.GetString("versionFileType"))
		l.LogFatalOnError(err)

		vs, err := semver.NewVersion(currentVersion)
		l.LogFatalOnError(err)

		nextVersionType := args[0]
		nextVersion, err := vs.SetNext(nextVersionType)
		l.LogFatalOnError(err)

		log.Println("new version will be: ", nextVersion)

		if versionCmdFlags.dryRun {
			log.Println("dry run finished...")
			os.Exit(1)
		}

		if viper.GetBool("pushChanges") {
			if isClean, err := gs.IsRepoClean(); !isClean || err != nil {
				log.Fatal(err)
			}
		}

		fs.WriteVersionFile(viper.GetString("versionFileType"), nextVersion)

		if viper.GetBool("pushChanges") {
			gs.AddVersionChanges(versionFilepath)
			gs.CommitVersionChanges(nextVersion)
		}

		var createGitTagError error
		if viper.GetBool("tagVersions") {
			createGitTagError = gs.CreateTag(nextVersion)
		}

		if viper.GetBool("pushChanges") && createGitTagError == nil {
			if err = gs.Push(); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}
	},
}

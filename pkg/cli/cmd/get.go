package cmd

import (
	"fmt"

	"github.com/meinto/git-semver"
	"github.com/meinto/git-semver/file"
	"github.com/meinto/git-semver/git"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var getCmdFlags struct {
	printRaw bool
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVarP(&getCmdFlags.printRaw, "raw", "r", false, "print only the plain version number")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get version number",
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

		if len(args) > 0 {
			nextVersionType := args[0]
			nextVersion, err := vs.Get(nextVersionType)
			l.LogFatalOnError(err)

			printNextVersion(nextVersionType, nextVersion, getCmdFlags.printRaw)
		} else {
			currentVersion, err := vs.Get("")
			l.LogFatalOnError(err)

			printCurrentVersion(currentVersion, getCmdFlags.printRaw)
		}
	},
}

func printNextVersion(nextVersionType, nextVersion string, raw bool) {
	printVersion(nextVersion, fmt.Sprintf("Next %s version: ", nextVersionType), raw)
}

func printCurrentVersion(currentVersion string, raw bool) {
	printVersion(currentVersion, "Current version: ", raw)
}

func printVersion(nextVersion, message string, raw bool) {
	if !raw {
		fmt.Print(message)
	}
	fmt.Println(nextVersion)
}

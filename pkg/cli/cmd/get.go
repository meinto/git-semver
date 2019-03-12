package cmd

import (
	"fmt"
	"log"

	"github.com/meinto/git-semver"
	"github.com/meinto/git-semver/file"
	"github.com/meinto/git-semver/git"

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

		gs := git.NewGitService(viper.GetString("gitPath"))
		repoPath, err := gs.GitRepoPath()
		if err != nil {
			log.Fatal(err)
		}

		versionFilepath := repoPath + "/" + viper.GetString("versionFile")
		fs := file.NewVersionFileService(versionFilepath)

		currentVersion, err := fs.ReadVersionFromFile(viper.GetString("versionFileType"))
		if err != nil {
			log.Fatal(err)
		}

		vs, err := semver.NewVersion(currentVersion)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) > 0 {
			nextVersionType := args[0]
			nextVersion, err := vs.Get(nextVersionType)
			if err != nil {
				log.Fatal(err)
			}

			printNextVersion(nextVersionType, nextVersion, getCmdFlags.printRaw)
		} else {
			currentVersion, err := vs.Get("")
			if err != nil {
				log.Fatal(err)
			}

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

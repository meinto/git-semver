package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/meinto/git-semver/file"
	"github.com/meinto/git-semver/git"

	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"

	// "github.com/pkg/errors"
	"github.com/spf13/viper"
)

var rootCmdFlags struct {
	gitPath         string
	verbose         bool
	push            bool
	createTag       bool
	versionFile     string
	versionFileType string
	sshFilePath     string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootCmdFlags.gitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdFlags.verbose, "verbose", "v", false, "more logs")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdFlags.push, "push", "P", false, "push git tags")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdFlags.createTag, "tag", "T", false, "create a git tag")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.versionFile, "versionFile", "f", "VERSION", "name of version file")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.versionFileType, "versionFileType", "t", "raw", "type of version file (json, raw)")

	defaultSSHFilePath, err := file.GetDefaultSSHFilePath()
	if err != nil {
		log.Println(err)
	}
	rootCmd.PersistentFlags().StringVar(&rootCmdFlags.sshFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")

	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
	viper.BindPFlag("versionFile", rootCmd.PersistentFlags().Lookup("versionFile"))
	viper.BindPFlag("versionFileType", rootCmd.PersistentFlags().Lookup("versionFileType"))
}

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	Run: func(cmd *cobra.Command, args []string) {

		// if flags.RootCmdFlags.CreateTag() {
		// 	gitRepoPath, err := filepath.Abs(flags.RootCmdFlags.RepoPath())
		// 	internal.LogFatalOnErr(errors.Wrap(err, "cannot resolve repo path"))

		// 	pathToVersionFile := internal.VersionFilePath(gitRepoPath, flags.RootCmdFlags.VersionFile())

		// 	_, err = os.Stat(pathToVersionFile)
		// 	internal.LogFatalOnErr(errors.Wrap(err, "version file doesn't exist"))

		// 	versionFile, err := os.Open(pathToVersionFile)
		// 	internal.LogFatalOnErr(errors.Wrap(err, fmt.Sprintf("cannot read %s", flags.RootCmdFlags.VersionFile())))
		// 	defer versionFile.Close()

		// 	byteValue, err := ioutil.ReadAll(versionFile)
		// 	internal.LogFatalOnErr(errors.Wrap(err, "cannot read file"))
		// 	currentVersion := internal.GetVersion(flags.RootCmdFlags.VersionFileFormat(), byteValue)

		// 	util.MakeGitTag(gitRepoPath, currentVersion)
		// }

		// if flags.RootCmdFlags.Push() {
		// 	if err := util.Push(flags.RootCmdFlags.RepoPath(), flags.RootCmdFlags.SSHFilePath()); err != nil {
		// 		log.Fatalf("cannot push tag: %s", err.Error())
		// 	}
		// }

		g := git.NewGitService(viper.GetString("gitPath"))
		repoPath, _ := g.GitRepoPath()

		box := packr.NewBox(repoPath + "/buildAssets")
		version, err := box.FindString("VERSION")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Version of git-semver: %s\n", version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

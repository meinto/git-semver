package cmd

import (
	"github.com/meinto/git-semver/git"
	// "log"
	"fmt"
	// "path/filepath"
	// "os"
	// "io/ioutil"
	// "github.com/meinto/git-semver/util" 

	"github.com/gobuffalo/packr"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal/flags"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal"
	"github.com/spf13/cobra" 
	// "github.com/pkg/errors" 
	"github.com/spf13/viper"
)

var rootCmdOptions struct {
	// Author                string
	GitPath               string
	// UseBuiltInGitBindings bool
}

func init() {
	flags.RootCmdFlags.Init(rootCmd) 
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
}

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "standalone tool to version your gitlab repo with semver",
	PreRun: func(cmd *cobra.Command, args []string) {
		flags.RootCmdFlags.PreRun(cmd)
	},
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

		box := packr.NewBox(repoPath+"/buildAssets")
		version, err := box.FindString("VERSION")
		internal.LogFatalOnErr(err)
		fmt.Printf("Version of git-semver: %s\n", version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	internal.LogFatalOnErr(err)
}

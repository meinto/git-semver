package flags

import (
	"log"

	"github.com/meinto/git-semver/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootCmdFlagsType struct {
	verbose         bool
	repoPath        string
	createTag       bool
	push            bool
	versionFile     string
	versionFileType string
	sshFilePath     string
}

var RootCmdFlags rootCmdFlagsType

func (fs *rootCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&RootCmdFlags.verbose, "verbose", "v", false, "more logs")
	cmd.Flags().StringVarP(&RootCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().BoolVarP(&RootCmdFlags.push, "push", "P", false, "push git tags")
	cmd.Flags().BoolVarP(&RootCmdFlags.createTag, "tag", "T", false, "create a git tag")
	cmd.Flags().StringVarP(&RootCmdFlags.versionFile, "versionFile", "f", "VERSION", "name of version file")
	cmd.Flags().StringVarP(&RootCmdFlags.versionFileType, "versionFileType", "t", "raw", "type of version file (json, raw)")

	defaultSSHFilePath, err := util.GetDefaultSSHFilePath()
	util.LogOnError(err, RootCmdFlags.Verbose())
	cmd.Flags().StringVar(&RootCmdFlags.sshFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")

	viper.SetConfigName("semver.config")
	viper.SetConfigType("json")
	repoPath, _ := cmd.Flags().GetString("path")
	viper.AddConfigPath(repoPath)
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("there is no semver.config file: ", err)
	}
}

func (fs *rootCmdFlagsType) PreRun(cmd *cobra.Command) {
	bindViperFlag("versionFile", cmd.Flags().Lookup("versionFile"))
	bindViperFlag("versionFileType", cmd.Flags().Lookup("versionFileType"))
}

func (fs *rootCmdFlagsType) Verbose() bool {
	return fs.verbose
}

func (fs *rootCmdFlagsType) RepoPath() string {
	return fs.repoPath
}

func (fs *rootCmdFlagsType) CreateTag() bool {
	return fs.createTag
}

func (fs *rootCmdFlagsType) Push() bool {
	return fs.push
}

func (fs *rootCmdFlagsType) VersionFile() string {
	return viper.GetString("versionFile")
}

func (fs *rootCmdFlagsType) VersionFileFormat() string {
	return viper.GetString("versionFileType")
}

func (fs *rootCmdFlagsType) SSHFilePath() string {
	return fs.sshFilePath
}

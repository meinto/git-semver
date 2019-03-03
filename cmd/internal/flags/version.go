package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/meinto/git-semver/util"
)

type versionCmdFlagsType struct {
	repoPath        string
	versionFile     string
	versionFileType string
	dryRun          bool
	createTag       bool
	push            bool
	author          string
	email           string
	sshFilePath     string
}

var VersionCmdFlags versionCmdFlagsType

func (fs *versionCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&VersionCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().StringVarP(&VersionCmdFlags.author, "author", "a", "semver", "name of the author")
	cmd.Flags().StringVarP(&VersionCmdFlags.email, "email", "e", "semver@no-reply.git", "email of the author")
	cmd.Flags().BoolVarP(&VersionCmdFlags.dryRun, "dryrun", "d", false, "only log how version number would change")
	cmd.Flags().BoolVarP(&VersionCmdFlags.createTag, "tag", "T", false, "create a git tag")
	cmd.Flags().BoolVarP(&VersionCmdFlags.push, "push", "P", false, "push git tags and version changes")
	versionFileFlags(cmd)

	bindViperFlag("tagVersions", cmd.Flags().Lookup("tag"))
	bindViperFlag("pushChanges", cmd.Flags().Lookup("push"))
	bindViperFlag("author", cmd.Flags().Lookup("author"))
	bindViperFlag("email", cmd.Flags().Lookup("email"))

	defaultSSHFilePath, err := util.GetDefaultSSHFilePath()
	util.LogOnError(err, RootCmdFlags.Verbose())
	cmd.Flags().StringVar(&VersionCmdFlags.sshFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")

}

func (fs *versionCmdFlagsType) RepoPath() string {
	return fs.repoPath
}

func (fs *versionCmdFlagsType) VersionFile() string {
	return viper.GetString("versionFile")
}

func (fs *versionCmdFlagsType) VersionFileType() string {
	return viper.GetString("versionFileType")
}

func (fs *versionCmdFlagsType) DryRun() bool {
	return fs.dryRun
}

func (fs *versionCmdFlagsType) CreateTag() bool {
	return viper.GetBool("tagVersions")
}

func (fs *versionCmdFlagsType) Push() bool {
	return viper.GetBool("pushChanges")
}

func (fs *versionCmdFlagsType) Author() string {
	return viper.GetString("author")
}

func (fs *versionCmdFlagsType) Email() string {
	return viper.GetString("email")
}

func (fs *versionCmdFlagsType) SSHFilePath() string {
	return fs.sshFilePath
}

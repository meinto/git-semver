package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/meinto/git-semver/util"
)

type versionCmdFlags struct {
	repoPath          string
	versionFile       string
	versionFileFormat string
	dryRun            bool
	createTag         bool
	push              bool
	author            string
	email             string
	sshFilePath       string
}

var VersionCmdFlags versionCmdFlags

func (fs *versionCmdFlags) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&VersionCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().StringVarP(&VersionCmdFlags.author, "author", "a", "semver", "name of the author")
	cmd.Flags().StringVarP(&VersionCmdFlags.email, "email", "e", "semver@no-reply.git", "email of the author")
	cmd.Flags().StringVarP(&VersionCmdFlags.versionFile, "outfile", "o", "semver.json", "name of version file")
	cmd.Flags().StringVarP(&VersionCmdFlags.versionFileFormat, "outfileFormat", "f", "json", "format of outfile (json, raw)")
	cmd.Flags().BoolVarP(&VersionCmdFlags.dryRun, "dryrun", "d", false, "only log how version number would change")
	cmd.Flags().BoolVarP(&VersionCmdFlags.createTag, "tag", "t", false, "create a git tag")
	cmd.Flags().BoolVarP(&VersionCmdFlags.push, "push", "P", false, "push git tags and version changes")

	bindViperFlag("versionFileName", cmd.Flags().Lookup("outfile"))
	bindViperFlag("versionFileType", cmd.Flags().Lookup("outfileFormat"))
	bindViperFlag("tagVersions", cmd.Flags().Lookup("tag"))
	bindViperFlag("pushChanges", cmd.Flags().Lookup("push"))
	bindViperFlag("author", cmd.Flags().Lookup("author"))
	bindViperFlag("email", cmd.Flags().Lookup("email"))

	defaultSSHFilePath, err := util.GetDefaultSSHFilePath()
	util.LogOnError(err, RootCmdFlags.Verbose())
	cmd.Flags().StringVar(&VersionCmdFlags.sshFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")

}

func (fs *versionCmdFlags) RepoPath() string {
	return fs.repoPath
}

func (fs *versionCmdFlags) VersionFile() string {
	return viper.GetString("versionFileName")
}

func (fs *versionCmdFlags) VersionFileFormat() string {
	return viper.GetString("versionFileType")
}

func (fs *versionCmdFlags) DryRun() bool {
	return fs.dryRun
}

func (fs *versionCmdFlags) CreateTag() bool {
	return viper.GetBool("tagVersions")
}

func (fs *versionCmdFlags) Push() bool {
	return viper.GetBool("pushChanges")
}

func (fs *versionCmdFlags) Author() string {
	return viper.GetString("author")
}

func (fs *versionCmdFlags) Email() string {
	return viper.GetString("email")
}

func (fs *versionCmdFlags) SSHFilePath() string {
	return fs.sshFilePath
}

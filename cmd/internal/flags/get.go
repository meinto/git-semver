package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getCmdFlagsType struct {
	repoPath          string
	versionFile       string
	versionFileFormat string
	printRaw          bool
}

var GetCmdFlags getCmdFlagsType

func (fs *getCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&GetCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().BoolVarP(&GetCmdFlags.printRaw, "raw", "r", false, "print only the plain version number")
	versionFileFlags(cmd)
}

func (fs *getCmdFlagsType) RepoPath() string {
	return fs.repoPath
}

func (fs *getCmdFlagsType) VersionFile() string {
	return viper.GetString("versionFileName")
}

func (fs *getCmdFlagsType) VersionFileFormat() string {
	return viper.GetString("versionFileType")
}

func (fs *getCmdFlagsType) PrintRaw() bool {
	return fs.printRaw
}

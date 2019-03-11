package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getCmdFlagsType struct {
	repoPath        string
	versionFile     string
	versionFileType string
	printRaw        bool
}

var GetCmdFlags getCmdFlagsType

func (fs *getCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&GetCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().BoolVarP(&GetCmdFlags.printRaw, "raw", "r", false, "print only the plain version number")
	cmd.Flags().StringVarP(&GetCmdFlags.versionFile, "versionFile", "f", "VERSION", "name of version file")
	cmd.Flags().StringVarP(&GetCmdFlags.versionFileType, "versionFileType", "t", "raw", "type of version file (json, raw)")
}

func (fs *getCmdFlagsType) PreRun(cmd *cobra.Command) {
	bindViperFlag("versionFile", cmd.Flags().Lookup("versionFile"))
	bindViperFlag("versionFileType", cmd.Flags().Lookup("versionFileType"))

	repoPath, _ := cmd.Flags().GetString("path")
	LaodViperConfig(repoPath)
}

func (fs *getCmdFlagsType) RepoPath() string {
	return fs.repoPath
}

func (fs *getCmdFlagsType) VersionFile() string {
	return viper.GetString("versionFile")
}

func (fs *getCmdFlagsType) VersionFileFormat() string {
	return viper.GetString("versionFileType")
}

func (fs *getCmdFlagsType) PrintRaw() bool {
	return fs.printRaw
}

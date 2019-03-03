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
	cmd.Flags().StringVarP(&GetCmdFlags.versionFile, "outfile", "o", "semver.json", "name of version file")
	cmd.Flags().StringVarP(&GetCmdFlags.versionFileFormat, "outfileFormat", "f", "json", "format of outfile (json, raw)")
	cmd.Flags().BoolVarP(&GetCmdFlags.printRaw, "raw", "r", false, "print only the plain version number")

	viper.BindPFlag("versionFileName", cmd.Flags().Lookup("outfile"))
	viper.BindPFlag("versionFileType", cmd.Flags().Lookup("outfileFormat"))
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

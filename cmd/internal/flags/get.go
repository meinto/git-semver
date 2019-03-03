package flags

import "github.com/spf13/cobra"

type getCmdFlagsType struct {
	repoPath string
	printRaw bool
}

var GetCmdFlags getCmdFlagsType

func (fs *getCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&GetCmdFlags.repoPath, "path", "p", ".", "path to git repository")
	cmd.Flags().BoolVarP(&GetCmdFlags.printRaw, "raw", "r", false, "print only the plain version number")
}

func (fs *getCmdFlagsType) RepoPath() string {
	return fs.repoPath
}

func (fs *getCmdFlagsType) PrintRaw() bool {
	return fs.printRaw
}

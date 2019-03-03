package flags

import (
	"github.com/spf13/cobra"
)

type rootCmdFlags struct {
	verbose bool
}

var RootCmdFlags rootCmdFlags

func (fs *rootCmdFlags) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&RootCmdFlags.verbose, "verbose", "v", false, "more logs")
}

func (fs *rootCmdFlags) Verbose() bool {
	return fs.verbose
}

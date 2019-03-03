package flags

import (
	"github.com/spf13/cobra"
)

type rootCmdFlagsType struct {
	verbose bool
}

var RootCmdFlags rootCmdFlagsType

func (fs *rootCmdFlagsType) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&RootCmdFlags.verbose, "verbose", "v", false, "more logs")
}

func (fs *rootCmdFlagsType) Verbose() bool {
	return fs.verbose
}

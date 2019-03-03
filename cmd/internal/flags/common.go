package flags

import (
	"github.com/spf13/cobra"
)

func versionFileFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&VersionCmdFlags.versionFile, "versionFile", "f", "VERSION", "name of version file")
	cmd.Flags().StringVarP(&VersionCmdFlags.versionFileType, "versionFileType", "t", "raw", "type of version file (json, raw)")

	bindViperFlag("versionFile", cmd.Flags().Lookup("versionFile"))
	bindViperFlag("versionFileType", cmd.Flags().Lookup("versionFileType"))
}

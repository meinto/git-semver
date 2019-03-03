package flags

import (
	"github.com/meinto/git-semver/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func bindViperFlag(name string, flag *pflag.Flag) {
	util.LogOnError(viper.BindPFlag(name, flag), RootCmdFlags.Verbose())
}

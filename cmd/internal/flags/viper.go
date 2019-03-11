package flags

import (
	"log"

	"github.com/meinto/git-semver/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func bindViperFlag(name string, flag *pflag.Flag) {
	util.LogOnError(viper.BindPFlag(name, flag), RootCmdFlags.Verbose())
}

func LaodViperConfig(repoPath string) {
	viper.SetConfigName("semver.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(repoPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("there is no semver.config file: ", err)
	}
}

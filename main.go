package main

import (
	"log"

	"github.com/meinto/git-semver/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("semver.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("there is no semver.config file: ", err)
	}

	cmd.Execute()
}

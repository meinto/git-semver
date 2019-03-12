package main

//go:generate ../../.circleci/generate-assets.sh

import (
	"log"

	"github.com/meinto/git-semver/git"
	"github.com/meinto/git-semver/pkg/cli/cmd"
	"github.com/spf13/viper"
)

func main() {
	g := git.NewGitService(viper.GetString("gitPath"))
	repoPath, _ := g.GitRepoPath()

	viper.SetConfigName("semver.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(repoPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("there is no semver.config file: ", err)
	}

	cmd.Execute()
}

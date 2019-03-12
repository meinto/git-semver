package main

//go:generate ./.circleci/generate-assets.sh

import (
	"github.com/meinto/git-semver/cmd"
)

func main() {
	cmd.Execute()
}

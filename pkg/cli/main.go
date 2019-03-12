package main

//go:generate ../../.circleci/generate-assets.sh

import (
	"github.com/meinto/git-semver/pkg/cli/cmd"
)

func main() {
	cmd.Execute()
}

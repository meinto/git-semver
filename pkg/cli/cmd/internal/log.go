package internal

import (
	"log"

	"github.com/meinto/git-semver/cmd/internal/flags"
	"github.com/meinto/git-semver/util"
)

func LogOnError(err error) {
	util.LogFatalOnErr(err, flags.RootCmdFlags.Verbose())
}

func LogFatalOnErr(err error) {
	util.LogFatalOnErr(err, flags.RootCmdFlags.Verbose())
}

func LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

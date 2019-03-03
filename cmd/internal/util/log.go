package util

import (
	"log"

	"github.com/meinto/git-semver/cmd/internal/flags"
)

func pattern() string {
	if flags.RootCmdFlags.Verbose {
		return "%+v\n"
	}
	return "%v\n"
}

func LogOnError(err error) {
	if err != nil {
		log.Printf(pattern(), err)
	}
}

func LogFatalOnErr(err error) {
	if err != nil {
		log.Fatalf(pattern(), err)
	}
}

func LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

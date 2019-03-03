package util

import (
	"log"
)

func pattern(verbose bool) string {
	if verbose {
		return "%+v\n"
	}
	return "%v\n"
}

func LogOnError(err error, verbose bool) {
	if err != nil {
		log.Printf(pattern(verbose), err)
	}
}

func LogFatalOnErr(err error, verbose bool) {
	if err != nil {
		log.Fatalf(pattern(verbose), err)
	}
}

func LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

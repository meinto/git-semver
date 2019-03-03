package util

import (
	"log"
)

func LogFatalOnErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

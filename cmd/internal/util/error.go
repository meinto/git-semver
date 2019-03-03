package util

import (
	"log"
)

func LogOnError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
	}
}

func LogFatalOnErr(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

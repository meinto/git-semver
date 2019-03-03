package util

import (
	"log"
)

func LogFatalOnErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

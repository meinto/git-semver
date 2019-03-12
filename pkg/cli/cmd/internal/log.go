package internal

import "log"

type Logger interface {
	LogFatalOnError(err error)
	LogFatalIfNotOk(ok bool, message string)
}

type logger struct {
	verbose bool
}

func NewLogger(verbose bool) Logger {
	return &logger{verbose}
}

func (l *logger) LogFatalOnError(err error) {
	if err != nil {
		log.Fatalf(pattern(l.verbose), err)
	}
}

func (l *logger) LogFatalIfNotOk(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

func pattern(verbose bool) string {
	if verbose {
		return "%+v\n"
	}
	return "%v\n"
}

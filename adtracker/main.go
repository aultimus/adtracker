package main

import (
	"github.com/aultimus/adtracker"
	"github.com/cocoonlife/timber"
)

func main() {
	timber.AddLogger(timber.ConfigLogger{
		LogWriter: new(timber.ConsoleWriter),
		Level:     timber.DEBUG,
		Formatter: timber.NewPatFormatter("[%D %T] [%L] %s %M"),
	})

	// TODO: make port configurable via flag
	adtracker.Run(5000)
}

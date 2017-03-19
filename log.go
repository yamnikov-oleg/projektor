package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	Verbose bool
)

func init() {
	flag.BoolVar(&Verbose, "v", false, "Verbose output")
}

func logf(format string, a ...interface{}) {
	if Verbose {
		stamp := time.Now().Format(time.StampMilli)
		fmt.Printf(stamp+" "+format, a...)
	}
}

func errduring(action string, err error, nextmove string, a ...interface{}) {
	line := "Error during " + action + ":\n"
	if err != nil {
		line += err.Error() + "\n"
	}
	if nextmove != "" {
		line += nextmove + "\n"
	}
	line += "\n"
	logf(line, a...)
}

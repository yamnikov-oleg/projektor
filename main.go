// +build linux

package main

import (
	"fmt"
	"os"

	"github.com/yamnikov-oleg/go-gtk/gtk"
)

const (
	Verbose = true
)

func logf(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

func errduring(action string, err error, nextmove string, a ...interface{}) {
	line := action + ":\n"
	if err != nil {
		line += err.Error() + "\n"
	}
	if nextmove != "" {
		line += nextmove + "\n"
	}
	logf(line, a...)
}

func main() {
	gtk.Init(&os.Args)

	ConstructWindow()

	gtk.Main()
}

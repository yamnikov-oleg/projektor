// +build linux

package main

import "fmt"

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
	reader := NewEntriesReader()
	for reader.Next() {
		fmt.Printf("%#v\n", reader.Entry)
	}
}

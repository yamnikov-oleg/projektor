// +build linux

package main

import "os"

var (
	HOME = os.Getenv("HOME")
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		Verbose = true
		logf("Verbose mode on\n")
	}
	KillIfRunning()
	IndexDesktopEntries()
	IndexAvailableCommands()
	SetupUi()
}

// +build linux

package main

import "os"

var (
	HOME = os.Getenv("HOME")
)

func main() {
	KillIfRunning()
	IndexDesktopEntries()
	SetupUi()
}

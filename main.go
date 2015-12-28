// +build linux

package main

import "os"

func main() {
	if KillIfRunning() {
		os.Exit(0)
	}
	IndexDesktopEntries()
	StartUi()
}

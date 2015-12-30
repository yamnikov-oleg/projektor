// +build linux

package main

func main() {
	KillIfRunning()
	IndexDesktopEntries()
	StartUi()
}

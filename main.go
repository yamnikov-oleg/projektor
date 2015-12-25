// +build linux

package main

import (
	"os"

	"github.com/yamnikov-oleg/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	ConstructWindow()

	gtk.Main()
}

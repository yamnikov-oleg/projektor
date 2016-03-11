package main

import (
	"os"
	"strings"

	"github.com/yamnikov-oleg/go-ps"
)

var (
	Pid     = os.Getpid()
	Cmdline = os.Args[0]
)

func KillIfOtherInstance(p ps.Process) {
	if p.Pid() == Pid {
		return
	}

	if !strings.Contains(p.Executable(), Cmdline) {
		return
	}

	osproc, err := os.FindProcess(p.Pid())
	if err != nil {
		return
	}

	osproc.Kill()
	os.Exit(0)
}

func KillIfRunning() {
	processes, err := ps.Processes()
	if err != nil {
		errduring("retreiving process list", err, "Skipping process scanning")
		return
	}

	for _, p := range processes {
		KillIfOtherInstance(p)
	}

	return
}

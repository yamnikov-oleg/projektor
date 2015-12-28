package main

import (
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

var (
	Pid     = os.Getpid()
	Cmdline = os.Args[0]
)

func KillIfOtherInstance(p ps.Process) bool {
	if p.Pid() == Pid {
		return false
	}

	if !strings.Contains(p.Executable(), Cmdline) {
		return false
	}

	op, err := os.FindProcess(p.Pid())
	if err != nil {
		return false
	}

	op.Kill()
	return true
}

func KillIfRunning() (ret bool) {
	processes, err := ps.Processes()
	if err != nil {
		errduring("retreiving process list", err, "Skipping process scanning")
		return false
	}

	for _, p := range processes {
		if KillIfOtherInstance(p) {
			ret = true
		}
	}

	return
}

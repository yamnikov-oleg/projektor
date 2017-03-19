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

func KillIfOtherInstance(p ps.Process) bool {
	if p.Pid() == Pid {
		return false
	}

	if !strings.Contains(p.Executable(), Cmdline) {
		return false
	}

	if strings.Contains(p.Cmdline(), "-dry") {
		return false
	}

	osproc, err := os.FindProcess(p.Pid())
	if err != nil {
		return false
	}

	osproc.Kill()
	return true
}

func FindAndKill() (ret bool) {
	processes, err := ps.Processes()
	if err != nil {
		errduring("retreiving process list", err, "Skipping process scanning")
		return
	}

	for _, p := range processes {
		ret = KillIfOtherInstance(p) || ret
	}
	return
}

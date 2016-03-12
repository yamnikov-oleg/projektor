// +build linux

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	HOME = os.Getenv("HOME")

	SIFlag         = "inst"
	SingleInstance bool
)

func init() {
	flag.BoolVar(&SingleInstance, SIFlag, false, "Run an instance of projektor, not a daemon")
}

func RunInstance() {
	logf("Running single instance of projektor\n")

	IndexDesktopEntries()
	IndexAvailableCommands()
	SetupUi()
}

func RunDaemon() {
	logf("Running key binding daemon of projektor\n")

	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(xu)
	cb := func(xu *xgbutil.XUtil, e xevent.KeyPressEvent) {
		if FindAndKill() {
			return
		}
		cmd := exec.Command(os.Args[0], "-"+SIFlag)
		err := cmd.Start()
		if err != nil {
			errduring("instance creation", err, "")
		}
		go cmd.Wait()
	}
	err = keybind.KeyPressFun(cb).Connect(xu, xu.RootWin(), "mod4-q", true)
	if err != nil {
		log.Fatal(err)
	}

	xevent.Main(xu)
}

func main() {
	flag.Parse()
	if SingleInstance {
		RunInstance()
	} else {
		RunDaemon()
	}
}

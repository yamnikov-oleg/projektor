// +build linux

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	HOME   = os.Getenv("HOME")
	AppDir = HOME + "/.projektor"

	SIFlag         = "inst"
	SingleInstance bool

	CPUProfile string
)

func init() {
	flag.BoolVar(&SingleInstance, SIFlag, false, "Run an instance of projektor, not a daemon")
	flag.StringVar(&CPUProfile, "cpuprofile", "", "Run CPU profiling and output results to the `file`")
}

func RunInstance() {
	logf("Running single instance of projektor\n")

	LoadHistory()
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

	if CPUProfile != "" {
		f, err := os.Create(CPUProfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if SingleInstance {
		RunInstance()
	} else {
		RunDaemon()
	}
}

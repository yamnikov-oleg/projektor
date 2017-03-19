// +build linux

package main

import (
	"flag"
	"fmt"
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

	Version     string = "v0.1+"
	ShowVersion bool

	DryRun bool
)

func init() {
	flag.BoolVar(&SingleInstance, SIFlag, false, "Run an instance of projektor, not a daemon")
	flag.StringVar(&CPUProfile, "cpuprofile", "", "Run CPU profiling and output results to the `file`")
	flag.BoolVar(&ShowVersion, "V", false, "Display version of current projektor build.\n\tPlus sign means that the build includes several more commits over the release.")
	flag.BoolVar(&DryRun, "dry", false, "Prepare to run projektor instance but do not run (useful to force kernel caching).")
}

func RunInstance(dry bool) {
	logf("Running single instance of projektor\n")

	logf("Loading history...\n")
	LoadHistory()
	logf("Indexing desktop entries...\n")
	IndexDesktopEntries()
	logf("Set up the UI.\n")
	SetupUi(dry)
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
	err = keybind.KeyPressFun(cb).Connect(xu, xu.RootWin(), Config.KeyBind, true)
	if err != nil {
		log.Fatal(err)
	}

	xevent.Main(xu)
}

func main() {
	flag.Parse()

	if ShowVersion {
		fmt.Printf("Projektor %v\n", Version)
		os.Exit(0)
	}

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
		RunInstance(DryRun)
	} else {
		RunDaemon()
	}
}

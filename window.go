package main

import (
	"fmt"
	"unsafe"

	"github.com/yamnikov-oleg/go-gtk/gdk"
	"github.com/yamnikov-oleg/go-gtk/glib"
	"github.com/yamnikov-oleg/go-gtk/gtk"
)

var Ui struct {
	Window *gtk.Window
}

func ConstructWindow() {
	wnd := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	Ui.Window = wnd

	wnd.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	wnd.SetDecorated(false)
	wnd.SetSkipTaskbarHint(true)
	wnd.SetBorderWidth(6)
	wnd.SetSizeRequest(400, 0)

	wnd.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		e := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		fmt.Printf("%vx%v\n", int(e.X), int(e.Y))
	})
	wnd.Connect("focus-out-event", gtk.MainQuit)
	wnd.Connect("destroy", gtk.MainQuit)

	rootbox := gtk.NewVBox(false, 6)
	wnd.Add(rootbox)

	textedit := gtk.NewEntry()
	textedit.GrabFocus()
	rootbox.PackStart(textedit, false, false, 0)

	wnd.ShowAll()
}

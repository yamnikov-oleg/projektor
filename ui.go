package main

import (
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"github.com/yamnikov-oleg/go-gtk/gdk"
	"github.com/yamnikov-oleg/go-gtk/gio"
	"github.com/yamnikov-oleg/go-gtk/glib"
	"github.com/yamnikov-oleg/go-gtk/gtk"
)

var Ui struct {
	Window      UiWindow
	RootBox     *gtk.VBox
	SearchEntry UiEntry
	ScrollWin   *gtk.ScrolledWindow
	ListStore   *gtk.ListStore
	TreeView    UiTreeView
	Pointer     UiPointer
}

type UiWindow struct {
	*gtk.Window
}

func (UiWindow) OnKeyPress(e *gdk.EventKey) {
	switch e.Keyval {
	case gdk.KEY_Down:
		Ui.TreeView.Selected().Inc().Select()

	case gdk.KEY_Up:
		Ui.TreeView.Selected().Dec().Select()

	case gdk.KEY_Return:
		Ui.TreeView.Selected().Execute()

	case gdk.KEY_Escape:
		gtk.MainQuit()
	}
}

func (UiWindow) OnButtonPress(e *gdk.EventButton) {
	var wid, hei, clickX, clickY int
	Ui.Window.GetSize(&wid, &hei)
	clickX, clickY = int(e.X), int(e.Y)

	if clickX < 0 || clickX > wid || clickY < 0 || clickY > hei {
		gtk.MainQuit()
	}
}

func (UiWindow) OnFocusIn() {
	Ui.Pointer.Grab()
}

type UiEntry struct {
	*gtk.Entry
}

func (UiEntry) OnChanged() {
	UpdateSearchResults()
}

type UiTreeIter struct {
	*gtk.TreeIter
}

var NilTreeIter = UiTreeIter{nil}

func NewTreeIter() UiTreeIter {
	return UiTreeIter{new(gtk.TreeIter)}
}

func (iter UiTreeIter) None() bool {
	return iter.TreeIter == nil
}

func (iter UiTreeIter) Execute() {
	if iter.None() {
		return
	}
	var val glib.GValue
	Ui.ListStore.GetValue(iter.TreeIter, 2, &val)
	cmd := strings.Fields(val.GetString())
	exec.Command(cmd[0], cmd[1:]...).Start()
	gtk.MainQuit()
}

func (iter UiTreeIter) Select() {
	if iter.None() {
		return
	}
	Ui.TreeView.GetSelection().SelectIter(iter.TreeIter)
	Ui.TreeView.ScrollToCell(Ui.ListStore.GetPath(iter.TreeIter), nil, false, 0, 0)
}

func (iter UiTreeIter) Inc() UiTreeIter {
	if iter.None() || !Ui.ListStore.IterNext(iter.TreeIter) {
		return NilTreeIter
	}
	return iter
}

func (iter UiTreeIter) Dec() UiTreeIter {
	if iter.None() || !Ui.ListStore.IterPrev(iter.TreeIter) {
		return NilTreeIter
	}
	return iter
}

type UiTreeView struct {
	*gtk.TreeView
}

func (UiTreeView) OnRowActivated() {
	Ui.TreeView.Selected().Execute()
}

func (UiTreeView) Selected() UiTreeIter {
	selection := Ui.TreeView.GetSelection()
	if selection.CountSelectedRows() == 0 {
		return NilTreeIter
	}

	iter := NewTreeIter()
	selection.GetSelected(iter.TreeIter)
	return iter
}

func (UiTreeView) First() UiTreeIter {
	iter := NewTreeIter()
	if !Ui.ListStore.GetIterFirst(iter.TreeIter) {
		return NilTreeIter
	}
	return iter
}

func (UiTreeView) Clear() {
	Ui.ListStore.Clear()
}

func (UiTreeView) AppendLaunchEntry(entry *LaunchEntry) {
	iter := NewTreeIter()
	Ui.ListStore.Append(iter.TreeIter)

	gicon, err := gio.NewIconForString(entry.Icon)
	if err != nil {
		errduring("appending entry to ListStore", err, "Skipping it")
		return
	}

	Ui.ListStore.Set(iter.TreeIter,
		0, gicon.GIcon,
		1, entry.MarkupName,
		2, entry.Cmdline,
	)
}

type UiPointer struct {
	*gdk.Device
}

func (UiPointer) Grab() {
	status := Ui.Pointer.Device.Grab(Ui.Window.GetWindow(), gdk.OWNERSHIP_APPLICATION, true, gdk.BUTTON_PRESS_MASK, nil, gdk.CURRENT_TIME)
	if status != gdk.GRAB_SUCCESS {
		errduring("pointer grabbing, grab status %v", nil, "", status)
	}
}

func UpdateSearchResults() {
	Ui.TreeView.Clear()
	text := strings.TrimSpace(Ui.SearchEntry.GetText())

	results := SearchAppEntries(text)
	for _, entry := range results {
		Ui.TreeView.AppendLaunchEntry(entry)
	}
	Ui.TreeView.First().Select()
}

func SetupUi() {
	gtk.Init(&os.Args)

	//
	// Constructors
	//
	Ui.Window = UiWindow{gtk.NewWindow(gtk.WINDOW_TOPLEVEL)}
	Ui.RootBox = gtk.NewVBox(false, 6)
	Ui.SearchEntry = UiEntry{gtk.NewSearchEntry()}
	Ui.ScrollWin = gtk.NewScrolledWindow(nil, nil)
	Ui.TreeView = UiTreeView{gtk.NewTreeView()}
	Ui.ListStore = gtk.NewListStore(gio.GetIconType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	Ui.Pointer = UiPointer{gdk.GetDefaultDisplay().GetDeviceManager().GetClientPointer()}

	//
	// Window
	//
	Ui.Window.SetPosition(gtk.WIN_POS_CENTER)
	Ui.Window.SetGravity(gdk.GRAVITY_SOUTH)
	Ui.Window.SetDecorated(false)
	Ui.Window.SetSkipTaskbarHint(true)
	Ui.Window.SetBorderWidth(6)
	Ui.Window.SetSizeRequest(400, 480)
	Ui.Window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		e := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		Ui.Window.OnKeyPress(e)
	})
	Ui.Window.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		e := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		Ui.Window.OnButtonPress(e)
	})
	Ui.Window.Connect("focus-in-event", Ui.Window.OnFocusIn)
	Ui.Window.Connect("destroy", gtk.MainQuit)

	//
	// SearchEntry
	//
	Ui.SearchEntry.Connect("changed", Ui.SearchEntry.OnChanged)

	//
	// TreeView
	//
	Ui.TreeView.SetHeadersVisible(false)
	cr := gtk.NewCellRendererPixbuf()
	glib.ObjectFromNative(unsafe.Pointer(cr.ToCellRenderer())).Set("stock-size", int(gtk.ICON_SIZE_DIALOG))
	Ui.TreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes2("Icon", cr, "gicon", 0))
	Ui.TreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Id", gtk.NewCellRendererText(), "markup", 1))
	Ui.TreeView.SetModel(Ui.ListStore)
	Ui.TreeView.Connect("row-activated", Ui.TreeView.OnRowActivated)

	//
	// Focus setup
	//
	Ui.RootBox.SetCanFocus(false)
	Ui.ScrollWin.SetCanFocus(false)
	Ui.TreeView.SetCanFocus(false)
	Ui.SearchEntry.GrabFocus()

	//
	// Packing
	//
	Ui.ScrollWin.Add(Ui.TreeView.TreeView)
	Ui.RootBox.PackStart(Ui.SearchEntry.Entry, false, false, 0)
	Ui.RootBox.PackEnd(Ui.ScrollWin, true, true, 0)
	Ui.Window.Add(Ui.RootBox)

	//
	// Stylesheet loading
	//
	provider := gtk.NewCssProvider()
	screen := gdk.GetDefaultDisplay().GetDefaultScreen()
	gtk.StyleContextAddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	err := provider.LoadFromData(CSS_CODE)
	if err != nil {
		errduring("CSS loading", err, "")
	}

	UpdateSearchResults()
	Ui.Window.ShowAll()
	gtk.Main()
}

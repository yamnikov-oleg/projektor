package main

import (
	"fmt"
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
	Window      *gtk.Window
	RootBox     *gtk.VBox
	SearchEntry *gtk.Entry
	ListStore   *gtk.ListStore
	ScrollWin   *gtk.ScrolledWindow
	TreeView    *gtk.TreeView
	Pointer     *gdk.Device
}

func setupSearchEntry() {
	Ui.SearchEntry = gtk.NewSearchEntry()
	Ui.SearchEntry.GrabFocus()
}

func setupAppList() {
	Ui.TreeView = gtk.NewTreeView()
	Ui.TreeView.SetCanFocus(false)
	Ui.TreeView.SetHeadersVisible(false)

	cr := gtk.NewCellRendererPixbuf()
	glib.ObjectFromNative(unsafe.Pointer(cr.ToCellRenderer())).Set("stock-size", int(gtk.ICON_SIZE_DIALOG))
	Ui.TreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes2("Icon", cr, "gicon", 0))
	Ui.TreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Id", gtk.NewCellRendererText(), "markup", 1))

	Ui.ListStore = gtk.NewListStore(gio.GetIconType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	Ui.TreeView.SetModel(Ui.ListStore)

	Ui.ScrollWin = gtk.NewScrolledWindow(nil, nil)
	Ui.ScrollWin.SetCanFocus(false)
	Ui.ScrollWin.Add(Ui.TreeView)

	Ui.TreeView.Connect("row-activated", OnTreeViewRowActivate)
}

func makeSearching() {
	Ui.ListStore.Clear()
	text := Ui.SearchEntry.GetText()
	text = strings.TrimSpace(text)
	loText := strings.ToLower(text)

	results := SearchDesktopEntries(text)
	for _, entry := range results {
		listStoreAppendEntry(entry, loText)
	}
	treeViewSelectFirst()
}

func setupSearchLogic() {
	Ui.SearchEntry.Connect("changed", makeSearching)
}

func setupUiElements() {
	Ui.RootBox = gtk.NewVBox(false, 6)
	Ui.RootBox.SetCanFocus(false)

	setupSearchEntry()
	Ui.RootBox.PackStart(Ui.SearchEntry, false, false, 0)

	setupAppList()
	Ui.RootBox.PackEnd(Ui.ScrollWin, true, true, 0)

	setupSearchLogic()
}

func setupPointerDevice() {
	Ui.Pointer = gdk.GetDefaultDisplay().GetDeviceManager().GetClientPointer()
}

func setupWindow() {
	Ui.Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	Ui.Window.SetPosition(gtk.WIN_POS_CENTER)
	Ui.Window.SetGravity(gdk.GRAVITY_SOUTH)
	Ui.Window.SetDecorated(false)
	Ui.Window.SetSkipTaskbarHint(true)
	Ui.Window.SetBorderWidth(6)
	Ui.Window.SetSizeRequest(400, 480)

	Ui.Window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		e := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		OnWindowKeyPress(e)
	})
	Ui.Window.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		e := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		OnWindowButtonPress(e)
	})
	Ui.Window.Connect("destroy", gtk.MainQuit)
	Ui.Window.Connect("focus-in-event", func() {
		pointerGrab()
	})

	setupUiElements()
	Ui.Window.Add(Ui.RootBox)

	setupPointerDevice()
}

func loadCss() {
	provider := gtk.NewCssProvider()
	screen := gdk.GetDefaultDisplay().GetDefaultScreen()
	gtk.StyleContextAddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	err := provider.LoadFromData(CSS_CODE)
	fmt.Printf("%#v\n", err)
}

func StartUi() {
	gtk.Init(&os.Args)

	setupWindow()
	loadCss()
	makeSearching()
	Ui.Window.ShowAll()

	gtk.Main()
}

func escapeAmp(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
}

func entryDisplayName(entry *DtEntry, query string) string {
	if query == "" {
		return escapeAmp(entry.Name)
	}
	ind := strings.Index(entry.LoCaseName, query)
	if ind < 0 {
		return escapeAmp(entry.Name)
	}
	return escapeAmp(fmt.Sprintf("%v<b>%v</b>%v", entry.Name[:ind], entry.Name[ind:ind+len(query)], entry.Name[ind+len(query):]))
}

func listStoreAppendEntry(entry *DtEntry, searchQuery string) {
	var iter gtk.TreeIter
	Ui.ListStore.Append(&iter)

	gicon, err := gio.NewIconForString(entry.Icon)
	if err != nil {
		errduring("appending entry to ListStore", err, "Skipping it")
		return
	}
	Ui.ListStore.Set(&iter,
		0, gicon.GIcon,
		1, entryDisplayName(entry, searchQuery),
		2, entry.Exec,
	)
}

func treeViewSelect(iter *gtk.TreeIter) {
	Ui.TreeView.GetSelection().SelectIter(iter)
	Ui.TreeView.ScrollToCell(Ui.ListStore.GetPath(iter), nil, false, 0, 0)
}

func treeViewSelectFirst() {
	var iter gtk.TreeIter
	if !Ui.ListStore.GetIterFirst(&iter) {
		return
	}
	treeViewSelect(&iter)
}

func pointerGrab() {
	status := Ui.Pointer.Grab(Ui.Window.GetWindow(), gdk.OWNERSHIP_APPLICATION, true, gdk.BUTTON_PRESS_MASK, nil, gdk.CURRENT_TIME)
	if status != gdk.GRAB_SUCCESS {
		errduring("pointer grabbing, grab status %v", nil, "", status)
	}
}

func runSelectedApp() {
	selection := Ui.TreeView.GetSelection()
	if selection.CountSelectedRows() == 0 {
		return
	}

	var iter gtk.TreeIter
	selection.GetSelected(&iter)

	var val glib.GValue
	Ui.ListStore.GetValue(&iter, 2, &val)
	cmd := strings.Fields(val.GetString())
	exec.Command(cmd[0], cmd[1:]...).Start()
	gtk.MainQuit()
}

func OnWindowKeyPress(e *gdk.EventKey) {
	switch e.Keyval {

	case gdk.KEY_Escape:
		gtk.MainQuit()

	case gdk.KEY_Up, gdk.KEY_Down:
		selection := Ui.TreeView.GetSelection()
		if selection.CountSelectedRows() == 0 {
			return
		}
		var iter gtk.TreeIter
		selection.GetSelected(&iter)
		if e.Keyval == gdk.KEY_Up {
			if !Ui.ListStore.IterPrev(&iter) {
				return
			}
		} else {
			if !Ui.ListStore.IterNext(&iter) {
				return
			}
		}
		treeViewSelect(&iter)

	case gdk.KEY_Return:
		runSelectedApp()
	}
}

func OnWindowButtonPress(e *gdk.EventButton) {
	var wid, hei int
	Ui.Window.GetSize(&wid, &hei)

	clickX := int(e.X)
	clickY := int(e.Y)

	if clickX < 0 || clickX > wid || clickY < 0 || clickY > hei {
		gtk.MainQuit()
	}
}

func OnTreeViewRowActivate() {
	runSelectedApp()
}

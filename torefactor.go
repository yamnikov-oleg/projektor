package main

import (
	"fmt"
	"strings"

	"github.com/yamnikov-oleg/go-gtk/gio"
	"github.com/yamnikov-oleg/go-gtk/gtk"
)

func makeSearching() {
	Ui.ListStore.Clear()
	text := Ui.SearchEntry.GetText()
	text = strings.TrimSpace(text)
	loText := strings.ToLower(text)

	results := SearchDesktopEntries(text)
	for _, entry := range results {
		listStoreAppendEntry(entry, loText)
	}
	Ui.TreeView.First().Select()
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

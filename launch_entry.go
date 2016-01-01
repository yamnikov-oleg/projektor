package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yamnikov-oleg/go-gtk/gio"
	"github.com/yamnikov-oleg/projektor/conf"
)

type LaunchEntry struct {
	// Clean name for an entry. E.g. "Atom Text Editor"
	Name string
	// Same as Name, but lowercased, e.g. "atom text editor"
	LoCaseName string
	// Formatted for display on a gtk widget, e.g. "<b>Ato</b>m Text Editor"
	MarkupName string
	// Name which is injected into search entry on Tab hit
	TabName string

	Icon string

	Cmdline string
	// Describes priority of an entry in results list. Lower index -> higher priority.
	QueryIndex int
}

func NewEntryFromDesktopFile(filepath string) (le *LaunchEntry, err error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer fd.Close()

	cf, err := conf.Read(fd)
	if err != nil {
		return
	}

	section := cf.Sections["Desktop Entry"]
	if section.Bool("Hidden") {
		return nil, errors.New("desktop entry hidden")
	}
	if section.Bool("NoDisplay") {
		return nil, errors.New("desktop entry not displayed")
	}
	if section.Has("OnlyShowIn") && !strings.Contains(section.Str("OnlyShowIn"), CurrentDesktop) {
		return nil, errors.New("desktop entry is hidden on current desktop")
	}
	if section.Has("NotShowIn") && strings.Contains(section.Str("NotShowIn"), CurrentDesktop) {
		return nil, errors.New("desktop entry is hidden on current desktop")
	}

	le = &LaunchEntry{}
	le.Name = section.Str("Name")
	le.LoCaseName = strings.ToLower(le.Name)
	le.Icon = section.Str("Icon")

	r := strings.NewReplacer(" %f", "", " %F", "", " %u", "", " %U", "")
	le.Cmdline = r.Replace(section.Str("Exec"))
	le.TabName = le.Cmdline

	return
}

func NewEntryFromCommand(command string) *LaunchEntry {
	return &LaunchEntry{
		Icon:       "application-default-icon",
		Name:       command,
		MarkupName: "<b>" + command + "</b>",
		TabName:    command,
		Cmdline:    command,
	}
}

func NewEntryForFile(path string, displayName string, tabName string) (*LaunchEntry, error) {
	gFileInfo, fiErr := gio.NewFileForPath(path).QueryInfo("standard::*", gio.FILE_QUERY_INFO_NONE, nil)
	if fiErr != nil {
		return nil, fiErr
	}

	icon := gFileInfo.GetIcon()
	return &LaunchEntry{
		Icon:       icon.ToString(),
		MarkupName: displayName,
		TabName:    tabName,
		Cmdline:    "xdg-open " + path,
	}, nil
}

func (le *LaunchEntry) UpdateMarkupName(index, length int) {
	index2 := index + length
	le.MarkupName = EscapeAmpersand(
		fmt.Sprintf("%v<b>%v</b>%v",
			le.Name[:index],
			le.Name[index:index2],
			le.Name[index2:],
		),
	)
}

type LaunchEntriesList []*LaunchEntry

func (list LaunchEntriesList) SortByName() {
	sort.Sort(lelSortableByName(list))
}

func (list LaunchEntriesList) SortByIndex() {
	sort.Sort(lelSortableByIndex(list))
}

type lelSortableByName LaunchEntriesList

func (list lelSortableByName) Len() int {
	return len(list)
}

func (list lelSortableByName) Less(i, j int) bool {
	return list[i].LoCaseName < list[j].LoCaseName
}

func (list lelSortableByName) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type lelSortableByIndex LaunchEntriesList

func (list lelSortableByIndex) Len() int {
	return len(list)
}

func (list lelSortableByIndex) Less(i, j int) bool {
	if list[i].QueryIndex == list[j].QueryIndex {
		return list[i].LoCaseName < list[j].LoCaseName
	}
	return list[i].QueryIndex < list[j].QueryIndex
}

func (list lelSortableByIndex) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type EntrySearchFunc func(string) LaunchEntriesList

func SearchEntries(query string, fns []EntrySearchFunc) (list LaunchEntriesList) {
	for _, fn := range fns {
		list = append(list, fn(query)...)
	}
	return
}

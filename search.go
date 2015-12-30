package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yamnikov-oleg/go-gtk/gio"
)

type SearchPair struct {
	Index int
	Entry *DtEntry
}

type SearchPairList []SearchPair

func (list SearchPairList) Len() int {
	return len(list)
}

func (list SearchPairList) Less(i, j int) bool {
	return list[i].Index < list[j].Index
}

func (list SearchPairList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func ExecCommandEntry(command string) *DtEntry {
	return &DtEntry{
		Icon: "application-default-icon",
		Name: fmt.Sprintf("\u2192 <b>%v</b>", command),
		Exec: command,
	}
}

func SpecialEntry(query string) *DtEntry {
	if query == "" {
		return nil
	}

	if query[0] != '/' && query[0] != '~' {
		return ExecCommandEntry(query)
	}

	if query[0] == '~' {
		query = HOME + query[1:]
	}

	stat, statErr := os.Stat(query)
	if statErr != nil {
		return nil
	}

	if !stat.IsDir() && (stat.Mode().Perm()&0111) != 0 {
		return ExecCommandEntry(query)
	}

	gFileInfo, fiErr := gio.NewFileForPath(query).QueryInfo("standard::*", gio.FILE_QUERY_INFO_NONE, nil)
	if fiErr != nil {
		return nil
	}

	icon := gFileInfo.GetIcon()
	return &DtEntry{
		Icon: icon.ToString(),
		Name: fmt.Sprintf("<b>%v</b>", query),
		Exec: "xdg-open " + query,
	}
}

func SearchDesktopEntries(query string) (entries []*DtEntry) {
	// if query == "" {
	// 	return
	// }

	var pairs SearchPairList
	loQuery := strings.ToLower(query)

	reader := NewEntriesInterator()
	for reader.Next() {
		entry := reader.Entry()
		index := strings.Index(entry.LoCaseName, loQuery)
		if index != -1 {
			pairs = append(pairs, SearchPair{index, entry})
		}
	}

	sort.Sort(pairs)

	entries = make([]*DtEntry, len(pairs))
	for i, p := range pairs {
		entries[i] = p.Entry
	}

	specEntry := SpecialEntry(query)
	if specEntry != nil {
		entries = append(entries, specEntry)
	}

	return
}

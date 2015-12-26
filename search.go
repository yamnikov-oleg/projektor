package main

import (
	"sort"
	"strings"
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

func SearchDesktopEntries(query string) (entries []*DtEntry) {
	if query == "" {
		return
	}

	var pairs SearchPairList

	reader := NewEntriesInterator()
	for reader.Next() {
		entry := reader.Entry()
		index := strings.Index(entry.LoCaseName, query)
		if index != -1 {
			pairs = append(pairs, SearchPair{index, entry})
		}
	}

	sort.Sort(pairs)

	entries = make([]*DtEntry, len(pairs))
	for i, p := range pairs {
		entries[i] = p.Entry
	}

	return
}

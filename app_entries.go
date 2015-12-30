package main

import (
	"io/ioutil"
	"os"
	"strings"
)

const (
	MaxAppEntrySize = 1 * 1024 * 1024 // 1Mb
)

var (
	SharedAppDir   = "/usr/share/applications"
	LocalAppDir    = HOME + "/.local/share/applications"
	CurrentDesktop = os.Getenv("XDG_CURRENT_DESKTOP")
)

type AppEntriesParser struct {
	files        []string
	currentIndex int
	Entry        *LaunchEntry
}

func NewAppEntriesParser() *AppEntriesParser {
	return &AppEntriesParser{nil, -1, nil}
}

func (parser *AppEntriesParser) AppendDirectory(dir string) error {
	efs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range efs {
		parser.files = append(parser.files, dir+"/"+fi.Name())
	}
	return nil
}

func (parser *AppEntriesParser) Good() bool {
	return parser.currentIndex >= 0 && parser.currentIndex < len(parser.files)
}

func (parser *AppEntriesParser) Next() bool {
	parser.currentIndex++
	if !parser.Good() {
		return false
	}

	filepath := parser.files[parser.currentIndex]
	file, err := os.Stat(filepath)
	if err != nil {
		errduring("reading entry file `%v`", err, "Skipping it", filepath)
		return parser.Next()
	}

	if file.IsDir() {
		return parser.Next()
	}
	if file.Size() > MaxAppEntrySize {
		errduring("reading .desktop file `%v`: size to big!", nil, "Skipping it", filepath)
		return parser.Next()
	}
	if !strings.HasSuffix(file.Name(), ".desktop") {
		return parser.Next()
	}

	le, err := NewEntryFromDesktopFile(filepath)
	if err != nil {
		return parser.Next()
	}
	parser.Entry = le

	return true
}

var ApplicationEntries LaunchEntriesList

func IndexDesktopEntries() {
	parser := NewAppEntriesParser()

	if err := parser.AppendDirectory(SharedAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", SharedAppDir)
	}
	if err := parser.AppendDirectory(LocalAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", LocalAppDir)
	}

	for parser.Next() {
		ApplicationEntries = append(ApplicationEntries, parser.Entry)
	}

	ApplicationEntries.SortByName()
}

type EntriesIterator struct {
	Index int
	List  LaunchEntriesList
}

func NewEntriesInterator(list LaunchEntriesList) *EntriesIterator {
	return &EntriesIterator{-1, list}
}

func (ei *EntriesIterator) Good() bool {
	return ei.Index >= 0 && ei.Index < len(ei.List)
}

func (ei *EntriesIterator) Next() bool {
	ei.Index++
	return ei.Good()
}

func (ei *EntriesIterator) Entry() *LaunchEntry {
	if !ei.Good() {
		return nil
	}
	return ei.List[ei.Index]
}

func SearchAppEntries(query string) LaunchEntriesList {
	loQuery := strings.ToLower(query)
	results := LaunchEntriesList{}

	iterator := NewEntriesInterator(ApplicationEntries)
	for iterator.Next() {
		entry := iterator.Entry()
		index := strings.Index(entry.LoCaseName, loQuery)
		if index != -1 {
			entry.QueryIndex = index
			entry.UpdateMarkupName(index, len(query))
			results = append(results, entry)
		}
	}

	results.SortByIndex()

	specEntry := SpecialEntry(query)
	if specEntry != nil {
		results = append(results, specEntry)
	}

	return results
}

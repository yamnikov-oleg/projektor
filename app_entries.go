package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
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

func collectDesktopEntries() (entries LaunchEntriesList) {
	parser := NewAppEntriesParser()

	if err := parser.AppendDirectory(SharedAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", SharedAppDir)
	}
	if err := parser.AppendDirectory(LocalAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", LocalAppDir)
	}

	for parser.Next() {
		entries = append(entries, parser.Entry)
	}

	entries.SortByName()

	return
}

var CacheFilePath = path.Join(AppDir, "desktop-launchers.cache")

func getCachedDesktopEntries() (LaunchEntriesList, error) {
	fd, err := os.Open(CacheFilePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	var entries LaunchEntriesList
	if err := yaml.Unmarshal(contents, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func cacheDesktopEntries(entries LaunchEntriesList) error {
	data, err := yaml.Marshal(entries)
	if err != nil {
		return err
	}

	fd, err := os.Create(CacheFilePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	if _, err := fd.Write(data); err != nil {
		return err
	}

	return nil
}

var ApplicationEntries LaunchEntriesList

func IndexDesktopEntries() {
	updateCache := func() LaunchEntriesList {
		updatedEntries := collectDesktopEntries()
		err := cacheDesktopEntries(updatedEntries)
		if err != nil {
			errduring("saving desktop entries cache", err, "Skipping caching procedure.")
		}
		logf("Updated cache with %v entries.\n", len(updatedEntries))
		return updatedEntries
	}

	logf("Performing desktop entries cache indexing...\n")

	entries, err := getCachedDesktopEntries()
	if err == nil {
		logf("Found cache of %v entries.\n", len(entries))
		logf("Running cache update in background...\n")
		go updateCache()
		ApplicationEntries = entries
	} else {
		logf("No cache found, collecting desktop entries in foreground...\n")
		ApplicationEntries = updateCache()
	}
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
	if query == "" {
		return nil
	}

	loQuery := strings.ToLower(query)
	results := LaunchEntriesList{}

	iterator := NewEntriesInterator(ApplicationEntries)
	for iterator.Next() {
		entry := iterator.Entry()
		if IsInHistory(entry.Cmdline) {
			continue
		}
		index := strings.Index(entry.LoCaseName, loQuery)
		if index != -1 {
			entry.QueryIndex = index
			entry.UpdateMarkupName(index, len(query))
			results = append(results, entry)
		}
	}

	results.SortByIndex()
	return results
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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
	return results
}

func ExpandQueryPath(query string) (isPath bool, path string) {
	if query == "" {
		return false, query
	}
	if query[0] != '/' && query[0] != '~' {
		return false, query
	}

	isPath = true
	path = query

	if path[0] == '~' {
		path = HOME + path[1:]
	}
	return
}

func SearchFileEntries(query string) (results LaunchEntriesList) {
	isPath, queryPath := ExpandQueryPath(query)
	if !isPath {
		return nil
	}

	stat, statErr := os.Stat(queryPath)
	if statErr == nil && (stat.IsDir() || (stat.Mode().Perm()&0111) == 0) {
		entry, err := NewEntryForFile(queryPath, "<b>"+query+"</b>", query)
		if err != nil {
			errduring("making file entry `%v`", err, "Skipping it", queryPath)
		} else {
			results = append(results, entry)
		}
	}

	dirPath := queryPath
	queryFileName := ""
	if statErr == nil && stat.IsDir() && !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	} else if lastSlashInd := strings.LastIndex(queryPath, "/"); lastSlashInd >= 0 {
		dirPath = queryPath[:lastSlashInd+1]
		queryFileName = queryPath[lastSlashInd+1:]
	}

	displayDirPath := query
	if statErr == nil && stat.IsDir() && !strings.HasSuffix(displayDirPath, "/") {
		displayDirPath += "/"
	} else if lastSlashInd := strings.LastIndex(query, "/"); lastSlashInd >= 0 {
		displayDirPath = query[:lastSlashInd+1]
	}

	dir, err := os.Open(dirPath)
	if err != nil {
		errduring("opening dir `%v`", err, "Skipping it", dirPath)
		return
	}
	dirStat, err := dir.Stat()
	if err != nil || !dirStat.IsDir() {
		errduring("retrieving dir stat `%v`", err, "Skipping it", dirPath)
		return
	}
	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		errduring("retrieving dirnames `%v`", err, "Skipping it", dirPath)
	}

	sort.Strings(filenames)

	queryFnLen := len(queryFileName)
	for _, name := range filenames {
		if !strings.HasPrefix(name, queryFileName) {
			continue
		}

		filePath := dirPath + name
		if filePath == queryPath {
			continue
		}

		tabFilePath := displayDirPath + name + "/"
		displayFilePath := fmt.Sprintf(".../<b>%v</b>%v", name[0:queryFnLen], name[queryFnLen:])

		entry, err := NewEntryForFile(filePath, displayFilePath, tabFilePath)
		if err != nil {
			errduring("file entry addition `%v`", err, "Skipping it", filePath)
			continue
		}
		results = append(results, entry)
	}

	return
}

func SearchCmdEntries(query string) LaunchEntriesList {
	if query == "" {
		return nil
	}

	isPath, path := ExpandQueryPath(query)
	if !isPath {
		return LaunchEntriesList{NewEntryFromCommand(query)}
	}

	stat, statErr := os.Stat(path)
	// not exists OR is a directory OR is not executable
	if statErr != nil || stat.IsDir() || (stat.Mode().Perm()&0111) == 0 {
		return nil
	}

	return LaunchEntriesList{NewEntryFromCommand(query)}
}

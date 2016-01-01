package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var AvailableCommands []string

func GetAllExecutablesFromDir(dir string) (execs []string) {
	fd, err := os.Open(dir)
	if err != nil {
		errduring("path directory `%v` open", err, "Skipping it", dir)
		return
	}
	defer fd.Close()

	filenames, err := fd.Readdirnames(-1)
	if err != nil {
		errduring("path directory `%v` filenames reading", err, "Skipping it", dir)
		return
	}

	for _, fname := range filenames {
		if stat, err := os.Stat(dir + "/" + fname); err == nil && IsExecutable(stat) {
			execs = append(execs, fname)
		}
	}
	return
}

func IndexAvailableCommands() {
	paths := strings.Split(os.Getenv("PATH"), fmt.Sprintf("%c", os.PathListSeparator))
	for _, pathDir := range paths {
		AvailableCommands = append(AvailableCommands, GetAllExecutablesFromDir(pathDir)...)
	}
	sort.Strings(AvailableCommands)
}

func SearchCmdEntries(query string) (list LaunchEntriesList) {
	if query == "" {
		return nil
	}

	isPath, path := ExpandPathString(query)
	queryLength := len(query)
	if !isPath {
		firstEntry := NewEntryFromCommand(query)
		firstEntry.QueryIndex = -1
		list = append(list, firstEntry)

		for _, cmd := range AvailableCommands {
			if cmd == "" {
				continue
			}
			if cmd == query {
				continue
			}

			ind := strings.Index(cmd, query)
			if ind < 0 {
				continue
			}

			entry := NewEntryFromCommand(cmd)
			entry.QueryIndex = ind
			entry.UpdateMarkupName(ind, queryLength)
			list = append(list, entry)
		}
		list.SortByIndex()
		return
		//return LaunchEntriesList{NewEntryFromCommand(query)}
	}

	stat, statErr := os.Stat(path)
	// not exists OR is a directory OR is not executable
	if statErr != nil || !IsExecutable(stat) {
		return nil
	}

	return LaunchEntriesList{NewEntryFromCommand(query)}
}

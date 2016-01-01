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
	query = ExpandEnvVars(query)

	if query == "" {
		return nil
	}

	queryCmd := query
	if ind := strings.Index(query, " "); ind > 0 {
		queryCmd = query[:ind]
	}
	isPath, path := ExpandPathString(queryCmd)

	if isPath {
		if stat, err := os.Stat(path); err != nil || !IsExecutable(stat) {
			return nil
		}
	}

	list = append(list, NewEntryFromCommand(query))
	list[0].QueryIndex = -1

	if isPath {
		// Nothing to do else
		// command lookup below
		return
	}

	qlen := len(query)
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
		entry.UpdateMarkupName(ind, qlen)
		list = append(list, entry)
	}
	list.SortByIndex()

	return
}

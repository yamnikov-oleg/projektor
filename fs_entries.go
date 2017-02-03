package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type PathQuery struct {
	QueryPath     string
	DirectoryPath string
	Filename      string
}

func NewPathQuery(q string) (isPath bool, pq *PathQuery) {
	pq = &PathQuery{}

	isPath, pq.QueryPath = ExpandPathString(q)
	if !isPath {
		return
	}

	pq.DirectoryPath = pq.QueryPath

	stat, err := os.Stat(pq.QueryPath)

	if err == nil && stat.IsDir() && !strings.HasSuffix(pq.DirectoryPath, "/") {
		pq.DirectoryPath += "/"
	} else {
		if ind := strings.LastIndex(pq.QueryPath, "/"); ind >= 0 {
			pq.DirectoryPath = pq.QueryPath[:ind+1]
			pq.Filename = pq.QueryPath[ind+1:]
		}
	}

	return
}

func (pq *PathQuery) MakeLaunchEntry() (*LaunchEntry, error) {
	stat, err := os.Stat(pq.QueryPath)
	if err != nil {
		return nil, err
	}
	if IsExecutable(stat) {
		return nil, fmt.Errorf("`%v` is executable", pq.QueryPath)
	}

	entry, err := NewEntryForFile(pq.QueryPath, "<b>"+pq.QueryPath+"</b>", pq.DirectoryPath)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (pq *PathQuery) DirFilenames() ([]string, error) {
	dir, err := os.Open(pq.DirectoryPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	stat, err := dir.Stat()
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("`%v` is not a directory", pq.DirectoryPath)
	}

	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	return filenames, nil
}

func SearchFileEntries(query string) (results LaunchEntriesList) {
	query = ExpandEnvVars(query)
	isPath, pq := NewPathQuery(query)
	if !isPath {
		return
	}

	if entry, err := pq.MakeLaunchEntry(); err != nil {
		errduring("making file entry `%v`", err, "Skipping it", pq.QueryPath)
	} else if !IsInHistory(entry.Cmdline) {
		entry.QueryIndex = -1
		results = append(results, entry)
	}

	filenames, err := pq.DirFilenames()
	if err != nil {
		errduring("retrieving dir `%v` filenames", err, "No file entries are retrieved", pq.DirectoryPath)
		return
	}
	sort.Strings(filenames)

	qflen := len(pq.Filename)
	pqLoaseFilename := strings.ToLower(pq.Filename)
	for _, name := range filenames {
		lcname := strings.ToLower(name)
		if !strings.HasPrefix(lcname, pqLoaseFilename) {
			continue
		}
		if name == pq.Filename {
			continue
		}

		path := pq.DirectoryPath + name
		isDir := false
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			isDir = true
		}
		tabPath := pq.DirectoryPath + name
		if isDir {
			tabPath += "/"
		}
		displayPath := fmt.Sprintf(".../<b>%v</b>%v", name[0:qflen], name[qflen:])

		entry, err := NewEntryForFile(path, displayPath, tabPath)
		if err != nil {
			errduring("file entry addition `%v`", err, "Skipping it", path)
			continue
		}
		if !isDir {
			entry.QueryIndex = 1
		}
		results = append(results, entry)
	}

	results.SortByIndex()
	return
}

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type PathQuery struct {
	OriginalQuery      string
	QueryPath          string
	DirectorySubstring string
	DirectoryPath      string
	Filename           string
}

func NewPathQuery(q string) (isPath bool, pq *PathQuery) {
	pq = &PathQuery{OriginalQuery: q}

	isPath, pq.QueryPath = ExpandPathString(q)
	if !isPath {
		return
	}

	pq.DirectoryPath = pq.QueryPath
	pq.DirectorySubstring = pq.OriginalQuery

	stat, err := os.Stat(pq.QueryPath)

	if err == nil && stat.IsDir() && !strings.HasSuffix(pq.DirectoryPath, "/") {
		pq.DirectoryPath += "/"
		pq.DirectorySubstring += "/"
	} else {
		if ind := strings.LastIndex(pq.QueryPath, "/"); ind >= 0 {
			pq.DirectoryPath = pq.QueryPath[:ind+1]
			pq.Filename = pq.QueryPath[ind+1:]
		}
		if ind := strings.LastIndex(pq.OriginalQuery, "/"); ind >= 0 {
			pq.DirectorySubstring = pq.OriginalQuery[:ind+1]
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
		return nil, fmt.Errorf("`%v` is executable", pq.OriginalQuery)
	}

	entry, err := NewEntryForFile(pq.QueryPath, "<b>"+pq.OriginalQuery+"</b>", pq.OriginalQuery)
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
	isPath, pq := NewPathQuery(query)
	if !isPath {
		return
	}

	if entry, err := pq.MakeLaunchEntry(); err != nil {
		errduring("making file entry `%v`", err, "Skipping it", pq.QueryPath)
	} else {
		results = append(results, entry)
	}

	filenames, err := pq.DirFilenames()
	if err != nil {
		errduring("retrieving dir `%v` filenames", err, "No file entries are retrieved", pq.DirectoryPath)
		return
	}
	sort.Strings(filenames)

	queryFnLen := len(pq.Filename)
	for _, name := range filenames {
		if !strings.HasPrefix(name, pq.Filename) {
			continue
		}

		filePath := pq.DirectoryPath + name
		if filePath == pq.QueryPath {
			continue
		}

		tabFilePath := pq.DirectorySubstring + name
		if stat, err := os.Stat(filePath); err == nil && stat.IsDir() {
			tabFilePath += "/"
		}
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

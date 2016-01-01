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

func SearchFileEntries(query string) (results LaunchEntriesList) {
	isPath, pq := NewPathQuery(query)
	if !isPath {
		return nil
	}

	if entry, err := pq.MakeLaunchEntry(); err != nil {
		errduring("making file entry `%v`", err, "Skipping it", pq.QueryPath)
	} else {
		results = append(results, entry)
	}

	dir, err := os.Open(pq.DirectoryPath)
	if err != nil {
		errduring("opening dir `%v`", err, "Skipping it", pq.DirectoryPath)
		return
	}
	dirStat, err := dir.Stat()
	if err != nil || !dirStat.IsDir() {
		errduring("retrieving dir stat `%v`", err, "Skipping it", pq.DirectoryPath)
		return
	}
	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		errduring("retrieving dirnames `%v`", err, "Skipping it", pq.DirectoryPath)
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

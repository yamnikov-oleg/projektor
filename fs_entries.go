package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func SearchFileEntries(query string) (results LaunchEntriesList) {
	isPath, queryPath := ExpandPathString(query)
	if !isPath {
		return nil
	}

	stat, statErr := os.Stat(queryPath)
	if statErr == nil && !IsExecutable(stat) {
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

		tabFilePath := displayDirPath + name
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

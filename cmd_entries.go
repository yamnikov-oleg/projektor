package main

import (
	"fmt"
	"os"
	"strings"
)

var PathEntries = strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

func SearchCmdEntries(query string) (list LaunchEntriesList) {
	if query != "" && !IsInHistory(query) && IsCommand(query) {
		return LaunchEntriesList{NewEntryFromCommand(query)}
	}
	return nil
}

func IsCommand(cmd string) bool {
	return IsPathExecutable(cmd) || IsCommandInPath(cmd)
}

func IsPathExecutable(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return IsExecutable(stat)
}

func IsCommandInPath(cmd string) bool {
	for _, path := range PathEntries {
		filePath := fmt.Sprintf("%v%v%v", path, string(os.PathSeparator), cmd)
		if IsPathExecutable(filePath) {
			return true
		}
	}
	return false
}

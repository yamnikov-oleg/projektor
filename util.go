package main

import (
	"os"
	"strings"
)

func EscapeAmpersand(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
}

func ExpandPathString(query string) (isPath bool, path string) {
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

func IsExecutable(info os.FileInfo) bool {
	if info == nil || info.IsDir() || (info.Mode().Perm()&0111) == 0 {
		return false
	}
	return true
}

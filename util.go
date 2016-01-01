package main

import (
	"os"
	"regexp"
	"strings"
)

var (
	EnvVarRegexp = regexp.MustCompile(`\$(\w+)`)
)

func EscapeAmpersand(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
}

func ExpandEnvVars(query string) string {
	matches := EnvVarRegexp.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		envVar := match[0]
		cleanEnvVar := match[1]
		query = strings.Replace(query, envVar, os.Getenv(cleanEnvVar), 1)
	}
	return query
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

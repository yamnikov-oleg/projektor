package main

import (
	"os"
	"regexp"
	"strings"
	"path"
)

var (
	EnvVarRegexp    = regexp.MustCompile(`\$(\w+)`)
	UrlSchemaRegexp = regexp.MustCompile(`\w+://`)
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

func ExpandPathString(query string) (bool, string) {
	if query == "" {
		return false, query
	}
	if query[0] != '/' && query[0] != '~' {
		return false, query
	}

	if query[0] == '~' {
		query = HOME + query[1:]
	}
	return true, path.Clean(query)
}

func IsExecutable(info os.FileInfo) bool {
	if info == nil || info.IsDir() || (info.Mode().Perm()&0111) == 0 {
		return false
	}
	return true
}

func IsUrl(query string) bool {
	inds := UrlSchemaRegexp.FindStringIndex(query)
	if len(inds) == 0 {
		return false
	}
	if inds[0] != 0 {
		return false
	}
	return true
}

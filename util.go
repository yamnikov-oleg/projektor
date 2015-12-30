package main

import "strings"

func EscapeAmpersand(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
}

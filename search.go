package main

import (
	"fmt"
	"os"

	"github.com/yamnikov-oleg/go-gtk/gio"
)

func ExecCommandEntry(command string) *LaunchEntry {
	return &LaunchEntry{
		Icon:       "application-default-icon",
		MarkupName: fmt.Sprintf("\u2192 <b>%v</b>", command),
		Cmdline:    command,
	}
}

func SpecialEntry(query string) *LaunchEntry {
	if query == "" {
		return nil
	}

	if query[0] != '/' && query[0] != '~' {
		return ExecCommandEntry(query)
	}

	if query[0] == '~' {
		query = HOME + query[1:]
	}

	stat, statErr := os.Stat(query)
	if statErr != nil {
		return nil
	}

	if !stat.IsDir() && (stat.Mode().Perm()&0111) != 0 {
		return ExecCommandEntry(query)
	}

	gFileInfo, fiErr := gio.NewFileForPath(query).QueryInfo("standard::*", gio.FILE_QUERY_INFO_NONE, nil)
	if fiErr != nil {
		return nil
	}

	icon := gFileInfo.GetIcon()
	return &LaunchEntry{
		Icon:       icon.ToString(),
		MarkupName: fmt.Sprintf("<b>%v</b>", query),
		Cmdline:    "xdg-open " + query,
	}
}

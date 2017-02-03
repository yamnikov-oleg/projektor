package main

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type HistRecord struct {
	// = LaunchEntry.Name
	Name string `yaml:"name"`
	// = LaunchEntry.TabName
	TabName string `yaml:"tabname"`
	// = LaunchEntry.Icon
	Icon string `yaml:"icon"`
	// = LaunchEntry.Cmdline
	Cmdline string `yaml:"cmdline"`
}

type HistoryWarehouse []HistRecord

const (
	HistCapacity = 40
)

var (
	HistFilepath = AppDir + "/history.yaml"
	// Older records are first, newer are at the end
	History HistoryWarehouse
)

func LoadHistory() {
	err := os.MkdirAll(AppDir, 0700)
	if err != nil {
		errduring("creating app directory at %q", err, "Exiting", AppDir)
		os.Exit(1)
	}
	f, err := os.Open(HistFilepath)
	if err != nil {
		errduring("opening history file %q", err, "Skipping history loading", HistFilepath)
		return
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		errduring("reading history file contents %q", err, "Skipping history loading", HistFilepath)
		return
	}
	err = yaml.Unmarshal(contents, &History)
	if err != nil {
		errduring("unmarshaling history from %q", err, "Skipping history loading", HistFilepath)
		return
	}

	// Truncate to proper size
	hl := len(History)
	if hl > HistCapacity {
		History = History[hl-HistCapacity:]
	}
}

func MakeHistRecord(entry HistRecord) {
	remi := -1
	for i, r := range History {
		if r.Cmdline == entry.Cmdline {
			remi = i
			break
		}
	}
	if remi >= 0 {
		History = append(History[:remi], History[remi+1:]...)
	}

	History = append(History, entry)

	buf, err := yaml.Marshal(History)
	if err != nil {
		errduring("history marshaling", err, "Skipping history saving")
		return
	}

	f, err := os.Create(HistFilepath)
	if err != nil {
		errduring("history file creation", err, "Skipping history saving")
		return
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		errduring("history file writing", err, "Skipping history saving")
		return
	}
}

func IsInHistory(cmdline string) bool {
	for _, r := range History {
		if r.Cmdline == cmdline {
			return true
		}
	}
	return false
}

func SearchHistEntries(query string) (list LaunchEntriesList) {
	expanded := ExpandEnvVars(query)
	if strings.HasPrefix(expanded, "~") {
		expanded = strings.Replace(expanded, "~", HOME, 1)
	}
	seps := []string{
		query,
		strings.ToLower(query),
		expanded,
		strings.ToLower(expanded),
	}

	findAny := func(s string, seps []string) (index, length int) {
		if query == "" {
			return 0, 0
		}
		for _, sep := range seps {
			if sep == "" {
				continue
			}
			index = strings.Index(s, sep)
			if index >= 0 {
				return index, len(sep)
			}
		}
		return -1, 0
	}

	for i, r := range History {
		lon := strings.ToLower(r.Name)
		ind, length := findAny(lon, seps)
		if ind < 0 {
			continue
		}

		entry := &LaunchEntry{
			Type:       HistEntry,
			Name:       r.Name,
			TabName:    r.TabName,
			Icon:       r.Icon,
			Cmdline:    r.Cmdline,
			QueryIndex: -i,
		}
		entry.UpdateMarkupName(ind, length)
		list = append(list, entry)
	}

	list.SortByIndex()
	return
}

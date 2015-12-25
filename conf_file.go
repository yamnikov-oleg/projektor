package main

import (
	"bufio"
	"io"
	"strings"
)

type Values map[string]string

func (v Values) Val(key string) string {
	return v[key]
}

func (v Values) Has(key string) bool {
	_, ok := v[key]
	return ok
}

type ConfFile struct {
	Sections map[string]Values
}

func ReadConfFile(reader io.Reader) (cf *ConfFile, err error) {
	scanner := bufio.NewScanner(reader)
	section := ""

	cf = &ConfFile{}
	cf.Sections = make(map[string]Values)
	cf.Sections[section] = make(Values)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue // Empty line, skipping it...
		}
		if line[0] == '[' {
			section = line[1 : len(line)-1]
			cf.Sections[section] = make(Values)
			continue
		}
		sepIndex := strings.Index(line, "=")
		if sepIndex < 0 {
			continue // Not a key-value pair, skipping it...
		}

		key := line[:sepIndex]
		value := line[sepIndex+1:]

		cf.Sections[section][key] = value
	}

	err = scanner.Err()
	return
}

func (cf *ConfFile) Sec(name string) Values {
	return cf.Sections[name]
}

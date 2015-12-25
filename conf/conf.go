package conf

import (
	"bufio"
	"io"
	"strings"
)

type Values map[string]string

func (v Values) Bool(key string) bool {
	return v[key] == "true"
}

func (v Values) Str(key string) string {
	return v[key]
}

func (v Values) Has(key string) bool {
	_, ok := v[key]
	return ok
}

type File struct {
	Sections map[string]Values
}

func Read(reader io.Reader) (cf *File, err error) {
	scanner := bufio.NewScanner(reader)
	section := ""

	cf = &File{}
	cf.Sections = make(map[string]Values)

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

		if section != "" {
			cf.Sections[section][key] = value
		}
	}

	err = scanner.Err()
	return
}

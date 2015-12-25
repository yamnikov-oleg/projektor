// +build linux

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yamnikov-oleg/projektor/conf"
)

var (
	sharedAppDir = "/usr/share/applications"
	localAppDir  = os.Getenv("HOME") + "/.local/share/applications"
)

const (
	Verbose      = true
	MaxEntrySize = 1 * 1024 * 1024 // 1Mb
)

func logf(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

func errduring(action string, err error, nextmove string, a ...interface{}) {
	line := action + ":\n"
	if err != nil {
		line += err.Error() + "\n"
	}
	if nextmove != "" {
		line += nextmove + "\n"
	}
	logf(line, a...)
}

type Entry struct {
	Name string
	Icon string
	Exec string
}

func ReadEntry(filename string) (en *Entry, err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fd.Close()

	cf, err := conf.Read(fd)
	if err != nil {
		return
	}

	section := cf.Sections["Desktop Entry"]
	if section.Bool("Hidden") {
		return nil, errors.New("desktop entry hidden")
	}
	if section.Bool("NoDisplay") {
		return nil, errors.New("desktop entry not displayed")
	}

	en = &Entry{}
	en.Name = section.Str("Name")
	en.Icon = section.Str("Icon")

	r := strings.NewReplacer(" %f", "", " %F", "", " %u", "", " %U", "")
	en.Exec = r.Replace(section.Str("Exec"))

	return
}

type EntriesReader struct {
	files        []string
	currentIndex int
	Entry        *Entry
}

func NewEntriesReader() *EntriesReader {
	er := &EntriesReader{nil, -1, nil}

	if err := er.AppendDirectory(sharedAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", sharedAppDir)
	}
	if err := er.AppendDirectory(localAppDir); err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", localAppDir)
	}

	return er
}

func (er *EntriesReader) AppendDirectory(dir string) error {
	efs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range efs {
		er.files = append(er.files, dir+"/"+fi.Name())
	}
	return nil
}

func (er *EntriesReader) Good() bool {
	return er.currentIndex >= 0 && er.currentIndex < len(er.files)
}

func (er *EntriesReader) Next() bool {
	er.currentIndex++
	if !er.Good() {
		return false
	}

	filepath := er.files[er.currentIndex]
	file, err := os.Stat(filepath)
	if err != nil {
		errduring("reading entry file `%v`", err, "Skipping it", filepath)
		return er.Next()
	}

	if file.IsDir() {
		return er.Next()
	}
	if file.Size() > MaxEntrySize {
		errduring("reading .desktop file `%v`: size to big!", nil, "Skipping it", filepath)
		return er.Next()
	}
	if !strings.HasSuffix(file.Name(), ".desktop") {
		return er.Next()
	}

	en, err := ReadEntry(filepath)
	if err != nil {
		return er.Next()
	}
	er.Entry = en

	return true
}

func main() {
	reader := NewEntriesReader()
	for reader.Next() {
		fmt.Printf("%#v\n", reader.Entry)
	}
}

// +build linux

package main

import (
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

	en = &Entry{}
	section := cf.Sections["Desktop Entry"]
	en.Name = section.Str("Name")
	en.Icon = section.Str("Icon")
	en.Exec = section.Str("Exec")

	en.Exec = strings.Replace(en.Exec, " %f", "", -1)
	en.Exec = strings.Replace(en.Exec, " %F", "", -1)
	en.Exec = strings.Replace(en.Exec, " %u", "", -1)
	en.Exec = strings.Replace(en.Exec, " %U", "", -1)
	return
}

type EntriesReader struct {
	files        []string
	currentIndex int
	Entry        *Entry
}

func NewEntriesReader() *EntriesReader {
	er := &EntriesReader{nil, -1, nil}

	efs, err := ioutil.ReadDir(sharedAppDir)
	if err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", sharedAppDir)
	} else {
		for _, fi := range efs {
			er.files = append(er.files, sharedAppDir+"/"+fi.Name())
		}
	}

	efs, err = ioutil.ReadDir(localAppDir)
	if err != nil {
		errduring("reading of the directory `%v`", err, "Skipping it", localAppDir)
	} else {
		for _, fi := range efs {
			er.files = append(er.files, localAppDir+"/"+fi.Name())
		}
	}

	return er
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

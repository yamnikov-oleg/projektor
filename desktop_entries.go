package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yamnikov-oleg/projektor/conf"
)

const (
	MaxEntrySize = 1 * 1024 * 1024 // 1Mb
)

var (
	sharedAppDir = "/usr/share/applications"
	localAppDir  = os.Getenv("HOME") + "/.local/share/applications"
)

type DtEntry struct {
	Name       string
	LoCaseName string
	Icon       string
	Exec       string
}

func DtEntryFromFile(filename string) (en *DtEntry, err error) {
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

	en = &DtEntry{}
	en.Name = section.Str("Name")
	en.LoCaseName = strings.ToLower(en.Name)
	en.Icon = section.Str("Icon")

	r := strings.NewReplacer(" %f", "", " %F", "", " %u", "", " %U", "")
	en.Exec = r.Replace(section.Str("Exec"))

	return
}

type EntriesReader struct {
	files        []string
	currentIndex int
	Entry        *DtEntry
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

	en, err := DtEntryFromFile(filepath)
	if err != nil {
		return er.Next()
	}
	er.Entry = en

	return true
}

var DesktopEntries []*DtEntry

func IndexDesktopEntries() {
	reader := NewEntriesReader()
	for reader.Next() {
		DesktopEntries = append(DesktopEntries, reader.Entry)
	}
}

type EntriesIterator struct {
	Index int
}

func NewEntriesInterator() *EntriesIterator {
	return &EntriesIterator{-1}
}

func (ei *EntriesIterator) Good() bool {
	return ei.Index >= 0 && ei.Index < len(DesktopEntries)
}

func (ei *EntriesIterator) Next() bool {
	ei.Index++
	return ei.Good()
}

func (ei *EntriesIterator) Entry() *DtEntry {
	if !ei.Good() {
		return nil
	}
	return DesktopEntries[ei.Index]
}

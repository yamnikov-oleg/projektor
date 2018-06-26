package main

import (
	"io/ioutil"
	"os"
	"path"
)

const CSS_CODE = `
  GtkEntry {
    background-image: none;
    border: none;
    box-shadow: none;
    font-size: 12pt;
  }
  GtkTreeView {
    background-color: transparent;
    font-size: 8pt;
  }
  GtkTreeView:selected {
    background-color: rgba(0,0,0,0.1);
  }
  GtkTreeView.cell {
    padding: 6px 3px;
  }
`

var (
	StylesFilePath = path.Join(AppDir, "styles.css")
)

func AppStyle() string {
	if _, err := os.Stat(StylesFilePath); !os.IsNotExist(err) {
		cssBytes, _ := ioutil.ReadFile(StylesFilePath)
		return string(cssBytes)
	}
	return CSS_CODE
}

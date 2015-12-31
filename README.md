# projektor
Fast application launcher for Gnome written in Go

![Screenshot](/screenshots/01.png?raw=true)

## Features

* Search and launch applications installed on your system
* Navigate through file system and open directories and files via ```xdg-open```
* Execute custom command lines

## Building

* Install gtk+-3.0 development files
* Install latest version of Golang from https://golang.org
* ```go get github.com/yamnikov-oleg/projektor```
* Make sure ```&GOPATH/bin``` is appended to ```&PATH```
* Done.

## Usage

Start ```projektor``` from console. Best practice would be to bind the command to a key shortcut (e.g. ```Super+Q```) using some utility.

* Navigate through launch entries using mouse wheel or arrow keys.
* Use top text box for searching.
* Press enter or double-click an entry to launch it.
* Start typing a path to navigate file system.
* Press Tab to insert selected entry's command line into the text box.
* Type any command line then select appropriate launch entry to execute the command.
* Press escape, click outside the window or start ```projektor``` again to close projektor window.

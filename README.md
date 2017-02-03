# projektor
Fast application launcher for Gnome written in Go

![Screenshot](/screenshots/03.png?raw=true)

## Features

* Search and launch applications installed on your system
* Navigate through file system, open directories and files
* Execute custom command lines in background
* Open urls in the default web browser

## Building

* Install gtk+-3.0 development files
* Install latest version of Golang from https://golang.org
* `go get github.com/yamnikov-oleg/projektor`
* Make sure `$GOPATH/bin` is appended to `$PATH` (optional, but convenient).
* Done.

## Usage

Execute `projektor` command to start a daemon, which will listen for `Super+Q` hotkey. Use `Super+Q` to launch projektor window any time. Add `projektor` to startup applications on your system.

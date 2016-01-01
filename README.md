# projektor
Fast application launcher for Gnome written in Go

![Screenshot](/screenshots/02.png?raw=true)

## Features

* [Demo video](https://youtu.be/-i69v6F41ps)
* Search and launch applications installed on your system
* Navigate through file system using Tab key, open directories and files
* Execute custom command lines
* Open urls in the default web browser

## Building

* Install gtk+-3.0 development files
* Install latest version of Golang from https://golang.org
* ```go get github.com/yamnikov-oleg/projektor```
* Make sure ```&GOPATH/bin``` is appended to ```&PATH```
* Done.

## Usage

Start ```projektor``` from console. Best practice would be to bind the command to a key shortcut (e.g. ```Super+Q```) using some utility.

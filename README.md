# projektor
Fast application launcher for Gnome written in Go

![Screenshot](/screenshots/03.png?raw=true)

## Features

* Search and launch any __app__ installed on you system:
  + Type _chrome_, _gimp_, _steam_ etc. and press Enter to launch the app.
* Run any __command__ which can be found in your `PATH`:
  + _gksu service mysql restart_, _cp ~/src/file1 ~/dst/_ etc.
  + Search for an app, select it with arrow keys and press Tab to copy its
    commandline into the input box. Now you can modify launch options and start
    it with Enter button.
* Open __files and directories__; navigate your __file system__ through launcher's
  UI without using the mouse.
  + Type path to a file and press Enter to open it in an appropriate app.
  + Type path to a directory and press Enter to open it in your file manager.
  + If your type directory path, Projektor will list all the files in that
    directory. Select it with arrow keys and press Enter to open.
  + Select a directory and press Tab to copy its path into the input box.
    Projektor will now walk into that dir and list its files.
  + Type .. (two dots) to walk into the parent directory.
* Open any _URLs_ and web addresses:
  + Type _github.com_ and Projektor will suggest you to open this website in your
    favourite web browser.
  + Type `ftp://admin@host.com` to launch your favourite FTP client.
  + Type any URL which is supported on your system.
* Perform __web searching__
  + Type any text and Projektor will show an entry to search it on the web with
    your default browser.
  + By default Projektor uses Google, but you can set any search engine in
    the config file.
* Calculate __math equations__:
  + `5^3 * (15/4 + 1.5)`
  + `sqrt(256) + pow(3, 2.5) / fact(5)`
  + `ln(e) * log2(16) * log10(1000) * log(3, 81)`
  + `sin(1) * cos(pi/2) + tan(3*pi/2) - cot(4*pi/3)`
* Finally, Projektor keeps __history__ of your recently launched entries.
  + When you're searching, Projektor will show history results at the top.
  + Open Projektor and simply press Enter to launch last used item.

## Installing

Please note, this launcher was only tested on Ubuntu-like linux distros. It might
not work on other systems as expected. Feel free to leave an issue, if you
experience problems using Projektor. Thank you :)

Projektor is built with Gtk 3 and requires the library to be installed on your system.

To install Projektor:

* Download latest build from the [releases page](https://github.com/yamnikov-oleg/projektor/releases).
  If build for you OS or architecture is not available, read Building section on
  how to build projektor yourself.
* Add projektor to your startup applications. Use `projektor &` to start the
  key binding daemon.
* Done! Use `Super+Q` to open the launcher. Happy projekting :)

## Building

* Install gtk+-3.0 development files (`sudo apt install libgtk-3-dev` if you're an Ubuntu user)
* Install latest version of Golang from https://golang.org
* Run `go get github.com/yamnikov-oleg/projektor` in your terminal.
* If the build succeeds, the binary will be located in the `$GOPATH/bin` directory.
  Now you can install using installing instructions in this Readme.

## Usage

When projektor keybind daemon is running, `Super+Q` keys will open the projektor
UI. Type what you want to launch, select the item with arrow keys, press enter
or double-click it with left mouse button to launch it. Also, try experimenting
with Tab key :)

For complete list of functions read the Features section.

## Configuration

Projektor configuration file is located at `~/.projektor/config.yaml`.

Here's annotation for every paramater of the default config:

```yaml
# Key bind, used by Projektor daemon. `mod4` is the Windows (Super) key.
# If your Super+Q shortcut doesn't work, edit this option.
# You can use modifier keys `shift` and `control`.
# To identify some complex keybind use `xbindkeys` tool.
keybind: mod4-q
# Searching categories, enabled for use.
# Disable a category by settings its flag to `false` and projektor will no longer
# offer you entries of that category.
enabledcategories:
  # Mathematical calculations
  calc: true
  # Entries history
  history: true
  # Installed applications
  apps: true
  # URLs and web addresses
  url: true
  # Custom command lines
  commands: true
  # File system
  files: true
  # Web searching
  websearch: true
# History category configuration
history:
  # How many last used items should Projektor remember for history?
  capacity: 40
# URL category configuration
url:
  # Icon for URL entries. Use a name of some gtk icon installed on your system,
  # e.g. `firefox`. You can specify absolute path to an image file as well,
  # e.g. `/home/me/.projektor/url-icon.png`
  icon: web-browser
# Web searching category configuration
websearch:
  # Template for search url. `%s` marker denotes, where the search query
  # should be inserted.
  engine: https://google.com/search?q=%s
  # Icon for web search entries. Use a name of some gtk icon installed
  # on your system. You can specify absolute path to an image file as well,
  # e.g. `/home/me/.projektor/google.png`
  icon: web-browser
```

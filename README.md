# ghnotify

Simple tool which polls the GitHub API to check if repos have been updated (i.e. commits have been pushed).

## Requirements

* Go (tested on 1.4.2, should work with 1+).

## Installation

If you haven't setup Go before, you need to first set a `GOPATH` (see [https://golang.org/doc/code.html#GOPATH](https://golang.org/doc/code.html#GOPATH)).

    $ go get github.com/dhowden/ghnotify

This will fetch the code and build the `ghnotify` command line tool and put it in `$GOPATH/bin` (assumed to be in your `PATH` already).

Now:

    $ ghnotify -config $GOPATH/src/github.com/dhowden/ghnotify/config.json

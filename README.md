# ghnotify
[![Build Status](https://travis-ci.org/dhowden/ghnotify.svg?branch=master)](https://travis-ci.org/dhowden/ghnotify)

Simple tool which polls the GitHub API to check if repos have been updated (i.e. commits have been pushed).

## Requirements

* [Go](http://golang.org/dl/) (tested on 1.4.2, should work with 1+).

## Installation

If you haven't setup Go before, you need to first set a `GOPATH` (see [https://golang.org/doc/code.html#GOPATH](https://golang.org/doc/code.html#GOPATH)).

    $ go get github.com/dhowden/ghnotify

This will fetch the code and build the `ghnotify` command line tool and put it in `$GOPATH/bin` (assumed to be in your `PATH` already).

Now:

    $ ghnotify -config $GOPATH/src/github.com/dhowden/ghnotify/config.json

## Slack Integration

First setup [Incoming WebHooks](https://api.slack.com/incoming-webhooks "Slack Incoming Webhooks") for your Slack account and you will get a URL which can be used to post messages into a Slack channel.  Then:

    $ ghnotify -slack-webhook-url YOUR_URL_HERE

## Flowdock Integration

First retrieve the Flow API token and then:

    $ ghnotify -flowdock-token 90a8b5e69bd9125f40...

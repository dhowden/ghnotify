// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"sort"
	"time"
)

// Notifier is an interface which defines the Notify method.
type Notifier interface {
	// Notify takes a map of repo -> updated time.  Returns an error if there
	// was a problem handling/passing on the notification.
	Notify(map[string]time.Time) error
}

type logNotifier struct{}

func (logNotifier) Notify(repos map[string]time.Time) error {
	keys := make([]string, 0, len(repos))
	for k := range repos {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	for _, k := range keys {
		log.Printf("%v\t: %v", k, repos[k])
	}
	return nil
}

// changesNotifier is a wrapper for a Notifier which only passes on Notify calls
// to the underlying Notifier if the repos or updated times have changed since the
// last call.
type changesNotifier struct {
	Notifier

	last map[string]time.Time
}

func (d changesNotifier) Notify(repos map[string]time.Time) error {
	changes := make(map[string]time.Time)
	for k, v := range repos {
		if d.last[k] != v {
			changes[k] = v
			d.last[k] = v
		}
	}

	if len(changes) == 0 {
		return nil
	}
	return d.Notifier.Notify(changes)
}

// NewMultiNotifier creates a new Notifier which will call Notify on every passed Notifier.
// NB: if any notifier returns an error, then it will be returned immediately and the remaing
// notifiers will not be called.
func NewMultiNotifier(notifiers ...Notifier) Notifier {
	return multiNotifier{
		notifiers: notifiers,
	}
}

type multiNotifier struct {
	notifiers []Notifier
}

func (m multiNotifier) Notify(repos map[string]time.Time) error {
	for _, n := range m.notifiers {
		err := n.Notify(repos)
		if err != nil {
			return err
		}
	}
	return nil
}

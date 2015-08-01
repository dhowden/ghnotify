package main

import (
	"log"
	"sort"
	"time"
)

type Notifier interface {
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
	return d.Notifier.Notify(changes)
}

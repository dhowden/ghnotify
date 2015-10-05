// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type githubPoller struct {
	remaining, limit int
	reset            time.Time
	minPoll          time.Duration

	errCh    chan error
	repolist []string

	notifier Notifier
}

func (g *githubPoller) poll() {
	timeout := time.After(0 * time.Second)
	for {
		select {
		case <-timeout:
			repos := make(map[string]time.Time, len(g.repolist))
			for _, repo := range g.repolist {
				u, err := g.fetchRepoUpdated(repo)
				if err != nil {
					g.errCh <- err
					break
				}
				repos[repo] = u
			}

			g.notifier.Notify(repos)

			polls := (g.remaining) / len(g.repolist)
			r := g.reset.Sub(time.Now())
			if polls > 0 {
				rerun := time.Duration(int64(r) / int64(polls))
				if g.minPoll > rerun {
					rerun = g.minPoll
				}
				timeout = time.After(rerun)
				break
			}
			timeout = time.After(g.reset.Sub(time.Now()))
		}
	}
}

func (g *githubPoller) fetchRepoUpdated(repo string) (result time.Time, err error) {
	r, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo, nil)
	if err != nil {
		return
	}
	r.Header.Add("Accept", "application/vnd.github.v3+json")

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = g.updateLimits(resp.Header)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result, err = getUpdatedAt(bodyBytes)
	return
}

func (g *githubPoller) updateLimits(h http.Header) (err error) {
	limit, err := strconv.Atoi(h.Get(limitHeaderLimit))
	if err != nil {
		return
	}
	g.limit = int(limit)

	remaining, err := strconv.Atoi(h.Get(limitHeaderRemaining))
	if err != nil {
		return
	}
	g.remaining = int(remaining)

	reset, err := strconv.ParseInt(h.Get(limitHeaderReset), 10, 0)
	if err != nil {
		return
	}
	g.reset = time.Unix(reset, 0)
	return
}

func getUpdatedAt(data []byte) (time.Time, error) {
	resp := struct {
		UpdatedAt time.Time `json:"updated_at"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return time.Time{}, err
	}
	return resp.UpdatedAt, nil
}

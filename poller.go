package main

import (
	"encoding/json"
	"fmt"
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
				}
				repos[repo] = u
			}

			fmt.Printf("Fetched repo data:\n%+v\n", repos)

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

	err = g.getHeaderInfo(resp.Header)
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

func (g *githubPoller) getHeaderInfo(h http.Header) (err error) {
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

func getUpdatedAt(data []byte) (result time.Time, err error) {
	resp := make(map[string]interface{})
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}

	u, ok := resp["updated_at"]
	if !ok {
		err = fmt.Errorf("`updated_at` not in returned resp")
		return
	}

	us, ok := u.(string)
	if !ok {
		err = fmt.Errorf("expected `updated_at` to be a string, got %+v", u)
		return
	}

	result, err = time.Parse(time.RFC3339, us)
	return
}

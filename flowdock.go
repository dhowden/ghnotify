// Copyright 2015, Dejan Golja
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type flowDockNotifier struct {
	Flow_token string
}

func (f flowDockNotifier) Notify(repos map[string]time.Time) error {
	keys := make([]string, 0, len(repos))
	for k := range repos {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	for _, k := range keys {
		line := fmt.Sprintf("Repo <https://github.com/%v|%v> was updated at %v", k, k, repos[k])
		lines = append(lines, line)
	}

	form := url.Values{}
	form.Set("content", strings.Join(lines, "\n"))
	form.Set("event", "comment")
	form.Set("external_user_name", "ghnotify")

	resp, err := http.Post("https://api.flowdock.com/messages/chat/"+f.Flow_token,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("error performing POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}
	return nil
}

// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type slackWebHookNotifier struct {
	url string
}

func (s slackWebHookNotifier) Notify(repos map[string]time.Time) error {
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

	payload := struct {
		Text string `json:"text"`
	}{
		Text: strings.Join(lines, "\n"),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshallaing payload: %v", err)
	}

	resp, err := http.Post(s.url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("error performing POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}
	return nil
}

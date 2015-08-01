package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	configFile string
)

var (
	limitHeaderLimit     string = "X-Ratelimit-Limit"
	limitHeaderRemaining        = "X-Ratelimit-Remaining"
	limitHeaderReset            = "X-Ratelimit-Reset"
)

type Config struct {
	Repos    []string `json:"repos"`
	DataFile string   `json:"data_file"`
	MinPoll  string   `json:"min_poll",omitempty`
}

func init() {
	flag.StringVar(&configFile, "config", "config.json", "Config file")
}

func main() {
	flag.Parse()

	if configFile == "" {
		fmt.Fprintf(os.Stderr, "no config specified\n")
		os.Exit(1)
		return
	}

	config, err := readConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %s\n", configFile, err)
		os.Exit(1)
		return
	}

	repoList := config.Repos
	if len(repoList) == 0 {
		fmt.Fprintf(os.Stderr, "no repos specified\n")
		os.Exit(1)
		return
	}

	var minPoll time.Duration
	if config.MinPoll != "" {
		minPoll, err = time.ParseDuration(config.MinPoll)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse min_poll")
			os.Exit(1)
			return
		}
	}

	errCh := make(chan error)
	go func() {
		for err := range errCh {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()

	n := changesNotifier{
		logNotifier{},
		make(map[string]time.Time),
	}

	poller := githubPoller{
		repolist: repoList,
		minPoll:  minPoll,
		errCh:    errCh,
		notifier: n,
	}
	poller.poll()
}

func readConfig(fileName string) (c Config, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &c)
	return
}

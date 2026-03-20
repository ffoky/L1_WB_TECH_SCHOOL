package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"L2.16/wget"
	"L2.16/wget/downloader"
)

const (
	errMissingUrl = "wget: missing URL"
	timeout       = 10 * time.Second
)

func main() {
	var opts wget.Options
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wget [OPTION]... [URL]...\n\n")
	}
	flag.BoolVar(&opts.Recursive, "r", false, "enable recursion")
	flag.IntVar(&opts.DepthLevel, "l", 1, "depth level of recursion, use with r flag")
	flag.IntVar(&opts.Workers, "n", 5, "number of parallel downloads")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, errMissingUrl)
		flag.Usage()
		os.Exit(1)
	}
	url := args[0]

	if !opts.Recursive {
		opts.DepthLevel = 0
	}

	client := downloader.NewClient(timeout)
	w := wget.NewWget(client, opts)
	w.LoadRobots(url)
	err := w.Run(url, opts.DepthLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

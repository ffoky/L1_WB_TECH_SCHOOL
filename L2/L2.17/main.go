package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"L2.17/telnet"
)

func main() {
	var opts telnet.Options
	timeout := 0
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: telnet [--timeout=10] host port\"\n\n")
	}

	flag.IntVar(&timeout, "timeout", 10, "connection timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	opts.Host = args[0]
	opts.Port = args[1]
	opts.Timeout = time.Duration(timeout) * time.Second
	if err := telnet.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

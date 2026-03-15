package main

import (
	"flag"
	"fmt"
	"os"

	"L2.12/grep"
)

func main() {
	var opts grep.Options
	var cContext int

	flag.IntVar(&opts.AContext, "A", 0, "print N lines of context after match")
	flag.IntVar(&opts.BContext, "B", 0, "print N lines of context before match")
	flag.IntVar(&cContext, "C", 0, "print N lines of context before and after match")
	flag.BoolVar(&opts.CountSame, "c", false, "count matching lines")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&opts.InvertFilter, "v", false, "invert filter")
	flag.BoolVar(&opts.ExactMatch, "F", false, "fixed string match")
	flag.BoolVar(&opts.StringNumber, "n", false, "print line numbers")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if cContext > 0 {
		opts.AContext = cContext
		opts.BContext = cContext
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	pattern := args[0]

	reader := os.Stdin
	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			}
		}()
		reader = f
	}

	if err := grep.Grep(reader, os.Stdout, pattern, opts); err != nil {
		fmt.Fprintf(os.Stderr, "grep: %v\n", err)
		os.Exit(1)
	}
}

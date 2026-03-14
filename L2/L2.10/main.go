package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"L2.10/sorting"
)

func main() {
	var opts sorting.Options

	flag.IntVar(&opts.Column, "k", 0, "sort by column N")
	flag.BoolVar(&opts.Numeric, "n", false, "sort by numeric values")
	flag.BoolVar(&opts.Reverse, "r", false, "sort reversed")
	flag.BoolVar(&opts.Unique, "u", false, "only distinct chars")
	flag.BoolVar(&opts.Month, "M", false, "sort by month")
	flag.BoolVar(&opts.TrimTailSpace, "b", false, "ignore trailing blanks")
	flag.BoolVar(&opts.IsSorted, "c", false, "check if sorted")
	flag.BoolVar(&opts.HumanReadableSize, "h", false, "sort by human readable sizes")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sort [OPTIONS] [FILE]\n")
		flag.PrintDefaults()
	}

	var reader io.Reader = os.Stdin
	if args := flag.Args(); len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "sort: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "sort: %v\n", err)
			}
		}()
		reader = file
	}

	if err := sorting.Sort(reader, os.Stdout, opts); err != nil {
		fmt.Fprintf(os.Stderr, "sort: %v\n", err)
		os.Exit(1)
	}
}

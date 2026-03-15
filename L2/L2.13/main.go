package main

import (
	"flag"
	"fmt"
	"os"

	"L2.13/cut"
)

func main() {
	var opts cut.Options

	var fieldsRaw string
	flag.StringVar(&fieldsRaw, "f", "", "fields to select")
	flag.StringVar(&opts.Delimiter, "d", "\t", "delimiter")
	flag.BoolVar(&opts.Separated, "s", false, "only lines with delimiter")

	flag.Parse()

	fields, err := cut.ParseFields(fieldsRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cut: %v\n", err)
		os.Exit(1)
	}
	opts.Fields = fields

	if err = cut.Cut(os.Stdin, os.Stdout, opts); err != nil {
		fmt.Fprintf(os.Stderr, "cut: %v\n", err)
		os.Exit(1)
	}
}

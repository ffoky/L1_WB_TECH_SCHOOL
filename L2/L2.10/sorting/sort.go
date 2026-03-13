package sorting

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Sort читает строки из reader, сортирует по opts и пишет в writer.
func Sort(reader io.Reader, writer io.Writer, opts Options) error {
	lines, err := readLines(reader)
	if err != nil {
		return err
	}

	cmp := makeCmpFunc(opts)

	if opts.IsSorted {
		return checkSorted(lines, cmp)
	}

	sort.SliceStable(lines, func(i, j int) bool {
		return cmp(lines[i], lines[j]) < 0
	})

	if opts.Unique {
		lines = deduplicate(lines, cmp)
	}

	return writeLines(writer, lines)
}

// readLines читает все строки.
func readLines(reader io.Reader) ([]string, error) {
	newReader := bufio.NewReader(reader)
	var lines []string
	for {
		line, err := newReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					lines = append(lines, strings.TrimRight(line, "\r\n"))
				}
				break
			}
			return nil, err
		}
		lines = append(lines, strings.TrimRight(line, "\r\n"))
	}
	return lines, nil
}

// checkSorted проверяет что строки отсортированы.
func checkSorted(lines []string, cmp func(a, b string) int) error {
	for i := 1; i < len(lines); i++ {
		if cmp(lines[i-1], lines[i]) > 0 {
			return fmt.Errorf("sort: -:%d: disorder: %s", i+1, lines[i])
		}
	}
	return nil
}

// deduplicate убирает  дубликаты.
func deduplicate(lines []string, cmp func(a, b string) int) []string {
	unique := lines[:1]
	for _, line := range lines[1:] {
		if cmp(line, unique[len(unique)-1]) != 0 {
			unique = append(unique, line)
		}
	}
	return unique
}

// writeLines записывает строки в writer.
func writeLines(writer io.Writer, lines []string) error {
	w := bufio.NewWriter(writer)
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return w.Flush()
}

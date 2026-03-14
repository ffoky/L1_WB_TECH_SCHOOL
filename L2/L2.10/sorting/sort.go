package sorting

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const parallelThreshold = 50_000

// без рефлексии из sort.SliceStable, сравниваем предвычисленные ключи напрямую
type keySorter struct {
	lines   []string
	keys    []float64
	reverse bool
}

func (s *keySorter) Len() int { return len(s.lines) }

func (s *keySorter) Less(i, j int) bool {
	if s.reverse {
		return s.keys[i] > s.keys[j]
	}
	return s.keys[i] < s.keys[j]
}

func (s *keySorter) Swap(i, j int) {
	s.lines[i], s.lines[j] = s.lines[j], s.lines[i]
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}

func Sort(reader io.Reader, writer io.Writer, opts Options) error {
	lines, err := readLines(reader)
	if err != nil {
		return err
	}

	cmpFunc := makeCmpFunc(opts)

	if opts.IsSorted {
		return checkSorted(lines, cmpFunc)
	}

	if len(lines) > parallelThreshold {
		sortParallel(lines, opts, cmpFunc)
	} else {
		sortSingle(lines, opts, cmpFunc)
	}

	if opts.Unique {
		lines = deduplicate(lines, cmpFunc)
	}

	return writeLines(writer, lines)
}

func sortSingle(lines []string, opts Options, cmpFunc func(a, b string) int) {
	if opts.Numeric || opts.Month || opts.HumanReadableSize {
		keys := make([]float64, len(lines))
		for i, line := range lines {
			keys[i] = parseKey(line, opts)
		}
		sort.Stable(&keySorter{lines: lines, keys: keys, reverse: opts.Reverse})
	} else {
		slices.SortStableFunc(lines, func(a, b string) int {
			return cmpFunc(a, b)
		})
	}
}

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

func checkSorted(lines []string, cmp func(a, b string) int) error {
	for i := 1; i < len(lines); i++ {
		if cmp(lines[i-1], lines[i]) > 0 {
			return fmt.Errorf("sort: -:%d: disorder: %s", i+1, lines[i])
		}
	}
	return nil
}

func deduplicate(lines []string, cmp func(a, b string) int) []string {
	unique := lines[:1]
	for _, line := range lines[1:] {
		if cmp(line, unique[len(unique)-1]) != 0 {
			unique = append(unique, line)
		}
	}
	return unique
}

func writeLines(writer io.Writer, lines []string) error {
	w := bufio.NewWriter(writer)
	for _, line := range lines {
		if _, err := w.WriteString(line); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}
	return w.Flush()
}

func parseKey(line string, opts Options) float64 {
	field := extractField(line, opts.Column)
	if opts.TrimTailSpace {
		field = strings.TrimRight(field, " \t")
	}
	switch {
	case opts.Numeric:
		f, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return 0
		}
		return f
	case opts.Month:
		return float64(monthIndex[field])
	case opts.HumanReadableSize:
		return float64(parseSize(field))
	default:
		return 0
	}
}

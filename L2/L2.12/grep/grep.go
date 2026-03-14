package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const noPrevShown = -2

func Grep(reader io.Reader, writer io.Writer, pattern string, opts Options) error {
	matcher, err := buildMatcher(pattern, opts)
	if err != nil {
		return err
	}

	lines, err := readLines(reader)
	if err != nil {
		return err
	}

	matches := findMatches(lines, matcher, opts.InvertFilter)

	if opts.CountSame {
		_, err = fmt.Fprintln(writer, len(matches))
		return err
	}

	output := collectOutput(lines, matches, opts)
	w := bufio.NewWriter(writer)
	for _, line := range output {
		if _, err = fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return w.Flush()
}

// собирает функцию проверки строки на совпадение,
// подстрока или регулярка
func buildMatcher(pattern string, opts Options) (func(string) bool, error) {
	if opts.ExactMatch {
		p := pattern
		if opts.IgnoreCase {
			p = strings.ToLower(p)
		}
		return func(line string) bool {
			l := line
			if opts.IgnoreCase {
				l = strings.ToLower(l)
			}
			return strings.Contains(l, p)
		}, nil
	}

	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	return re.MatchString, nil
}

func readLines(reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func findMatches(lines []string, matcher func(string) bool, invertFilter bool) []int {
	var matches []int
	for i, line := range lines {
		matched := matcher(line)
		if invertFilter {
			matched = !matched
		}
		if matched {
			matches = append(matches, i)
		}
	}
	return matches
}

// формирует вывод с учётом контекста и нумерации строк,
// если между группами совпадений есть разрыв, ставит "--"
func collectOutput(lines []string, matches []int, opts Options) []string {
	if len(matches) == 0 {
		return nil
	}

	ToShow := make([]bool, len(lines))
	isMatch := make([]bool, len(lines))
	for _, idx := range matches {
		isMatch[idx] = true
		start := idx - opts.BContext
		if start < 0 {
			start = 0
		}
		end := idx + opts.AContext
		if end >= len(lines) {
			end = len(lines) - 1
		}
		for j := start; j <= end; j++ {
			ToShow[j] = true
		}
	}

	var result []string
	prevShown := noPrevShown

	for i, line := range lines {
		if !ToShow[i] {
			continue
		}

		hasContext := opts.AContext > 0 || opts.BContext > 0
		if hasContext && prevShown >= 0 && i-prevShown > 1 {
			result = append(result, "--")
		}

		if opts.StringNumber {
			prefix := "-"
			if isMatch[i] {
				prefix = ":"
			}
			line = fmt.Sprintf("%d%s%s", i+1, prefix, line)
		}

		result = append(result, line)
		prevShown = i
	}

	return result
}

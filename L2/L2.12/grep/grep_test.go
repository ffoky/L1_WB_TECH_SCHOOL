package grep_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"L2.12/grep"
)

func runGrep(t *testing.T, input, pattern string, opts grep.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := grep.Grep(strings.NewReader(input), &buf, pattern, opts)
	require.NoError(t, err)
	return buf.String()
}

func TestGrep(t *testing.T) {
	input := "apple\nbanana\ncherry\ndate\nelderberry\n"

	tests := []struct {
		name     string
		pattern  string
		opts     grep.Options
		expected string
	}{
		{
			name:     "простое совпадение",
			pattern:  "ban",
			expected: "banana\n",
		},
		{
			name:     "нет совпадений",
			pattern:  "xyz",
			expected: "",
		},
		{
			name:     "несколько совпадений",
			pattern:  "e",
			expected: "apple\ncherry\ndate\nelderberry\n",
		},
		{
			name:     "-i игнорировать регистр",
			pattern:  "BANANA",
			opts:     grep.Options{IgnoreCase: true},
			expected: "banana\n",
		},
		{
			name:     "-v инвертировать",
			pattern:  "e",
			opts:     grep.Options{InvertFilter: true},
			expected: "banana\n",
		},
		{
			name:     "-c подсчёт",
			pattern:  "e",
			opts:     grep.Options{CountSame: true},
			expected: "4\n",
		},
		{
			name:     "-n номера строк",
			pattern:  "banana",
			opts:     grep.Options{StringNumber: true},
			expected: "2:banana\n",
		},
		{
			name:     "-F точное совпадение подстроки",
			pattern:  "an",
			opts:     grep.Options{ExactMatch: true},
			expected: "banana\n",
		},
		{
			name:     "-F с -i",
			pattern:  "AN",
			opts:     grep.Options{ExactMatch: true, IgnoreCase: true},
			expected: "banana\n",
		},
		{
			name:     "регулярное выражение",
			pattern:  "^[a-c]",
			expected: "apple\nbanana\ncherry\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runGrep(t, input, tt.pattern, tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGrepContext(t *testing.T) {
	input := "1\n2\n3\n4\n5\n6\n7\n8\n9\n"

	tests := []struct {
		name     string
		pattern  string
		opts     grep.Options
		expected string
	}{
		{
			name:     "-A 2 контекст после",
			pattern:  "3",
			opts:     grep.Options{AContext: 2},
			expected: "3\n4\n5\n",
		},
		{
			name:     "-B 2 контекст до",
			pattern:  "5",
			opts:     grep.Options{BContext: 2},
			expected: "3\n4\n5\n",
		},
		{
			name:     "-C 1 контекст вокруг",
			pattern:  "5",
			opts:     grep.Options{AContext: 1, BContext: 1},
			expected: "4\n5\n6\n",
		},
		{
			name:     "-A 2 в конце файла",
			pattern:  "8",
			opts:     grep.Options{AContext: 5},
			expected: "8\n9\n",
		},
		{
			name:     "-B 2 в начале файла",
			pattern:  "1",
			opts:     grep.Options{BContext: 5},
			expected: "1\n",
		},
		{
			name:     "пересекающийся контекст объединяется",
			pattern:  "[35]",
			opts:     grep.Options{AContext: 1, BContext: 1},
			expected: "2\n3\n4\n5\n6\n",
		},
		{
			name:     "без контекста нет разделителя",
			pattern:  "[19]",
			expected: "1\n9\n",
		},
		{
			name:     "разделитель между группами с контекстом",
			pattern:  "[19]",
			opts:     grep.Options{AContext: 1},
			expected: "1\n2\n--\n9\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runGrep(t, input, tt.pattern, tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGrepContextWithLineNumbers(t *testing.T) {
	input := "aaa\nbbb\nccc\nddd\neee\n"

	result := runGrep(t, input, "ccc", grep.Options{
		BContext:     1,
		AContext:     1,
		StringNumber: true,
	})

	// совпавшая строка с ":", контекстные с "-"
	assert.Equal(t, "2-bbb\n3:ccc\n4-ddd\n", result)
}

func TestGrepInvalidRegex(t *testing.T) {
	var buf bytes.Buffer
	err := grep.Grep(strings.NewReader("test"), &buf, "[invalid", grep.Options{})
	assert.Error(t, err)
}

func TestGrepEmptyInput(t *testing.T) {
	result := runGrep(t, "", "pattern", grep.Options{})
	assert.Equal(t, "", result)
}

func TestGrepCountWithInvert(t *testing.T) {
	input := "a\nb\nc\n"
	result := runGrep(t, input, "a", grep.Options{CountSame: true, InvertFilter: true})
	assert.Equal(t, "2\n", result)
}

func TestGrepFromFile(t *testing.T) {
	inputData, err := os.ReadFile("testdata/input.txt")
	require.NoError(t, err)
	input := string(inputData)

	tests := []struct {
		name     string
		pattern  string
		opts     grep.Options
		wantFile string
	}{
		{
			name:     "простой поиск fmt",
			pattern:  "fmt",
			wantFile: "testdata/grep_simple.txt",
		},
		{
			name:     "-n номера строк",
			pattern:  "fmt",
			opts:     grep.Options{StringNumber: true},
			wantFile: "testdata/grep_n.txt",
		},
		{
			name:     "-i игнор регистра",
			pattern:  "PRINTLN",
			opts:     grep.Options{IgnoreCase: true},
			wantFile: "testdata/grep_i.txt",
		},
		{
			name:     "-v инверсия",
			pattern:  "fmt",
			opts:     grep.Options{InvertFilter: true},
			wantFile: "testdata/grep_v.txt",
		},
		{
			name:     "-c подсчёт",
			pattern:  "fmt",
			opts:     grep.Options{CountSame: true},
			wantFile: "testdata/grep_c.txt",
		},
		{
			name:     "-F точное совпадение",
			pattern:  "fmt.Println",
			opts:     grep.Options{ExactMatch: true},
			wantFile: "testdata/grep_F.txt",
		},
		{
			name:     "-A 1 контекст после",
			pattern:  "func ",
			opts:     grep.Options{AContext: 1},
			wantFile: "testdata/grep_A1.txt",
		},
		{
			name:     "-B 1 контекст до",
			pattern:  "func ",
			opts:     grep.Options{BContext: 1},
			wantFile: "testdata/grep_B1.txt",
		},
		{
			name:     "-C 1 контекст вокруг",
			pattern:  "func ",
			opts:     grep.Options{AContext: 1, BContext: 1},
			wantFile: "testdata/grep_C1.txt",
		},
		{
			name:     "-n -i -C 1 комбинация",
			pattern:  "println",
			opts:     grep.Options{StringNumber: true, IgnoreCase: true, AContext: 1, BContext: 1},
			wantFile: "testdata/grep_niC1.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			result := runGrep(t, input, tt.pattern, tt.opts)
			assert.Equal(t, string(want), result)
		})
	}
}

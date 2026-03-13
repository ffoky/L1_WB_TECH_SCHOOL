package sorting_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"L2.10/sorting"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     sorting.Options
		expected string
		wantErr  bool
	}{
		{
			name:     "happy лексикографическая",
			input:    "banana\napple\ncherry\n",
			opts:     sorting.Options{},
			expected: "apple\nbanana\ncherry\n",
		},
		{
			name:     "happy пустой ввод",
			input:    "",
			opts:     sorting.Options{},
			expected: "",
		},
		{
			name:     "happy одна строка",
			input:    "apple\n",
			opts:     sorting.Options{},
			expected: "apple\n",
		},

		{
			name:     "happy -n",
			input:    "10\n9\n100\n2\n",
			opts:     sorting.Options{Numeric: true},
			expected: "2\n9\n10\n100\n",
		},
		{
			name:     "sad -n",
			input:    "foo\n2\nbar\n1\n",
			opts:     sorting.Options{Numeric: true},
			expected: "foo\nbar\n1\n2\n",
		},

		{
			name:     "happy -r",
			input:    "banana\napple\ncherry\n",
			opts:     sorting.Options{Reverse: true},
			expected: "cherry\nbanana\napple\n",
		},

		{
			name:     "happy -u есть дубликаты",
			input:    "banana\napple\nbanana\ncherry\napple\n",
			opts:     sorting.Options{Unique: true},
			expected: "apple\nbanana\ncherry\n",
		},
		{
			name:     "happy -u без дубликатов",
			input:    "apple\nbanana\ncherry\n",
			opts:     sorting.Options{Unique: true},
			expected: "apple\nbanana\ncherry\n",
		},

		{
			name:     "happy -M",
			input:    "Mar\nJan\nDec\nFeb\n",
			opts:     sorting.Options{Month: true},
			expected: "Jan\nFeb\nMar\nDec\n",
		},
		{
			name:     "sad -M неизвестный месяц",
			input:    "Mar\nFoo\nJan\n",
			opts:     sorting.Options{Month: true},
			expected: "Foo\nJan\nMar\n",
		},

		{
			name:     "happy -h",
			input:    "1G\n500M\n2K\n100\n",
			opts:     sorting.Options{HumanReadableSize: true},
			expected: "100\n2K\n500M\n1G\n",
		},
		{
			name:     "sad -h",
			input:    "1G\nfoo\n2K\n",
			opts:     sorting.Options{HumanReadableSize: true},
			expected: "foo\n2K\n1G\n",
		},
		{
			name:     "happy -k 2",
			input:    "banana\t5\napple\t2\ncherry\t1\n",
			opts:     sorting.Options{Column: 2},
			expected: "cherry\t1\napple\t2\nbanana\t5\n",
		},
		{
			name:     "sad -k 3, но столбца не существует",
			input:    "banana\t5\napple\ncherry\t1\n",
			opts:     sorting.Options{Column: 3},
			expected: "banana\t5\napple\ncherry\t1\n", // порядок не меняется — все поля пустые
		},
		{
			name:     "happy -b",
			input:    "banana  \napple\ncherry\t\n",
			opts:     sorting.Options{TrimTailSpace: true},
			expected: "apple\nbanana  \ncherry\t\n",
		},

		{
			name:    "happy -c отсортировано",
			input:   "apple\nbanana\ncherry\n",
			opts:    sorting.Options{IsSorted: true},
			wantErr: false,
		},
		{
			name:    "sad -c не отсортировано",
			input:   "banana\napple\ncherry\n",
			opts:    sorting.Options{IsSorted: true},
			wantErr: true,
		},

		{
			name:     "happy -n -r",
			input:    "10\n9\n100\n2\n",
			opts:     sorting.Options{Numeric: true, Reverse: true},
			expected: "100\n10\n9\n2\n",
		},
		{
			name:     "happy -k 2 -n",
			input:    "banana\t10\napple\t9\ncherry\t100\n",
			opts:     sorting.Options{Column: 2, Numeric: true},
			expected: "apple\t9\nbanana\t10\ncherry\t100\n",
		},
		{
			name:     "happy -n -u",
			input:    "2\n1\n2\n3\n1\n",
			opts:     sorting.Options{Numeric: true, Unique: true},
			expected: "1\n2\n3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var buf bytes.Buffer

			err := sorting.Sort(reader, &buf, tt.opts)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestSortFromFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		opts     sorting.Options
	}{
		{
			name:     "happy -n из файла",
			input:    "testdata/n.txt",
			expected: "testdata/n_want.txt",
			opts:     sorting.Options{Numeric: true},
		},
		{
			name:     "happy -M из файла",
			input:    "testdata/m.txt",
			expected: "testdata/m_want.txt",
			opts:     sorting.Options{Month: true},
		},
		{
			name:     "happy -h из файла",
			input:    "testdata/h.txt",
			expected: "testdata/h_want.txt",
			opts:     sorting.Options{HumanReadableSize: true},
		},
		{
			name:     "happy crlf Windows line endings",
			input:    "testdata/crlf.txt",
			expected: "testdata/crlf_want.txt",
			opts:     sorting.Options{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.input)
			require.NoError(t, err)
			defer func() {
				err = f.Close()
				require.NoError(t, err)
			}()

			expected, err := os.ReadFile(tt.expected)
			require.NoError(t, err)

			var buf bytes.Buffer
			err = sorting.Sort(f, &buf, tt.opts)
			require.NoError(t, err)

			assert.Equal(t, string(expected), buf.String())
		})
	}
}

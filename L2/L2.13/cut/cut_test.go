package cut

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runCut(t *testing.T, input string, opts Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := Cut(strings.NewReader(input), &buf, opts)
	require.NoError(t, err)
	return buf.String()
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"single field", "2", []int{2}, false},
		{"multiple fields", "1,3,5", []int{1, 3, 5}, false},
		{"range", "2-4", []int{2, 3, 4}, false},
		{"mixed", "1,3-5,8", []int{1, 3, 4, 5, 8}, false},
		{"empty string", "", nil, true},
		{"garbage", "abc", nil, true},
		{"bad range", "1-abc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFields(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCut(t *testing.T) {
	tests := []struct {
		name  string
		input string
		opts  Options
		want  string
	}{
		{
			name:  "tab basic",
			input: "a\tb\tc\n",
			opts:  Options{Fields: []int{1, 3}, Delimiter: "\t"},
			want:  "a\tc\n",
		},
		{
			name:  "custom delimiter",
			input: "a:b:c\n",
			opts:  Options{Fields: []int{2}, Delimiter: ":"},
			want:  "b\n",
		},
		{
			name:  "out of bounds",
			input: "a\tb\n",
			opts:  Options{Fields: []int{1, 5}, Delimiter: "\t"},
			want:  "a\n",
		},
		{
			name:  "-s skips lines without delimiter",
			input: "no delim\na:b\n",
			opts:  Options{Fields: []int{1}, Delimiter: ":", Separated: true},
			want:  "a\n",
		},
		{
			name:  "no -s prints lines as is",
			input: "no delim\na:b\n",
			opts:  Options{Fields: []int{1}, Delimiter: ":"},
			want:  "no delim\na\n",
		},
		{
			name:  "empty input",
			input: "",
			opts:  Options{Fields: []int{1}, Delimiter: "\t"},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runCut(t, tt.input, tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCutFromFile(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		opts      Options
		wantFile  string
	}{
		{
			name:      "tab f1,3",
			inputFile: "testdata/input_tab.txt",
			opts:      Options{Fields: []int{1, 3}, Delimiter: "\t"},
			wantFile:  "testdata/tab_f1_3.txt",
		},
		{
			name:      "tab f2-4",
			inputFile: "testdata/input_tab.txt",
			opts:      Options{Fields: []int{2, 3, 4}, Delimiter: "\t"},
			wantFile:  "testdata/tab_f2_4.txt",
		},
		{
			name:      "tab out of bounds",
			inputFile: "testdata/input_tab.txt",
			opts:      Options{Fields: []int{1, 10}, Delimiter: "\t"},
			wantFile:  "testdata/tab_oob.txt",
		},
		{
			name:      "colon f1",
			inputFile: "testdata/input_colon.txt",
			opts:      Options{Fields: []int{1}, Delimiter: ":"},
			wantFile:  "testdata/colon_f1.txt",
		},
		{
			name:      "colon f1,6-7",
			inputFile: "testdata/input_colon.txt",
			opts:      Options{Fields: []int{1, 6, 7}, Delimiter: ":"},
			wantFile:  "testdata/colon_f1_6_7.txt",
		},
		{
			name:      "-s skip no delimiter",
			inputFile: "testdata/input_mixed.txt",
			opts:      Options{Fields: []int{1}, Delimiter: ":", Separated: true},
			wantFile:  "testdata/mixed_s.txt",
		},
		{
			name:      "no -s prints all",
			inputFile: "testdata/input_mixed.txt",
			opts:      Options{Fields: []int{1}, Delimiter: ":"},
			wantFile:  "testdata/mixed_no_s.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(tt.inputFile)
			require.NoError(t, err)

			want, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			got := runCut(t, string(input), tt.opts)
			assert.Equal(t, string(want), got)
		})
	}
}


package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCutRedirect(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		symbol     string
		wantBefore string
		wantAfter  string
	}{
		{"выходной редирект", "echo hello > file.txt", ">", "echo hello", "file.txt"},
		{"входной редирект", "wc -l < input.txt", "<", "wc -l", "input.txt"},
		{"без редиректа", "echo hello", ">", "echo hello", ""},
		{"пробелы вокруг", "echo hello  >  file.txt", ">", "echo hello", "file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before, after := cutRedirect(tt.input, tt.symbol)
			assert.Equal(t, tt.wantBefore, before)
			assert.Equal(t, tt.wantAfter, after)
		})
	}
}

func TestParseRedirects(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd string
		wantIn  string
		wantOut string
	}{
		{"без редиректов", "echo hello", "echo hello", "", ""},
		{">", "echo hello > out.txt", "echo hello", "", "out.txt"},
		{"<", "wc -l < in.txt", "wc -l", "in.txt", ""},
		{"< и >", "sort < in.txt > out.txt", "sort", "in.txt", "out.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := parseRedirects(tt.input)
			assert.Equal(t, tt.wantCmd, res.cmd)
			assert.Equal(t, tt.wantIn, res.inFile)
			assert.Equal(t, tt.wantOut, res.outFile)
		})
	}
}

func TestSplitCmdsByConds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"echo hello", "echo hello", []string{"echo hello"}},
		{"&&", "echo a && echo b", []string{"echo a", "&&", "echo b"}},
		{"||", "echo a || echo b", []string{"echo a", "||", "echo b"}},
		{
			"echo a && echo b || echo c",
			"echo a && echo b || echo c",
			[]string{"echo a", "&&", "echo b", "||", "echo c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCmdsByConds(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

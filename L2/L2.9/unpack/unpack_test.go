package unpack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"a0bc2d5e", "bccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{`qwe\4\5`, "qwe45", false},
		{`qwe\45`, "qwe44444", false},
		{"a", "a", false},
		{"a3", "aaa", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Unpack(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

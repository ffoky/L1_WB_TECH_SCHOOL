package wget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"без http", "example.com", "http://example.com"},
		{"с http", "http://example.com", "http://example.com"},
		{"с https", "https://example.com", "https://example.com"},
		{"с / на конце", "https://example.com/", "https://example.com"},
		{"без htpps и с / на конце", "example.com/", "http://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validate(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsSameDomain(t *testing.T) {
	tests := []struct {
		name string
		base string
		link string
		want bool
	}{
		{"тот же домен", "https://example.com", "https://example.com/about", true},
		{"другой домен", "https://example.com", "https://other.com/page", false},
		{"тот же с путём", "https://example.com/page", "https://example.com/other", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSameDomain(tt.base, tt.link)
			assert.Equal(t, tt.want, got)
		})
	}
}

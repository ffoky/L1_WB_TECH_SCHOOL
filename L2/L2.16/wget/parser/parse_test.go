package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		html string
		want []string
	}{
		{
			"ссылки href",
			`<html><body><a href="/about">About</a><a href="/contact">Contact</a></body></html>`,
			[]string{"/about", "/contact"},
		},
		{
			"картинки src",
			`<html><body><img src="/logo.png"><script src="/app.js"></script></body></html>`,
			[]string{"/logo.png", "/app.js"},
		},
		{
			"href и src вместе",
			`<html><head><link href="/style.css"></head><body><a href="/page">P</a><img src="/img.jpg"></body></html>`,
			[]string{"/style.css", "/page", "/img.jpg"},
		},
		{
			"без ссылок",
			`<html><body><p>hello</p></body></html>`,
			[]string{},
		},
		{
			"абсолютные ссылки",
			`<html><body><a href="https://example.com/page">Link</a></body></html>`,
			[]string{"https://example.com/page"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, err := Parse(strings.NewReader(tt.html))
			require.NoError(t, err)
			assert.Equal(t, tt.want, links)
		})
	}
}

func TestParseEmpty(t *testing.T) {
	links, err := Parse(strings.NewReader(""))
	require.NoError(t, err)
	assert.Empty(t, links)
}

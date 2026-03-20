package saver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUrlPath(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"без расширения", "https://example.com/about", "example.com/about/index.html"},
		{"корень", "https://example.com/", "example.com/index.html"},
		{"без пути", "https://example.com", "example.com/index.html"},
		{"вложенный путь", "https://example.com/docs/guide", "example.com/docs/guide/index.html"},
		{"с расширением html", "https://example.com/page.html", "example.com/page.html"},
		{"с расширением jpg", "https://example.com/img/photo.jpg", "example.com/img/photo.jpg"},
		{"с расширением css", "https://example.com/css/style.css", "example.com/css/style.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUrlPath(tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(dir)

	body := []byte("<html>hello</html>")
	path, err := Save("https://example.com/page.html", body)
	require.NoError(t, err)
	assert.Equal(t, "example.com/page.html", path)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, body, data)
}

func TestSaveIndex(t *testing.T) {
	dir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(dir)

	body := []byte("<html>index</html>")
	path, err := Save("https://example.com/", body)
	require.NoError(t, err)
	assert.Equal(t, "example.com/index.html", path)
}

package saver

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

const (
	filePerm = 0o644
	dirPerm  = 0o755
)

func Save(URL string, body []byte) (string, error) {
	filePath, err := ParseUrlPath(URL)
	if err != nil {
		return "", fmt.Errorf("save url: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err = os.MkdirAll(dir, dirPerm); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}

	if err = os.WriteFile(filePath, body, filePerm); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return filePath, nil
}

func ParseUrlPath(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("parse url:%s: %w", URL, err)
	}
	path := u.Path
	host := u.Host
	if path == "" || path == "/" {
		path = "/index.html"
	}

	// если путь без расширения, это страница, добавляем /index.html
	if filepath.Ext(path) == "" {
		path = path + "/index.html"
	}

	return filepath.Join(host, path), nil
}

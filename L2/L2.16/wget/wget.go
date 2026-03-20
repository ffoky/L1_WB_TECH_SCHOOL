package wget

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"L2.16/wget/downloader"
	"L2.16/wget/parser"
	"L2.16/wget/robots"
	"L2.16/wget/saver"
)

type Wget struct {
	client  *downloader.Client
	visited map[string]bool
	robots  *robots.Rules
	sem     chan struct{}
	mu      sync.Mutex
}

func NewWget(client *downloader.Client, opts Options) *Wget {
	return &Wget{
		client:  client,
		visited: make(map[string]bool),
		sem:     make(chan struct{}, opts.Workers),
	}
}

func (w *Wget) Run(rawUrl string, depth int) error {
	URL := validate(rawUrl)
	w.mu.Lock()
	if w.visited[URL] {
		w.mu.Unlock()
		return nil
	}
	w.visited[URL] = true
	w.mu.Unlock()

	if w.robots != nil && !w.robots.IsAllowed(URL) {
		_, _ = fmt.Fprintf(os.Stderr, "blocked by robots.txt: %s\n", URL)
		return nil
	}

	body, err := w.client.GetHTML(URL)
	if err != nil {
		return fmt.Errorf("wget: %w", err)
	}

	links, err := parser.Parse(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("wget: %w", err)
	}

	body = replaceLinksToLocal(URL, links, body)

	path, err := saver.Save(URL, body)
	if err != nil {
		return fmt.Errorf("wget: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "saved to %s\n", path)

	if depth <= 0 {
		return nil
	}

	var wg sync.WaitGroup
	for _, link := range w.filterLinks(URL, links) {
		w.sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-w.sem }()
			if err := w.Run(link, depth-1); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}
	wg.Wait()
	return nil
}

func (w *Wget) filterLinks(baseURL string, links []string) []string {
	var result []string
	for _, link := range links {
		full, err := convertToAbsoluteURL(baseURL, link)
		if err != nil {
			continue
		}
		full = validate(full)
		if isSameDomain(baseURL, full) {
			result = append(result, full)
		}
	}
	return result
}

func (w *Wget) LoadRobots(baseURL string) {
	w.robots = robots.Fetch(w.client, validate(baseURL))
}

func convertToAbsoluteURL(baseURL, link string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("convert: %w", err)
	}
	parsed, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("convert: %w", err)
	}
	return base.ResolveReference(parsed).String(), nil
}

func validate(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}
	rawURL = strings.TrimRight(rawURL, "/")

	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	u.Fragment = ""
	return u.String()
}

func isSameDomain(baseURL, linkURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	link, err := url.Parse(linkURL)
	if err != nil {
		return false
	}
	return base.Host == link.Host
}

func replaceLinksToLocal(baseURL string, links []string, body []byte) []byte {
	currentPath, err := saver.ParseUrlPath(baseURL)
	if err != nil {
		return body
	}
	currentDir := filepath.Dir(currentPath)

	var fullLink string
	for _, link := range links {
		fullLink, err = convertToAbsoluteURL(baseURL, link)
		if err != nil {
			continue
		}
		targetPath, err := saver.ParseUrlPath(fullLink)
		if err != nil {
			continue
		}
		relativePath, err := filepath.Rel(currentDir, targetPath)
		if err != nil {
			continue
		}
		body = bytes.ReplaceAll(body, []byte(link), []byte(relativePath))
	}
	return body
}

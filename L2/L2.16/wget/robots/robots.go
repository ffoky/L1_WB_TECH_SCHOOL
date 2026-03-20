package robots

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"L2.16/wget/downloader"
)

type Rules struct {
	disallowed []string
}

func Fetch(client *downloader.Client, baseURL string) *Rules {
	u, err := url.Parse(baseURL)
	if err != nil {
		return &Rules{}
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)
	body, err := client.GetHTML(robotsURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "no robots.txt, skipping")
		return &Rules{}
	}

	return parse(string(body))
}

func parse(content string) *Rules {
	rules := &Rules{}
	scanner := bufio.NewScanner(strings.NewReader(content))

	appliesToAllBots := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if isCommentOrEmpty(line) {
			continue
		}

		if isUserAgentLine(line) {
			appliesToAllBots = extractValue(line) == "*"
			continue
		}

		if appliesToAllBots && isDisallowLine(line) {
			path := extractValue(line)
			if path != "" {
				rules.disallowed = append(rules.disallowed, path)
			}
		}
	}

	return rules
}

func isCommentOrEmpty(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func isUserAgentLine(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), "user-agent:")
}

func isDisallowLine(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), "disallow:")
}

func extractValue(line string) string {
	return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
}

func (r *Rules) IsAllowed(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	for _, path := range r.disallowed {
		if strings.HasPrefix(u.Path, path) {
			return false
		}
	}
	return true
}

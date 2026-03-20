package parser

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func Parse(body io.Reader) ([]string, error) {
	node, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}
	links := make([]string, 0)
	collectLinks(node, &links)
	return links, nil
}

func collectLinks(n *html.Node, links *[]string) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "href" || attr.Key == "src" {
				*links = append(*links, attr.Val)
			}
		}
	}

	collectLinks(n.FirstChild, links)
	collectLinks(n.NextSibling, links)
}

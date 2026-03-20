package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		&http.Client{Timeout: timeout},
	}
}

func (c *Client) GetHTML(url string) ([]byte, error) {
	response, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download from %s: %w", url, err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}(response.Body)

	return io.ReadAll(response.Body)
}

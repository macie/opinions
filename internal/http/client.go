// Package http implements cancellable REST requests with custom User-Agent.
package http

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

// Selected HTTP status codes as registered with IANA. See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusOK = 200
)

// A Client is an HTTP client.
type Client struct {
	defaultClient http.Client
	buildOS       string
	AppVersion    string
}

// Get issues a GET to the specified URL with given context and custom
// User-Agent. It is a replacement for net/http module Get function.
func (c Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent())

	return c.defaultClient.Do(req)
}

// UserAgent constructs User-Agent string in format:
// `opinions/<version_number> (<os>; +https://github.com/macie/opinions)`.
func (c Client) UserAgent() string {
	version := c.AppVersion
	if version == "" {
		version = "local-dev"
	}
	os := c.buildOS
	if os == "" {
		os = runtime.GOOS
	}

	return fmt.Sprintf("opinions/%s (%s; +https://github.com/macie/opinions)", version, os)
}

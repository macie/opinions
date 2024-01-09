// Package opinions helps finding links to discussions about given topics/URLs
// on social news websites.
package opinions

import (
	"context"
	"fmt"
	"net/http"
)

// Discussion is a representation of discussion inside social media service.
type Discussion struct {
	Service  string
	URL      string
	Title    string
	Source   string
	Comments int
}

// String returns string representation of discussion metadata.
func (d Discussion) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s", d.Service, d.URL, d.Title, d.Source)
}

// GetRequester is an interface for sending HTTP requests with custom User-Agent.
type GetRequester interface {
	// Get issues a GET to the specified URL and follows redirects.
	//
	// When err is nil, resp always contains a non-nil resp.Body. Caller should close resp.Body when done reading from it.
	//
	// It is modelled after net/http module Get function, see:
	// https://pkg.go.dev/net/http#Client.Get
	Get(ctx context.Context, url string) (resp *http.Response, err error)
}

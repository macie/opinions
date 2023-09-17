// Package http implements cancellable REST requests with custom User-Agent.
package http

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

// DefaultClient is the default Client and is used by Get.
var DefaultClient = &http.Client{}

// Get issues a GET to the specified URL with given context and custom
// User-Agent. It is a replacement for net/http module Get function.
func Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// FIXME: use real version number
	req.Header.Set("User-Agent", fmt.Sprintf("opinions/%s (%s; +https://github.com/macie/opinions)", "dev", runtime.GOOS))

	return DefaultClient.Do(req)
}

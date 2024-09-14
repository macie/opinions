package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/macie/opinions/internal/ensure"
)

func ExampleGet() {
	type HttpbinResponse struct {
		Headers struct {
			UserAgent string `json:"User-Agent"`
		} `json:"headers"`
	}
	var response HttpbinResponse

	c := Client{buildOS: "anyOS", AppVersion: "0.0.0-local"}
	URL := "https://httpbin.org/get"
	raw := ensure.MustReturn(c.Get(context.TODO(), URL))
	defer func() {
		_ = raw.Body.Close() // ignore closing errors during testing
	}()

	ensure.Must(json.NewDecoder(raw.Body).Decode(&response))

	fmt.Println(response.Headers.UserAgent)
	// Output:
	// opinions/0.0.0-local (anyOS; +https://github.com/macie/opinions)
}

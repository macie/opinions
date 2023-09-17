package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/macie/opinions/testing"
)

func ExampleGet() {
	type HttpbinResponse struct {
		Headers struct {
			UserAgent string `json:"User-Agent"`
		} `json:"headers"`
	}
	var response HttpbinResponse

	URL := "https://httpbin.org/get"
	raw := testing.MustReturn(Get(context.TODO(), URL))
	defer raw.Body.Close()

	body := testing.MustReturn(io.ReadAll(raw.Body))
	testing.Must(json.Unmarshal(body, &response))

	fmt.Println(response.Headers.UserAgent)
	// Output:
	// opinions/dev (openbsd; +https://github.com/macie/opinions)
}

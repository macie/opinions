package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// must safeguards happy path of function call.
func must[T any](val T, err error) T {
	if err != nil {
		log.Fatalf("function call returns error %v", err)
	}
	return val
}

func ExampleGet() {
	type HttpbinResponse struct {
		Headers struct {
			UserAgent string `json:"User-Agent"`
		} `json:"headers"`
	}
	var response HttpbinResponse

	URL := "https://httpbin.org/get"
	raw := must(Get(context.TODO(), URL))
	defer raw.Body.Close()

	body := must(io.ReadAll(raw.Body))
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("function call returns error %v", err)
	}

	fmt.Println(response.Headers.UserAgent)
	// Output:
	// opinions/dev (openbsd; +https://github.com/macie/opinions)
}

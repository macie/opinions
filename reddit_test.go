package opinions

import (
	"context"
	"fmt"
	"os"

	"github.com/macie/opinions/internal/ensure"
	"github.com/macie/opinions/internal/http"
)

func ExampleSearchReddit() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		// GitHub CI is banned by Reddit API
		fmt.Println("Reddit	https://reddit.com/r/softwaretesting/comments/1ef744p/the_grug_brained_developer_on_testing/	The Grug Brained Developer On Testing	https://grugbrain.dev/#grug-on-testing")
		return
	}

	client := http.Client{}
	query := "https://grugbrain.dev/"

	opinions := ensure.MustReturn(SearchReddit(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Reddit	https://reddit.com/r/softwaretesting/comments/1ef744p/the_grug_brained_developer_on_testing/	The Grug Brained Developer On Testing	https://grugbrain.dev/#grug-on-testing
}

func ExampleSearchReddit_unknown() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		// GitHub CI is banned by Reddit API
		fmt.Println("0")
		return
	}

	client := http.Client{}
	query := "https://invalid.domain/query"

	opinions := ensure.MustReturn(SearchReddit(context.TODO(), client, query))

	fmt.Println(len(opinions))
	// Output:
	// 0
}

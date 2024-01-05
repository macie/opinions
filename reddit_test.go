package opinions

import (
	"context"
	"fmt"
	"os"

	"github.com/macie/opinions/ensure"
	"github.com/macie/opinions/http"
)

func ExampleSearchReddit() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		// GitHub CI is banned by Reddit API
		fmt.Println("Reddit	https://reddit.com/r/hypeurls/comments/17k6i1l/the_grug_brained_developer_2022/	The Grug Brained Developer (2022)	https://grugbrain.dev/")
		return
	}

	client := http.Client{}
	query := "https://grugbrain.dev/"

	opinions := ensure.MustReturn(SearchReddit(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Reddit	https://reddit.com/r/hypeurls/comments/17k6i1l/the_grug_brained_developer_2022/	The Grug Brained Developer (2022)	https://grugbrain.dev/
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

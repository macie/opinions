package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/ensure"
	"github.com/macie/opinions/http"
)

func ExampleSearchHackerNews() {
	client := http.Client{}
	query := "grugbrain.dev"

	opinions := ensure.MustReturn(SearchHackerNews(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/
}

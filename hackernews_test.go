package main

import (
	"context"
	"fmt"

	"github.com/macie/opinions/ensure"
)

func ExampleSearchHackerNews() {
	query := "grugbrain.dev"
	opinions := ensure.MustReturn(SearchHackerNews(context.TODO(), query))

	fmt.Println(opinions[0])
	// Output:
	// Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/
}

package main

import (
	"context"
	"fmt"

	"github.com/macie/opinions/testing"
)

func ExampleSearchHackerNews() {
	query := "grugbrain.dev"
	opinions := testing.MustReturn(SearchHackerNews(context.TODO(), query))

	fmt.Println(opinions[0].String())
	// Output:
	// Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/
}

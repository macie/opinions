package main

import (
	"context"
	"testing"
)

func TestSearchHackerNews(t *testing.T) {
	query := "grugbrain.dev"
	want := "Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/"
	opinions, err := SearchHackerNews(context.TODO(), query)
	if err != nil {
		t.Errorf("SearchHackerNews(ctx, \"%v\") returns error %v", query, err)
	}

	get := ""
	if len(opinions) > 0 {
		get = opinions[0].String()
	}
	if get != want {
		t.Errorf("SearchHackerNews(ctx, \"%v\") got \"%v\" want \"%v\"", query, get, want)
	}
}

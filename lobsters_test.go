package main

import (
	"context"
	"fmt"

	"github.com/macie/opinions/testing"
)

func ExampleSearchLobsters() {
	query := "https://grugbrain.dev"
	opinions := testing.MustReturn(SearchLobsters(context.TODO(), query))

	fmt.Println(opinions[0].String())
	// Output:
	// Lobsters	https://lobste.rs/s/ifaar4/grug_brained_developer	The Grug Brained Developer	https://grugbrain.dev/
}

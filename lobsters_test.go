package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/internal/ensure"
	"github.com/macie/opinions/internal/http"
)

func ExampleSearchLobsters() {
	client := http.Client{}
	query := "https://grugbrain.dev"

	opinions := ensure.MustReturn(SearchLobsters(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Lobsters	https://lobste.rs/s/ifaar4/grug_brained_developer	The Grug Brained Developer	https://grugbrain.dev/
}

func ExampleSearchLobsters_unknown() {
	client := http.Client{}
	query := "https://invalid.domain/query"

	opinions := ensure.MustReturn(SearchLobsters(context.TODO(), client, query))

	fmt.Println(len(opinions))
	// Output:
	// 0
}

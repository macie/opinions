package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/ensure"
)

func ExampleSearchLobsters() {
	query := "https://grugbrain.dev"
	opinions := ensure.MustReturn(SearchLobsters(context.TODO(), query))

	fmt.Println(opinions[0])
	// Output:
	// Lobsters	https://lobste.rs/s/ifaar4/grug_brained_developer	The Grug Brained Developer	https://grugbrain.dev/
}

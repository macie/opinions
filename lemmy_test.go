package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/internal/ensure"
	"github.com/macie/opinions/internal/http"
)

func ExampleSearchLemmy() {
	client := http.Client{}
	query := "https://grugbrain.dev/"

	opinions := ensure.MustReturn(SearchLemmy(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Lemmy	https://lemmy.world/post/24591938	The Grug Brained Developer - A manifesto	https://grugbrain.dev/
}

func ExampleSearchLemmy_unknown() {
	client := http.Client{}
	query := "https://invalid.domain/query"

	opinions := ensure.MustReturn(SearchLemmy(context.TODO(), client, query))

	fmt.Println(len(opinions))
	// Output:
	// 0
}

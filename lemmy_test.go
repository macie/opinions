package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/ensure"
	"github.com/macie/opinions/http"
)

func ExampleSearchLemmy() {
	client := http.Client{}
	query := "https://grugbrain.dev/"

	opinions := ensure.MustReturn(SearchLemmy(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Lemmy	https://lemmy.world/post/7563451	The Grug Brained Developer (2022)	https://grugbrain.dev/
}

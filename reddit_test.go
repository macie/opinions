package opinions

import (
	"context"
	"fmt"

	"github.com/macie/opinions/ensure"
	"github.com/macie/opinions/http"
)

func ExampleSearchReddit() {
	client := http.Client{}
	query := "https://grugbrain.dev/"

	opinions := ensure.MustReturn(SearchReddit(context.TODO(), client, query))

	fmt.Println(opinions[0])
	// Output:
	// Reddit	https://reddit.com/r/hypeurls/comments/17k6i1l/the_grug_brained_developer_2022/	The Grug Brained Developer (2022)	https://grugbrain.dev/
}

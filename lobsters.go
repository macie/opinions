package opinions

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/macie/opinions/html"
	"github.com/macie/opinions/http"
)

// SearchLobsters query Lobsters search engine for given prompt. It returns list
// of discussions sorted by relevance.
func SearchLobsters(ctx context.Context, client http.Client, query string) ([]Discussion, error) {
	discussions := make([]Discussion, 0)
	searchURL := "https://lobste.rs/search?what=stories&order=relevance&q="

	// queries with URL must be prefixed for more accurate results
	_, err := url.Parse(query)
	if err != nil {
		query = "domain:" + query
	}

	r, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return discussions, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return discussions, fmt.Errorf("cannot search Lobsters: `GET %s` responded with status code %d", r.Request.URL, r.StatusCode)
	}

	body, err := html.Parse(r.Body)
	if err != nil {
		return discussions, err
	}

	for _, listItem := range html.FindAll(body, "ol > li") {
		srcNode := html.First(listItem, ".link > a")
		commentNode := html.First(listItem, ".mobile_comments")
		comments, err := strconv.Atoi(html.Text(html.First(commentNode, "span")))
		if err != nil {
			comments = 0
		}

		discussions = append(discussions, Discussion{
			Service:  "Lobsters",
			URL:      "https://lobste.rs" + html.Attr(commentNode, "href"),
			Title:    html.Text(srcNode),
			Source:   html.Attr(srcNode, "href"),
			Comments: comments,
		})
	}

	return discussions, nil
}

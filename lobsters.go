package opinions

import (
	"context"
	"net/url"
	"strconv"

	"github.com/macie/opinions/html"
	"github.com/macie/opinions/http"
)

// SearchLobsters query lobster search engine for given prompt. It returns list
// of stories which contains comments sorted by relevance.
func SearchLobsters(ctx context.Context, client http.Client, query string) ([]Discussion, error) {
	discussions := make([]Discussion, 0)
	searchURL := "https://lobste.rs/search?what=stories&order=relevance&q="

	// queries with URL must be prefixed for more accurate results
	_, err := url.Parse(query)
	if err != nil {
		query = "domain:" + query
	}

	raw, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return discussions, err
	}
	defer raw.Body.Close()

	doc, err := html.Parse(raw.Body)
	if err != nil {
		return discussions, err
	}

	for _, d := range html.FindAll(doc, "ol > li") {
		srcNode := html.First(d, ".link > a")
		commentNode := html.First(d, ".mobile_comments")
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

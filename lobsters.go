package opinions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/macie/opinions/internal/html"
)

// SearchLobsters query Lobsters search engine for given prompt. It returns list
// of discussions sorted by relevance.
func SearchLobsters(ctx context.Context, client GetRequester, query string) (discussions []Discussion, err error) {
	searchURL := "https://lobste.rs/search?what=stories&order=relevance&q="

	// queries with URL must be prefixed for more accurate results
	_, err = url.Parse(query)
	if err != nil {
		query = "domain:" + query
	}

	r, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return noDiscussions, err
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if r.StatusCode != http.StatusOK {
		return noDiscussions, fmt.Errorf("`GET %s` responded with unexpected status code %d", r.Request.URL, r.StatusCode)
	}

	body, err := html.Parse(r.Body)
	if err != nil {
		return noDiscussions, err
	}
	items := html.FindAll(body, "ol > li")

	discussions = make([]Discussion, 0, len(items))
	for _, listItem := range items {
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

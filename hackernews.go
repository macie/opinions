package opinions

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"time"

	"github.com/macie/opinions/http"
)

// HackerNewsResponse represents some interesting fields of response from
// HN Search API.
//
// See: https://hn.algolia.com/api
type HackerNewsResponse struct {
	Hits []struct {
		CreatedAt   time.Time `json:"created_at"`
		Title       string    `json:"title"`
		URL         string    `json:"url"`
		NumComments int       `json:"num_comments"`
		ObjectID    string    `json:"objectID"`
	} `json:"hits"`
}

// SearchHackerNews query HN Search API for given prompt. It returns list of
// discussions sorted by relevance, then popularity, then number of comments.
//
// See: https://hn.algolia.com/api
func SearchHackerNews(ctx context.Context, client http.Client, query string) ([]Discussion, error) {
	searchURL := "http://hn.algolia.com/api/v1/search?tags=story&query="
	discussions := make([]Discussion, 0)

	r, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return discussions, err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return discussions, err
	}

	var response HackerNewsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return discussions, err
	}

	for _, entry := range response.Hits {
		discussions = append(discussions, Discussion{
			Service:  "Hacker News",
			URL:      "https://news.ycombinator.com/item?id=" + url.QueryEscape(entry.ObjectID),
			Title:    entry.Title,
			Source:   entry.URL,
			Comments: entry.NumComments,
		})
	}

	return discussions, nil
}

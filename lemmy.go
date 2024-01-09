package opinions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/macie/opinions/internal/http"
)

// LemmyResponse represents some interesting fields of response from
// Lemmy Search API.
//
// See: https://join-lemmy.org/api/interfaces/SearchResponse.html
type LemmyResponse struct {
	Posts []struct {
		Post struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"post"`
		Counts struct {
			Comments int `json:"comments"`
		} `json:"counts"`
	} `json:"posts"`
}

// SearchLemmy query Lemmy Search API for given prompt. It returns list
// of discussions sorted by rank based on the score and time of the latest
// comment, with decay over time.
//
// See: https://join-lemmy.org/docs/users/03-votes-and-ranking.html
func SearchLemmy(ctx context.Context, client http.Client, query string) ([]Discussion, error) {
	discussions := make([]Discussion, 0)
	searchURL := "https://lemmy.world/api/v3/search?listingType=All&sort=Active"

	_, err := url.Parse(query)
	switch {
	case err == nil:
		searchURL += "&type_=Url&q="
		// URL query must contain trailing slash for more accurate results
		if !strings.HasSuffix(query, "/") {
			query += "/"
		}
	default:
		searchURL += "&type_=Posts&q="
	}

	r, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return discussions, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return discussions, fmt.Errorf("cannot search Lemmy: `GET %s` responded with status code %d", r.Request.URL, r.StatusCode)
	}

	var response LemmyResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return discussions, err
	}

	for _, entry := range response.Posts {
		discussions = append(discussions, Discussion{
			Service:  "Lemmy",
			URL:      fmt.Sprintf("https://lemmy.world/post/%d", entry.Post.ID),
			Title:    entry.Post.Name,
			Source:   entry.Post.URL,
			Comments: entry.Counts.Comments,
		})
	}

	return discussions, nil
}

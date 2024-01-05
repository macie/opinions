package opinions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/macie/opinions/html"
	"github.com/macie/opinions/http"
)

// RedditResponse represents some interesting fields of response from Reddit API.
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				ID          string `json:"permalink"`
				Title       string `json:"title"`
				URL         string `json:"url"`
				NumComments int    `json:"num_comments"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// UnmarshalJSON deserialize inconsistent JSON responses to RedditResponse.
// Reddit returns empty object ("{}") when there are no search results.
func (r *RedditResponse) UnmarshalJSON(b []byte) error {
	isEmptyResponse := len(b) == 4 && string(b) == "\"{}\""
	if isEmptyResponse {
		return nil
	}

	// new type prevents recursive calls to RedditResponse.UnmarshalJSON()
	type resp *RedditResponse
	return json.Unmarshal(b, resp(r))
}

// SearchReddit searches Reddit for given query and returns list of discussions
// sorted by relevance.
//
// See: https://www.reddit.com/dev/api#GET_search
func SearchReddit(ctx context.Context, client http.Client, query string) ([]Discussion, error) {
	discussions := make([]Discussion, 0)
	searchURL := "https://www.reddit.com/search.json?sort=relevance&t=all&q="

	r, err := client.Get(ctx, searchURL+url.QueryEscape(query))
	if err != nil {
		return discussions, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		if r.Header.Get("X-Ratelimit-Remaining") == "0" { // https://support.reddithelp.com/hc/en-us/articles/16160319875092-Reddit-Data-API-Wiki
			return discussions, fmt.Errorf("cannot search Reddit: too many requests. Wait %s seconds", r.Header.Get("X-Ratelimit-Reset"))
		}

		if r.StatusCode == 403 {
			var details string
			body, err := html.Parse(r.Body)
			if err != nil {
				details = ""
			}
			details = html.Text(html.First(body, "p"))

			return discussions, fmt.Errorf("cannot search Reddit: your IP address seems to be banned by Reddit: `GET %s` responded with `%s`", r.Request.URL, details)
		}

		return discussions, fmt.Errorf("cannot search Reddit: `GET %s` responded with status code %d", r.Request.URL, r.StatusCode)
	}

	var response RedditResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return discussions, err
	}

	for _, entry := range response.Data.Children {
		discussions = append(discussions, Discussion{
			Service:  "Reddit",
			URL:      "https://reddit.com" + entry.Data.ID,
			Title:    entry.Data.Title,
			Source:   entry.Data.URL,
			Comments: entry.Data.NumComments,
		})
	}

	return discussions, nil
}

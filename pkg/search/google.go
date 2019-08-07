package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"net/url"
)

const (
	googleDefaultCX    = "002017775112634544492:7y5bpl2sn78"
	googleSearchAPIURL = "https://www.googleapis.com/customsearch/v1"
)

// Google implements Searcher for Google custom search engines.
type Google struct {
	APIKey string // api key enabled for the custom search api
	CX     string // custom search engine id. If empty a default is used.
}

type googleCustomSearchResult struct {
	Items []*struct {
		Link string
	} `json:"items"`
}

// Search implements Searcher for google.
// At most 100 results can be returned (limit by Google).
func (s *Google) Search(ctx context.Context, query string) <-chan Result {
	const itemCount = 10

	cx := s.CX
	if cx == "" {
		cx = googleDefaultCX
	}

	searchURL := fmt.Sprintf(googleSearchAPIURL+
		"?q=%s"+
		"&key=%s"+
		"&cx=%s"+
		"&fields=items(link)"+
		"&num=%d&start=%%d",
		url.QueryEscape(query), s.APIKey, cx, itemCount)

	urlChan := make(chan Result, itemCount)

	go func() {
		defer close(urlChan)

		for i := 1; i <= 100; i += itemCount {
			req := request.NewWithContext(ctx, fmt.Sprintf(searchURL, i))
			resp, err := req.Response()
			if err != nil {
				return
			}

			var data googleCustomSearchResult

			err = json.NewDecoder(resp.Body).Decode(&data)
			_ = resp.Body.Close()

			if err != nil {
				return
			}

			for _, item := range data.Items {
				if item.Link == "" {
					continue
				}

				select {
				case <-ctx.Done():
					return
				case urlChan <- Result{URL: item.Link}:
				}
			}
		}
	}()

	return urlChan
}

package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"log"
	"net/url"
)

const (
	googleDefaultCX    = "002017775112634544492:7y5bpl2sn78"
	googleSearchAPIURL = "https://www.googleapis.com/customsearch/v1"
)

type GoogleSearcher struct {
	APIKey string
	CX     string
}

type googleCustomSearchResult struct {
	Items []*struct {
		Link string
	} `json:"items"`
}

func (s *GoogleSearcher) Search(ctx context.Context, query string) <-chan Result {
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
				log.Print(err)
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

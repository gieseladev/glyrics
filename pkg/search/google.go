package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"log"
	"net/url"
	"strconv"
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

	searchURL, err := url.Parse(fmt.Sprintf(googleSearchAPIURL+
		"?q=%s"+
		"&key=%s"+
		"&cx=%s"+
		"&fields=items(link)"+
		"&num=%d",
		url.QueryEscape(query), s.APIKey, cx, itemCount))
	if err != nil {
		panic(err)
	}

	urlChan := make(chan Result, itemCount)

	go func() {
		defer close(urlChan)

		for i := 1; i <= 100; i += itemCount {
			searchURL.Query().Set("start", strconv.Itoa(i))
			req := request.NewWithContext(ctx, searchURL)
			body, err := req.Body()
			if err != nil {
				log.Print("glyrics/google: couldn't get body")
				return
			}

			var data googleCustomSearchResult
			err = json.NewDecoder(body).Decode(&data)

			_ = req.Close()

			if err != nil {
				log.Print("glyrics/google: couldn't decode search response")
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

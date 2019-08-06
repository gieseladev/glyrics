package search

import (
	"encoding/json"
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/requests"
	"net/url"
)

const (
	searchAPIURL = "https://www.googleapis.com/customsearch/v1"
)

type googleCustomSearchResult struct {
	Items []*struct {
		Link string
	} `json:"items"`
}

// GoogleSearch performs a google custom search with a search engine strongly
// optimised for lyrics. It returns a channel which yields all urls of the search
// results in order and a channel which can be used to stop the search.
// If not stopped the search channel will yield 100 search results.
func GoogleSearch(query string, apiKey string) (<-chan string, chan<- struct{}) {
	itemCount := 10

	searchURL := fmt.Sprintf(searchAPIURL+
		"?q=%s"+
		"&key=%s"+
		"&cx=002017775112634544492:7y5bpl2sn78"+
		"&fields=items(link)"+
		"&num=%d&start=%%d",
		url.QueryEscape(query), apiKey, itemCount)

	urlChan := make(chan string, itemCount)
	stopSignal := make(chan struct{})

	go func() {
		defer close(urlChan)

	SearchLoop:
		for i := 1; i <= 100; i += itemCount {
			req := requests.NewRequest(fmt.Sprintf(searchURL, i))
			resp, err := req.Response()
			if err != nil {
				panic(err)
			}

			var data googleCustomSearchResult

			err = json.NewDecoder(resp.Body).Decode(&data)
			_ = resp.Body.Close()

			if err != nil {
				break SearchLoop
			}

			for _, item := range data.Items {
				if item.Link == "" {
					continue
				}

				select {
				case <-stopSignal:
					break SearchLoop
				case urlChan <- item.Link:
				}
			}
		}
	}()

	return urlChan, stopSignal
}

package lyricsfinder

import (
	"encoding/json"
	"fmt"
	"net/http"
	url2 "net/url"
)

type googleCustomSearchResult struct {
	Items []*struct {
		Link string
	} `json:"items"`
}

func GoogleSearch(query string, apiKey string) (<-chan string, chan<- bool) {
	itemCount := 10

	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?"+
		"q=%s"+
		"&key=%s"+
		"&cx=002017775112634544492:7y5bpl2sn78"+
		"&fields=items(link)"+
		"&num=%d&start=%%d", url2.QueryEscape(query), apiKey, itemCount)

	urlChan := make(chan string, 3)
	stopSignal := make(chan bool)

	go func() {
		defer close(urlChan)

	SearchLoop:
		for i := 1; i <= 100; i += itemCount {
			// FIXME http.Get uses http.DefaultClient which doesn't have any timeout
			resp, err := http.Get(fmt.Sprintf(url, i))
			if err != nil {
				panic(err)
			}

			var data googleCustomSearchResult

			err = json.NewDecoder(resp.Body).Decode(&data)
			resp.Body.Close()
			if err != nil {
				continue
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

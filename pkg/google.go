package lyricsfinder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type googleCustomSearchResult struct {
	Items []*struct {
		Link string
	} `json:"items"`
}

func GoogleSearch(query string, apiKey string, ch chan string) {
	itemCount := 10

	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?"+
		"q=%s"+
		"&key=%s"+
		"&cx=002017775112634544492:7y5bpl2sn78"+
		"&fields=items(link)"+
		"&num=%%d", query, apiKey)

	for i := 1; i <= 100; i += itemCount {
		// FIXME http.Get uses http.DefaultClient which doesn't have any timeout
		resp, err := http.Get(fmt.Sprintf(url, i))
		if err != nil {
			panic(err)
		}

		var data googleCustomSearchResult

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		for _, item := range data.Items {
			ch <- item.Link
		}

		resp.Body.Close()
	}

	close(ch)
}

package lyricsfinder

import (
	"errors"
	"fmt"
	"github.com/gieseladev/lyricsfinder/pkg/extractors"
	"github.com/gieseladev/lyricsfinder/pkg/models"
)

func ExtractLyricsFromRequest(request models.Request) (*models.Lyrics, error) {
	for _, extractor := range extractors.Extractors {
		if extractor.CanHandle(request) {
			lyrics, err := extractor.ExtractLyrics(request)
			if err != nil {
				continue
			}

			return lyrics, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("No extractor could extract %+v", request))
}

func ExtractLyrics(url string) (*models.Lyrics, error) {
	return ExtractLyricsFromRequest(models.Request{Url: url})
}

func extractLyricsToChannel(url string, ch chan models.Lyrics) {
	lyrics, err := ExtractLyrics(url)
	if err == nil {
		ch <- *lyrics
	}
}

func SearchLyrics(query string, apiKey string, ch chan models.Lyrics) {
	urlChan := make(chan string, 2) // "preload" the next url (speed-up in case of new api request)
	go GoogleSearch(query, apiKey, urlChan)

	for url := range urlChan {
		go extractLyricsToChannel(url, ch)
	}

	close(ch)
}

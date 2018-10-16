package lyricsfinder

import (
	"errors"
	"fmt"
	"github.com/gieseladev/lyricsfindergo/pkg/extractors"
	"github.com/gieseladev/lyricsfindergo/pkg/models"
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

func SearchLyrics(query string, apiKey string, ch chan models.Lyrics, stopSignal chan bool) {
	urlChan := make(chan string, 2) // "preload" the next url (speed-up in case of new api request)
	go GoogleSearch(query, apiKey, urlChan, stopSignal)

	for url := range urlChan {
		extractLyricsToChannel(url, ch)
	}

	close(ch)
}

func SearchFirstLyrics(query string, apiKey string) models.Lyrics {
	lyricsChan := make(chan models.Lyrics)
	stopChan := make(chan bool)
	go SearchLyrics(query, apiKey, lyricsChan, stopChan)

	lyrics := <-lyricsChan
	stopChan <- true

	return lyrics
}

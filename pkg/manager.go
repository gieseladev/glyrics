package lyricsfinder

import (
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

	return nil, fmt.Errorf("no extractor could extract %+v", request)
}

func ExtractLyrics(url string) (*models.Lyrics, error) {
	return ExtractLyricsFromRequest(models.Request{Url: url})
}

func SearchLyrics(query string, apiKey string) (<-chan models.Lyrics, chan<- bool) {
	lyricsChan := make(chan models.Lyrics)
	urlChan, stopChan := GoogleSearch(query, apiKey)

	go func() {
		for url := range urlChan {
			lyrics, err := ExtractLyrics(url)
			if err == nil {
				lyricsChan <- *lyrics
			}
		}

		close(lyricsChan)
	}()

	return lyricsChan, stopChan
}

func SearchFirstLyrics(query string, apiKey string) models.Lyrics {
	lyricsChan, stopChan := SearchLyrics(query, apiKey)

	lyrics := <-lyricsChan
	stopChan <- true

	return lyrics
}

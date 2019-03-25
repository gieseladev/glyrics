// Package gLyrics is a library for extracting lyrics from lyrics websites.
// It also comes with a lyrics search function which uses google custom
// search to provide the most accurate results possible.
package glyrics

import (
	"fmt"
	"github.com/gieseladev/glyrics/pkg/extractors"
	"github.com/gieseladev/glyrics/pkg/models"
)

// ExtractLyricsFromRequest tries to extract lyrics from
// the provided Request.
// It tries all extractors from extractors.Extractors
// and returns the first one that was successful.
// Errors from extracting lyrics are ignored. The only error
// returned by this function is when no extractor was able to
// extract any lyrics.
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

// ExtractLyrics wraps the provided url in a Request and performs
// ExtractLyricsFromRequest.
func ExtractLyrics(url string) (*models.Lyrics, error) {
	return ExtractLyricsFromRequest(models.Request{Url: url})
}

// SearchLyrics uses GoogleSearch to search for lyrics websites based
// on the given query and turns them into Lyrics. Like the GoogleSearch
// method it returns two channels.
func SearchLyrics(query string, apiKey string) (<-chan models.Lyrics, chan<- bool) {
	lyricsChan := make(chan models.Lyrics)
	urlChan, stopChan := GoogleSearch(query, apiKey)

	go func() {
		defer close(lyricsChan)

		for url := range urlChan {
			lyrics, err := ExtractLyrics(url)
			if err == nil {
				lyricsChan <- *lyrics
			}
		}
	}()

	return lyricsChan, stopChan
}

// SearchNLyrics uses SearchLyrics to search for lyrics
// and returns at most the specified amount of Lyrics.
func SearchNLyrics(query, apiKey string, amount int) []models.Lyrics {
	lyrics := make([]models.Lyrics, amount)

	lyricsChan, stopChan := SearchLyrics(query, apiKey)

	for len(lyrics) < amount {
		ly, ok := <-lyricsChan
		if !ok {
			break
		}

		lyrics = append(lyrics, ly)
	}

	stopChan <- true

	return lyrics
}

// SearchFirstLyrics uses SearchLyrics to search for lyrics and
// returns the first result.
func SearchFirstLyrics(query string, apiKey string) models.Lyrics {
	lyricsChan, stopChan := SearchLyrics(query, apiKey)

	lyrics := <-lyricsChan
	stopChan <- true

	return lyrics
}

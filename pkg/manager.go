package lyricsfinder

import (
	"errors"
	"fmt"
)

var Extractors = make([]Extractor, 0)

func RegisterExtractor(extractor Extractor) {
	Extractors = append(Extractors, extractor)
}

func ExtractLyricsFromRequest(request Request) (*Lyrics, error) {
	for _, extractor := range Extractors {
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

func ExtractLyrics(url string) (*Lyrics, error) {
	return ExtractLyricsFromRequest(Request{Url: url})
}

func extractLyricsToChannel(url string, ch chan Lyrics) {
	lyrics, err := ExtractLyrics(url)
	if err == nil {
		ch <- *lyrics
	}
}

func SearchLyrics(query string, apiKey string, ch chan Lyrics) {
	urlChan := make(chan string, 3)
	go GoogleSearch(query, apiKey, urlChan)

	for url := range urlChan {
		go extractLyricsToChannel(url, ch)
	}

	close(ch)
}

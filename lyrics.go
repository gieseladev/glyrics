package glyrics

import (
	"fmt"
	"github.com/gieseladev/glyrics/v3/extractors"
	"github.com/gieseladev/glyrics/v3/pkg/requests"
	"github.com/gieseladev/glyrics/v3/pkg/search"
	"sync"
)

type LyricsInfo = extractors.LyricsInfo
type LyricsOrigin = extractors.LyricsOrigin

// ExtractLyricsFromRequest tries to extract lyrics from the provided Request.
// Errors from extracting lyrics are ignored. The only error
// returned by this function is when no extractor was able to
// extract any lyrics.
func ExtractLyricsFromRequest(req *requests.Request) (*LyricsInfo, error) {
	for _, e := range extractors.GetExtractorsForRequest(req) {
		lyrics, err := e.ExtractLyrics(req)
		if err != nil {
			continue
		}

		return lyrics, nil
	}

	return nil, fmt.Errorf("no extractor could extract from %s", req.Url)
}

// ExtractLyrics wraps the provided url in a Request and performs
// ExtractLyricsFromRequest.
func ExtractLyrics(url string) (*LyricsInfo, error) {
	return ExtractLyricsFromRequest(requests.NewRequest(url))
}

// SearchLyrics uses GoogleSearch to search for lyrics websites based
// on the given query and turns them into LyricsInfo. Like the GoogleSearch
// method it returns two channels.
func SearchLyrics(query string, apiKey string) (<-chan *LyricsInfo, chan<- struct{}) {
	lyricsChan := make(chan *LyricsInfo)
	urlChan, stopChan := search.GoogleSearch(query, apiKey)

	go func() {
		defer close(lyricsChan)

		for url := range urlChan {
			lyrics, err := ExtractLyrics(url)
			if err == nil {
				lyricsChan <- lyrics
			}
		}
	}()

	return lyricsChan, stopChan
}

// SearchNLyrics uses SearchLyrics to search for lyrics
// and returns at most the specified amount of LyricsInfo.
func SearchNLyrics(query, apiKey string, amount int) []LyricsInfo {
	lyrics := make([]LyricsInfo, 0, amount)
	var mut sync.Mutex

	lyricsChan, stopChan := SearchLyrics(query, apiKey)

	for len(lyrics) < amount {
		ly, ok := <-lyricsChan
		if !ok {
			break
		}

		mut.Lock()
		lyrics = append(lyrics, *ly)
		mut.Unlock()
	}

	close(stopChan)

	return lyrics
}

// SearchFirstLyrics uses SearchLyrics to search for lyrics and
// returns the first result.
func SearchFirstLyrics(query string, apiKey string) *LyricsInfo {
	lyricsChan, stopChan := SearchLyrics(query, apiKey)

	lyrics := <-lyricsChan
	close(stopChan)

	return lyrics
}

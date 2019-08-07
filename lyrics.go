/*
Package glyrics provides tools for extracting lyrics.
*/
package glyrics

import (
	"context"
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"github.com/gieseladev/glyrics/v3/pkg/search"
	"github.com/gieseladev/glyrics/v3/pkg/sources"
)

// LyricsInfo is an alias for lyrics.Info
type LyricsInfo = lyrics.Info

// LyricsOrigin is an alias for lyrics.Origin
type LyricsOrigin = lyrics.Origin

// ExtractFromRequest tries to extract lyrics from the provided Request.
// Errors from extracting lyrics are ignored. The only error returned by this
// function is when no extractor was able to extract any lyrics.
func ExtractFromRequest(req *request.Request) (*LyricsInfo, error) {
	for _, e := range sources.GetExtractorsForRequest(req) {
		if err := req.Context().Err(); err != nil {
			return nil, err
		}

		info, err := e.ExtractLyrics(req)
		if err != nil {
			continue
		}

		return info, nil
	}

	return nil, fmt.Errorf("no extractor could extract from %s", req.Url)
}

// ExtractWithContext extracts the lyrics from the url using the context.
func ExtractWithContext(ctx context.Context, url string) (*LyricsInfo, error) {
	return ExtractFromRequest(request.NewWithContext(ctx, url))
}

// Extract wraps the provided url in a Request and performs
// ExtractFromRequest.
func Extract(url string) (*LyricsInfo, error) {
	return ExtractWithContext(context.Background(), url)
}

// SearchLyricsWithContext uses the searcher to search for lyrics based on the
// query. It returns a channel which sends lyrics infos. To stop sending, cancel
// the context.
func Search(ctx context.Context, searcher search.Searcher, query string) <-chan *LyricsInfo {
	lyricsChan := make(chan *LyricsInfo)
	results := searcher.Search(ctx, query)

	go func() {
		defer close(lyricsChan)

		for result := range results {
			req := request.NewWithContext(ctx, result.URL)
			info, err := ExtractFromRequest(req)
			if err == nil {
				lyricsChan <- info
			}
		}
	}()

	return lyricsChan
}

// SearchN returns a slice with at most amount lyrics infos in it.
func SearchN(ctx context.Context, searcher search.Searcher,
	query string, amount int) []LyricsInfo {
	infos := make([]LyricsInfo, 0, amount)

	ctx, cancel := context.WithCancel(ctx)
	lyricsChan := Search(ctx, searcher, query)

	for len(infos) < amount {
		ly, ok := <-lyricsChan
		if !ok {
			break
		}

		infos = append(infos, *ly)
	}

	cancel()

	return infos
}

// SearchFirst returns the first search result from the searcher for the query.
// Might return nil if the context is cancelled or no results are found.
func SearchFirst(ctx context.Context, searcher search.Searcher, query string) *LyricsInfo {
	ctx, cancel := context.WithCancel(ctx)
	lyricsChan := Search(ctx, searcher, query)

	info := <-lyricsChan
	cancel()

	return info
}

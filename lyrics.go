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
	"golang.org/x/sync/semaphore"
	"net/url"
	"sync"
)

// LyricsInfo is an alias for lyrics.Info
type LyricsInfo = lyrics.Info

// LyricsOrigin is an alias for lyrics.Origin
type LyricsOrigin = lyrics.Origin

// ExtractFromRequest tries to extract lyrics from the provided Request.
// Errors from extracting lyrics are ignored. The only error returned by this
// function is when no extractor was able to extract any lyrics.
func ExtractFromRequest(req request.Requester) (*LyricsInfo, error) {
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

	return nil, fmt.Errorf("no extractor could extract from %s", req.URL().String())
}

// ExtractWithContext extracts the lyrics from the url using the context.
func ExtractWithContext(ctx context.Context, rawurl string) (*LyricsInfo, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req := request.NewWithContext(ctx, u)
	defer func() { _ = req.Close() }()
	return ExtractFromRequest(req)
}

// Extract wraps the provided url in a Request and performs
// ExtractFromRequest.
func Extract(u string) (*LyricsInfo, error) {
	return ExtractWithContext(context.Background(), u)
}

type workerCountKey struct{}

// WithWorkers sets the amount of workers to use for Search.
// The default is 3.
func WithWorkers(ctx context.Context, workers int64) context.Context {
	return context.WithValue(ctx, workerCountKey{}, workers)
}

// Search uses the searcher to search for lyrics based on the
// query. It returns a channel which sends lyrics infos. To stop sending, cancel
// the context.
func Search(ctx context.Context, searcher search.Searcher, query string) <-chan *LyricsInfo {
	infos := make(chan *LyricsInfo)
	results := searcher.Search(ctx, query)

	workers, ok := ctx.Value(workerCountKey{}).(int64)
	if !ok {
		workers = 3
	}

	sem := semaphore.NewWeighted(workers)

	type lyricsWithIndex struct {
		Index int
		Info  *LyricsInfo
	}

	unorderedLyrics := make(chan lyricsWithIndex)

	// Uses `workers` amount of goroutines to perform extractions. As the
	// results come in unordered, they aren't returned directly. Instead,
	// another goroutine handles the
	go func() {
		defer close(unorderedLyrics)

		var wg sync.WaitGroup

		index := 0
		for result := range results {
			if err := sem.Acquire(ctx, 1); err != nil {
				break
			}

			wg.Add(1)
			go func(index int, result search.Result) {
				defer sem.Release(1)
				defer wg.Done()

				u, _ := url.Parse(result.URL)
				req := request.NewWithContext(ctx, u)
				info, _ := ExtractFromRequest(req)
				_ = req.Close()

				unorderedLyrics <- lyricsWithIndex{Index: index, Info: info}
			}(index, result)

			index++
		}

		wg.Wait()
	}()

	// receives the unordered results from the extraction and buffers them
	// to return them in order.
	go func() {
		defer close(infos)

		buf := make(map[int]*LyricsInfo)
		index := 0

		for result := range unorderedLyrics {
			if result.Index != index {
				buf[result.Index] = result.Info
				continue
			}

			if info := result.Info; info != nil {
				infos <- info
			}
			index++

			for ; ; index++ {
				info, ok := buf[index]
				if !ok {
					break
				}
				delete(buf, index)

				if info != nil {
					infos <- info
				}

			}
		}
	}()

	return infos
}

// SearchN returns a slice with at most amount lyrics infos in it.
func SearchN(ctx context.Context, searcher search.Searcher,
	query string, amount int) []LyricsInfo {
	infos := make([]LyricsInfo, 0, amount)

	ctx, cancel := context.WithCancel(ctx)
	if amount < 10 {
		ctx = WithWorkers(ctx, int64(amount))
	}
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
	ctx = WithWorkers(ctx, 1)
	lyricsChan := Search(ctx, searcher, query)

	info := <-lyricsChan
	cancel()

	return info
}

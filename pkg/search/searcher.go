package search

import "context"

// Result represents a search result.
type Result struct {
	URL string // url of the search result.
}

// Searcher provides an interface for query-based searching.
type Searcher interface {
	// Search searches for results using the query.
	// It returns a channel which sends the results.
	// To stop searching, cancel the context.
	Search(ctx context.Context, query string) <-chan Result
}

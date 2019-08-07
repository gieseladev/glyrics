package glyrics

import (
	"context"
	"github.com/gieseladev/glyrics/v3/pkg/search"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getGoogleSearcher(t *testing.T) search.Searcher {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Fatal("GOOGLE_API_KEY not set!")
	}

	return &search.Google{APIKey: apiKey}
}

func TestSearchN(t *testing.T) {
	searcher := getGoogleSearcher(t)
	lyrics := SearchN(context.Background(), searcher, "Hello World", 3)
	assert.Len(t, lyrics, 3)
}

func TestSearchFirstLyrics(t *testing.T) {
	searcher := getGoogleSearcher(t)

	lyrics := SearchFirst(context.Background(), searcher, "I dunno what to lookup anymore")
	if lyrics == nil {
		t.Error("Didn't get any lyrics!")
	}
}

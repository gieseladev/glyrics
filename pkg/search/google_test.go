package search

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGoogleSearcher(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Fatal("GOOGLE_API_KEY not set!")
	}

	searcher := Google{APIKey: apiKey}
	results := searcher.Search(context.Background(), "test")

	select {
	case result := <-results:
		link := result.URL
		if link == "" {
			t.Error("Didn't get any links!")
		}
	case <-time.After(5 * time.Second):
		t.Error("Google Search timed out")
	}
}

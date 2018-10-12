package lyricsfinder

import (
	"os"
	"testing"
	"time"
)

func TestGoogleSearch(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Fatal("GOOGLE_API_KEY not set!")
	}

	ch := make(chan string)
	go GoogleSearch("test", apiKey, ch)

	select {
	case link := <-ch:
		if link == "" {
			t.Error("Didn't get any links!")
		}
	case <-time.After(5 * time.Second):
		t.Error("Google Search timed out")
	}
}

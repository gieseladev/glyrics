package glyrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLyrics_MarshalJSON(t *testing.T) {
	releaseDate := time.Now()
	lyrics := LyricsInfo{Url: "Url", Title: "Title", Artist: "Artist", Lyrics: "LyricsInfo",
		ReleaseDate: releaseDate,
		Origin:      LyricsOrigin{Name: "SourceName", Url: "SourceURL"},
	}

	rep, err := json.Marshal(lyrics)
	if err != nil {
		t.Fatal("Couldn't serialise LyricsInfo object")
	}

	dateRep, _ := releaseDate.MarshalJSON()
	expectedRep := []byte(fmt.Sprintf(`{"url":"Url","title":"Title","artist":"Artist","lyrics":"LyricsInfo",`+
		`"release_date":%s,"origin":{"name":"SourceName","url":"SourceURL"}}`, dateRep))

	if !bytes.Equal(rep, expectedRep) {
		t.Errorf("Rep didn't match expectations:\n%q\n====\n%q", rep, expectedRep)
	}
}

func TestExtractor(t *testing.T) {
	urls := []string{
		"https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules",
		"https://www.animelyrics.com/anime/haruhi/harehareyukaiemiri.htm",
		"https://genius.com/Ed-sheeran-the-a-team-lyrics",
		"https://www.lyrical-nonsense.com/lyrics/radwimps/zen-zen-zense",
	}

	for _, url := range urls {
		_, err := ExtractLyrics(url)
		if err != nil {
			t.Error("url:", url, "error:", err)
		}
	}
}

func getGoogleAPIKey(t *testing.T) string {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Fatal("GOOGLE_API_KEY not set!")
	}

	return apiKey
}

func TestSearchNLyrics(t *testing.T) {
	apiKey := getGoogleAPIKey(t)
	lyrics := SearchNLyrics("Hello World", apiKey, 3)
	if len(lyrics) > 3 {
		t.Error("found more than the requested amount of lyrics")
	}
}

func TestSearchFirstLyrics(t *testing.T) {
	apiKey := getGoogleAPIKey(t)

	lyrics := SearchFirstLyrics("I dunno what to lookup anymore", apiKey)
	if lyrics == nil {
		t.Error("Didn't get any lyrics!")
	}
}

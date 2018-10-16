package lyricsfinder

import (
	"github.com/gieseladev/lyricsfindergo/pkg/extractors"
	"github.com/gieseladev/lyricsfindergo/pkg/models"
	"os"
	"testing"
)

func TestExtractors(t *testing.T) {
	if len(extractors.Extractors) == 0 {
		t.Error("Didn't load any Extractors!")
	}
}

func TestExtractor(t *testing.T) {
	urls := []string{
		"https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules",
		"http://www.animelyrics.com/anime/haruhi/harehareyukaiemiri.htm",
		"https://genius.com/Ed-sheeran-the-a-team-lyrics",
		"http://www.lyrical-nonsense.com/lyrics/radwimps/zen-zen-zense",
	}

	for _, url := range urls {
		_, err := ExtractLyrics(url)
		if err != nil {
			t.Error("url:", url, "error:", err)
		}
	}
}

func TestSearchFirstLyrics(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Fatal("GOOGLE_API_KEY not set!")
	}

	lyrics := SearchFirstLyrics("The a Team", apiKey)
	if lyrics == (models.Lyrics{}) {
		t.Error("Didn't get any lyrics!")
	}
}

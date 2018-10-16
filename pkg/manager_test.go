package lyricsfinder

import (
	"github.com/gieseladev/lyricsfindergo/pkg/extractors"
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

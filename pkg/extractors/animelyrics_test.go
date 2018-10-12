package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"testing"
	"time"
)

func TestAnimeLyrics_ExtractLyricsTranslated(t *testing.T) {
	lyrics, err := AnimeLyricsExtractor.ExtractLyrics(
		models.Request{Url: "http://www.animelyrics.com/anime/kmb/fnknhnh.htm"})
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "Futari no Kimochi no Honto no Himitsu", "Yasuna", time.Time{},
		"../../test/testdata/lyrics/animelyrics-yasuna-fnknhnh.txt")

	lyrics, err = AnimeLyricsExtractor.ExtractLyrics(
		models.Request{Url: "http://www.animelyrics.com/anime/akamegakill/liarmask.htm"},
	)
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "Liar Mask", "Rika Mayama", time.Time{},
		"../../test/testdata/lyrics/animelyrics-rika_mayama-liar_mask.txt")
}

func TestAnimeLyrics_ExtractLyricsUntranslated(t *testing.T) {
	lyrics, err := AnimeLyricsExtractor.ExtractLyrics(
		models.Request{Url: "https://www.animelyrics.com/anime/accelworld/chasetheworld.htm"})
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "Chase the world", "May'n", time.Time{},
		"../../test/testdata/lyrics/animelyrics-mayn-chase_the_world.txt")
}

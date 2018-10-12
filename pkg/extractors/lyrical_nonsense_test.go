package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"testing"
	"time"
)

func TestLyricalNonsense_ExtractLyrics(t *testing.T) {
	lyrics, err := LyricalNonsenseExtractor.ExtractLyrics(models.Request{
		Url: "http://www.lyrical-nonsense.com/lyrics/radwimps/zen-zen-zense",
	})
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "Zenzenzense", "RADWIMPS", time.Time{},
		"../../test/testdata/lyrics/lyrical_nonsense-radwimps-zenzenzense.txt")

	lyrics, err = LyricalNonsenseExtractor.ExtractLyrics(models.Request{
		Url: "https://www.lyrical-nonsense.com/lyrics/himouto-umaru-chan-r-theme-songs/umarun-taisou-sisters",
	})
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "Umarun Taisou", "SisterS", time.Time{},
		"../../test/testdata/lyrics/lyrical_nonsense-sisters-umarun_taisou.txt")
}

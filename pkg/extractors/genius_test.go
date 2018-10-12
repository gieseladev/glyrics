package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"testing"
	"time"
)

func TestGeniusLyrics_ExtractLyrics(t *testing.T) {
	lyrics, err := GeniusLyricsExtractor.ExtractLyrics(models.Request{
		Url: "https://genius.com/Ed-sheeran-the-a-team-lyrics",
	})
	if err != nil {
		t.Error(err)
	}

	ExpectLyricsFile(t, lyrics, "The A Team", "Ed Sheeran",
		time.Date(2011, 6, 12, 0, 0, 0, 0, time.UTC),
		"../../test/testdata/lyrics/genius-ed_sheeran-the_a_team.txt")
}

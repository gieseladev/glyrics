package lyrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestLyrics_MarshalJSON(t *testing.T) {
	releaseDate := time.Now()
	lyrics := Info{Url: "Website", Title: "Title", Artist: "Artist", Lyrics: "LyricsInfo",
		ReleaseDate: releaseDate,
		Origin:      Origin{Name: "SourceName", Website: "SourceURL"},
	}

	rep, err := json.Marshal(lyrics)
	if err != nil {
		t.Fatal("Couldn't serialise LyricsInfo object")
	}

	dateRep, _ := releaseDate.MarshalJSON()
	expectedRep := []byte(fmt.Sprintf(`{"url":"Website","title":"Title","artist":"Artist","lyrics":"LyricsInfo",`+
		`"release_date":%s,"origin":{"name":"SourceName","url":"SourceURL"}}`, dateRep))

	if !bytes.Equal(rep, expectedRep) {
		t.Errorf("Rep didn't match expectations:\n%q\n====\n%q", rep, expectedRep)
	}
}

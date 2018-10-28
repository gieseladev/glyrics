package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestLyrics_MarshalJSON(t *testing.T) {
	releaseDate := time.Now()
	lyrics := Lyrics{Url: "Url", Title: "Title", Artist: "Artist", Lyrics: "Lyrics",
		ReleaseDate: releaseDate,
		Origin:      &LyricsOrigin{Name: "SourceName", Url: "SourceURL"},
	}

	rep, err := json.Marshal(lyrics)
	if err != nil {
		t.Fatal("Couldn't serialise Lyrics object")
	}

	dateRep, _ := releaseDate.MarshalJSON()
	expectedRep := []byte(fmt.Sprintf(`{"url":"Url","title":"Title","artist":"Artist","lyrics":"Lyrics",`+
		`"release_date":%s,"origin":{"name":"SourceName","url":"SourceURL"}}`, dateRep))

	if !bytes.Equal(rep, expectedRep) {
		t.Errorf("Rep didn't match expectations:\n%q\n====\n%q", rep, expectedRep)
	}
}

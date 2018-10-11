package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestMusixMatch_CanHandle(t *testing.T) {
	req := lyricsfinder.Request{Url: "https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules"}
	defer req.Close()

	if !MusixMatchExtractor.CanHandle(req) {
		t.Errorf("Extractor didn't accept %s even though it should've", req.Url)
	}
}

func TestMusixMatch_ExtractLyrics(t *testing.T) {
	req := lyricsfinder.Request{Url: "https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules"}
	defer req.Close()

	lyrics, err := MusixMatchExtractor.ExtractLyrics(req)
	if err != nil {
		t.Error(err)
		return
	}

	rawExpectedLyrics, err := ioutil.ReadFile("../../test/testdata/lyrics/musixmatch-dua_lipa-new_rules.txt")
	if err != nil {
		t.Fatal(err)
	}

	expectedLyrics := strings.Replace(string(rawExpectedLyrics), "\r\n", "\n", -1)

	if lyrics.Title != "New Rules" {
		t.Errorf("Title \"%s\" didn't match", lyrics.Title)
	}
	if lyrics.Artist != "Dua Lipa" {
		t.Errorf("Artist \"%s\" didn't match", lyrics.Artist)
	}
	if lyrics.ReleaseDate != time.Date(2017, time.June, 2, 0, 0, 0, 0, time.UTC) {
		t.Errorf("Date \"%s\" didn't match", lyrics.ReleaseDate)
	}
	if lyrics.Lyrics != expectedLyrics {
		t.Errorf("Lyrics didn't match:\n====\n%s\n====\nVS\n====\n%s\n====", lyrics.Lyrics, expectedLyrics)
	}
}

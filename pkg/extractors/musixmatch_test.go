package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func ExpectLyrics(t *testing.T, lyrics *models.Lyrics, title, artist string, releaseDate time.Time, lyricsText string) {
	if lyrics.Title != title {
		t.Errorf("Title %q didn't match: %q", lyrics.Title, title)
	}
	if lyrics.Artist != artist {
		t.Errorf("Artist %q didn't match: %q", lyrics.Artist, artist)
	}
	if lyrics.ReleaseDate != releaseDate {
		t.Errorf("Date %s didn't match: %s", lyrics.ReleaseDate, releaseDate)
	}
	if lyrics.Lyrics != lyricsText {
		t.Errorf("Lyrics didn't match:\n====\n%q\n====\nVS\n====\n%q\n====", lyrics.Lyrics, lyricsText)
	}
}

func ExpectLyricsFile(t *testing.T, lyrics *models.Lyrics, title, artist string, releaseDate time.Time, lyricsFile string) {
	rawExpectedLyrics, err := ioutil.ReadFile(lyricsFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedLyrics := strings.Replace(string(rawExpectedLyrics), "\r\n", "\n", -1)
	ExpectLyrics(t, lyrics, title, artist, releaseDate, expectedLyrics)
}

func TestMusixMatch_CanHandle(t *testing.T) {
	req := models.Request{Url: "https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules"}
	defer req.Close()

	if !MusixMatchExtractor.CanHandle(req) {
		t.Errorf("Extractor didn't accept %s even though it should've", req.Url)
	}
}

func TestMusixMatch_ExtractLyrics(t *testing.T) {
	req := models.Request{Url: "https://www.musixmatch.com/lyrics/Dua-Lipa/New-Rules"}
	defer req.Close()

	lyrics, err := MusixMatchExtractor.ExtractLyrics(req)
	if err != nil {
		t.Error(err)
		return
	}

	ExpectLyricsFile(t, lyrics, "New Rules", "Dua Lipa",
		time.Date(2017, time.June, 2, 0, 0, 0, 0, time.UTC),
		"../../test/testdata/lyrics/musixmatch-dua_lipa-new_rules.txt")
}

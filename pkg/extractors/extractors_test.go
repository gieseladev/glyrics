package extractors

import (
	"fmt"
	"github.com/gieseladev/lyricsfindergo/pkg/models"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type lyricsTestCase struct {
	SkipIf      string
	Url         string
	Extractor   string
	Title       string
	Artist      string
	ReleaseDate time.Time `yaml:"release_date"`
	Lyrics      string
}

func (test lyricsTestCase) String() string {
	return fmt.Sprintf("(%-16s) %s - %s", test.Extractor, test.Artist, test.Title)
}

func (test *lyricsTestCase) ShouldSkip(t *testing.T) bool {
	if test.SkipIf == "travis" && os.Getenv("TRAVIS") == "true" {
		return true
	}

	return false
}

func (test *lyricsTestCase) Test(t *testing.T, lyrics *models.Lyrics) {
	if test.Title != lyrics.Title {
		t.Errorf("Title %q didn't match: %q", lyrics.Title, test.Title)
	}
	if test.Artist != lyrics.Artist {
		t.Errorf("Artist %q didn't match: %q", lyrics.Artist, test.Artist)
	}
	if test.ReleaseDate != lyrics.ReleaseDate {
		t.Errorf("Date %s didn't match: %s", lyrics.ReleaseDate, test.ReleaseDate)
	}
	if test.Lyrics != lyrics.Lyrics {
		t.Errorf("Lyrics didn't match:\n====\n%q\n====\nVS\n====\n%q\n====", lyrics.Lyrics, test.Lyrics)
	}

	if lyrics.Origin == nil {
		t.Errorf("Lyrics don't have an origin: %s", lyrics)
	}
	if test.Extractor != lyrics.Origin.Name {
		t.Errorf("Origin %q didn't match: %q", lyrics.Origin.Name, test.Extractor)
	}
}

func gatherTestCases(t *testing.T) []lyricsTestCase {
	pattern := filepath.Join("..", "..", "test", "testdata", "lyrics", "*.yml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Error(err)
	}

	cases := make([]lyricsTestCase, len(files))

	for i, file := range files {
		var testCase lyricsTestCase
		file, err := os.Open(file)
		if err != nil {
			t.Log(err)
			continue
		}

		err = yaml.NewDecoder(file).Decode(&testCase)
		if err != nil {
			t.Log(err)
			continue
		}

		cases[i] = testCase
	}

	return cases
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func findExtractor(name string) Extractor {
	name = strings.Replace(name, " ", "", -1)

	for _, extractor := range Extractors {
		extractorName := getType(extractor)
		if strings.HasPrefix(strings.ToLower(extractorName), strings.ToLower(name)) {
			return extractor
		}
	}
	return nil
}

func TestExtractors(t *testing.T) {
	cases := gatherTestCases(t)

	t.Logf("Testing %d case(s)", len(cases))

	for _, testCase := range cases {
		if testCase == (lyricsTestCase{}) {
			t.Error("Empty test case, skipping!")
			continue
		}

		t.Log(testCase.String())
		if testCase.ShouldSkip(t) {
			t.Log("> Skipped")
			continue
		}

		extractor := findExtractor(testCase.Extractor)
		if extractor == nil {
			t.Error("ERROR: Couldn't find extractor")
			continue
		}

		req := models.Request{Url: testCase.Url}

		if !extractor.CanHandle(req) {
			t.Errorf("ERROR: Extractor %s can't handle url %s", extractor, req.Url)
		}

		lyrics, err := extractor.ExtractLyrics(req)
		if err != nil {
			t.Error(err)
			continue
		}

		testCase.Test(t, lyrics)
	}
}

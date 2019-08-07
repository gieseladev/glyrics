package sources

import (
	"fmt"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func allTrue(values ...bool) bool {
	for _, value := range values {
		if !value {
			return false
		}
	}

	return true
}

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

func (test *lyricsTestCase) Check(t *testing.T, info *lyrics.Info) bool {
	a := assert.New(t)

	return allTrue(
		a.Equal(test.Title, info.Title, "title didn't match"),
		a.Equal(test.Artist, info.Artist, "artist didn't match"),
		a.Equal(test.ReleaseDate, info.ReleaseDate, "release date didn't match"),
		a.Equal(test.Lyrics, info.Lyrics, "lyrics didn't match"),

		a.Equal(test.Extractor, info.Origin.Name, "origin name didn't match"),
	)
}

func (test *lyricsTestCase) Test(t *testing.T) {
	r := require.New(t)

	r.NotNil(test, "empty test case")

	if test.ShouldSkip(t) {
		t.Log("> Skipped")
		return
	}

	req := request.New(test.Url)

	extractors := GetExtractorsForRequest(req)
	r.Len(extractors, 1, "must return exactly one extractor")
	extractor := extractors[0]

	lyrics, err := extractor.ExtractLyrics(req)
	r.NoError(err, "extractor returned error")

	test.Check(t, lyrics)
}

func gatherTestCases(t *testing.T) []lyricsTestCase {
	pattern := filepath.FromSlash("../../test/data/lyrics/*.yml")
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

func TestExtractors(t *testing.T) {
	cases := gatherTestCases(t)

	t.Logf("Testing %d case(s)", len(cases))

	for _, testCase := range cases {
		t.Run(testCase.String(), testCase.Test)
	}
}

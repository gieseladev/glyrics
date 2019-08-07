package sources

import (
	"errors"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
)

var (
	// AZLyricsOrigin is the lyrics origin for AZLyrics.
	AZLyricsOrigin = lyrics.Origin{Name: "AZLyrics", Website: "azlyrics.com"}

	// AZLyricsExtractor is an extractor for AZLyrics
	AZLyricsExtractor = ExtractorFunc(extractAZLyricsLyrics)
)

func extractAZLyricsLyrics(req *request.Request) (*lyrics.Info, error) {
	const (
		// format: `"<title>" lyrics`
		titlePrefixLen = len(`"`)
		titleSuffixLen = len(`" lyrics`)
		titleMinLen    = titlePrefixLen + titleSuffixLen

		// format: `<artist> Lyrics`
		artistPrefixLen = 0
		artistSuffixLen = len(` Lyrics`)
		artistMinLen    = artistPrefixLen + artistSuffixLen
	)

	req.Request().Header.Set(
		"user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
	)

	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	// check for the presence of the sorting buttons
	if sortingButtons := doc.Find(".btn.sorting"); sortingButtons.Length() != 0 {
		return nil, errors.New("not a lyrics page")
	}

	center := doc.Find("div.text-center:not(.noprint)")
	if center.Length() == 0 {
		return nil, errors.New("no lyrics found")
	}

	title := strings.TrimSpace(center.Find("h1").Text())
	// using <= because a literally empty title is just as implausible
	if len(title) <= titleMinLen {
		return nil, errors.New("no title found, suspecting no lyrics page")
	}

	title = title[titlePrefixLen : len(title)-titleSuffixLen]

	artist := doc.Find("div.lyricsh h2 b").Text()
	if len(artist) > artistMinLen {
		artist = artist[artistPrefixLen : len(artist)-artistSuffixLen]
	} else {
		artist = ""
	}

	lyricsText := strings.TrimSpace(center.Find("div:not([class])").First().Text())

	return &lyrics.Info{URL: req.URL, Title: title, Artist: artist, Lyrics: lyricsText,
		Origin: AZLyricsOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?azlyrics.com/.+`)),
		AZLyricsExtractor,
	))
}

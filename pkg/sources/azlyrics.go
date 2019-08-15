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

func extractAZLyricsLyrics(req request.Requester) (*lyrics.Info, error) {
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

	return &lyrics.Info{URL: req.URL().String(), Title: title, Artist: artist, Lyrics: lyricsText,
		Origin: AZLyricsOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexCanExtractTeller(regexp.MustCompile(`https?://(?:www.)?azlyrics.com/.+`)),
		AZLyricsExtractor,
	))
}

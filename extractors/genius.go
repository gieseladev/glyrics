package extractors

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/v3/pkg/requests"
	"regexp"
	"strings"
	"time"
)

// GeniusOrigin is the LyricsOrigin for Genius.
var GeniusOrigin = LyricsOrigin{Name: "Genius", Url: "genius.com"}

func ExtractGeniusLyrics(req *requests.Request) (*LyricsInfo, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	title := doc.Find("h1.header_with_cover_art-primary_info-title").Text()
	artist := doc.Find("a.header_with_cover_art-primary_info-primary_artist").Text()

	var rawDate string
	doc.Find("div.metadata_unit").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.Find("span.metadata_unit-label").Text() == "Release Date" {
			rawDate = selection.Find("span.metadata_unit-info").Text()
			return false
		}
		return true
	})

	releaseDate, _ := time.Parse("January 2, 2006", rawDate)

	lyrics := strings.TrimSpace(doc.Find("div.lyrics").First().Text())

	if lyrics == "" {
		return nil, errors.New("no lyrics found")
	}

	return &LyricsInfo{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		ReleaseDate: releaseDate,
		Origin:      GeniusOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?genius.com/.*`)),
		ExtractorFunc(ExtractGeniusLyrics),
	))
}

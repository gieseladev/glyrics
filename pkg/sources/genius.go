package sources

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
	"time"
)

var (
	// GeniusOrigin is the Origin for Genius.
	GeniusOrigin = lyrics.Origin{Name: "Genius", Website: "genius.com"}

	// GeniusExtractor is an extractor for Genius
	GeniusExtractor = ExtractorFunc(extractGeniusLyrics)
)

func extractGeniusLyrics(req *request.Request) (*lyrics.Info, error) {
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

	lyricsText := strings.TrimSpace(doc.Find("div.lyrics").First().Text())

	if lyricsText == "" {
		return nil, errors.New("no lyrics found")
	}

	return &lyrics.Info{Url: req.Url, Title: title, Artist: artist, Lyrics: lyricsText,
		ReleaseDate: releaseDate,
		Origin:      GeniusOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?genius.com/.*`)),
		GeniusExtractor,
	))
}

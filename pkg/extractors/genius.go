package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/pkg/models"
	"regexp"
	"strings"
	"time"
)

// GeniusOrigin is the models.LyricsOrigin for Genius.
var GeniusOrigin = models.LyricsOrigin{Name: "Genius", Url: "genius.com"}

type geniusLyrics struct {
	RegexCanHandle
}

func (extractor *geniusLyrics) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
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

	return &models.Lyrics{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		ReleaseDate: releaseDate,
		Origin:      &GeniusOrigin}, nil
}

// GeniusLyricsExtractor is the Extractor instance used for Genius
var GeniusLyricsExtractor = geniusLyrics{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?genius.com/.*`),
}}

func init() {
	RegisterExtractor(&GeniusLyricsExtractor)
}

package extractors

import (
	"errors"
	"github.com/gieseladev/lyricsfindergo/pkg/models"
	"regexp"
	"strings"
)

var AZLyricsOrigin = models.LyricsOrigin{Name: "AZLyrics", Url: "azlyrics.com"}

type azLyrics struct {
	RegexCanHandle
}

func (extractor *azLyrics) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	req.Request().Header.Set(
		"user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
	)

	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	center := doc.Find("div.text-center:not(.noprint)")
	if center.Length() == 0 {
		return nil, errors.New("no lyrics found")
	}

	title := strings.TrimSpace(center.Find("h1").Text())
	title = title[1 : len(title)-8]

	artist := doc.Find("div.lyricsh h2 b").Text()
	artist = artist[:len(artist)-7]

	lyrics := strings.TrimSpace(center.Find("div:not([class])").First().Text())

	return &models.Lyrics{Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &AZLyricsOrigin}, nil
}

var AZLyricsExtractor = azLyrics{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?azlyrics.com/.*`),
}}

func init() {
	RegisterExtractor(&AZLyricsExtractor)
}

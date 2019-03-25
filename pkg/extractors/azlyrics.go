package extractors

import (
	"errors"
	"github.com/gieseladev/glyrics/pkg/models"
	"regexp"
	"strings"
)

// AZLyricsOrigin is the models.LyricsOrigin for AZLyrics.
var AZLyricsOrigin = models.LyricsOrigin{Name: "AZLyrics", Url: "azlyrics.com"}

type azLyrics struct {
	RegexCanHandle
}

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

func (extractor *azLyrics) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
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

	lyrics := strings.TrimSpace(center.Find("div:not([class])").First().Text())

	return &models.Lyrics{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &AZLyricsOrigin}, nil
}

// AZLyricsExtractor is the Extractor instance used for AZLyrics
var AZLyricsExtractor = azLyrics{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?azlyrics.com/.*`),
}}

func init() {
	RegisterExtractor(&AZLyricsExtractor)
}

package extractors

import (
	"errors"
	"github.com/gieseladev/glyrics/pkg/models"
	"regexp"
	"strings"
)

// LyricsModeOrigin is the models.LyricsOrigin for LyricsMode.
var LyricsModeOrigin = models.LyricsOrigin{Name: "LyricsMode", Url: "lyricsmode.com"}

type lyricsMode struct {
	RegexCanHandle
}

var artistTitleSplit = regexp.MustCompile(`\s*(?P<artist>.+?)\s+–\s+(?P<title>.+) (?:lyrics)?`)

func (extractor *lyricsMode) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	header := doc.Find("h1.song_name.fs32").Text()

	var artist, title string
	match := artistTitleSplit.FindStringSubmatch(header)
	if len(match) >= 2 {
		artist = match[1]
		title = match[2]
	} else {
		return nil, errors.New("couldn't find title and artist")
	}

	lyricsContainer := doc.Find("#lyrics_text").First()
	lyricsContainer.Children().RemoveFiltered("div.hide")

	lyrics := strings.TrimSpace(lyricsContainer.Text())

	return &models.Lyrics{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &LyricsModeOrigin}, nil
}

// LyricsModeExtractor is the Extractor instance used for LyricsMode
var LyricsModeExtractor = lyricsMode{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?lyricsmode.com/.*`),
}}

func init() {
	RegisterExtractor(&LyricsModeExtractor)
}

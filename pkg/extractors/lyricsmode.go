package extractors

import (
	"errors"
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"regexp"
	"strings"
)

var LyricsModeOrigin = models.LyricsOrigin{Name: "LyricsMode", Url: "lyricsmode.com"}

type lyricsMode struct {
	RegexCanHandle
}

var ArtistTitleSplit = regexp.MustCompile(`\s*(?P<artist>.+?)\s+–\s+(?P<title>.+) (?:lyrics)?`)

func (extractor *lyricsMode) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	header := doc.Find("h1.song_name.fs32").Text()

	var artist, title string
	match := ArtistTitleSplit.FindStringSubmatch(header)
	if len(match) >= 2 {
		artist = match[1]
		title = match[2]
	} else {
		return nil, errors.New("couldn't find title and artist")
	}

	lyrics := strings.TrimSpace(doc.Find("p#lyrics_text").First().Text())

	return &models.Lyrics{Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &LyricsModeOrigin}, nil
}

var LyricsModeExtractor = lyricsMode{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?lyricsmode.com/.*`),
}}

func init() {
	RegisterExtractor(&LyricsModeExtractor)
}

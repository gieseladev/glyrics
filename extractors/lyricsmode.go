package extractors

import (
	"errors"
	"github.com/gieseladev/glyrics/v3/pkg/requests"
	"regexp"
	"strings"
)

// LyricsModeOrigin is the glyrics.LyricsOrigin for LyricsMode.
var LyricsModeOrigin = LyricsOrigin{Name: "LyricsMode", Url: "lyricsmode.com"}

var artistTitleSplit = regexp.MustCompile(`\s*(?P<artist>.+?)\s+–\s+(?P<title>.+) (?:lyrics)?`)

func ExtractLyricsModeLyrics(req *requests.Request) (*LyricsInfo, error) {
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

	return &LyricsInfo{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		Origin: LyricsModeOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?lyricsmode.com/.*`)),
		ExtractorFunc(ExtractLyricsModeLyrics),
	))
}

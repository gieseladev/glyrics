package sources

import (
	"errors"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
)

var (
	// LyricsModeOrigin is the glyrics.Origin for LyricsMode.
	LyricsModeOrigin = lyrics.Origin{Name: "LyricsMode", Website: "lyricsmode.com"}

	// LyricsModeExtractor is an extractor for LyricsMode
	LyricsModeExtractor = ExtractorFunc(extractLyricsModeLyrics)
)

var lyricsModeHeaderMatcher = regexp.MustCompile(`\s*(?P<artist>.+?)\s+–\s+(?P<title>.+) (?:lyrics)?`)

func extractLyricsModeLyrics(req request.Requester) (*lyrics.Info, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	header := doc.Find("h1.song_name.fs32").Text()

	var artist, title string
	match := lyricsModeHeaderMatcher.FindStringSubmatch(header)
	if len(match) >= 2 {
		artist = match[1]
		title = match[2]
	} else {
		return nil, errors.New("couldn't find title and artist")
	}

	lyricsContainer := doc.Find("#lyrics_text").First()
	lyricsContainer.Children().RemoveFiltered("div.hide")

	lyricsText := strings.TrimSpace(lyricsContainer.Text())

	return &lyrics.Info{URL: req.URL().String(), Title: title, Artist: artist, Lyrics: lyricsText,
		Origin: LyricsModeOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexCanExtractTeller(regexp.MustCompile(`https?://(?:www.)?lyricsmode.com/.*`)),
		LyricsModeExtractor,
	))
}

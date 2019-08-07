package sources

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
)

var (
	// LyricalNonsenseOrigin is the lyrics origin for Lyrical Nonsense.
	LyricalNonsenseOrigin = lyrics.Origin{Name: "Lyrical Nonsense", Website: "lyrical-nonsense.com"}

	// LyricalNonsenseExtractor is an extractor for Lyrical Nonsense
	LyricalNonsenseExtractor = ExtractorFunc(extractLyricalNonsenseLyrics)
)

func extractLyricalNonsenseLyrics(req *request.Request) (*lyrics.Info, error) {
	const (
		titleSuffixLen = len(` 歌詞`)
		titleMinLen    = titleSuffixLen
	)

	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	headerContents := doc.Find("h1").Contents()
	if headerContents.Length() != 2 {
		return nil, errors.New("couldn't decode header information")
	}

	rawTitle := strings.TrimSpace(headerContents.Get(0).Data)
	artist := headerContents.Next().Text()
	if len(rawTitle) < titleMinLen {
		return nil, errors.New("no title found")
	}

	title := rawTitle[:len(rawTitle)-titleSuffixLen]

	lyricsSel := doc.Find("div#Romaji div.olyrictext").Add("div#LyricsInfo div.olyrictext")
	lyricsBuilder := strings.Builder{}

	lyricsSel.First().Find("p").Each(func(i int, selection *goquery.Selection) {
		for _, line := range strings.Split(selection.Text(), "\n") {
			line = strings.TrimSpace(line)
			lyricsBuilder.WriteString(line + "\n")
		}
		lyricsBuilder.WriteString("\n")
	})

	lyricsText := strings.TrimSpace(lyricsBuilder.String())

	return &lyrics.Info{Url: req.Url, Title: title, Artist: artist, Lyrics: lyricsText,
		Origin: LyricalNonsenseOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?lyrical-nonsense.com/lyrics/.*`)),
		LyricalNonsenseExtractor,
	))
}

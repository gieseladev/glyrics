package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/pkg/requests"
	"regexp"
	"strings"
)

// LyricalNonsenseOrigin is the glyrics.LyricsOrigin for Lyrical Nonsense.
var LyricalNonsenseOrigin = LyricsOrigin{Name: "Lyrical Nonsense", Url: "lyrical-nonsense.com"}

func ExtractLyricalNonsenseLyrics(req *requests.Request) (*LyricsInfo, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	titleSel := doc.Find("span.titletext2new").Add("div.titlelyricblocknew h1")
	title := titleSel.First().Text()

	artistSel := doc.Find("div.artistcontainer span.artisttext2new").Add("div.artistcontainer h2")
	artist := artistSel.First().Text()

	lyricsSel := doc.Find("div#Romaji div.olyrictext").Add("div#LyricsInfo div.olyrictext")
	lyricsBuilder := strings.Builder{}

	lyricsSel.First().Find("p").Each(func(i int, selection *goquery.Selection) {
		for _, line := range strings.Split(selection.Text(), "\n") {
			line = strings.TrimSpace(line)
			lyricsBuilder.WriteString(line + "\n")
		}
		lyricsBuilder.WriteString("\n")
	})

	lyrics := strings.TrimSpace(lyricsBuilder.String())

	return &LyricsInfo{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		Origin: LyricalNonsenseOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?lyrical-nonsense.com/lyrics/.*`)),
		ExtractorFunc(ExtractLyricalNonsenseLyrics),
	))
}

package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"regexp"
	"strings"
)

var LyricalNonsenseOrigin = models.LyricsOrigin{Name: "LyricalNonsense", Url: "lyrical-nonsense.com"}

type lyricalNonsense struct {
	RegexCanHandle
}

func (extractor *lyricalNonsense) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	titleSel := doc.Find("span.titletext2new").Add("div.titlelyricblocknew h1")
	title := titleSel.First().Text()

	artistSel := doc.Find("div.artistcontainer span.artisttext2new").Add("div.artistcontainer h2")
	artist := artistSel.First().Text()

	lyricsSel := doc.Find("div#Romaji div.olyrictext").Add("div#Lyrics div.olyrictext")
	lyricsBuilder := strings.Builder{}

	lyricsSel.First().Find("p").Each(func(i int, selection *goquery.Selection) {
		lyricsBuilder.WriteString(selection.Text() + "\n\n")
	})

	lyrics := strings.TrimSpace(lyricsBuilder.String())

	return &models.Lyrics{Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &LyricalNonsenseOrigin}, nil
}

var LyricalNonsenseExtractor = lyricalNonsense{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?lyrical-nonsense.com/lyrics/.*`),
}}

func init() {
	RegisterExtractor(&LyricalNonsenseExtractor)
}

package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/pkg/models"
	"regexp"
	"strings"
)

// AnimeLyricsOrigin is the models.LyricsOrigin for AnimeLyrics.
var AnimeLyricsOrigin = models.LyricsOrigin{Name: "Animelyrics", Url: "animelyrics.com"}

type animeLyrics struct {
	RegexCanHandle
}

var artistMatcher = regexp.MustCompile(`Performed by:? (?P<artist>[\w' ]+)\b`)

func (extractor *animeLyrics) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	var artist, lyrics string

	title := strings.TrimSpace(doc.Find("div~h1").First().Contents().First().Text())

	artistSearchDoc := doc.Clone()
	artistSearchDoc.Find("br").ReplaceWithHtml("\n")

	artistMatch := artistMatcher.FindStringSubmatch(artistSearchDoc.Text())
	if len(artistMatch) > 1 {
		artist = artistMatch[1]
	}

	if window := doc.Find(`table[cellspacing="0"][border="0"]`); window.Length() > 0 {
		window.Find("dt").Remove()
		window.Find("span").AfterHtml("\n\n")
		lyrics = window.Find("td.romaji").Text()
	} else {
		center := doc.Find("div.centerbox")
		passedDt := false
		lyricsBuilder := strings.Builder{}

		center.Contents().EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if goquery.NodeName(selection) == "#text" {
				if passedDt {
					lyricsBuilder.WriteString(selection.Text())
				}
			} else {
				if selection.Is("dt") {
					passedDt = true
				} else if !selection.Is("br") && passedDt {
					return false
				}
			}
			return true
		})

		lyrics = lyricsBuilder.String()
	}

	lyrics = strings.TrimSpace(strings.Replace(lyrics, "\u00a0", " ", -1))

	return &models.Lyrics{Url: req.Url, Title: title, Artist: artist, Lyrics: lyrics,
		Origin: &AnimeLyricsOrigin}, nil
}

// AnimeLyricsExtractor is the Extractor instance used for AnimeLyrics
var AnimeLyricsExtractor = animeLyrics{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?animelyrics.com/.*`),
}}

func init() {
	RegisterExtractor(&AnimeLyricsExtractor)
}

package sources

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
)

var (
	// AnimeLyricsOrigin is the lyrics origin for Animelyrics
	AnimeLyricsOrigin = lyrics.Origin{Name: "Animelyrics", Website: "animelyrics.com"}

	// AnimeLyricsExtractor is an extractor for Animelyrics
	AnimeLyricsExtractor = ExtractorFunc(extractAnimeLyricsLyrics)
)

var animeLyricsArtistMatcher = regexp.MustCompile(`Performed by:? (?P<artist>[\w' ]+)\b`)

func extractAnimeLyricsLyrics(req request.Requester) (*lyrics.Info, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	var artist, lyricsText string

	title := strings.TrimSpace(doc.Find("div~h1").First().Contents().First().Text())

	artistSearchDoc := doc.Clone()
	artistSearchDoc.Find("br").ReplaceWithHtml("\n")

	artistMatch := animeLyricsArtistMatcher.FindStringSubmatch(artistSearchDoc.Text())
	if len(artistMatch) > 1 {
		artist = artistMatch[1]
	}

	if window := doc.Find(`table[cellspacing="0"][border="0"]`); window.Length() > 0 {
		window.Find("dt").Remove()
		window.Find("span").AfterHtml("\n\n")
		lyricsText = window.Find("td.romaji").Text()
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

		lyricsText = lyricsBuilder.String()
	}

	lyricsText = strings.TrimSpace(strings.Replace(lyricsText, "\u00a0", " ", -1))

	return &lyrics.Info{URL: req.URL().String(),
		Title: title, Artist: artist,
		Lyrics: lyricsText,
		Origin: AnimeLyricsOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexCanExtractTeller(regexp.MustCompile(`https?://(?:www.)?animelyrics.com/.*`)),
		AnimeLyricsExtractor,
	))
}

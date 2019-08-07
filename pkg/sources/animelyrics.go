package sources

import (
	"github.com/PuerkitoBio/goquery"
	lyrics2 "github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"strings"
)

// AnimeLyricsOrigin is the lyrics origin for Animelyrics
var AnimeLyricsOrigin = lyrics2.Origin{Name: "Animelyrics", Website: "animelyrics.com"}

var animeLyricsArtistMatcher = regexp.MustCompile(`Performed by:? (?P<artist>[\w' ]+)\b`)

func ExtractAnimeLyricsLyrics(req *request.Request) (*lyrics2.Info, error) {
	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	var artist, lyrics string

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

	return &lyrics2.Info{Url: req.Url,
		Title: title, Artist: artist,
		Lyrics: lyrics,
		Origin: AnimeLyricsOrigin}, nil
}

func init() {
	RegisterExtractor(CreateMaybeExtractor(
		RegexExtractorTeller(regexp.MustCompile(`https?://(?:www.)?animelyrics.com/.*`)),
		ExtractorFunc(ExtractAnimeLyricsLyrics),
	))
}

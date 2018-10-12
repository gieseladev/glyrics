package extractors

import (
	"errors"
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"regexp"
	"strings"
	"time"
)

var MusixMatchOrigin = models.LyricsOrigin{Name: "musixMatch", Url: "musixmatch.com"}

type musixMatch struct {
	RegexCanHandle
}

func (extractor *musixMatch) ExtractLyrics(req models.Request) (*models.Lyrics, error) {
	req.Request().Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0")

	doc, err := req.Document()
	if err != nil {
		return nil, err
	}

	if doc.Find(`div.mxm-empty-state[data-reactid="87"]`).Length() > 0 {
		return nil, errors.New("no lyrics")
	}

	window := doc.Find("div.mxm-lyrics span")
	if window.Length() == 0 {
		return nil, errors.New("no lyrics")
	}

	window.Find("script").ReplaceWithHtml("\n\n")

	lyrics := window.First().Text()
	title := strings.TrimSpace(doc.Find("h1.mxm-track-title__track").First().Text())[6:]
	artist := doc.Find("a.mxm-track-title__artist").First().Text()

	rawDate := doc.Find("div.mxm-track-footer__album h3.mui-cell__subtitle").Text()

	date, _ := parseOrdinalDate("Jan 2 2006", rawDate)

	return &models.Lyrics{Title: title, Artist: artist, ReleaseDate: date, Lyrics: lyrics,
		Origin: &MusixMatchOrigin}, nil
}

var MusixMatchExtractor = musixMatch{RegexCanHandle{
	UrlMatch: regexp.MustCompile(`https?://(?:www.)?musixmatch.com/lyrics/.*`),
}}

func init() {
	RegisterExtractor(&MusixMatchExtractor)
}

var dayOrdinals = map[string]string{ // map[ordinal]cardinal
	"1st": "1", "2nd": "2", "3rd": "3", "4th": "4", "5th": "5",
	"6th": "6", "7th": "7", "8th": "8", "9th": "9", "10th": "10",
	"11th": "11", "12th": "12", "13th": "13", "14th": "14", "15th": "15",
	"16th": "16", "17th": "17", "18th": "18", "19th": "19", "20th": "20",
	"21st": "21", "22nd": "22", "23rd": "23", "24th": "24", "25th": "25",
	"26th": "26", "27th": "27", "28th": "28", "29th": "29", "30th": "30",
	"31st": "31",
}

func parseOrdinalDate(layout, value string) (time.Time, error) {
	const ( // day number
		cardMinLen = len("1")
		cardMaxLen = len("31")
		ordSfxLen  = len("th")
		ordMinLen  = cardMinLen + ordSfxLen
	)

	for k := 0; k < len(value)-ordMinLen; {
		// i number start
		for ; k < len(value) && (value[k] > '9' || value[k] < '0'); k++ {
		}
		i := k
		// j cardinal end
		for ; k < len(value) && (value[k] <= '9' && value[k] >= '0'); k++ {
		}
		j := k
		if j-i > cardMaxLen || j-i < cardMinLen {
			continue
		}
		// k ordinal end
		// ASCII Latin (uppercase | 0x20) = lowercase
		for ; k < len(value) && (value[k]|0x20 >= 'a' && value[k]|0x20 <= 'z'); k++ {
		}
		if k-j != ordSfxLen {
			continue
		}

		// day ordinal to cardinal
		for ; i < j-1 && (value[i] == '0'); i++ {
		}
		o := strings.ToLower(value[i:k])
		c, ok := dayOrdinals[o]
		if ok {
			value = value[:i] + c + value[k:]
			break
		}
	}

	return time.Parse(layout, value)
}

package extractors

import (
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"regexp"
)

var Extractors = make([]Extractor, 0)

func RegisterExtractor(extractor Extractor) {
	Extractors = append(Extractors, extractor)
}

type Extractor interface {
	CanHandle(req models.Request) bool

	ExtractLyrics(req models.Request) (*models.Lyrics, error)
}

type RegexCanHandle struct {
	UrlMatch *regexp.Regexp
}

func (extractor *RegexCanHandle) CanHandle(req models.Request) bool {
	return extractor.UrlMatch.MatchString(req.Url)
}

package lyricsfinder

import "regexp"

type Extractor interface {
	CanHandle(req Request) bool

	ExtractLyrics(req Request) (*Lyrics, error)
}

type RegexCanHandle struct {
	UrlMatch *regexp.Regexp
}

func (extractor *RegexCanHandle) CanHandle(req Request) bool {
	return extractor.UrlMatch.MatchString(req.Url)
}

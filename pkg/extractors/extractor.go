// Package extractors contains the various lyrics extractors for gLyrics.
package extractors

import (
	"github.com/gieseladev/glyrics/pkg/models"
	"regexp"
)

// Extractors is a slice of all the loaded extractors.
// Use RegisterExtractor to register a new extractor.
var Extractors = make([]Extractor, 0)

// RegisterExtractor adds a new Extractor to the Extractors slice.
// Use this function to add a new extractor to be used by gLyrics
func RegisterExtractor(extractor Extractor) {
	Extractors = append(Extractors, extractor)
}

// Extractor extracts lyrics from a Request.
type Extractor interface {
	// CanHandle performs simple checks to determine whether the extractor
	// has any chance of extracting lyrics from the Request.
	// It does not guarantee that Extractor.ExtractLyrics will be successful.
	CanHandle(req models.Request) bool

	// ExtractLyrics performs the actual extraction.
	ExtractLyrics(req models.Request) (*models.Lyrics, error)
}

// RegexCanHandle is a mixin for Extractor which uses a Regexp to check
// the Request url against in CanHandle.
type RegexCanHandle struct {
	UrlMatch *regexp.Regexp
}

// CanHandle checks whether the Request url matches the provided Regexp.
func (extractor *RegexCanHandle) CanHandle(req models.Request) bool {
	return extractor.UrlMatch.MatchString(req.Url)
}

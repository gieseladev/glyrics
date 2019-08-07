/*
Package sources contains the extractors for gLyrics.
*/
package sources

import (
	"github.com/gieseladev/glyrics/v3/pkg/lyrics"
	"github.com/gieseladev/glyrics/v3/pkg/request"
	"regexp"
	"sync"
)

var (
	registeredExtractors = make([]MaybeExtractor, 0)
	extractorsMux        sync.RWMutex
)

// RegisterExtractor adds a new Extractor to the registered sources.
// Registered extractors are returned by GetExtractorsForRequest.
func RegisterExtractor(e MaybeExtractor) {
	extractorsMux.Lock()
	defer extractorsMux.Unlock()

	registeredExtractors = append(registeredExtractors, e)
}

// GetExtractorsForRequest returns a slice of all extractor which can extract
// lyrics from the given request.
func GetExtractorsForRequest(req *request.Request) []Extractor {
	extractorsMux.RLock()
	defer extractorsMux.RUnlock()

	var extractors []Extractor
	for _, e := range registeredExtractors {
		if e.CanExtract(req) {
			extractors = append(extractors, e)
		}
	}

	return extractors
}

// Extractor extracts lyrics from a Request.
type Extractor interface {
	// ExtractLyrics performs the actual extraction.
	ExtractLyrics(req *request.Request) (*lyrics.Info, error)
}

// CanExtractTeller can tell whether lyrics can be extracted from the given
// request.
type CanExtractTeller interface {
	// CanExtract performs simple checks to determine whether the extractor
	// has any chance of extracting lyrics from the Request.
	CanExtract(req *request.Request) bool
}

// MaybeExtractor combines Extractor with CanExtractTeller.
type MaybeExtractor interface {
	CanExtractTeller
	Extractor
}

type maybeExtractor struct {
	CanExtractTeller
	Extractor
}

// CreateMaybeExtractor combines a CanExtractTeller and an Extractor to a
// MaybeExtractor.
func CreateMaybeExtractor(teller CanExtractTeller, extractor Extractor) MaybeExtractor {
	return maybeExtractor{
		CanExtractTeller: teller,
		Extractor:        extractor,
	}
}

// ExtractorFunc is a function which implements the Extractor interface.
type ExtractorFunc func(req *request.Request) (*lyrics.Info, error)

func (e ExtractorFunc) ExtractLyrics(req *request.Request) (*lyrics.Info, error) {
	return e(req)
}

type regexExtractorTeller struct {
	*regexp.Regexp
}

// RegexExtractorTeller wraps the regular expression in a struct that implements
// CanExtractTeller.
func RegexExtractorTeller(re *regexp.Regexp) CanExtractTeller {
	return &regexExtractorTeller{re}
}

func (e *regexExtractorTeller) CanExtract(req *request.Request) bool {
	return e.MatchString(req.Url)
}

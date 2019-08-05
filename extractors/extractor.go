/*
Package extractors contains the extractors for gLyrics.
*/
package extractors

import (
	"github.com/gieseladev/glyrics/pkg/requests"
	"regexp"
	"sync"
	"time"
)

// LyricsOrigin contains the details regarding the origin of lyrics.
type LyricsOrigin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// LyricsInfo represents a song's lyrics and metadata.
type LyricsInfo struct {
	Url         string       `json:"url"`
	Title       string       `json:"title"`
	Artist      string       `json:"artist"`
	Lyrics      string       `json:"lyrics"`
	ReleaseDate time.Time    `json:"release_date,omitempty"`
	Origin      LyricsOrigin `json:"origin,omitempty"`
}

var (
	registeredExtractors = make([]MaybeExtractor, 0)
	extractorsMux        sync.RWMutex
)

// RegisterExtractor adds a new Extractor to the registered extractors.
// Registered extractors are returned by GetExtractorsForRequest.
func RegisterExtractor(e MaybeExtractor) {
	extractorsMux.Lock()
	defer extractorsMux.Unlock()

	registeredExtractors = append(registeredExtractors, e)
}

// GetExtractorsForRequest returns a slice of all extractors which can extract
// lyrics from the given request.
func GetExtractorsForRequest(req *requests.Request) []Extractor {
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
	ExtractLyrics(req *requests.Request) (*LyricsInfo, error)
}

// CanExtractTeller can tell whether lyrics can be extracted from the given
// request.
type CanExtractTeller interface {
	// CanExtract performs simple checks to determine whether the extractor
	// has any chance of extracting lyrics from the Request.
	CanExtract(req *requests.Request) bool
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
type ExtractorFunc func(req *requests.Request) (*LyricsInfo, error)

func (e ExtractorFunc) ExtractLyrics(req *requests.Request) (*LyricsInfo, error) {
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

func (e *regexExtractorTeller) CanExtract(req *requests.Request) bool {
	return e.MatchString(req.Url)
}

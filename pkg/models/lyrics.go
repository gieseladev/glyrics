package models

import (
	"time"
)

// LyricsOrigin contains metadata regarding the extractor
// which extracted the lyrics.
type LyricsOrigin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Lyrics represents a song's lyrics and metadata.
// The object is JSONifiable.
type Lyrics struct {
	Url         string        `json:"url"`
	Title       string        `json:"title"`
	Artist      string        `json:"artist"`
	Lyrics      string        `json:"lyrics"`
	ReleaseDate time.Time     `json:"release_date,omitempty"`
	Origin      *LyricsOrigin `json:"origin,omitempty"`
}

package models

import (
	"time"
)

type LyricsOrigin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Lyrics struct {
	Url         string        `json:"url"`
	Title       string        `json:"title"`
	Artist      string        `json:"artist"`
	Lyrics      string        `json:"lyrics"`
	ReleaseDate time.Time     `json:"release_date,omitempty"`
	Origin      *LyricsOrigin `json:"origin,omitempty"`
}

/*
Package lyrics provides the lyrics information and origin models.
*/
package lyrics

import "time"

// Origin contains the details regarding the origin of lyrics.
type Origin struct {
	Name    string `json:"name"`
	Website string `json:"url"`
}

// Info represents a song's lyrics and metadata.
type Info struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Lyrics      string    `json:"lyrics"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Origin      Origin    `json:"origin,omitempty"`
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gieseladev/glyrics/v3"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var errlog = log.New(os.Stderr, "", 0)

func writeHeadline(w io.Writer, lyrics *glyrics.LyricsInfo) int {
	var total int

	n, _ := fmt.Fprint(w, lyrics.Title)
	total += n

	if lyrics.Artist != "" {
		n, _ := fmt.Fprintf(w, " by %s", lyrics.Artist)
		total += n
	}

	if lyrics.ReleaseDate != (time.Time{}) {
		n, _ := fmt.Fprintf(w, " (%d)", lyrics.ReleaseDate.Year())
		total += n
	}

	return total
}

func printLyrics(lyrics *glyrics.LyricsInfo) {
	var headlineBuilder strings.Builder
	writeHeadline(&headlineBuilder, lyrics)
	headline := headlineBuilder.String()
	underline := strings.Repeat("=", utf8.RuneCountInString(headline))

	fmt.Printf("%s\n%s\n\n%s\n\nfrom %s",
		headline, underline, lyrics.Lyrics, lyrics.Origin.Website,
	)
}

func printLyricsFormat(c *cli.Context, lyrics *glyrics.LyricsInfo) {
	format := c.GlobalString("format")

	switch format {
	case "json":
		_ = json.NewEncoder(os.Stdout).Encode(lyrics)
	default:
		printLyrics(lyrics)
	}
}

func printMultipleLyricsOverview(infos []glyrics.LyricsInfo) {
	w := os.Stdout

	var avgWritten float64

	for i, info := range infos {
		_, _ = fmt.Fprintf(w, "%2d. ", i+1)
		written := writeHeadline(w, &info)

		if writtenF64 := float64(written); writtenF64 > avgWritten || avgWritten-writtenF64 > 50 {
			avgWritten = writtenF64
		}

		padding := int(avgWritten) - written
		if padding > 0 {
			_, _ = w.WriteString(strings.Repeat(" ", padding))
		}

		_, _ = fmt.Fprintf(w, " (%s)\n", info.URL)
	}
}

func printMultipleLyricsFormat(c *cli.Context, infos []glyrics.LyricsInfo) {
	format := c.GlobalString("format")

	switch format {
	case "json":
		_ = json.NewEncoder(os.Stdout).Encode(infos)
	default:
		printMultipleLyricsOverview(infos)
	}
}

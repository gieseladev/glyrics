package main

import (
	"fmt"
	"github.com/gieseladev/lyricsfinder/internal"
	"github.com/gieseladev/lyricsfinder/pkg"
	"github.com/gieseladev/lyricsfinder/pkg/models"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
	"time"
)

func printLyrics(lyrics *models.Lyrics) {
	headlineBuilder := strings.Builder{}
	headlineBuilder.WriteString(lyrics.Title)
	if lyrics.Artist != "" {
		headlineBuilder.WriteString(" by " + lyrics.Artist)
	}
	if lyrics.ReleaseDate != (time.Time{}) {
		headlineBuilder.WriteString(" (" + string(lyrics.ReleaseDate.Year()) + ")")
	}

	headline := headlineBuilder.String()
	underline := strings.Repeat("=", len(headline))

	lyricsText := fmt.Sprintf("%s\n%s\n\n%s\n\nfrom %s",
		headline, underline, lyrics.Lyrics, lyrics.Origin.Url,
	)

	log.Print(lyricsText)
}

func searchLyrics(c *cli.Context) {
	query := strings.Join(c.Args(), " ")
	apiKey := c.String("token")
	if apiKey == "" {
		config, err := internal.GetConfig()
		if err != nil {
			log.Fatal("No token passed and couldn't load config file", err)
		}

		apiKey = config.GoogleApiKey
	}

	ch := make(chan models.Lyrics, 1)
	go lyricsfinder.SearchLyrics(query, apiKey, ch)

	lyrics := <-ch
	if lyrics != (models.Lyrics{}) {
		printLyrics(&lyrics)
	} else {
		log.Fatal("Couldn't find any results!")
	}
}

func extractLyrics(c *cli.Context) {
	url := c.Args().First()

	lyrics, err := lyricsfinder.ExtractLyrics(url)
	if err != nil {
		log.Fatal("Couldn't extract lyrics", err)
	}

	printLyrics(lyrics)
}

func main() {
	app := cli.NewApp()
	app.Name = "lyricsfinder"
	app.Version = "2.1.1"

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search for lyrics",
			Action:  searchLyrics,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "token",
					Usage:  "Google api key for custom search",
					EnvVar: "GOOGLE_API_KEY",
					Value:  "",
				},
			},
		},
		{
			Name:    "extract",
			Aliases: []string{"e"},
			Usage:   "Extract lyrics from url",
			Action:  extractLyrics,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
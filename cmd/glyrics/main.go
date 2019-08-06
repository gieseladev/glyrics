package main

import (
	"fmt"
	"github.com/gieseladev/glyrics/v3"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
	"time"
)

func printLyrics(lyrics *glyrics.LyricsInfo) {
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

	config, err := GetConfig()
	if apiKey == "" {
		if err != nil {
			log.Fatal("No token passed and couldn't load config file: ", err)
		}

		apiKey = config.GoogleApiKey
	} else if err == nil {
		config.GoogleApiKey = apiKey
		_ = config.SaveConfig()
	} else {
		_ = CliConfig{GoogleApiKey: apiKey}.SaveConfig()
	}

	lyrics := glyrics.SearchFirstLyrics(query, apiKey)

	if lyrics != nil {
		printLyrics(lyrics)
	} else {
		log.Fatal("Couldn't find any results!")
	}
}

func extractLyrics(c *cli.Context) {
	url := c.Args().First()

	lyrics, err := glyrics.ExtractLyrics(url)
	if err != nil {
		log.Fatal("Couldn't extract lyrics: ", err)
	}

	printLyrics(lyrics)
}

func main() {
	app := cli.NewApp()
	app.Name = "gLyrics"
	app.Usage = "find the lyrics you've always wanted to find"
	app.Description = "This is a command line tool to access the power of gLyrics."
	app.Version = "2.2.1"

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search for lyrics",
			Description: "Uses google custom search to find the most accurate lyrics for you. " +
				"Requires an api key with access to the custom search api. " +
				"This token is only required the first time.",
			ArgsUsage: "query",
			Action:    searchLyrics,
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
			Name:      "extract",
			Aliases:   []string{"e"},
			Usage:     "Extract lyrics from url",
			ArgsUsage: "url",
			Action:    extractLyrics,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

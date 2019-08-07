package main

import (
	"context"
	"fmt"
	"github.com/gieseladev/glyrics/v3"
	"github.com/gieseladev/glyrics/v3/pkg/search"
	"github.com/urfave/cli"
	"os"
	"strings"
	"time"
	"unicode/utf8"
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
	underline := strings.Repeat("=", utf8.RuneCountInString(headline))

	lyricsText := fmt.Sprintf("%s\n%s\n\n%s\n\nfrom %s",
		headline, underline, lyrics.Lyrics, lyrics.Origin.Website,
	)

	fmt.Print(lyricsText)
}

func searchLyrics(c *cli.Context) {
	query := strings.Join(c.Args(), " ")
	apiKey := c.String("token")

	config, err := GetConfig()
	if apiKey == "" {
		if err != nil {
			fmt.Print("No token passed and couldn't load config file: ", err)
			os.Exit(1)
		}

		apiKey = config.GoogleApiKey
	} else if err == nil {
		config.GoogleApiKey = apiKey
		_ = config.SaveConfig()
	} else {
		_ = CliConfig{GoogleApiKey: apiKey}.SaveConfig()
	}

	searcher := &search.GoogleSearcher{APIKey: apiKey}
	lyrics := glyrics.SearchFirst(context.Background(), searcher, query)

	if lyrics != nil {
		printLyrics(lyrics)
	} else {
		fmt.Print("Couldn't find any results!")
	}
}

func extractLyrics(c *cli.Context) {
	url := c.Args().First()

	lyrics, err := glyrics.Extract(url)
	if err != nil {
		fmt.Print("Couldn't extract lyrics: ", err)
		os.Exit(1)
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
		fmt.Print(err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gieseladev/glyrics/v3"
	"github.com/gieseladev/glyrics/v3/pkg/search"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getConfigLocation(c *cli.Context) (string, error) {
	if c.GlobalBool("no-config") {
		return "", errors.New("config suppressed")
	}

	return c.GlobalString("config"), nil
}

func getConfig(c *cli.Context) (*CliConfig, error) {
	location, err := getConfigLocation(c)
	if err != nil {
		return nil, err
	}

	return LoadConfig(location)
}

func putConfig(c *cli.Context, config *CliConfig) {
	if location, err := getConfigLocation(c); err == nil {
		if err := SaveConfig(location, config); err != nil {
			errlog.Print("couldn't save config: ", err)
		}
	}
}

func getApiKey(c *cli.Context) (string, error) {
	apiKey := c.String("api-key")

	config, err := getConfig(c)

	if apiKey == "" {
		if err != nil {
			return "", errors.New("no api specified and couldn't load config: " + err.Error())
		}

		return config.GoogleApiKey, nil
	}

	if config == nil {
		config = &CliConfig{GoogleApiKey: apiKey}
	} else {
		config.GoogleApiKey = apiKey
	}

	putConfig(c, config)

	return apiKey, nil
}

func findLyricsCommand(c *cli.Context) {
	query := strings.Join(c.Args(), " ")
	apiKey, err := getApiKey(c)
	if err != nil {
		errlog.Fatal(err)
	}

	searcher := &search.Google{APIKey: apiKey}
	lyrics := glyrics.SearchFirst(context.Background(), searcher, query)

	if lyrics != nil {
		printLyricsFormat(c, lyrics)
	} else {
		fmt.Print("Couldn't find any results!")
	}
}

func searchLyricsCommand(c *cli.Context) {
	query := strings.Join(c.Args(), " ")
	apiKey, err := getApiKey(c)
	if err != nil {
		errlog.Fatal(err)
	}

	searcher := &search.Google{APIKey: apiKey}
	infos := glyrics.SearchN(context.Background(), searcher, query, c.Int("results"))
	printMultipleLyricsFormat(c, infos)
}

func extractLyricsCommand(c *cli.Context) {
	url := c.Args().First()

	lyrics, err := glyrics.Extract(url)
	if err != nil {
		errlog.Fatal("Couldn't extract lyrics: ", err)
	}

	printLyricsFormat(c, lyrics)
}

// Modified regex from: https://stackoverflow.com/a/3809435
var urlRegex = regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)`)

func smartCommand(c *cli.Context) {
	query := strings.Join(c.Args(), "")
	if urlRegex.MatchString(query) {
		extractLyricsCommand(c)
	} else {
		findLyricsCommand(c)
	}
}

func GetApp() *cli.App {
	app := cli.NewApp()
	app.Name = "gLyrics"
	app.Usage = "extract and find lyrics"
	app.Description = "Command line tool to access gLyrics."
	app.Version = "3.1.1"

	var defaultConfigLocation string
	if homedir, err := os.UserHomeDir(); err == nil {
		defaultConfigLocation = filepath.Join(homedir, ".glyrics")
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "config file location",
			EnvVar: "GLYRICS_CONFIG",
			Value:  defaultConfigLocation,
		},
		cli.BoolFlag{
			Name:  "no-config",
			Usage: "suppress the usage and creation of a config file",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "output format",
		},
	}

	apiKeyFlag := cli.StringFlag{
		Name:   "api-key",
		Usage:  "google api key with access to the custom search api",
		EnvVar: "GOOGLE_API_KEY",
	}

	app.Commands = []cli.Command{
		{
			Name:      "extract",
			Aliases:   []string{"e"},
			Usage:     "Extract lyrics from url",
			ArgsUsage: "url",
			Action:    extractLyricsCommand,
		},
		{
			Name:    "find",
			Aliases: []string{"f"},
			Usage:   "Search for lyrics and print the first result",
			Description: "Uses the google custom search engine to find the most accurate lyrics for you. " +
				"Requires an api key with access to the custom search api. " +
				"This token is stored in the config file and so is only required the first time.",
			ArgsUsage: "query",
			Action:    findLyricsCommand,
			Flags: []cli.Flag{
				apiKeyFlag,
			},
		},
		{
			Name:      "search",
			Aliases:   []string{"s"},
			Usage:     "Search for lyrics and show an overview of the results",
			ArgsUsage: "query",
			Action:    searchLyricsCommand,
			Flags: []cli.Flag{
				apiKeyFlag,
				cli.IntFlag{
					Name:  "results",
					Usage: "amount of results to display",
					Value: 5,
				},
			},
		},
	}

	app.Action = smartCommand

	return app
}

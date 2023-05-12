package main

import (
	"errors"
	"log"
	"os"

	_ "embed"

	"github.com/jamesjarvis/mappyboi/v2/pkg/base"
	"github.com/jamesjarvis/mappyboi/v2/pkg/input/google"
	"github.com/jamesjarvis/mappyboi/v2/pkg/parser"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version string

var (
	baseFileFlag              = "base_file"
	googleLocationHistoryFlag = "google_location_history"
)

func mustCreateFileIfNotExists(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		return
	}
	if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return
}

func app(c *cli.Context) error {
	log.Println("Mappyboi " + c.App.Version)

	mustCreateFileIfNotExists(c.Path(baseFileFlag))

	// Setup base.
	log.Printf("Loading Base file from %s...\n", c.Path(baseFileFlag))
	baseLocationHistory, err := base.ReadBase(c.Path(baseFileFlag))
	if err != nil {
		return err
	}
	log.Printf("Loaded Base file from %s, %d entries\n", c.Path(baseFileFlag), len(baseLocationHistory.Data))

	// Parse additional files and fold back into base.
	{
		var parsers []parser.Parser
		if c.IsSet(googleLocationHistoryFlag) {
			parsers = append(parsers, &google.LocationHistory{
				Filepath: c.Path(googleLocationHistoryFlag),
			})
		}
		log.Printf("Parsing supplied location files...\n")
		parsedLocationHistory, err := parser.ParseAll(parsers...)
		if err != nil {
			return err
		}
		log.Printf("Parsed %d entries from supplied location files\n", len(parsedLocationHistory.Data))
		baseLocationHistory.Insert(parsedLocationHistory.Data...)
		log.Printf("Combined all locations into %d entries\n", len(baseLocationHistory.Data))
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "mappyboi v2",
		Usage: "Store all Google Takeout / Apple Health exports and transform to custom outputs",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      baseFileFlag,
				Aliases:   []string{"base"},
				Usage:     "Base location history append only `FILE`",
				TakesFile: true,
				Required:  true,
			},
			&cli.PathFlag{
				Name:      googleLocationHistoryFlag,
				Aliases:   []string{"glh"},
				Usage:     "Google Takeout Location History `FILE`",
				TakesFile: true,
			},
			cli.VersionFlag,
		},
		Action: app,
	}

	app.Version = version
	app.EnableBashCompletion = true
	app.Authors = []*cli.Author{
		{
			Name:  "James Jarvis",
			Email: "git@jamesjarvis.io",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

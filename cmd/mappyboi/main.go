package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jamesjarvis/mappyboi/pkg/models"
	"github.com/jamesjarvis/mappyboi/pkg/parser"
	"github.com/urfave/cli/v2"
)

const (
	googleFlag = "location_history"
	gpxFlag    = "gpx_folder"
)

func PrintStats(data *models.Data) {
	fmt.Printf("Total gps points: %d\n", len(data.GoLocations))
}

func main() {
	app := &cli.App{
		Name:  "mappyboi",
		Usage: "Make a heatmap out of Google Takeout / Apple Health exports",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    googleFlag,
				Aliases: []string{"lh"},
				Usage:   "Load Google Location History from `FILE`",
			},
			&cli.PathFlag{
				Name:    gpxFlag,
				Aliases: []string{"gpx"},
				Usage:   "Load GPX files from from `FOLDER`",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Mappyboi v0.1")

			googleParser := &parser.GoogleLocationHistory{}
			if c.IsSet(googleFlag) {
				googleParser = &parser.GoogleLocationHistory{
					Filepath: c.Path(googleFlag),
				}
			}

			gpxParsers := []*parser.GPXFile{}
			if c.IsSet(gpxFlag) {
				gpxs, err := parser.FindGPXFiles(c.Path(gpxFlag))
				if err != nil {
					return err
				}
				gpxParsers = gpxs
			}

			// TODO: Clean this up, couldn't get the variadic thing to stop complaining.
			parsers := []parser.Parser{}
			parsers = append(parsers, googleParser)
			for _, p := range gpxParsers {
				parsers = append(parsers, p)
			}

			allData, err := parser.ParseAll(parsers...)
			if err != nil {
				return err
			}

			PrintStats(allData)

			return nil
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

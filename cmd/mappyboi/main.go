package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jamesjarvis/mappyboi/pkg/maptemplate"
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

			parsers := []parser.Parser{}

			if c.IsSet(googleFlag) {
				parsers = append(parsers, &parser.GoogleLocationHistory{
					Filepath: c.Path(googleFlag),
				})
			}

			if c.IsSet(gpxFlag) {
				gpxs, err := parser.FindGPXFiles(c.Path(gpxFlag))
				if err != nil {
					return err
				}
				for _, p := range gpxs {
					parsers = append(parsers, p)
				}
			}

			if len(parsers) == 0 {
				return nil
			}

			allData, err := parser.ParseAll(parsers...)
			if err != nil {
				return err
			}

			PrintStats(allData)

			return maptemplate.GenerateHTML(allData)
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

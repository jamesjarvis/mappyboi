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
	googleFlag        = "location_history"
	gpxFlag           = "gpx_folder"
	outputFlag        = "output_file"
	defaultOutputFile = "map.html"
	version 					= "v1.0.1"
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
				Name:      googleFlag,
				Aliases:   []string{"lh"},
				Usage:     "Load Google Location History from `FILE`",
				TakesFile: true,
			},
			&cli.StringSliceFlag{
				Name:    gpxFlag,
				Aliases: []string{"gpx"},
				Usage:   "Load GPX files from `FOLDER`",
			},
			&cli.PathFlag{
				Name:      outputFlag,
				Aliases:   []string{"o"},
				Usage:     "Output `FILE` to export heatmap to. Must be .html format",
				TakesFile: true,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Mappyboi " + version)

			// Locate data from flags
			parsers := []parser.Parser{}
			if c.IsSet(googleFlag) {
				parsers = append(parsers, &parser.GoogleLocationHistory{
					Filepath: c.Path(googleFlag),
				})
			}
			if c.IsSet(gpxFlag) {
				gpxs := c.StringSlice(gpxFlag)
				for _, g := range gpxs {
					gpxs, err := parser.FindGPXFiles(g)
					if err != nil {
						return err
					}
					for _, p := range gpxs {
						parsers = append(parsers, p)
					}
				}
			}

			// Exit early if no data passed in
			if len(parsers) == 0 {
				fmt.Println("You need some data for a heatmap, silly!")
				return nil
			}

			// Get the output file path
			outputPath := defaultOutputFile
			if c.IsSet(outputFlag) {
				outputPath = c.Path(outputFlag)
			}

			// mmm consume the data
			allData, err := parser.ParseAll(parsers...)
			if err != nil {
				return err
			}

			PrintStats(allData)

			return maptemplate.GenerateHTML(outputPath, allData)
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

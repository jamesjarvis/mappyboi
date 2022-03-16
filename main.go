package main

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/jamesjarvis/mappyboi/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/pkg/maptemplate"
	"github.com/jamesjarvis/mappyboi/pkg/models"
	"github.com/jamesjarvis/mappyboi/pkg/parser"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version string

const (
	googleFlag        = "location_history"
	gpxFlag           = "gpx_folder"
	outputFlag        = "output_file"
	defaultOutputFile = "map.html"
	minDistanceFlag   = "min_distance"
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
			&cli.Float64Flag{
				Name:    minDistanceFlag,
				Aliases: []string{"min"},
				Usage:   "If you struggle to open the file in a browser due to too many points, reduce the number of points by increasing this value.",
			},
			cli.VersionFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Mappyboi " + c.App.Version)

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
						return fmt.Errorf("error finding gpx files for %s: %w", g, err)
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
				return fmt.Errorf("error parsing data: %w", err)
			}

			// Simplify routes to minimise number of points.
			// Unfortunately leaflet will stack overflow after around 600k points :'(

			if c.IsSet(minDistanceFlag) {
				minDistance := c.Float64(minDistanceFlag)
				allData, err = conversions.ReducePoints(allData, minDistance)
				if err != nil {
					return fmt.Errorf("failed to reduce points to %f: %w", minDistance, err)
				}
			}

			PrintStats(allData)

			return maptemplate.GenerateHTML(outputPath, allData)
		},
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

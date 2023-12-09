package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/jamesjarvis/mappyboi/v2/pkg/base"
	"github.com/jamesjarvis/mappyboi/v2/pkg/input/fit"
	"github.com/jamesjarvis/mappyboi/v2/pkg/input/google"
	"github.com/jamesjarvis/mappyboi/v2/pkg/input/gpx"
	"github.com/jamesjarvis/mappyboi/v2/pkg/input/polarsteps"
	"github.com/jamesjarvis/mappyboi/v2/pkg/maptemplate"
	"github.com/jamesjarvis/mappyboi/v2/pkg/parser"
	"github.com/jamesjarvis/mappyboi/v2/pkg/transform"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version string

var (
	baseFileFlag = "base_file"
	// Input
	googleLocationHistoryFlag = "google_location_history"
	gpxDirectoryFlag          = "gpx_directory"
	fitDirectoryFlag          = "fit_directory"
	polarstepDirectoryFlag    = "polarstep_directory"
	// Output
	outputTypeFlag = "output_type"
	outputFileFlag = "output_file"
	// Transformations
	outputTransformReducePointsFlag = "output_reduce_points"
	outputRandomisePoints           = "output_randomise_points"
)

type output string

const (
	output_UNKNOWN output = "UNKNOWN"
	output_MAP     output = "MAP"
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
}

func mustConvertOutputType(o string) output {
	switch o {
	case string(output_MAP):
		return output_MAP
	}
	return output_UNKNOWN
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
	originalLength := len(baseLocationHistory.Data)
	log.Printf("Loaded Base file from %s, %d entries\n", c.Path(baseFileFlag), originalLength)

	// Parse additional files and fold back into base.
	var parsers []parser.Parser
	if c.IsSet(googleLocationHistoryFlag) {
		parsers = append(parsers, &google.LocationHistory{
			Filepath: c.Path(googleLocationHistoryFlag),
		})
	}
	if c.IsSet(gpxDirectoryFlag) {
		gpxs, err := gpx.FindGPXFiles(c.Path(gpxDirectoryFlag))
		if err != nil {
			return err
		}
		for _, p := range gpxs {
			parsers = append(parsers, p)
		}
	}
	if c.IsSet(fitDirectoryFlag) {
		fits, err := fit.FindFitFiles(c.Path(fitDirectoryFlag))
		if err != nil {
			return err
		}
		for _, p := range fits {
			parsers = append(parsers, p)
		}
	}
	if c.IsSet(polarstepDirectoryFlag) {
		psteps, err := polarsteps.FindPolarstepsFiles(c.Path(polarstepDirectoryFlag))
		if err != nil {
			return err
		}
		for _, p := range psteps {
			parsers = append(parsers, p)
		}
	}
	if len(parsers) > 0 {
		log.Printf("Parsing supplied location files...\n")
		parsedLocationHistory, err := parser.ParseAll(parsers...)
		if err != nil {
			return err
		}
		log.Printf("Parsed %d entries from supplied location files\n", len(parsedLocationHistory.Data))
		baseLocationHistory.Insert(parsedLocationHistory.Data...)
		log.Printf("Combined all locations into %d entries (%d new)\n", len(baseLocationHistory.Data), len(baseLocationHistory.Data)-originalLength)
	}

	// Cleanup location history.
	err = baseLocationHistory.Cleanup()
	if err != nil {
		return err
	}

	// Write to base.
	if len(baseLocationHistory.Data) > originalLength {
		log.Printf("Writing %d entries to Base file %s...", len(baseLocationHistory.Data), c.Path(baseFileFlag))
		err = base.WriteBase(c.Path(baseFileFlag), baseLocationHistory)
		if err != nil {
			return err
		}
		log.Printf("Completed writing to Base file %s", c.Path(baseFileFlag))
	} else {
		log.Printf("Skipping write, as data is unchanged\n")
	}

	// Run output transformations.

	// Simplify routes to minimise number of points.
	// Unfortunately leaflet will stack overflow after around 600k points :'(
	if c.IsSet(outputTransformReducePointsFlag) {
		minDistance := c.Float64(outputTransformReducePointsFlag)
		baseLocationHistory, err = transform.ReducePoints(baseLocationHistory, minDistance)
		if err != nil {
			return fmt.Errorf("failed to reduce points to %f: %w", minDistance, err)
		}
	}
	// Randomise output.
	if c.IsSet(outputRandomisePoints) {
		baseLocationHistory, err = transform.RandomisePoints(baseLocationHistory)
		if err != nil {
			return fmt.Errorf("failed to shuffle points: %w", err)
		}
	}

	// Generate output.
	if c.IsSet(outputTypeFlag) && c.IsSet(outputFileFlag) {
		log.Printf("Generating output of type %s to %s...\n", c.String(outputTypeFlag), c.Path(outputFileFlag))
		outputType := mustConvertOutputType(c.String(outputTypeFlag))
		switch outputType {
		case output_UNKNOWN:
			return fmt.Errorf("invalid output type %s", c.String(outputTypeFlag))
		case output_MAP:
			err = maptemplate.GenerateHTML(c.Path(outputFileFlag), baseLocationHistory)
			if err != nil {
				return err
			}
		}
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
				Usage:     "Base location history append only `FILE`, in .json or .json.gz",
				TakesFile: true,
				Required:  true,
			},
			&cli.PathFlag{
				Name:      googleLocationHistoryFlag,
				Aliases:   []string{"glh"},
				Usage:     "Google Takeout Location History `FILE`",
				TakesFile: true,
			},
			&cli.PathFlag{
				Name:      gpxDirectoryFlag,
				Aliases:   []string{"gpxd"},
				Usage:     "GPX `DIRECTORY` to load .gpx files from",
				TakesFile: false,
			},
			&cli.PathFlag{
				Name:      fitDirectoryFlag,
				Aliases:   []string{"fitd"},
				Usage:     "FIT `DIRECTORY` to load .fit files from",
				TakesFile: false,
			},
			&cli.PathFlag{
				Name:      polarstepDirectoryFlag,
				Aliases:   []string{"pstepd"},
				Usage:     "Polarstep `DIRECTORY` to load locations.json files from",
				TakesFile: false,
			},
			&cli.StringFlag{
				Name:    outputTypeFlag,
				Aliases: []string{"ot"},
				Usage:   "Output format, must be one of [MAP]",
			},
			&cli.PathFlag{
				Name:      outputFileFlag,
				Aliases:   []string{"of"},
				Usage:     "Output `FILE` to write to",
				TakesFile: true,
			},
			&cli.Float64Flag{
				Name:    outputTransformReducePointsFlag,
				Aliases: []string{"rp"},
				Usage:   "If you struggle to open the file in a browser due to too many points, reduce the number of points by increasing this value.",
			},
			&cli.BoolFlag{
				Name:    outputRandomisePoints,
				Aliases: []string{"rand"},
				Usage:   "If you want to export the view of the points, but otherwise randomise the data to prevent perfect tracking, this will randomise the order.",
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

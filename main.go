package main

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version string

var (
	baseFileFlag = "base_file"
)

func app(c *cli.Context) error {
	fmt.Println("Mappyboi " + c.App.Version)

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

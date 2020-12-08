package parser

import (
	"encoding/json"

	"github.com/jamesjarvis/mappyboi/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/pkg/models"
	"github.com/tkrajina/gpxgo/gpx"
)

type Parser interface {
	String() string
	Parse() (*models.Data, error)
}

type GoogleLocationHistory struct {
	Filepath string
}

func (p *GoogleLocationHistory) String() string {
	return p.Filepath
}

func (p *GoogleLocationHistory) Parse() (*models.Data, error) {
	var data models.GoogleData

	file, err := Load(p.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return conversions.GoogleDataToData(&data)
}

type GPXFile struct {
	Filepath string
}

func (p *GPXFile) String() string {
	return p.Filepath
}

func (p *GPXFile) Parse() (*models.Data, error) {
	g, err := gpx.ParseFile(p.Filepath)
	if err != nil {
		return nil, err
	}

	// Reduce number of points in GPX track.
	g.ReduceGpxToSingleTrack()
	g.SimplifyTracks(1.5)

	// Add to data object.
	data := &models.Data{
		GoLocations: make([]*models.GoLocation, 0, g.GetTrackPointsNo()),
	}
	g.ExecuteOnAllPoints(func(gpxPoint *gpx.GPXPoint) {
		data.GoLocations = append(data.GoLocations, conversions.GPXPointToGoLocation(gpxPoint))
	})

	return data, nil
}

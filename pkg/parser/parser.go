package parser

import (
	"github.com/jamesjarvis/mappyboi/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/pkg/models"
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tkrajina/gpxgo/gpx"
)

type Parser interface {
	String() string
	Parse() (types.LocationHistory, error)
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

	// Add to data object.
	data := &models.Data{
		GoLocations: make([]*models.GoLocation, 0, g.GetTrackPointsNo()),
	}
	g.ExecuteOnAllPoints(func(gpxPoint *gpx.GPXPoint) {
		data.GoLocations = append(data.GoLocations, conversions.GPXPointToGoLocation(gpxPoint))
	})

	return data, nil
}

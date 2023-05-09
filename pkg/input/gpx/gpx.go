package gpx

import (
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tkrajina/gpxgo/gpx"
)

type GPXFile struct {
	Filepath string
}

func (p *GPXFile) String() string {
	return p.Filepath
}

func (p *GPXFile) Parse() (types.LocationHistory, error) {
	g, err := gpx.ParseFile(p.Filepath)
	if err != nil {
		return types.LocationHistory{}, err
	}

	// Reduce number of points in GPX track.
	g.ReduceGpxToSingleTrack()

	// Add to data object.
	data := types.LocationHistory{
		Data: make([]types.Location, 0, g.GetTrackPointsNo()),
	}
	g.ExecuteOnAllPoints(func(gpxPoint *gpx.GPXPoint) {
		data.Insert(GPXPointToGoLocation(gpxPoint))
	})

	return data, nil
}

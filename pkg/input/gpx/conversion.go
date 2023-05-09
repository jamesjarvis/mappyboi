package gpx

import (
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tkrajina/gpxgo/gpx"
)

func GPXPointToGoLocation(gpxPoint *gpx.GPXPoint) types.Location {
	return types.Location{
		Time:             gpxPoint.Timestamp,
		Latitude:         gpxPoint.Latitude,
		Longitude:        gpxPoint.Longitude,
		Accuracy:         gpxPoint.HorizontalDilution.Value(),
		Altitude:         gpxPoint.Elevation.Value(),
		VerticalAccuracy: gpxPoint.VerticalDilution.Value(),
	}
}

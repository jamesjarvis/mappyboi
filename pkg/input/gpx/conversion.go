package gpx

import (
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tkrajina/gpxgo/gpx"
)

func gpxPointToGoLocation(gpxPoint *gpx.GPXPoint) types.Location {
	return types.Location{
		Time:             gpxPoint.Timestamp,
		Latitude:         gpxPoint.Latitude,
		Longitude:        gpxPoint.Longitude,
		Accuracy:         int(gpxPoint.HorizontalDilution.Value()),
		Altitude:         int(gpxPoint.Elevation.Value()),
		VerticalAccuracy: int(gpxPoint.VerticalDilution.Value()),
	}
}

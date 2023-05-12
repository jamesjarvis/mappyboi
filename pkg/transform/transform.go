package transform

import (
	"log"
	"sort"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tkrajina/gpxgo/gpx"
)

func gPXPointToGoLocation(gpxPoint gpx.GPXPoint) types.Location {
	return types.Location{
		Time:             gpxPoint.Timestamp,
		Latitude:         gpxPoint.Latitude,
		Longitude:        gpxPoint.Longitude,
		Altitude:         int(gpxPoint.Elevation.Value()),
		Accuracy:         int(gpxPoint.HorizontalDilution.Value()),
		VerticalAccuracy: int(gpxPoint.VerticalDilution.Value()),
	}
}

func goLocationToGPXPoint(goLocation types.Location) gpx.GPXPoint {
	return gpx.GPXPoint{
		Timestamp:          goLocation.Time,
		HorizontalDilution: *gpx.NewNullableFloat64(float64(goLocation.Accuracy)),
		VerticalDilution:   *gpx.NewNullableFloat64(float64(goLocation.VerticalAccuracy)),
		Point: gpx.Point{
			Latitude:  goLocation.Latitude,
			Longitude: goLocation.Longitude,
			Elevation: *gpx.NewNullableFloat64(float64(goLocation.Altitude)),
		},
	}
}

// ReducePoints reduces the number of points from locs, by enforcing a minimum distance between points.
// Note that this only reduces the points that are chronological, if there are multiple points at the same
// place on different days this will not be reduced.
func ReducePoints(locs types.LocationHistory, minDistance float64) (types.LocationHistory, error) {
	log.Printf("Starting point reduction, %d points with a new mininum distance of %.1fm\n", len(locs.Data), minDistance)

	points := make([]gpx.GPXPoint, 0, len(locs.Data))
	for _, p := range locs.Data {
		points = append(points, goLocationToGPXPoint(p))
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})

	gpxTrack := &gpx.GPXTrackSegment{
		Points: points,
	}
	gpxTrack.ReduceTrackPoints(minDistance)

	log.Printf("Reduced minimum distance between points to %.1fm, %d points (%.0f%% reduction)\n", minDistance, len(gpxTrack.Points), (1-(float64(len(gpxTrack.Points))/float64(len(locs.Data))))*100)

	newData := types.LocationHistory{
		Data: make([]types.Location, 0, len(gpxTrack.Points)),
	}
	for _, p := range gpxTrack.Points {
		newData.Insert(gPXPointToGoLocation(p))
	}

	return newData, nil
}

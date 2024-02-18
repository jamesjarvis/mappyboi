package transform

import (
	"log"
	"math/rand"
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

// Transformer receives a location history, applies a given transformation, and returns
// an error if encountered.
type Transformer func(types.LocationHistory) error

// ProcessPoints applies the provided transformers to the provided location history points.
// It returns the transformed location history, or an error if encountered.
func ProcessPoints(
	locs types.LocationHistory,
	transformers ...Transformer,
) error {
	var err error
	for _, transformer := range transformers {
		err = transformer(locs)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithMinimumDistance returns a transformer that reduces the location history points, by enforcing a minimum distance between points.
// Note that this only reduces chronological points at the moment, and multiple overlapping points that are not sequential will not be
// reduced.
func WithMinimumDistance(minDistance float64) Transformer {
	return func(lh types.LocationHistory) error {
		log.Printf("Starting point reduction, %d points with a new mininum distance of %.1fm\n", len(locs.Data), minDistance)

		points := make([]gpx.GPXPoint, 0, len(lh.Data))
		for _, p := range lh.Data {
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

		lh = newData

		return nil
	}
}

// WithRandomOrder randomises the order of the provided points.
func WithRandomOrder() Transformer {
	return func(lh types.LocationHistory) error {
		log.Printf("Starting shuffle of %d points", len(lh.Data))
		rand.Shuffle(len(lh.Data), func(i, j int) {
			lh.Data[i], lh.Data[j] = lh.Data[j], lh.Data[i]
		})
		log.Printf("Shuffled order of %d points", len(lh.Data))
		return nil
	}
}

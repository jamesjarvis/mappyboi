package conversions

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"log"

	"github.com/jamesjarvis/mappyboi/pkg/models"
	"github.com/tkrajina/gpxgo/gpx"
)

func E7ToStandard(e7Pos int64) float64 {
	return float64(e7Pos) / float64(1e7)
}

func GoogleDataToData(gd *models.GoogleData) (*models.Data, error) {
	data := &models.Data{
		GoLocations: make([]*models.GoLocation, 0, len(gd.Locations)),
	}

	for _, gloc := range gd.Locations {
		var t time.Time

		if gloc.TimestampMs != "" {
			ms, err := strconv.ParseInt(gloc.TimestampMs, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestampMs: %w", err)
			}
			t = time.Unix(0, ms*int64(time.Millisecond))
		} else {
			t = gloc.Timestamp
		}

		data.GoLocations = append(data.GoLocations, &models.GoLocation{
			Time:             t,
			Latitude:         E7ToStandard(gloc.LatitudeE7),
			Longitude:        E7ToStandard(gloc.LongitudeE7),
			Accuracy:         float64(gloc.Accuracy),
			Altitude:         float64(gloc.Altitude),
			VerticalAccuracy: float64(gloc.VerticalAccuracy),
		})
	}

	return data, nil
}

func GPXPointToGoLocation(gpxPoint *gpx.GPXPoint) *models.GoLocation {
	return &models.GoLocation{
		Time:             gpxPoint.Timestamp,
		Latitude:         gpxPoint.Latitude,
		Longitude:        gpxPoint.Longitude,
		Accuracy:         gpxPoint.HorizontalDilution.Value(),
		Altitude:         gpxPoint.Elevation.Value(),
		VerticalAccuracy: gpxPoint.VerticalDilution.Value(),
	}
}

func GoLocationToGPXPoint(goLocation *models.GoLocation) gpx.GPXPoint {
	return gpx.GPXPoint{
		Timestamp:          goLocation.Time,
		HorizontalDilution: *gpx.NewNullableFloat64(goLocation.Accuracy),
		VerticalDilution:   *gpx.NewNullableFloat64(goLocation.VerticalAccuracy),
		Point: gpx.Point{
			Latitude:  goLocation.Latitude,
			Longitude: goLocation.Longitude,
			Elevation: *gpx.NewNullableFloat64(goLocation.Altitude),
		},
	}
}

func ReducePoints(data *models.Data, minDistance float64) (*models.Data, error) {
	log.Printf("Starting point reduction, %d points\n", len(data.GoLocations))

	points := make([]gpx.GPXPoint, 0, len(data.GoLocations))
	for _, p := range data.GoLocations {
		points = append(points, GoLocationToGPXPoint(p))
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})

	gpxTrack := &gpx.GPXTrackSegment{
		Points: points,
	}
	gpxTrack.ReduceTrackPoints(minDistance)

	log.Printf("Reduced minimum distance between points to %.1fm, %d points (%.0f%% reduction)\n", minDistance, len(gpxTrack.Points), (1-(float64(len(gpxTrack.Points))/float64(len(data.GoLocations))))*100)

	newData := &models.Data{
		GoLocations: make([]*models.GoLocation, 0, len(gpxTrack.Points)),
	}
	for _, p := range gpxTrack.Points {
		newData.GoLocations = append(newData.GoLocations, GPXPointToGoLocation(&p))
	}
	data = nil
	data = newData

	return data, nil
}

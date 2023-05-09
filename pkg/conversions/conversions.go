package conversions

// import (
// 	"sort"

// 	"log"

// 	"github.com/jamesjarvis/mappyboi/pkg/models"
// 	"github.com/tkrajina/gpxgo/gpx"
// )

// func GoLocationToGPXPoint(goLocation *models.GoLocation) gpx.GPXPoint {
// 	return gpx.GPXPoint{
// 		Timestamp:          goLocation.Time,
// 		HorizontalDilution: *gpx.NewNullableFloat64(goLocation.Accuracy),
// 		VerticalDilution:   *gpx.NewNullableFloat64(goLocation.VerticalAccuracy),
// 		Point: gpx.Point{
// 			Latitude:  goLocation.Latitude,
// 			Longitude: goLocation.Longitude,
// 			Elevation: *gpx.NewNullableFloat64(goLocation.Altitude),
// 		},
// 	}
// }

// func ReducePoints(data *models.Data, minDistance float64) (*models.Data, error) {
// 	log.Printf("Starting point reduction, %d points\n", len(data.GoLocations))

// 	points := make([]gpx.GPXPoint, 0, len(data.GoLocations))
// 	for _, p := range data.GoLocations {
// 		points = append(points, GoLocationToGPXPoint(p))
// 	}

// 	sort.Slice(points, func(i, j int) bool {
// 		return points[i].Timestamp.Before(points[j].Timestamp)
// 	})

// 	gpxTrack := &gpx.GPXTrackSegment{
// 		Points: points,
// 	}
// 	gpxTrack.ReduceTrackPoints(minDistance)

// 	log.Printf("Reduced minimum distance between points to %.1fm, %d points (%.0f%% reduction)\n", minDistance, len(gpxTrack.Points), (1-(float64(len(gpxTrack.Points))/float64(len(data.GoLocations))))*100)

// 	newData := &models.Data{
// 		GoLocations: make([]*models.GoLocation, 0, len(gpxTrack.Points)),
// 	}
// 	for _, p := range gpxTrack.Points {
// 		newData.GoLocations = append(newData.GoLocations, GPXPointToGoLocation(&p))
// 	}
// 	data = nil
// 	data = newData

// 	return data, nil
// }

package conversions

import (
	"strconv"
	"time"

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
		ms, err := strconv.ParseInt(gloc.TimestampMs, 10, 64)
		if err != nil {
			return nil, err
		}

		data.GoLocations = append(data.GoLocations, &models.GoLocation{
			Time:             time.Unix(0, ms*int64(time.Millisecond)),
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

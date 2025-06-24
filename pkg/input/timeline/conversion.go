package timeline

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jamesjarvis/mappyboi/v2/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

func timelineDataToData(locations []GoogleLocation) (*types.LocationHistory, error) {
	lh := &types.LocationHistory{
		Data: make([]types.Location, 0, len(locations)),
	}

	for _, gloc := range locations {
		startTime := gloc.StartTime.UTC()
		if gloc.Visit.Probability != "" {
			loc := gloc.Visit.TopCandidate.PlaceLocation
			lat, lon, err := conversions.GeocodeURIToWGS84(loc)
			if err != nil {
				return nil, fmt.Errorf("error converting geocode URI '%s' to WGS84: %w", loc, err)
			}
			lh.Insert(types.Location{
				Time:      startTime,
				Latitude:  lat,
				Longitude: lon,
			})
			continue
		}
		for _, tp := range gloc.TimelinePath {
			loc := tp.Point
			lat, lon, err := conversions.GeocodeURIToWGS84(loc)
			if err != nil {
				return nil, fmt.Errorf("error converting geocode URI '%s' to WGS84: %w", loc, err)
			}
			durationOffset, err := strconv.Atoi(tp.DurationMinutesOffsetFromStart)
			if err != nil {
				return nil, fmt.Errorf("error parsing duration minutes offset '%s': %w", tp.DurationMinutesOffsetFromStart, err)
			}
			lh.Insert(types.Location{
				Time:      startTime.Add(time.Minute * time.Duration(durationOffset)),
				Latitude:  lat,
				Longitude: lon,
			})
		}
	}

	return lh, nil
}

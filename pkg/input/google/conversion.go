package google

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jamesjarvis/mappyboi/v2/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

func googleDataToData(gd GoogleData) (*types.LocationHistory, error) {
	lh := &types.LocationHistory{
		Data: make([]types.Location, 0, len(gd.Locations)),
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
			t = gloc.Timestamp.UTC()
		}

		lh.Insert(types.Location{
			Time:             t,
			Latitude:         conversions.E7ToWGS84(gloc.LatitudeE7),
			Longitude:        conversions.E7ToWGS84(gloc.LongitudeE7),
			Altitude:         gloc.Altitude,
			Accuracy:         gloc.Accuracy,
			VerticalAccuracy: gloc.VerticalAccuracy,
		})
	}

	return lh, nil
}

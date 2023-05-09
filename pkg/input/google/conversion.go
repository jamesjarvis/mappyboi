package google

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

func E7ToStandard(e7Pos int64) float64 {
	return float64(e7Pos) / float64(1e7)
}

func GoogleDataToData(gd GoogleData) (types.LocationHistory, error) {
	lh := types.LocationHistory{
		Data: make([]types.Location, 0, len(gd.Locations)),
	}

	for _, gloc := range gd.Locations {
		var t time.Time

		if gloc.TimestampMs != "" {
			ms, err := strconv.ParseInt(gloc.TimestampMs, 10, 64)
			if err != nil {
				return types.LocationHistory{}, fmt.Errorf("error parsing timestampMs: %w", err)
			}
			t = time.Unix(0, ms*int64(time.Millisecond))
		} else {
			t = gloc.Timestamp
		}

		lh.Insert(types.Location{
			Time:             t,
			Latitude:         E7ToStandard(gloc.LatitudeE7),
			Longitude:        E7ToStandard(gloc.LongitudeE7),
			Accuracy:         float64(gloc.Accuracy),
			Altitude:         float64(gloc.Altitude),
			VerticalAccuracy: float64(gloc.VerticalAccuracy),
		})
	}

	return lh, nil
}

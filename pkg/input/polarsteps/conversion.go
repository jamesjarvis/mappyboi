package polarsteps

import (
	"math"
	"time"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

func polarstepDataToLocationHistory(pd PolarstepData) (*types.LocationHistory, error) {
	lh := &types.LocationHistory{
		Data: make([]types.Location, 0, len(pd.Locations)),
	}

	for _, psteploc := range pd.Locations {
		sec, dec := math.Modf(psteploc.Time)
		t := time.Unix(int64(sec), int64(dec*(1e9))).UTC()

		lh.Insert(types.Location{
			Time:      t,
			Latitude:  psteploc.Lat,
			Longitude: psteploc.Lon,
		})
	}

	return lh, nil
}

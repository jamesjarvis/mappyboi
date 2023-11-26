package fit

import (
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tormoder/fit"
)

func fitRecordToGoLocation(fitRecord *fit.RecordMsg) types.Location {
	return types.Location{
		Time:      fitRecord.Timestamp,
		Latitude:  fitRecord.PositionLat.Degrees(),
		Longitude: fitRecord.PositionLong.Degrees(),
		Accuracy:  int(fitRecord.GpsAccuracy),
		Altitude:  int(fitRecord.GetAltitudeScaled()),
	}
}

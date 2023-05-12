/*
Package base contains encodable data models and methods for reading and writing
from the base location history file.
*/
package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jamesjarvis/mappyboi/v2/pkg/conversions"
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

// encodableLocationHistory is the wire encodable form of the location history.
type encodableLocationHistory struct {
	Locations []encodableLocation `json:"locations"`
}

// encodableLocation is the wire encodable form of the Location structure used within mappyboi.
type encodableLocation struct {
	// Timestamp (as an ISO 8601 string) of the record.
	Timestamp string `json:"timestamp"`
	// WGS84 Latitude and Longitude coordinates of the location.
	// Degrees multiplied by 10^7 and rounded to the nearest integer, in the range
	// -1800000000 to +1800000000 (divide value by 10^7 for the usual range -180° to +180°).
	LatitudeE7  int64 `json:"latitude"`
	LongitudeE7 int64 `json:"longitude"`
	// Altitude above the WGS84 reference ellipsoid, in meters.
	Altitude int64 `json:"altitude"`
	// Approximate accuracy radius of the location measurement, in meters.
	// A lower value means better precision.
	Accuracy         int64 `json:"accuracy"`
	VerticalAccuracy int64 `json:"verticalaccuracy"`
}

func convertFromEncodable(loc encodableLocation) types.Location {
	timestamp, err := time.Parse(time.RFC3339, loc.Timestamp)
	if err != nil {
		panic(err)
	}
	return types.Location{
		Time:             timestamp,
		Latitude:         conversions.E7ToWGS84(loc.LatitudeE7),
		Longitude:        conversions.E7ToWGS84(loc.LongitudeE7),
		Altitude:         int(loc.Altitude),
		Accuracy:         int(loc.Accuracy),
		VerticalAccuracy: int(loc.VerticalAccuracy),
	}
}

func convertToEncodable(loc types.Location) encodableLocation {
	return encodableLocation{
		Timestamp:        loc.Time.Format(time.RFC3339),
		LatitudeE7:       conversions.WGS84ToE7(loc.Latitude),
		LongitudeE7:      conversions.WGS84ToE7(loc.Longitude),
		Altitude:         int64(loc.Altitude),
		Accuracy:         int64(loc.Accuracy),
		VerticalAccuracy: int64(loc.VerticalAccuracy),
	}
}

// ReadBase reads the base file into memory, within types.LocationHistory
func ReadBase(filePath string) (types.LocationHistory, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return types.LocationHistory{}, fmt.Errorf("failed to read base: %w", err)
	}
	defer f.Close()

	jsonDecoder := json.NewDecoder(f)

	locationHistory := encodableLocationHistory{}
	err = jsonDecoder.Decode(&locationHistory)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return types.LocationHistory{}, nil
		}
		return types.LocationHistory{}, fmt.Errorf("failed to decode json: %w", err)
	}

	convertedLocationHistory := types.LocationHistory{}
	for _, location := range locationHistory.Locations {
		convertedLocationHistory.Insert(
			convertFromEncodable(location),
		)
	}

	return convertedLocationHistory, nil
}

func WriteBase(filePath string, locationHistory types.LocationHistory) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open base: %w", err)
	}
	defer f.Close()

	jsonEncoder := json.NewEncoder(f)

	convertedLocationHistory := encodableLocationHistory{}
	for _, location := range locationHistory.Data {
		convertedLocationHistory.Locations = append(
			convertedLocationHistory.Locations,
			convertToEncodable(location),
		)
	}

	err = jsonEncoder.Encode(convertedLocationHistory)
	if err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}

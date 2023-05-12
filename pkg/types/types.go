/*
Package types holds generic data models to represent objects used within mappyboi.
*/
package types

import (
	"sort"
	"time"
)

// Location describes a single point in time and space.
// It is designed to be hashable, for equality comparisons
// (there is no point holding the same data twice!).
type Location struct {
	Time time.Time
	// GPS Coordinates.
	Latitude  float64
	Longitude float64
	// Altitude above the WGS84 reference ellipsoid, in meters.
	Altitude int
	// Approximate accuracy radius of the location measurement, in meters.
	// A lower value means better precision.
	Accuracy         int
	VerticalAccuracy int
}

// cleanupLocation returns the location truncated to the nearest level of accuracy we care about:
// - time: 1 second
func cleanupLocation(loc Location) Location {
	return Location{
		Time:             loc.Time.Truncate(time.Second).UTC(),
		Latitude:         loc.Latitude,
		Longitude:        loc.Longitude,
		Altitude:         loc.Altitude,
		Accuracy:         loc.Accuracy,
		VerticalAccuracy: loc.VerticalAccuracy,
	}
}

// LocationHistory stores a structured set of location history.
type LocationHistory struct {
	Data []Location // Ordered by Time ASC.
	seen map[time.Time]struct{}
}

// Insert modifies the receiver LocationHistory object by combining it
// with the incoming Location objects. If at item with the same time already exists
// within the map, it will be skipped. To maintain chronological ordering,
// one must call .Cleanup() afterwards.
func (lh *LocationHistory) Insert(data ...Location) {
	if lh.seen == nil {
		lh.seen = map[time.Time]struct{}{}
	}
	for _, v := range data {
		cleanValue := cleanupLocation(v)
		if _, exists := lh.seen[cleanValue.Time]; exists {
			continue
		}
		lh.seen[cleanValue.Time] = struct{}{}
		lh.Data = append(lh.Data, cleanValue)
	}
}

// Cleanup performs cleanup operations on the data, including sorting.
func (lh *LocationHistory) Cleanup() error {
	sort.SliceStable(lh.Data, func(i, j int) bool {
		return lh.Data[i].Time.Before(lh.Data[j].Time)
	})
	return nil
}

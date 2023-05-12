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

// locationKey exists only to serve as a map key.
type locationKey struct {
	time      time.Time
	latitude  float64
	longitude float64
}

// LocationHistory stores a structured set of location history.
type LocationHistory struct {
	Data []Location // Ordered by Time ASC.
	seen map[locationKey]struct{}
}

// Insert modifies the receiver LocationHistory object by combining it
// with the incoming Location objects. The receiver object will maintain
// chronological sorting. If the item already exists within the map, it
// will be skipped.
func (lh LocationHistory) Insert(data ...Location) {
	// TODO(jamesjarvis): perf on this sucks, do better.
	for _, v := range data {
		key := locationKey{
			time:      v.Time,
			latitude:  v.Latitude,
			longitude: v.Longitude,
		}
		if _, exists := lh.seen[key]; exists {
			continue
		}
		lh.seen[key] = struct{}{}
	}

	sort.SliceStable(lh.Data, func(i, j int) bool {
		return lh.Data[i].Time.Before(lh.Data[j].Time)
	})
}

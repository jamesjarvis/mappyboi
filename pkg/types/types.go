package types

import "time"

// Location describes a single point in time and space.
// It is designed to be hashable, for equality comparisons
// (there is no point holding the same data twice!).
type Location struct {
	Time             time.Time
	Latitude         float64
	Longitude        float64
	Accuracy         float64
	Altitude         float64
	VerticalAccuracy float64
}

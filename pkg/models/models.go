package models

import "time"

type Data struct {
	GoLocations []*GoLocation
}

type GoLocation struct {
	Time             time.Time
	Latitude         float64
	Longitude        float64
	Accuracy         float64
	Altitude         float64
	VerticalAccuracy float64
}

package models

import "time"

type GoogleData struct {
	Locations []GoogleLocation `json:"locations"`
}

type GoogleLocation struct {
	TimestampMs      string            `json:"timestampMs"` // deprecated, changed sometime between January and March 2022.
	Timestamp        time.Time         `json:"timestamp"`   // added sometime between January and March 2022.
	LatitudeE7       int64             `json:"latitudeE7"`
	LongitudeE7      int64             `json:"longitudeE7"`
	Accuracy         int               `json:"accuracy"`
	Altitude         int               `json:"altitude"`
	VerticalAccuracy int               `json:"verticalAccuracy"`
	Activities       []ActivityWrapper `json:"activity"`
}

type ActivityWrapper struct {
	TimestampMs string     `json:"timestampMs"`
	Activities  []Activity `json:"activity"`
}

type Activity struct {
	Type       string `json:"type"`
	Confidence int    `json:"confidence"`
}

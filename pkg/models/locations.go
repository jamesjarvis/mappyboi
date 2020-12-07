package models

type GoogleData struct {
	Locations []*Location `json:"locations"`
}

type Location struct {
	TimestampMs      string            `json:"timestampMs"`
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

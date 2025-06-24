package timeline

import "time"

type GoogleLocation struct {
	EndTime   time.Time `json:"endTime"`
	StartTime time.Time `json:"startTime"`
	// Either Visit or TimelinePath will be populated.
	Visit        Visit               `json:"visit"`
	TimelinePath []TimelinePathPoint `json:"timelinePath"`
}

type Visit struct {
	TopCandidate TopCandidate `json:"topCandidate"`
	Probability  string       `json:"probability"`
}

type TopCandidate struct {
	PlaceID       string `json:"placeID"`
	PlaceLocation string `json:"placeLocation"`
	Probability   string `json:"probability"`
}

type TimelinePathPoint struct {
	Point                          string `json:"point"`
	DurationMinutesOffsetFromStart string `json:"durationMinutesOffsetFromStartTime"`
}

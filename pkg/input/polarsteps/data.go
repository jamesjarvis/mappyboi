package polarsteps

type PolarstepData struct {
	Locations []PolarstepLocation `json:"locations"`
}

type PolarstepLocation struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Time float64 `json:"time"` // Unix timestamp in seconds.
}

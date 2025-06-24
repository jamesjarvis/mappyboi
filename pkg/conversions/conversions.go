/*
Package conversions contains useful converter functions used in various parts of mappyboi.
*/
package conversions

import "fmt"

func E7ToWGS84(e7Pos int64) float64 {
	return float64(e7Pos) / float64(1e7)
}

func WGS84ToE7(wgs84pos float64) int64 {
	return int64(wgs84pos * 1e7)
}

func GeocodeURIToWGS84(uri string) (float64, float64, error) {
	var lat, lon float64
	_, err := fmt.Sscanf(uri, "geo:%f,%f", &lat, &lon)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocode URI '%s': %w", uri, err)
	}
	return lat, lon, nil
}

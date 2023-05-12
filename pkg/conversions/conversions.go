/*
Package conversions contains useful converter functions used in various parts of mappyboi.
*/
package conversions

func E7ToWGS84(e7Pos int64) float64 {
	return float64(e7Pos) / float64(1e7)
}

func WGS84ToE7(wgs84pos float64) int64 {
	return int64(wgs84pos * 1e7)
}

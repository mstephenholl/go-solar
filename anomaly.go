// Package solar provides functions for calculating sunrise, sunset, and solar position
// for any location on Earth. It uses astronomical algorithms to compute solar elevation
// angles and times when the sun crosses specific elevation thresholds.
package solar

import (
	"math"
)

// MeanAnomaly calculates the solar mean anomaly, which is the angle of the sun
// relative to the earth for the specified Julian day.
func MeanAnomaly(d float64) float64 {
	v := math.Remainder(SolarMeanAnomalyBase+SolarMeanAnomalyRate*(d-J2000), FullCircleDegrees)
	if v < 0 {
		v += FullCircleDegrees
	}
	return v
}

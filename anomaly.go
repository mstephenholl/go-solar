package sunrise

import (
	"math"
)

// SolarMeanAnomaly calculates the angle of the sun relative to the earth for
// the specified Julian day.
func SolarMeanAnomaly(d float64) float64 {
	v := math.Remainder(SolarMeanAnomalyBase+SolarMeanAnomalyRate*(d-J2000), FullCircleDegrees)
	if v < 0 {
		v += FullCircleDegrees
	}
	return v
}

package solar

import (
	"math"
)

// EclipticLongitude calculates the angular distance of the earth along the
// ecliptic.
func EclipticLongitude(solarAnomaly, equationOfCenter, d float64) float64 {
	return math.Mod(solarAnomaly+equationOfCenter+HalfCircleDegrees+ArgumentOfPerihelion(d), FullCircleDegrees)
}

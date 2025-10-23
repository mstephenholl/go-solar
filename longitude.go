package solar

import (
	"math"
)

// eclipticLongitude calculates the angular distance of the earth along the
// ecliptic.
func eclipticLongitude(solarAnomaly, equationOfCenter, d float64) float64 {
	return math.Mod(solarAnomaly+equationOfCenter+HalfCircleDegrees+argumentOfPerihelion(d), FullCircleDegrees)
}

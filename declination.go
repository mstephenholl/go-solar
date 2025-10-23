package solar

import (
	"math"
)

// declination calculates one of the two angles required to locate a point on
// the celestial sphere in the equatorial coordinate system. The ecliptic
// longitude parameter must be in degrees.
func declination(eclipticLongitude float64) float64 {
	return math.Asin(math.Sin(eclipticLongitude*Degree)*SinDeclinationCoefficient) / Degree
}

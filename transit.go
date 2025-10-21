package sunrise

import (
	"math"
)

// SolarTransit calculates the Julian data for the local true solar transit.
func SolarTransit(d, solarAnomaly, eclipticLongitude float64) float64 {
	equationOfTime := EquationOfTimeC1*math.Sin(solarAnomaly*Degree) -
		EquationOfTimeC2*math.Sin(2*eclipticLongitude*Degree)
	return d + equationOfTime
}

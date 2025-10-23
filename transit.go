package solar

import (
	"math"
)

// transit calculates the Julian date for the local true solar transit (solar noon).
func transit(d, meanAnomaly, eclipticLongitude float64) float64 {
	equationOfTime := EquationOfTimeC1*math.Sin(meanAnomaly*Degree) -
		EquationOfTimeC2*math.Sin(2*eclipticLongitude*Degree)
	return d + equationOfTime
}

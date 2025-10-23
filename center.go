package solar

import (
	"math"
)

// equationOfCenter calculates the angular difference between the position of
// the earth in its elliptical orbit and the position it would occupy in a
// circular orbit for the given mean anomaly.
func equationOfCenter(solarAnomaly float64) float64 {
	var (
		anomalyInRad = solarAnomaly * Degree
		anomalySin   = math.Sin(anomalyInRad)
		anomaly2Sin  = math.Sin(2 * anomalyInRad)
		anomaly3Sin  = math.Sin(3 * anomalyInRad)
	)
	return EquationOfCenterC1*anomalySin + EquationOfCenterC2*anomaly2Sin + EquationOfCenterC3*anomaly3Sin
}

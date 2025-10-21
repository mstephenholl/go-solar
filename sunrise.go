package solar

import (
	"math"
	"time"
)

// Sunrise calculates when the sun will rise on the given day at the specified location.
// Returns time.Time{} if the sun does not rise (e.g., polar night).
func Sunrise(latitude, longitude float64, year int, month time.Month, day int) time.Time {
	rise, _ := sunriseSunsetInternal(latitude, longitude, year, month, day)
	return rise
}

// Sunset calculates when the sun will set on the given day at the specified location.
// Returns time.Time{} if the sun does not set (e.g., midnight sun).
func Sunset(latitude, longitude float64, year int, month time.Month, day int) time.Time {
	_, set := sunriseSunsetInternal(latitude, longitude, year, month, day)
	return set
}

// SunriseSunset calculates when the sun will rise and when it will set on the
// given day at the specified location.
// Returns time.Time{} if the sun does not rise or set.
func SunriseSunset(latitude, longitude float64, year int, month time.Month, day int) (time.Time, time.Time) {
	return sunriseSunsetInternal(latitude, longitude, year, month, day)
}

// sunriseSunsetInternal is the internal implementation shared by all public functions.
func sunriseSunsetInternal(latitude, longitude float64, year int, month time.Month, day int) (time.Time, time.Time) {
	var (
		d                 = MeanSolarNoon(longitude, year, month, day)
		meanAnomaly       = MeanAnomaly(d)
		equationOfCenter  = EquationOfCenter(meanAnomaly)
		eclipticLongitude = EclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = Transit(d, meanAnomaly, eclipticLongitude)
		declination       = Declination(eclipticLongitude)
		hourAngle         = HourAngle(latitude, declination)
		frac              = hourAngle / FullCircleDegrees
		sunrise           = transit - frac
		sunset            = transit + frac
	)

	// Check for no sunrise, no sunset
	if hourAngle == math.MaxFloat64 || hourAngle == -1*math.MaxFloat64 {
		return time.Time{}, time.Time{}
	}

	return JulianDayToTime(sunrise), JulianDayToTime(sunset)
}

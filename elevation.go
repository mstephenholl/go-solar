package solar

import (
	"math"
	"time"
)

// timeOfElevationInternal is the internal implementation with old signature
func timeOfElevationInternal(latitude, longitude, elevation float64, year int, month time.Month, day int) (morning, evening time.Time) {
	var (
		d                 = meanSolarNoonInternal(longitude, year, month, day)
		meanAnomaly       = meanAnomaly(d)
		equationOfCenter  = equationOfCenter(meanAnomaly)
		eclipticLongitude = eclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = transit(d, meanAnomaly, eclipticLongitude)
		declination       = declination(eclipticLongitude)
		// https://solarsena.com/solar-elevation-angle-altitude/
		numerator   = math.Sin(elevation*Degree) - (math.Sin(latitude*Degree) * math.Sin(declination*Degree))
		denominator = math.Cos(latitude*Degree) * math.Cos(declination*Degree)
		hourAngle   = math.Acos(numerator / denominator)
		frac        = hourAngle / (2 * math.Pi)
		morningJD   = transit - frac
		eveningJD   = transit + frac
	)

	// Check for cases where the sun never reaches the given elevation.
	if math.IsNaN(hourAngle) {
		return time.Time{}, time.Time{}
	}

	morning = JulianDayToTime(morningJD)
	evening = JulianDayToTime(eveningJD)
	return morning, evening
}

// TimeOfElevation calculates the times of day when the sun is at a given elevation
// above the horizon on a given day at the specified location.
//
// All times are calculated and returned in UTC. Useful for calculating twilight times,
// golden hour, etc.
//
// Common elevation angles:
//   - -0.833°: Official sunrise/sunset (accounts for atmospheric refraction)
//   - -6°: Civil twilight
//   - -12°: Nautical twilight
//   - -18°: Astronomical twilight
//   - 6°: Golden hour
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - elevation: Solar elevation angle in degrees (negative for below horizon)
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//
// Returns:
//   - morning: Time in UTC when sun reaches elevation in the morning (time.Time{} if never reached)
//   - evening: Time in UTC when sun reaches elevation in the evening (time.Time{} if never reached)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	// Calculate civil twilight times
//	morning, evening := solar.TimeOfElevation(loc, -6.0, t)
func TimeOfElevation(loc Location, elevation float64, t Time) (morning, evening time.Time) {
	return timeOfElevationInternal(loc.Latitude(), loc.Longitude(), elevation, t.Year(), t.Month(), t.Day())
}

// elevationInternal is the internal implementation with old signature
func elevationInternal(latitude, longitude float64, when time.Time) float64 {
	var (
		d                 = meanSolarNoonInternal(longitude, when.Year(), when.Month(), when.Day())
		meanAnomaly       = meanAnomaly(d)
		equationOfCenter  = equationOfCenter(meanAnomaly)
		eclipticLongitude = eclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = transit(d, meanAnomaly, eclipticLongitude)
		declination       = declination(eclipticLongitude)
		frac              = transit - TimeToJulianDay(when)
		hourAngle         = 2 * math.Pi * frac
		// https://solarsena.com/solar-elevation-angle-altitude/
		firstPart  = math.Sin(latitude*Degree) * math.Sin(declination*Degree)
		secondPart = math.Cos(latitude*Degree) * math.Cos(declination*Degree) * math.Cos(hourAngle)
	)

	return math.Asin(firstPart+secondPart) / Degree
}

// Elevation calculates the angle of the sun above the horizon at a given moment
// at the specified location.
//
// The time parameter should be in UTC. If you have a local time, convert it to UTC first
// using time.UTC() or time.In(time.UTC).
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - when: The moment in time to calculate elevation (in UTC)
//
// Returns:
//   - Solar elevation angle in degrees (positive above horizon, negative below)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	elevation := solar.Elevation(loc, time.Now().UTC())
//	// elevation is in degrees, positive above horizon, negative below
func Elevation(loc Location, when time.Time) float64 {
	return elevationInternal(loc.Latitude(), loc.Longitude(), when)
}

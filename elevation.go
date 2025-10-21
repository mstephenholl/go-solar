package solar

import (
	"math"
	"time"
)

// TimeOfElevation calculates the times of day when the sun is at a given elevation
// above the horizon on a given day at the specified location.
//
// All times are calculated and returned in UTC. The input date (year, month, day) is
// interpreted as UTC. Useful for calculating twilight times, golden hour, etc.
//
// Common elevation angles:
//   - -0.833°: Official sunrise/sunset (accounts for atmospheric refraction)
//   - -6°: Civil twilight
//   - -12°: Nautical twilight
//   - -18°: Astronomical twilight
//   - 6°: Golden hour
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - elevation: Solar elevation angle in degrees (negative for below horizon)
//   - year, month, day: The date in UTC for which to calculate times
//
// Returns:
//   - morning: Time in UTC when sun reaches elevation in the morning (time.Time{} if never reached)
//   - evening: Time in UTC when sun reaches elevation in the evening (time.Time{} if never reached)
//
// Example:
//
//	// Calculate civil twilight times
//	morning, evening := solar.TimeOfElevation(40.7128, -74.0060, -6.0, 2024, time.June, 21)
func TimeOfElevation(latitude, longitude, elevation float64, year int, month time.Month, day int) (morning, evening time.Time) {
	var (
		d                 = MeanSolarNoon(longitude, year, month, day)
		meanAnomaly       = MeanAnomaly(d)
		equationOfCenter  = EquationOfCenter(meanAnomaly)
		eclipticLongitude = EclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = Transit(d, meanAnomaly, eclipticLongitude)
		declination       = Declination(eclipticLongitude)
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

// Elevation calculates the angle of the sun above the horizon at a given moment
// at the specified location.
//
// The time parameter should be in UTC. If you have a local time, convert it to UTC first
// using time.UTC() or time.In(time.UTC).
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - when: The moment in time to calculate elevation (in UTC)
//
// Returns:
//   - Solar elevation angle in degrees (positive above horizon, negative below)
//
// Example:
//
//	// Calculate current sun elevation
//	elevation := solar.Elevation(40.7128, -74.0060, time.Now().UTC())
//	fmt.Printf("Sun is %.2f degrees above horizon\n", elevation)
func Elevation(latitude, longitude float64, when time.Time) float64 {
	var (
		d                 = MeanSolarNoon(longitude, when.Year(), when.Month(), when.Day())
		meanAnomaly       = MeanAnomaly(d)
		equationOfCenter  = EquationOfCenter(meanAnomaly)
		eclipticLongitude = EclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = Transit(d, meanAnomaly, eclipticLongitude)
		declination       = Declination(eclipticLongitude)
		frac              = transit - TimeToJulianDay(when)
		hourAngle         = 2 * math.Pi * frac
		// https://solarsena.com/solar-elevation-angle-altitude/
		firstPart  = math.Sin(latitude*Degree) * math.Sin(declination*Degree)
		secondPart = math.Cos(latitude*Degree) * math.Cos(declination*Degree) * math.Cos(hourAngle)
	)

	return math.Asin(firstPart+secondPart) / Degree
}

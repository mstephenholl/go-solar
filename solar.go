package solar

import (
	"math"
	"time"
)

// Sunrise calculates when the sun will rise on the given day at the specified location.
//
// All times are calculated and returned in UTC. The input date (year, month, day) is
// interpreted as UTC. To convert to local time, use time.In() with the appropriate
// timezone.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate sunrise
//
// Returns:
//   - Sunrise time in UTC
//   - time.Time{} if the sun does not rise on this day (e.g., polar night)
//
// Example:
//
//	sunrise := solar.Sunrise(40.7128, -74.0060, 2024, time.June, 21)
//	fmt.Println(sunrise.Format("15:04:05 MST"))  // UTC time
//	localTime := sunrise.In(time.Local)           // Convert to local time
func Sunrise(latitude, longitude float64, year int, month time.Month, day int) time.Time {
	rise, _ := sunriseSunsetInternal(latitude, longitude, year, month, day)
	return rise
}

// Sunset calculates when the sun will set on the given day at the specified location.
//
// All times are calculated and returned in UTC. The input date (year, month, day) is
// interpreted as UTC. To convert to local time, use time.In() with the appropriate
// timezone.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate sunset
//
// Returns:
//   - Sunset time in UTC
//   - time.Time{} if the sun does not set on this day (e.g., midnight sun)
//
// Example:
//
//	sunset := solar.Sunset(40.7128, -74.0060, 2024, time.June, 21)
//	fmt.Println(sunset.Format("15:04:05 MST"))  // UTC time
//	localTime := sunset.In(time.Local)          // Convert to local time
func Sunset(latitude, longitude float64, year int, month time.Month, day int) time.Time {
	_, set := sunriseSunsetInternal(latitude, longitude, year, month, day)
	return set
}

// SunriseSunset calculates when the sun will rise and when it will set on the
// given day at the specified location.
//
// All times are calculated and returned in UTC. The input date (year, month, day) is
// interpreted as UTC. To convert to local time, use time.In() with the appropriate
// timezone.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate sunrise and sunset
//
// Returns:
//   - sunrise: Sunrise time in UTC (time.Time{} if sun does not rise)
//   - sunset: Sunset time in UTC (time.Time{} if sun does not set)
//
// Example:
//
//	sunrise, sunset := solar.SunriseSunset(40.7128, -74.0060, 2024, time.June, 21)
//	fmt.Printf("Sunrise: %s, Sunset: %s\n", sunrise.Format("15:04 MST"), sunset.Format("15:04 MST"))
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

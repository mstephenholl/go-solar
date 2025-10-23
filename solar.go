package solar

import (
	"errors"
	"math"
	"time"
)

var (
	// ErrSunNeverRises is returned when the sun never rises at the given location and date (polar night).
	ErrSunNeverRises = errors.New("sun never rises at this location on this date")
	// ErrSunNeverSets is returned when the sun never sets at the given location and date (midnight sun).
	ErrSunNeverSets = errors.New("sun never sets at this location on this date")
)

// Sunrise calculates when the sun will rise on the given day at the specified location.
//
// All times are calculated and returned in UTC. To convert to local time, use time.In()
// with the appropriate timezone.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//
// Returns:
//   - Sunrise time in UTC
//   - error if the sun does not rise on this day (e.g., polar night)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	sunrise, err := solar.Sunrise(loc, t)
//	if err != nil {
//	    // Handle polar night or other errors
//	}
//	// sunrise is in UTC - convert to local time if needed
//	localTime := sunrise.In(time.Local)
func Sunrise(loc Location, t Time) (time.Time, error) {
	rise, _, err := sunriseSunsetInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day())
	return rise, err
}

// Sunset calculates when the sun will set on the given day at the specified location.
//
// All times are calculated and returned in UTC. To convert to local time, use time.In()
// with the appropriate timezone.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//
// Returns:
//   - Sunset time in UTC
//   - error if the sun does not set on this day (e.g., midnight sun)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	sunset, err := solar.Sunset(loc, t)
//	if err != nil {
//	    // Handle midnight sun or other errors
//	}
//	// sunset is in UTC - convert to local time if needed
//	localTime := sunset.In(time.Local)
func Sunset(loc Location, t Time) (time.Time, error) {
	_, set, err := sunriseSunsetInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day())
	return set, err
}

// SunriseSunset calculates when the sun will rise and when it will set on the
// given day at the specified location.
//
// All times are calculated and returned in UTC. To convert to local time, use time.In()
// with the appropriate timezone.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//
// Returns:
//   - sunrise: Sunrise time in UTC
//   - sunset: Sunset time in UTC
//   - error if the sun does not rise or set (e.g., polar night or midnight sun)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	sunrise, sunset, err := solar.SunriseSunset(loc, t)
//	if err != nil {
//	    // Handle polar night or midnight sun
//	}
//	// Both times are in UTC - convert to local time if needed
func SunriseSunset(loc Location, t Time) (time.Time, time.Time, error) {
	return sunriseSunsetInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day())
}

// sunriseSunsetInternal is the internal implementation shared by all public functions.
func sunriseSunsetInternal(latitude, longitude float64, year int, month time.Month, day int) (time.Time, time.Time, error) {
	var (
		d                 = meanSolarNoonInternal(longitude, year, month, day)
		meanAnomaly       = meanAnomaly(d)
		equationOfCenter  = equationOfCenter(meanAnomaly)
		eclipticLongitude = eclipticLongitude(meanAnomaly, equationOfCenter, d)
		transit           = transit(d, meanAnomaly, eclipticLongitude)
		declination       = declination(eclipticLongitude)
		hourAngle         = hourAngle(latitude, declination)
		frac              = hourAngle / FullCircleDegrees
		sunrise           = transit - frac
		sunset            = transit + frac
	)

	// Check for no sunrise, no sunset
	if hourAngle == math.MaxFloat64 {
		// Polar night - sun never rises
		return time.Time{}, time.Time{}, ErrSunNeverRises
	}
	if hourAngle == -1*math.MaxFloat64 {
		// Midnight sun - sun never sets
		return time.Time{}, time.Time{}, ErrSunNeverSets
	}

	return JulianDayToTime(sunrise), JulianDayToTime(sunset), nil
}

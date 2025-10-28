package solar

import (
	"time"
)

// meanSolarNoonInternal calculates the time at which the sun is at its highest
// altitude. The returned time is in Julian days. This is an internal function.
func meanSolarNoonInternal(longitude float64, year int, month time.Month, day int) float64 {
	t := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	return TimeToJulianDay(t) - longitude/LongitudeDivisor
}

// MeanSolarNoon calculates the time at which the sun is at its highest altitude
// (solar noon) for the given location and date.
//
// All times are calculated and returned in UTC. To convert to local time, use
// time.In() with the appropriate timezone.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//
// Returns:
//   - Solar noon time in UTC
//
// Example:
//
//	loc := solar.NewLocation(43.65, -79.38)
//	t := solar.NewTime(2024, time.June, 21)
//	noon := solar.MeanSolarNoon(loc, t)
//	// noon is in UTC - convert to local time if needed
func MeanSolarNoon(loc Location, t Time) time.Time {
	jd := meanSolarNoonInternal(loc.Longitude(), t.Year(), t.Month(), t.Day())
	return JulianDayToTime(jd)
}

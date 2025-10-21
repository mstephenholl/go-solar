package solar

import (
	"time"
)

const (
	secondsInADay      = 86400
	unixEpochJulianDay = 2440587.5
)

// TimeToJulianDay converts a time.Time into a Julian day number.
//
// The input time should be in UTC for accurate astronomical calculations.
// The timezone of the input time is preserved in the calculation via Unix timestamp.
//
// Parameters:
//   - t: Time to convert (should be in UTC for astronomical calculations)
//
// Returns:
//   - Julian day number as float64
//
// Example:
//
//	jd := solar.TimeToJulianDay(time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC))
//	// Returns 2451545.0 (J2000.0 epoch)
func TimeToJulianDay(t time.Time) float64 {
	return float64(t.Unix())/secondsInADay + unixEpochJulianDay
}

// JulianDayToTime converts a Julian day number into a time.Time.
//
// The returned time is always in UTC timezone. This is the standard for
// astronomical calculations.
//
// Parameters:
//   - d: Julian day number as float64
//
// Returns:
//   - Time in UTC corresponding to the Julian day
//
// Example:
//
//	t := solar.JulianDayToTime(2451545.0)  // J2000.0 epoch
//	// Returns 2000-01-01 12:00:00 +0000 UTC
func JulianDayToTime(d float64) time.Time {
	return time.Unix(int64((d-unixEpochJulianDay)*secondsInADay), 0).UTC()
}

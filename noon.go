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

// MeanSolarNoonFromNMEA calculates the time at which the sun is at its highest
// altitude for the location and date encoded in an NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned time is also in UTC.
// Convert to local time using time.In() if needed.
//
// Supported NMEA sentence types:
//   - GGA (Global Positioning System Fix Data) - requires external date
//   - RMC (Recommended Minimum Specific GPS/Transit Data) - includes date
//
// For GGA sentences, you must provide the date via the optional date parameter.
// For RMC sentences, the date parameter is ignored as the date is in the sentence.
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (e.g., "$GPGGA,..." or "$GPRMC,...")
//   - date: Optional date for GGA sentences (year, month, day). Pass 0, 0, 0 for RMC.
//
// Returns:
//   - The solar noon time in UTC at the location specified in the NMEA sentence
//   - An error if the sentence is invalid or unsupported
//
// Example:
//
//	// Using RMC sentence (includes date)
//	solarNoon, err := solar.MeanSolarNoonFromNMEA("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71", 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Using GGA sentence (requires external date)
//	solarNoon, err = solar.MeanSolarNoonFromNMEA("$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*5C", 2024, time.June, 21)
func MeanSolarNoonFromNMEA(nmea string, year int, month time.Month, day int) (time.Time, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	return MeanSolarNoon(loc, t), nil
}

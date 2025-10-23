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

// ElevationFromNMEA calculates the solar elevation angle at the location and time
// encoded in an NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned elevation is calculated for
// the time encoded in the NMEA sentence.
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
//   - Solar elevation angle in degrees (positive above horizon, negative below)
//   - An error if the sentence is invalid or unsupported
//
// Example:
//
//	// Using RMC sentence (includes date and time)
//	elevation, err := solar.ElevationFromNMEA("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71", 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Using GGA sentence (requires external date)
//	elevation, err = solar.ElevationFromNMEA("$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*5C", 2024, time.June, 21)
func ElevationFromNMEA(nmea string, year int, month time.Month, day int) (float64, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return 0, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return 0, err
	}

	elevation := Elevation(loc, t.DateTime())
	return elevation, nil
}

// TimeOfElevationFromNMEA calculates when the sun reaches a specified elevation angle
// on the date encoded in an NMEA sentence, at the location from the NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned times are in UTC.
// Convert to local time using time.In() if needed.
//
// Supported NMEA sentence types:
//   - GGA (Global Positioning System Fix Data) - requires external date
//   - RMC (Recommended Minimum Specific GPS/Transit Data) - includes date
//
// For GGA sentences, you must provide the date via the optional date parameter.
// For RMC sentences, the date parameter is ignored as the date is in the sentence.
//
// Common elevation angles:
//   - -0.833°: Official sunrise/sunset (accounts for atmospheric refraction)
//   - -6°: Civil twilight
//   - -12°: Nautical twilight
//   - -18°: Astronomical twilight
//   - 6°: Golden hour
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (e.g., "$GPGGA,..." or "$GPRMC,...")
//   - elevation: Solar elevation angle in degrees (negative for below horizon)
//   - date: Optional date for GGA sentences (year, month, day). Pass 0, 0, 0 for RMC.
//
// Returns:
//   - morning: Time in UTC when sun reaches elevation in the morning (time.Time{} if never reached)
//   - evening: Time in UTC when sun reaches elevation in the evening (time.Time{} if never reached)
//   - error: An error if the sentence is invalid or unsupported
//
// Example:
//
//	// Using RMC sentence - calculate civil twilight times
//	morning, evening, err := solar.TimeOfElevationFromNMEA(
//	    "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71",
//	    -6.0,  // Civil twilight
//	    0, 0, 0,
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Using GGA sentence (requires external date)
//	morning, evening, err = solar.TimeOfElevationFromNMEA(
//	    "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*5C",
//	    -6.0,
//	    2024, time.June, 21,
//	)
func TimeOfElevationFromNMEA(nmea string, elevation float64, year int, month time.Month, day int) (time.Time, time.Time, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	morning, evening := TimeOfElevation(loc, elevation, t)
	return morning, evening, nil
}

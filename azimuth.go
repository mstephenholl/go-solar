package solar

import (
	"math"
	"time"
)

// azimuthInternal calculates the solar azimuth angle (internal function with old signature).
func azimuthInternal(latitude, longitude float64, when time.Time) float64 {
	// Calculate mean solar noon for the given day
	meanNoon := meanSolarNoonInternal(longitude, when.Year(), when.Month(), when.Day())

	// Calculate solar mean anomaly
	solarAnomaly := meanAnomaly(meanNoon)

	// Calculate equation of center
	equationOfCenter := equationOfCenter(solarAnomaly)

	// Calculate ecliptic longitude
	eclipticLongitude := eclipticLongitude(solarAnomaly, equationOfCenter, meanNoon)

	// Calculate solar transit
	solarTransit := transit(meanNoon, solarAnomaly, eclipticLongitude)

	// Calculate solar declination
	declination := declination(eclipticLongitude)

	// Calculate hour angle from the time
	frac := TimeToJulianDay(when) - solarTransit
	hourAngle := 2.0 * math.Pi * frac

	// Get current solar elevation
	elevation := elevationInternal(latitude, longitude, when)

	// Calculate azimuth using spherical trigonometry
	// Formula: cos(Az) = (sin(δ) * cos(φ) - cos(δ) * sin(φ) * cos(H)) / cos(h)
	// Where:
	//   Az = azimuth angle
	//   δ = declination
	//   φ = latitude
	//   H = hour angle
	//   h = elevation
	latRad := latitude * Degree
	declRad := declination * Degree
	elevRad := elevation * Degree

	firstPart := math.Sin(declRad) * math.Cos(latRad)
	secondPart := math.Cos(declRad) * math.Sin(latRad) * math.Cos(hourAngle)
	cosAzimuth := (firstPart - secondPart) / math.Cos(elevRad)

	// Calculate azimuth in degrees
	azimuth := math.Acos(cosAzimuth) / Degree

	// Adjust for afternoon (hour angle >= 0 means afternoon/evening)
	if hourAngle >= 0 {
		azimuth = FullCircleDegrees - azimuth
	}

	return azimuth
}

// Azimuth calculates the solar azimuth angle at a specific time and location.
// The azimuth is the sun's compass direction, measured clockwise from true north.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - when: The datetime at which to calculate the azimuth (use time.Date, time.Now(), or Time.DateTime())
//
// Returns:
//   - The solar azimuth angle in degrees (0° = North, 90° = East, 180° = South, 270° = West)
//
// The calculation uses the solar hour angle and declination to determine the sun's
// position in the sky. The azimuth is measured clockwise from north, ranging from
// 0° to 360°.
//
// Note: All calculations assume UTC time. Ensure the input time is in UTC timezone.
//
// Example:
//
//	loc := solar.NewLocation(43.65, -79.38)
//	when := time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC)
//	azimuth := solar.Azimuth(loc, when)
//	// azimuth is in degrees: 0°=North, 90°=East, 180°=South, 270°=West
func Azimuth(loc Location, when time.Time) float64 {
	return azimuthInternal(loc.Latitude(), loc.Longitude(), when)
}

// AzimuthFromNMEA calculates the solar azimuth angle from an NMEA GPS sentence.
// The azimuth is the sun's compass direction, measured clockwise from true north.
//
// This function parses the NMEA sentence to extract location and time, then calculates
// the solar azimuth at that moment.
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (GGA or RMC format)
//   - year: Year (ignored for RMC sentences that include date)
//   - month: Month (ignored for RMC sentences that include date)
//   - day: Day of month (ignored for RMC sentences that include date)
//
// Returns:
//   - The solar azimuth angle in degrees (0° = North, 90° = East, 180° = South, 270° = West)
//   - An error if the NMEA sentence is invalid or cannot be parsed
//
// Supported NMEA sentence types:
//   - GGA (GPS Fix Data): Requires year, month, day parameters
//   - RMC (Recommended Minimum): Date is extracted from sentence, external date parameters ignored
//
// Note: All calculations assume UTC time.
//
// Example:
//
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"
//	azimuth, err := solar.AzimuthFromNMEA(nmea, 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
func AzimuthFromNMEA(nmea string, year int, month time.Month, day int) (float64, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return 0, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return 0, err
	}

	// Calculate azimuth using the parsed location and time
	azimuth := Azimuth(loc, t.DateTime())
	return azimuth, nil
}

package solar

import (
	"math"
	"time"
)

// Azimuth calculates the solar azimuth angle at a specific time and location.
// The azimuth is the sun's compass direction, measured clockwise from true north.
//
// Parameters:
//   - latitude: The latitude of the location in decimal degrees (positive for North, negative for South)
//   - longitude: The longitude of the location in decimal degrees (positive for East, negative for West)
//   - when: The time at which to calculate the azimuth (in UTC)
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
//	// Calculate azimuth for Toronto at noon on January 1, 2000
//	when := time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC)
//	azimuth := solar.Azimuth(43.65, -79.38, when)
//	fmt.Printf("Sun azimuth: %.2f degrees\n", azimuth)
func Azimuth(latitude, longitude float64, when time.Time) float64 {
	// Calculate mean solar noon for the given day
	meanNoon := MeanSolarNoon(longitude, when.Year(), when.Month(), when.Day())

	// Calculate solar mean anomaly
	solarAnomaly := MeanAnomaly(meanNoon)

	// Calculate equation of center
	equationOfCenter := EquationOfCenter(solarAnomaly)

	// Calculate ecliptic longitude
	eclipticLongitude := EclipticLongitude(solarAnomaly, equationOfCenter, meanNoon)

	// Calculate solar transit
	solarTransit := Transit(meanNoon, solarAnomaly, eclipticLongitude)

	// Calculate solar declination
	declination := Declination(eclipticLongitude)

	// Calculate hour angle from the time
	frac := TimeToJulianDay(when) - solarTransit
	hourAngle := 2.0 * math.Pi * frac

	// Get current solar elevation
	elevation := Elevation(latitude, longitude, when)

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
//	fmt.Printf("Sun azimuth: %.2f degrees\n", azimuth)
func AzimuthFromNMEA(nmea string, year int, month time.Month, day int) (float64, error) {
	// Parse the NMEA sentence to get position and time
	pos, err := parseNMEA(nmea, year, month, day)
	if err != nil {
		return 0, err
	}

	// Calculate azimuth using the parsed location and time
	azimuth := Azimuth(pos.Latitude, pos.Longitude, pos.Time)
	return azimuth, nil
}

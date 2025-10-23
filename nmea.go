package solar

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrInvalidNMEA is returned when the NMEA sentence is malformed or unsupported.
	ErrInvalidNMEA = errors.New("invalid NMEA sentence")

	// ErrInvalidChecksum is returned when the NMEA checksum validation fails.
	ErrInvalidChecksum = errors.New("invalid NMEA checksum")

	// ErrUnsupportedSentence is returned when the NMEA sentence type is not supported.
	ErrUnsupportedSentence = errors.New("unsupported NMEA sentence type")

	// ErrInvalidPosition is returned when latitude or longitude cannot be parsed.
	ErrInvalidPosition = errors.New("invalid position data")

	// ErrInvalidDate is returned when date/time cannot be parsed from NMEA sentence.
	ErrInvalidDate = errors.New("invalid date/time data")
)

// nmeaPosition holds the parsed position and time data from an NMEA sentence.
type nmeaPosition struct {
	Latitude  float64
	Longitude float64
	Time      time.Time
}

// SunriseFromNMEA calculates the sunrise time for the location and date encoded
// in an NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned sunrise time is also in UTC.
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
//   - The sunrise time in UTC at the location specified in the NMEA sentence
//   - An error if the sentence is invalid or unsupported
//   - time.Time{} if the sun does not rise (e.g., polar night)
//
// Example:
//
//	// Using RMC sentence (includes date)
//	sunrise, err := solar.SunriseFromNMEA("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A", 0, 0, 0)
//
//	// Using GGA sentence (requires date)
//	sunrise, err := solar.SunriseFromNMEA("$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47", 2024, time.June, 21)
func SunriseFromNMEA(nmea string, year int, month time.Month, day int) (time.Time, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	return Sunrise(loc, t)
}

// SunsetFromNMEA calculates the sunset time for the location and date encoded
// in an NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned sunset time is also in UTC.
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
//   - The sunset time in UTC at the location specified in the NMEA sentence
//   - An error if the sentence is invalid or unsupported
//   - time.Time{} if the sun does not set (e.g., midnight sun)
//
// Example:
//
//	// Using RMC sentence (includes date)
//	sunset, err := solar.SunsetFromNMEA("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A", 0, 0, 0)
//
//	// Using GGA sentence (requires date)
//	sunset, err := solar.SunsetFromNMEA("$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47", 2024, time.June, 21)
func SunsetFromNMEA(nmea string, year int, month time.Month, day int) (time.Time, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	return Sunset(loc, t)
}

// SunriseSunsetFromNMEA calculates both sunrise and sunset times for the location
// and date encoded in an NMEA sentence.
//
// NMEA sentences contain GPS data in UTC. The returned sunrise and sunset times are also
// in UTC. Convert to local time using time.In() if needed.
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
//   - sunrise: The sunrise time in UTC at the location specified in the NMEA sentence
//   - sunset: The sunset time in UTC at the location specified in the NMEA sentence
//   - error: An error if the sentence is invalid or unsupported
//   - time.Time{} for sunrise/sunset if the sun does not rise/set
//
// Example:
//
//	// Using RMC sentence (includes date)
//	sunrise, sunset, err := solar.SunriseSunsetFromNMEA("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A", 0, 0, 0)
//
//	// Using GGA sentence (requires date)
//	sunrise, sunset, err := solar.SunriseSunsetFromNMEA("$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47", 2024, time.June, 21)
func SunriseSunsetFromNMEA(nmea string, year int, month time.Month, day int) (time.Time, time.Time, error) {
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return SunriseSunset(loc, t)
}

// parseNMEA parses an NMEA sentence and extracts position and date information.
func parseNMEA(nmea string, year int, month time.Month, day int) (nmeaPosition, error) {
	// Remove leading/trailing whitespace
	nmea = strings.TrimSpace(nmea)

	// NMEA sentences must start with $
	if !strings.HasPrefix(nmea, "$") {
		return nmeaPosition{}, fmt.Errorf("%w: missing $ prefix", ErrInvalidNMEA)
	}

	// Split into sentence and checksum
	parts := strings.Split(nmea[1:], "*")
	if len(parts) != 2 {
		return nmeaPosition{}, fmt.Errorf("%w: missing or invalid checksum", ErrInvalidNMEA)
	}

	sentence := parts[0]
	checksumStr := parts[1]

	// Validate checksum
	if err := validateChecksum(sentence, checksumStr); err != nil {
		return nmeaPosition{}, err
	}

	// Split sentence into fields
	fields := strings.Split(sentence, ",")
	if len(fields) < 2 {
		return nmeaPosition{}, fmt.Errorf("%w: insufficient fields", ErrInvalidNMEA)
	}

	// Determine sentence type (last 3 characters of talker+sentence ID)
	sentenceType := fields[0]
	if len(sentenceType) < 3 {
		return nmeaPosition{}, fmt.Errorf("%w: invalid sentence type", ErrInvalidNMEA)
	}
	sentenceType = sentenceType[len(sentenceType)-3:]

	// Parse based on sentence type
	switch sentenceType {
	case "GGA":
		return parseGGA(fields, year, month, day)
	case "RMC":
		return parseRMC(fields)
	default:
		return nmeaPosition{}, fmt.Errorf("%w: %s (supported: GGA, RMC)", ErrUnsupportedSentence, sentenceType)
	}
}

// validateChecksum validates the NMEA sentence checksum.
func validateChecksum(sentence, checksumStr string) error {
	// Calculate checksum
	var checksum byte
	for i := 0; i < len(sentence); i++ {
		checksum ^= sentence[i]
	}

	// Parse expected checksum
	expected, err := strconv.ParseUint(checksumStr, 16, 8)
	if err != nil {
		return fmt.Errorf("%w: cannot parse checksum", ErrInvalidChecksum)
	}

	if checksum != byte(expected) {
		return fmt.Errorf("%w: calculated %02X, expected %02X", ErrInvalidChecksum, checksum, expected)
	}

	return nil
}

// parseGGA parses a GGA (GPS Fix Data) sentence.
// Format: $--GGA,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x,xx,x.x,x.x,M,x.x,M,x.x,xxxx
func parseGGA(fields []string, year int, month time.Month, day int) (nmeaPosition, error) {
	if len(fields) < 7 {
		return nmeaPosition{}, fmt.Errorf("%w: GGA sentence too short", ErrInvalidNMEA)
	}

	// GGA requires external date
	if year == 0 || month == 0 || day == 0 {
		return nmeaPosition{}, fmt.Errorf("%w: GGA sentence requires date parameter", ErrInvalidDate)
	}

	// Parse time (field 1)
	timeStr := fields[1]
	parsedTime, err := parseNMEATime(timeStr, year, month, day)
	if err != nil {
		return nmeaPosition{}, err
	}

	// Parse latitude (fields 2-3)
	lat, err := parseLatitude(fields[2], fields[3])
	if err != nil {
		return nmeaPosition{}, err
	}

	// Parse longitude (fields 4-5)
	lon, err := parseLongitude(fields[4], fields[5])
	if err != nil {
		return nmeaPosition{}, err
	}

	return nmeaPosition{
		Latitude:  lat,
		Longitude: lon,
		Time:      parsedTime,
	}, nil
}

// parseRMC parses an RMC (Recommended Minimum) sentence.
// Format: $--RMC,hhmmss.ss,A,llll.ll,a,yyyyy.yy,a,x.x,x.x,ddmmyy,x.x,a
func parseRMC(fields []string) (nmeaPosition, error) {
	if len(fields) < 10 {
		return nmeaPosition{}, fmt.Errorf("%w: RMC sentence too short", ErrInvalidNMEA)
	}

	// Check status (field 2) - should be 'A' for valid
	if fields[2] != "A" {
		return nmeaPosition{}, fmt.Errorf("%w: invalid GPS fix (status: %s)", ErrInvalidNMEA, fields[2])
	}

	// Parse date (field 9) - ddmmyy format
	dateStr := fields[9]
	if len(dateStr) != 6 {
		return nmeaPosition{}, fmt.Errorf("%w: invalid date format", ErrInvalidDate)
	}

	day, err := strconv.Atoi(dateStr[0:2])
	if err != nil {
		return nmeaPosition{}, fmt.Errorf("%w: invalid day", ErrInvalidDate)
	}

	monthInt, err := strconv.Atoi(dateStr[2:4])
	if err != nil {
		return nmeaPosition{}, fmt.Errorf("%w: invalid month", ErrInvalidDate)
	}

	year, err := strconv.Atoi(dateStr[4:6])
	if err != nil {
		return nmeaPosition{}, fmt.Errorf("%w: invalid year", ErrInvalidDate)
	}
	// Convert 2-digit year to 4-digit
	// Years 00-49 are 2000-2049, years 50-99 are 1950-1999
	if year < 50 {
		year += 2000
	} else {
		year += 1900
	}

	// Parse time (field 1)
	timeStr := fields[1]
	parsedTime, err := parseNMEATime(timeStr, year, time.Month(monthInt), day)
	if err != nil {
		return nmeaPosition{}, err
	}

	// Parse latitude (fields 3-4)
	lat, err := parseLatitude(fields[3], fields[4])
	if err != nil {
		return nmeaPosition{}, err
	}

	// Parse longitude (fields 5-6)
	lon, err := parseLongitude(fields[5], fields[6])
	if err != nil {
		return nmeaPosition{}, err
	}

	return nmeaPosition{
		Latitude:  lat,
		Longitude: lon,
		Time:      parsedTime,
	}, nil
}

// parseNMEATime parses NMEA time format (hhmmss.ss) and combines with date.
func parseNMEATime(timeStr string, year int, month time.Month, day int) (time.Time, error) {
	if len(timeStr) < 6 {
		return time.Time{}, fmt.Errorf("%w: time string too short", ErrInvalidDate)
	}

	hour, err := strconv.Atoi(timeStr[0:2])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: invalid hour", ErrInvalidDate)
	}

	minute, err := strconv.Atoi(timeStr[2:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: invalid minute", ErrInvalidDate)
	}

	second, err := strconv.Atoi(timeStr[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: invalid second", ErrInvalidDate)
	}

	// Parse fractional seconds if present
	nanosecond := 0
	if len(timeStr) > 7 && timeStr[6] == '.' {
		fracStr := timeStr[7:]
		// Pad to 9 digits for nanoseconds
		for len(fracStr) < 9 {
			fracStr += "0"
		}
		if len(fracStr) > 9 {
			fracStr = fracStr[:9]
		}
		nanosecond, err = strconv.Atoi(fracStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("%w: invalid fractional seconds", ErrInvalidDate)
		}
	}

	return time.Date(year, month, day, hour, minute, second, nanosecond, time.UTC), nil
}

// parseLatitude parses NMEA latitude format (ddmm.mmmm,N/S).
func parseLatitude(latStr, nsStr string) (float64, error) {
	if latStr == "" || nsStr == "" {
		return 0, fmt.Errorf("%w: empty latitude fields", ErrInvalidPosition)
	}

	// Find decimal point
	dotIdx := strings.Index(latStr, ".")
	if dotIdx < 2 {
		return 0, fmt.Errorf("%w: invalid latitude format", ErrInvalidPosition)
	}

	// Extract degrees (everything before last 2 digits of whole part)
	degreesStr := latStr[:dotIdx-2]
	degrees, err := strconv.ParseFloat(degreesStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid latitude degrees", ErrInvalidPosition)
	}

	// Extract minutes (last 2 digits + decimal part)
	minutesStr := latStr[dotIdx-2:]
	minutes, err := strconv.ParseFloat(minutesStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid latitude minutes", ErrInvalidPosition)
	}

	// Convert to decimal degrees
	lat := degrees + minutes/60.0

	// Apply N/S indicator
	if nsStr == "S" {
		lat = -lat
	} else if nsStr != "N" {
		return 0, fmt.Errorf("%w: invalid N/S indicator: %s", ErrInvalidPosition, nsStr)
	}

	return lat, nil
}

// parseLongitude parses NMEA longitude format (dddmm.mmmm,E/W).
func parseLongitude(lonStr, ewStr string) (float64, error) {
	if lonStr == "" || ewStr == "" {
		return 0, fmt.Errorf("%w: empty longitude fields", ErrInvalidPosition)
	}

	// Find decimal point
	dotIdx := strings.Index(lonStr, ".")
	if dotIdx < 2 {
		return 0, fmt.Errorf("%w: invalid longitude format", ErrInvalidPosition)
	}

	// Extract degrees (everything before last 2 digits of whole part)
	degreesStr := lonStr[:dotIdx-2]
	degrees, err := strconv.ParseFloat(degreesStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid longitude degrees", ErrInvalidPosition)
	}

	// Extract minutes (last 2 digits + decimal part)
	minutesStr := lonStr[dotIdx-2:]
	minutes, err := strconv.ParseFloat(minutesStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid longitude minutes", ErrInvalidPosition)
	}

	// Convert to decimal degrees
	lon := degrees + minutes/60.0

	// Apply E/W indicator
	if ewStr == "W" {
		lon = -lon
	} else if ewStr != "E" {
		return 0, fmt.Errorf("%w: invalid E/W indicator: %s", ErrInvalidPosition, ewStr)
	}

	return lon, nil
}

package solar

import (
	"fmt"
	"time"
)

// Location represents a geographical location.
// It can be created from direct latitude/longitude coordinates
// or parsed from an NMEA GPS sentence.
type Location struct {
	latitude  float64
	longitude float64
}

// NewLocation creates a Location from latitude and longitude coordinates.
//
// Parameters:
//   - latitude: Decimal degrees, positive north, negative south (-90 to +90)
//   - longitude: Decimal degrees, positive east, negative west (-180 to +180)
//
// Example:
//
//	loc := solar.NewLocation(43.65, -79.38) // Toronto, Canada
func NewLocation(latitude, longitude float64) Location {
	return Location{
		latitude:  latitude,
		longitude: longitude,
	}
}

// NewLocationFromNMEA creates a Location from an NMEA GPS sentence.
//
// Supported NMEA sentence types:
//   - GGA (Global Positioning System Fix Data)
//   - RMC (Recommended Minimum Specific GPS/Transit Data)
//
// The year, month, and day parameters are required for GGA sentences
// (which don't include date information). For RMC sentences, these
// parameters are ignored as the date is parsed from the sentence.
//
// Parameters:
//   - nmea: NMEA sentence string (e.g., "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47")
//   - year: Year (e.g., 2025) - ignored for RMC sentences
//   - month: Month (e.g., time.January) - ignored for RMC sentences
//   - day: Day of month (e.g., 15) - ignored for RMC sentences
//
// Returns:
//   - Location: The parsed location
//   - error: Any error encountered during parsing
//
// Example:
//
//	// From GGA sentence (requires date)
//	nmea := "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47"
//	loc, err := solar.NewLocationFromNMEA(nmea, 2025, time.January, 15)
//
//	// From RMC sentence (includes date, parameters ignored)
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A"
//	loc, err := solar.NewLocationFromNMEA(nmea, 0, 0, 0)
func NewLocationFromNMEA(nmea string, year int, month time.Month, day int) (Location, error) {
	pos, err := parseNMEA(nmea, year, month, day)
	if err != nil {
		return Location{}, err
	}

	return Location{
		latitude:  pos.Latitude,
		longitude: pos.Longitude,
	}, nil
}

// Latitude returns the latitude in decimal degrees.
// Positive values are north, negative values are south.
func (l Location) Latitude() float64 {
	return l.latitude
}

// Longitude returns the longitude in decimal degrees.
// Positive values are east, negative values are west.
func (l Location) Longitude() float64 {
	return l.longitude
}

// String returns a string representation of the Location.
func (l Location) String() string {
	latDir := "N"
	if l.latitude < 0 {
		latDir = "S"
	}
	lonDir := "E"
	if l.longitude < 0 {
		lonDir = "W"
	}

	return fmt.Sprintf("%.4f°%s, %.4f°%s", Abs(l.latitude), latDir, Abs(l.longitude), lonDir)
}

// Time represents a specific date for solar calculations.
// It can be created from individual date components or from a time.Time object.
// All times are treated as UTC.
type Time struct {
	when time.Time
}

// NewTime creates a Time from individual date components.
// The time is set to midnight (00:00:00) UTC.
//
// Parameters:
//   - year: Year (e.g., 2025)
//   - month: Month (e.g., time.January)
//   - day: Day of month (e.g., 15)
//
// Example:
//
//	t := solar.NewTime(2025, time.January, 15)
func NewTime(year int, month time.Month, day int) Time {
	return Time{
		when: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

// NewTimeFromDateTime creates a Time from a time.Time object.
// The time is converted to UTC if it isn't already.
//
// Parameters:
//   - when: A time.Time object
//
// Example:
//
//	now := time.Now()
//	t := solar.NewTimeFromDateTime(now)
func NewTimeFromDateTime(when time.Time) Time {
	return Time{
		when: when.UTC(),
	}
}

// NewTimeFromNMEA creates a Time from an NMEA GPS sentence.
// The time is extracted from the NMEA sentence and combined with the provided date.
//
// For RMC sentences, the date is parsed from the sentence and the provided
// year, month, day parameters are ignored. For GGA sentences, the provided
// date parameters are used.
//
// Parameters:
//   - nmea: NMEA sentence string
//   - year: Year (e.g., 2025) - ignored for RMC sentences
//   - month: Month (e.g., time.January) - ignored for RMC sentences
//   - day: Day of month (e.g., 15) - ignored for RMC sentences
//
// Returns:
//   - Time: The parsed time
//   - error: Any error encountered during parsing
//
// Example:
//
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A"
//	t, err := solar.NewTimeFromNMEA(nmea, 0, 0, 0)
func NewTimeFromNMEA(nmea string, year int, month time.Month, day int) (Time, error) {
	pos, err := parseNMEA(nmea, year, month, day)
	if err != nil {
		return Time{}, err
	}

	return Time{
		when: pos.Time,
	}, nil
}

// DateTime returns the underlying time.Time value in UTC.
func (t Time) DateTime() time.Time {
	return t.when
}

// Year returns the year component.
func (t Time) Year() int {
	return t.when.Year()
}

// Month returns the month component.
func (t Time) Month() time.Month {
	return t.when.Month()
}

// Day returns the day of month component.
func (t Time) Day() int {
	return t.when.Day()
}

// String returns a string representation of the Time.
func (t Time) String() string {
	return t.when.Format("2006-01-02")
}

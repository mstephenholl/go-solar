package solar

import "time"

// TwilightType represents the type of twilight for dawn/dusk calculations.
// Twilight is the period between daylight and darkness (or vice versa) when
// the sun is below the horizon but its light is still visible.
type TwilightType int

const (
	// Civil twilight occurs when the sun is between 0° and 6° below the horizon.
	// This is the most commonly used definition for dawn and dusk in everyday contexts.
	// During civil twilight, there is enough natural light for most outdoor activities
	// without artificial lighting. The brightest stars and planets are visible.
	Civil TwilightType = iota

	// Nautical twilight occurs when the sun is between 6° and 12° below the horizon.
	// During nautical twilight, the horizon is still visible at sea, allowing sailors
	// to take reliable star sights for navigation. General ground outlines are visible,
	// but detailed outdoor work is difficult.
	Nautical

	// Astronomical twilight occurs when the sun is between 12° and 18° below the horizon.
	// During astronomical twilight, the sky is dark enough for most astronomical
	// observations, though the Sun's light may still interfere with observing
	// extremely faint objects. Beyond astronomical twilight is true night.
	Astronomical
)

// twilightAngle returns the solar elevation angle for the given twilight type.
func twilightAngle(t TwilightType) float64 {
	switch t {
	case Nautical:
		return NauticalTwilightAngle
	case Astronomical:
		return AstronomicalTwilightAngle
	default: // Civil is the default
		return CivilTwilightAngle
	}
}

// Dawn calculates the dawn time for a given location and date.
// Dawn is the beginning of morning twilight, when the sun reaches the specified
// angle below the horizon and natural light begins to appear.
//
// By default, civil twilight (-6°) is used, which is the most common definition
// of dawn in everyday contexts. You can optionally specify Nautical or Astronomical
// twilight types for specialized applications.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate dawn
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dawn time in UTC (time.Time{} if the sun never reaches the twilight angle on this day)
//
// Example:
//
//	// Calculate civil dawn (default)
//	dawn := solar.Dawn(40.7128, -74.0060, 2024, time.June, 21)
//
//	// Calculate nautical dawn
//	dawn := solar.Dawn(40.7128, -74.0060, 2024, time.June, 21, solar.Nautical)
//
//	// Calculate astronomical dawn
//	dawn := solar.Dawn(40.7128, -74.0060, 2024, time.June, 21, solar.Astronomical)
func Dawn(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) time.Time {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate dawn using TimeOfElevation with the appropriate angle
	dawn, _ := TimeOfElevation(latitude, longitude, twilightAngle(tt), year, month, day)
	return dawn
}

// Dusk calculates the dusk time for a given location and date.
// Dusk is the end of evening twilight, when the sun reaches the specified
// angle below the horizon and natural light fades to darkness.
//
// By default, civil twilight (-6°) is used, which is the most common definition
// of dusk in everyday contexts. You can optionally specify Nautical or Astronomical
// twilight types for specialized applications.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate dusk
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dusk time in UTC (time.Time{} if the sun never reaches the twilight angle on this day)
//
// Example:
//
//	// Calculate civil dusk (default)
//	dusk := solar.Dusk(40.7128, -74.0060, 2024, time.June, 21)
//
//	// Calculate nautical dusk
//	dusk := solar.Dusk(40.7128, -74.0060, 2024, time.June, 21, solar.Nautical)
//
//	// Calculate astronomical dusk
//	dusk := solar.Dusk(40.7128, -74.0060, 2024, time.June, 21, solar.Astronomical)
func Dusk(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) time.Time {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate dusk using TimeOfElevation with the appropriate angle
	_, dusk := TimeOfElevation(latitude, longitude, twilightAngle(tt), year, month, day)
	return dusk
}

// DawnDusk calculates both dawn and dusk times for a given location and date.
// This is more efficient than calling Dawn() and Dusk() separately.
//
// By default, civil twilight (-6°) is used. You can optionally specify Nautical
// or Astronomical twilight types for specialized applications.
//
// Parameters:
//   - latitude: Latitude in decimal degrees (-90 to +90, negative for South)
//   - longitude: Longitude in decimal degrees (-180 to +180, negative for West)
//   - year, month, day: The date in UTC for which to calculate dawn and dusk
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - dawn: Dawn time in UTC (time.Time{} if never occurs)
//   - dusk: Dusk time in UTC (time.Time{} if never occurs)
//
// Example:
//
//	// Calculate civil dawn and dusk (default)
//	dawn, dusk := solar.DawnDusk(40.7128, -74.0060, 2024, time.June, 21)
//
//	// Calculate nautical dawn and dusk
//	dawn, dusk := solar.DawnDusk(40.7128, -74.0060, 2024, time.June, 21, solar.Nautical)
func DawnDusk(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) (dawn, dusk time.Time) {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate both times using TimeOfElevation with the appropriate angle
	return TimeOfElevation(latitude, longitude, twilightAngle(tt), year, month, day)
}

// DawnFromNMEA calculates the dawn time from an NMEA GPS sentence.
//
// By default, civil twilight (-6°) is used. You can optionally specify Nautical
// or Astronomical twilight types for specialized applications.
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (GGA or RMC format)
//   - year: Year (ignored for RMC sentences that include date)
//   - month: Month (ignored for RMC sentences that include date)
//   - day: Day of month (ignored for RMC sentences that include date)
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dawn time in UTC (time.Time{} if never occurs)
//   - An error if the NMEA sentence is invalid or cannot be parsed
//
// Example:
//
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"
//	dawn, err := solar.DawnFromNMEA(nmea, 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
func DawnFromNMEA(nmea string, year int, month time.Month, day int, twilightType ...TwilightType) (time.Time, error) {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Parse the NMEA sentence to get position and time
	pos, err := parseNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	// Calculate dawn using the parsed location and date
	dawn := Dawn(pos.Latitude, pos.Longitude, pos.Time.Year(), pos.Time.Month(), pos.Time.Day(), tt)
	return dawn, nil
}

// DuskFromNMEA calculates the dusk time from an NMEA GPS sentence.
//
// By default, civil twilight (-6°) is used. You can optionally specify Nautical
// or Astronomical twilight types for specialized applications.
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (GGA or RMC format)
//   - year: Year (ignored for RMC sentences that include date)
//   - month: Month (ignored for RMC sentences that include date)
//   - day: Day of month (ignored for RMC sentences that include date)
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dusk time in UTC (time.Time{} if never occurs)
//   - An error if the NMEA sentence is invalid or cannot be parsed
//
// Example:
//
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"
//	dusk, err := solar.DuskFromNMEA(nmea, 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
func DuskFromNMEA(nmea string, year int, month time.Month, day int, twilightType ...TwilightType) (time.Time, error) {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Parse the NMEA sentence to get position and time
	pos, err := parseNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	// Calculate dusk using the parsed location and date
	dusk := Dusk(pos.Latitude, pos.Longitude, pos.Time.Year(), pos.Time.Month(), pos.Time.Day(), tt)
	return dusk, nil
}

// DawnDuskFromNMEA calculates both dawn and dusk times from an NMEA GPS sentence.
// This is more efficient than calling DawnFromNMEA() and DuskFromNMEA() separately.
//
// By default, civil twilight (-6°) is used. You can optionally specify Nautical
// or Astronomical twilight types for specialized applications.
//
// Parameters:
//   - nmea: An NMEA 0183 sentence string (GGA or RMC format)
//   - year: Year (ignored for RMC sentences that include date)
//   - month: Month (ignored for RMC sentences that include date)
//   - day: Day of month (ignored for RMC sentences that include date)
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - dawn: Dawn time in UTC (time.Time{} if never occurs)
//   - dusk: Dusk time in UTC (time.Time{} if never occurs)
//   - An error if the NMEA sentence is invalid or cannot be parsed
//
// Example:
//
//	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"
//	dawn, dusk, err := solar.DawnDuskFromNMEA(nmea, 0, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
func DawnDuskFromNMEA(nmea string, year int, month time.Month, day int, twilightType ...TwilightType) (dawn, dusk time.Time, err error) {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Parse the NMEA sentence to get position and time
	pos, parseErr := parseNMEA(nmea, year, month, day)
	if parseErr != nil {
		return time.Time{}, time.Time{}, parseErr
	}

	// Calculate both times using the parsed location and date
	dawn, dusk = DawnDusk(pos.Latitude, pos.Longitude, pos.Time.Year(), pos.Time.Month(), pos.Time.Day(), tt)
	return dawn, dusk, nil
}

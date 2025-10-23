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

// dawnInternal is the internal implementation with old signature
func dawnInternal(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) time.Time {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate dawn using timeOfElevationInternal with the appropriate angle
	dawn, _ := timeOfElevationInternal(latitude, longitude, twilightAngle(tt), year, month, day)
	return dawn
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
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dawn time in UTC (time.Time{} if the sun never reaches the twilight angle on this day)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	// Calculate civil dawn (default)
//	dawn := solar.Dawn(loc, t)
//
//	// Calculate nautical dawn
//	dawn := solar.Dawn(loc, t, solar.Nautical)
//
//	// Calculate astronomical dawn
//	dawn := solar.Dawn(loc, t, solar.Astronomical)
func Dawn(loc Location, t Time, twilightType ...TwilightType) time.Time {
	return dawnInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day(), twilightType...)
}

// duskInternal is the internal implementation with old signature
func duskInternal(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) time.Time {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate dusk using timeOfElevationInternal with the appropriate angle
	_, dusk := timeOfElevationInternal(latitude, longitude, twilightAngle(tt), year, month, day)
	return dusk
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
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - Dusk time in UTC (time.Time{} if the sun never reaches the twilight angle on this day)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	// Calculate civil dusk (default)
//	dusk := solar.Dusk(loc, t)
//
//	// Calculate nautical dusk
//	dusk := solar.Dusk(loc, t, solar.Nautical)
//
//	// Calculate astronomical dusk
//	dusk := solar.Dusk(loc, t, solar.Astronomical)
func Dusk(loc Location, t Time, twilightType ...TwilightType) time.Time {
	return duskInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day(), twilightType...)
}

// dawnDuskInternal is the internal implementation with old signature
func dawnDuskInternal(latitude, longitude float64, year int, month time.Month, day int, twilightType ...TwilightType) (dawn, dusk time.Time) {
	// Determine twilight type (default to Civil)
	var tt TwilightType
	if len(twilightType) > 0 {
		tt = twilightType[0]
	} else {
		tt = Civil
	}

	// Calculate both times using timeOfElevationInternal with the appropriate angle
	return timeOfElevationInternal(latitude, longitude, twilightAngle(tt), year, month, day)
}

// DawnDusk calculates both dawn and dusk times for a given location and date.
// This is more efficient than calling Dawn() and Dusk() separately.
//
// By default, civil twilight (-6°) is used. You can optionally specify Nautical
// or Astronomical twilight types for specialized applications.
//
// Parameters:
//   - loc: Location created via NewLocation() or NewLocationFromNMEA()
//   - t: Time created via NewTime(), NewTimeFromDateTime(), or NewTimeFromNMEA()
//   - twilightType: Optional twilight type (Civil, Nautical, or Astronomical). Defaults to Civil.
//
// Returns:
//   - dawn: Dawn time in UTC (time.Time{} if never occurs)
//   - dusk: Dusk time in UTC (time.Time{} if never occurs)
//
// Example:
//
//	loc := solar.NewLocation(40.7128, -74.0060)
//	t := solar.NewTime(2024, time.June, 21)
//	// Calculate civil dawn and dusk (default)
//	dawn, dusk := solar.DawnDusk(loc, t)
//
//	// Calculate nautical dawn and dusk
//	dawn, dusk := solar.DawnDusk(loc, t, solar.Nautical)
func DawnDusk(loc Location, t Time, twilightType ...TwilightType) (dawn, dusk time.Time) {
	return dawnDuskInternal(loc.Latitude(), loc.Longitude(), t.Year(), t.Month(), t.Day(), twilightType...)
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
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	return Dawn(loc, t, twilightType...), nil
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
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, err
	}

	return Dusk(loc, t, twilightType...), nil
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
	loc, err := NewLocationFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	t, err := NewTimeFromNMEA(nmea, year, month, day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	dawn, dusk = DawnDusk(loc, t, twilightType...)
	return dawn, dusk, nil
}

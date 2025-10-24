package solar_test

import (
	"fmt"
	"time"

	"github.com/mstephenholl/go-solar"
)

// ExampleSunrise demonstrates calculating just the sunrise time
// for Toronto, Canada on January 1, 2000.
func ExampleSunrise() {
	// Create location for Toronto
	loc := solar.NewLocation(43.65, -79.38)

	// Create time for January 1, 2000
	t := solar.NewTime(2000, time.January, 1)

	// Calculate sunrise
	rise, err := solar.Sunrise(loc, t)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
	// Output:
	// Sunrise: 12:50:59 UTC
}

// ExampleSunset demonstrates calculating just the sunset time
// for Toronto, Canada on January 1, 2000.
func ExampleSunset() {
	// Create location for Toronto
	loc := solar.NewLocation(43.65, -79.38)

	// Create time for January 1, 2000
	t := solar.NewTime(2000, time.January, 1)

	// Calculate sunset
	set, err := solar.Sunset(loc, t)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
	// Output:
	// Sunset: 21:50:37 UTC
}

// ExampleSunriseSunset demonstrates basic sunrise and sunset calculation
// for Toronto, Canada on January 1, 2000.
func ExampleSunriseSunset() {
	// Create location for Toronto
	loc := solar.NewLocation(43.65, -79.38)

	// Create time for January 1, 2000
	t := solar.NewTime(2000, time.January, 1)

	// Calculate sunrise and sunset
	rise, set, err := solar.SunriseSunset(loc, t)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
	fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
	// Output:
	// Sunrise: 12:50:59 UTC
	// Sunset: 21:50:37 UTC
}

// ExampleSunriseSunset_polarNight demonstrates the case where the sun
// never rises (polar night).
func ExampleSunriseSunset_polarNight() {
	// Create location for Igloolik, Nunavut
	loc := solar.NewLocation(69.3321443, -81.6781126)

	// Create time for June 25, 2020 (midnight sun period)
	t := solar.NewTime(2020, time.June, 25)

	// Calculate sunrise and sunset
	_, _, err := solar.SunriseSunset(loc, t)

	// Check for error (sun never rises or sets)
	if err != nil {
		fmt.Println("The sun does not rise or set on this day")
	}
	// Output:
	// The sun does not rise or set on this day
}

// ExampleElevation demonstrates calculating the sun's elevation angle
// at a specific time and location.
func ExampleElevation() {
	// Create location for New York City
	loc := solar.NewLocation(40.7128, -74.0060)

	// Check sun elevation at noon UTC on summer solstice
	when := time.Date(2022, time.June, 21, 12, 0, 0, 0, time.UTC)
	elevation := solar.Elevation(loc, when)

	fmt.Printf("Sun elevation: %.1f degrees\n", elevation)
	// Output:
	// Sun elevation: 26.5 degrees
}

// ExampleTimeOfElevation demonstrates finding when the sun reaches
// a specific elevation angle.
func ExampleTimeOfElevation() {
	// Create location for London
	loc := solar.NewLocation(51.5072, -0.1276)
	t := solar.NewTime(2022, time.June, 21)

	// Find when sun is at 10 degrees above horizon
	morning, evening := solar.TimeOfElevation(loc, 10.0, t)

	fmt.Printf("Morning: %s\n", morning.Format("15:04 MST"))
	fmt.Printf("Evening: %s\n", evening.Format("15:04 MST"))
	// Output:
	// Morning: 05:06 UTC
	// Evening: 18:58 UTC
}

// ExampleAbs demonstrates the generic Abs function with different types.
func ExampleAbs() {
	// Works with int
	fmt.Println(solar.Abs(-5))

	// Works with int64
	fmt.Println(solar.Abs(int64(-100)))

	// Works with float64
	fmt.Println(solar.Abs(-3.14))

	// Output:
	// 5
	// 100
	// 3.14
}

// ExampleAlmostEqual demonstrates floating-point comparison with tolerance.
func ExampleAlmostEqual() {
	a := 1.0
	b := 1.00001

	// These are almost equal within tolerance
	fmt.Println(solar.AlmostEqual(a, b, 0.001))

	// But not with stricter tolerance
	fmt.Println(solar.AlmostEqual(a, b, 0.00001))

	// Output:
	// true
	// false
}

// ExampleAzimuth demonstrates calculating the solar azimuth angle.
// The azimuth is the sun's compass direction measured clockwise from north.
func ExampleAzimuth() {
	// Create location for Toronto
	loc := solar.NewLocation(43.65, -79.38)

	// Calculate azimuth for January 1, 2000 at 5:00 PM UTC (noon local time)
	when := time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC)
	azimuth := solar.Azimuth(loc, when)

	fmt.Printf("Azimuth: %.1f degrees (South)\n", azimuth)
	// Output:
	// Azimuth: 174.8 degrees (South)
}

// ExampleDawn demonstrates calculating civil dawn (beginning of morning twilight).
func ExampleDawn() {
	// Toronto coordinates
	loc := solar.NewLocation(43.65, -79.38)
	t := solar.NewTime(2000, time.January, 1)

	// Calculate civil dawn (default)
	dawn := solar.Dawn(loc, t)

	fmt.Printf("Civil dawn: %s\n", dawn.Format("15:04 MST"))
	// Output:
	// Civil dawn: 12:18 UTC
}

// ExampleDusk demonstrates calculating civil dusk (end of evening twilight).
func ExampleDusk() {
	// Toronto coordinates
	loc := solar.NewLocation(43.65, -79.38)
	t := solar.NewTime(2000, time.January, 1)

	// Calculate civil dusk (default)
	dusk := solar.Dusk(loc, t)

	fmt.Printf("Civil dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Civil dusk: 22:23 UTC
}

// ExampleDawnDusk demonstrates calculating both dawn and dusk times.
func ExampleDawnDusk() {
	// Toronto coordinates
	loc := solar.NewLocation(43.65, -79.38)
	t := solar.NewTime(2000, time.January, 1)

	// Calculate civil dawn and dusk
	dawn, dusk := solar.DawnDusk(loc, t)

	fmt.Printf("Dawn: %s\n", dawn.Format("15:04 MST"))
	fmt.Printf("Dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Dawn: 12:18 UTC
	// Dusk: 22:23 UTC
}

// ExampleDawn_nautical demonstrates calculating nautical dawn.
func ExampleDawn_nautical() {
	// Toronto coordinates
	loc := solar.NewLocation(43.65, -79.38)
	t := solar.NewTime(2000, time.January, 1)

	// Calculate nautical dawn (sun at -12° below horizon)
	nauticalDawn := solar.Dawn(loc, t, solar.Nautical)

	fmt.Printf("Nautical dawn: %s\n", nauticalDawn.Format("15:04 MST"))
	// Output:
	// Nautical dawn: 11:42 UTC
}

// ExampleDawnDusk_astronomical demonstrates calculating astronomical twilight.
func ExampleDawnDusk_astronomical() {
	// Toronto coordinates
	loc := solar.NewLocation(43.65, -79.38)
	t := solar.NewTime(2000, time.January, 1)

	// Calculate astronomical dawn and dusk (sun at -18° below horizon)
	dawn, dusk := solar.DawnDusk(loc, t, solar.Astronomical)

	fmt.Printf("Astronomical dawn: %s\n", dawn.Format("15:04 MST"))
	fmt.Printf("Astronomical dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Astronomical dawn: 11:07 UTC
	// Astronomical dusk: 23:34 UTC
}

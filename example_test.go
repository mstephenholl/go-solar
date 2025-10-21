package solar_test

import (
	"fmt"
	"time"

	"github.com/mstephenholl/go-solar"
)

// ExampleSunrise demonstrates calculating just the sunrise time
// for Toronto, Canada on January 1, 2000.
func ExampleSunrise() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate sunrise for January 1, 2000
	rise := solar.Sunrise(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
	// Output:
	// Sunrise: 12:50:59 UTC
}

// ExampleSunset demonstrates calculating just the sunset time
// for Toronto, Canada on January 1, 2000.
func ExampleSunset() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate sunset for January 1, 2000
	set := solar.Sunset(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
	// Output:
	// Sunset: 21:50:37 UTC
}

// ExampleSunriseSunset demonstrates basic sunrise and sunset calculation
// for Toronto, Canada on January 1, 2000.
func ExampleSunriseSunset() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate sunrise and sunset for January 1, 2000
	rise, set := solar.SunriseSunset(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
	fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
	// Output:
	// Sunrise: 12:50:59 UTC
	// Sunset: 21:50:37 UTC
}

// ExampleSunriseSunset_polarNight demonstrates the case where the sun
// never rises (polar night).
func ExampleSunriseSunset_polarNight() {
	// Igloolik, Nunavut during polar night
	latitude := 69.3321443
	longitude := -81.6781126

	rise, set := solar.SunriseSunset(latitude, longitude, 2020, time.June, 25)

	// Check for no sunrise/sunset
	if rise.IsZero() && set.IsZero() {
		fmt.Println("The sun does not rise or set on this day")
	}
	// Output:
	// The sun does not rise or set on this day
}

// ExampleElevation demonstrates calculating the sun's elevation angle
// at a specific time and location.
func ExampleElevation() {
	// New York City coordinates
	latitude := 40.7128
	longitude := -74.0060

	// Check sun elevation at noon UTC on summer solstice
	when := time.Date(2022, time.June, 21, 12, 0, 0, 0, time.UTC)
	elevation := solar.Elevation(latitude, longitude, when)

	fmt.Printf("Sun elevation: %.1f degrees\n", elevation)
	// Output:
	// Sun elevation: 26.5 degrees
}

// ExampleTimeOfElevation demonstrates finding when the sun reaches
// a specific elevation angle.
func ExampleTimeOfElevation() {
	// London coordinates
	latitude := 51.5072
	longitude := -0.1276

	// Find when sun is at 10 degrees above horizon
	morning, evening := solar.TimeOfElevation(latitude, longitude, 10.0, 2022, time.June, 21)

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

// ExampleMin demonstrates the generic Min function.
func ExampleMin() {
	// Works with int
	fmt.Println(solar.Min(5, 10))

	// Works with float64
	fmt.Println(solar.Min(3.14, 2.71))

	// Output:
	// 5
	// 2.71
}

// ExampleMax demonstrates the generic Max function.
func ExampleMax() {
	// Works with int
	fmt.Println(solar.Max(5, 10))

	// Works with float64
	fmt.Println(solar.Max(3.14, 2.71))

	// Output:
	// 10
	// 3.14
}

// ExampleClamp demonstrates restricting a value to a range.
func ExampleClamp() {
	// Value within range
	fmt.Println(solar.Clamp(5, 0, 10))

	// Value below minimum
	fmt.Println(solar.Clamp(-5, 0, 10))

	// Value above maximum
	fmt.Println(solar.Clamp(15, 0, 10))

	// Output:
	// 5
	// 0
	// 10
}

// ExampleDegreesToRadians demonstrates angle conversion.
func ExampleDegreesToRadians() {
	radians := solar.DegreesToRadians(180.0)
	fmt.Printf("%.6f\n", radians)
	// Output:
	// 3.141593
}

// ExampleRadiansToDegrees demonstrates angle conversion.
func ExampleRadiansToDegrees() {
	degrees := solar.RadiansToDegrees(3.14159265)
	fmt.Printf("%.1f\n", degrees)
	// Output:
	// 180.0
}

// ExampleAzimuth demonstrates calculating the solar azimuth angle.
// The azimuth is the sun's compass direction measured clockwise from north.
func ExampleAzimuth() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate azimuth for January 1, 2000 at 5:00 PM UTC (noon local time)
	when := time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC)
	azimuth := solar.Azimuth(latitude, longitude, when)

	fmt.Printf("Azimuth: %.1f degrees (South)\n", azimuth)
	// Output:
	// Azimuth: 174.8 degrees (South)
}

// ExampleDawn demonstrates calculating civil dawn (beginning of morning twilight).
func ExampleDawn() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate civil dawn for January 1, 2000 (default)
	dawn := solar.Dawn(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Civil dawn: %s\n", dawn.Format("15:04 MST"))
	// Output:
	// Civil dawn: 12:18 UTC
}

// ExampleDusk demonstrates calculating civil dusk (end of evening twilight).
func ExampleDusk() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate civil dusk for January 1, 2000 (default)
	dusk := solar.Dusk(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Civil dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Civil dusk: 22:23 UTC
}

// ExampleDawnDusk demonstrates calculating both dawn and dusk times.
func ExampleDawnDusk() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate civil dawn and dusk for January 1, 2000
	dawn, dusk := solar.DawnDusk(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Dawn: %s\n", dawn.Format("15:04 MST"))
	fmt.Printf("Dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Dawn: 12:18 UTC
	// Dusk: 22:23 UTC
}

// ExampleDawn_nautical demonstrates calculating nautical dawn.
func ExampleDawn_nautical() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate nautical dawn (sun at -12° below horizon)
	nauticalDawn := solar.Dawn(latitude, longitude, 2000, time.January, 1, solar.Nautical)

	fmt.Printf("Nautical dawn: %s\n", nauticalDawn.Format("15:04 MST"))
	// Output:
	// Nautical dawn: 11:42 UTC
}

// ExampleDawnDusk_astronomical demonstrates calculating astronomical twilight.
func ExampleDawnDusk_astronomical() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate astronomical dawn and dusk (sun at -18° below horizon)
	dawn, dusk := solar.DawnDusk(latitude, longitude, 2000, time.January, 1, solar.Astronomical)

	fmt.Printf("Astronomical dawn: %s\n", dawn.Format("15:04 MST"))
	fmt.Printf("Astronomical dusk: %s\n", dusk.Format("15:04 MST"))
	// Output:
	// Astronomical dawn: 11:07 UTC
	// Astronomical dusk: 23:34 UTC
}

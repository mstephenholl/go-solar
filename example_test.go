package sunrise_test

import (
	"fmt"
	"time"

	sunrise "github.com/mstephenholl/go-solar"
)

// ExampleSunriseSunset demonstrates basic sunrise and sunset calculation
// for Toronto, Canada on January 1, 2000.
func ExampleSunriseSunset() {
	// Toronto coordinates
	latitude := 43.65
	longitude := -79.38

	// Calculate sunrise and sunset for January 1, 2000
	rise, set := sunrise.SunriseSunset(latitude, longitude, 2000, time.January, 1)

	fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
	fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
	// Output:
	// Sunrise: 12:51:00 UTC
	// Sunset: 21:50:36 UTC
}

// ExampleSunriseSunset_polarNight demonstrates the case where the sun
// never rises (polar night).
func ExampleSunriseSunset_polarNight() {
	// Igloolik, Nunavut during polar night
	latitude := 69.3321443
	longitude := -81.6781126

	rise, set := sunrise.SunriseSunset(latitude, longitude, 2020, time.June, 25)

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
	elevation := sunrise.Elevation(latitude, longitude, when)

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
	morning, evening := sunrise.TimeOfElevation(latitude, longitude, 10.0, 2022, time.June, 21)

	fmt.Printf("Morning: %s\n", morning.Format("15:04 MST"))
	fmt.Printf("Evening: %s\n", evening.Format("15:04 MST"))
	// Output:
	// Morning: 05:06 UTC
	// Evening: 18:58 UTC
}

// ExampleAbs demonstrates the generic Abs function with different types.
func ExampleAbs() {
	// Works with int
	fmt.Println(sunrise.Abs(-5))

	// Works with int64
	fmt.Println(sunrise.Abs(int64(-100)))

	// Works with float64
	fmt.Println(sunrise.Abs(-3.14))

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
	fmt.Println(sunrise.AlmostEqual(a, b, 0.001))

	// But not with stricter tolerance
	fmt.Println(sunrise.AlmostEqual(a, b, 0.00001))

	// Output:
	// true
	// false
}

// ExampleMin demonstrates the generic Min function.
func ExampleMin() {
	// Works with int
	fmt.Println(sunrise.Min(5, 10))

	// Works with float64
	fmt.Println(sunrise.Min(3.14, 2.71))

	// Output:
	// 5
	// 2.71
}

// ExampleMax demonstrates the generic Max function.
func ExampleMax() {
	// Works with int
	fmt.Println(sunrise.Max(5, 10))

	// Works with float64
	fmt.Println(sunrise.Max(3.14, 2.71))

	// Output:
	// 10
	// 3.14
}

// ExampleClamp demonstrates restricting a value to a range.
func ExampleClamp() {
	// Value within range
	fmt.Println(sunrise.Clamp(5, 0, 10))

	// Value below minimum
	fmt.Println(sunrise.Clamp(-5, 0, 10))

	// Value above maximum
	fmt.Println(sunrise.Clamp(15, 0, 10))

	// Output:
	// 5
	// 0
	// 10
}

// ExampleDegreesToRadians demonstrates angle conversion.
func ExampleDegreesToRadians() {
	radians := sunrise.DegreesToRadians(180.0)
	fmt.Printf("%.6f\n", radians)
	// Output:
	// 3.141593
}

// ExampleRadiansToDegrees demonstrates angle conversion.
func ExampleRadiansToDegrees() {
	degrees := sunrise.RadiansToDegrees(3.14159265)
	fmt.Printf("%.1f\n", degrees)
	// Output:
	// 180.0
}

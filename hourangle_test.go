package solar

import (
	"math"
	"testing"
)

var dataHourAngle = []struct {
	inLatitude    float64
	inDeclination float64
	out           float64
}{
	// Normal cases - sun rises and sets
	// 1970-01-01 - prime meridian
	{0, -22.97753, 90.904793},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{43.65, -23.01689, 67.453649},
	// 2004-04-01 - (52° N, 5° E)
	{52, 4.75374, 97.477355},
}

func TestHourAngle(t *testing.T) {
	for _, tt := range dataHourAngle {
		v := hourAngle(tt.inLatitude, tt.inDeclination)
		if Round(v, DefaultPlaces) != Round(tt.out, DefaultPlaces) {
			t.Fatalf("%f != %f", v, tt.out)
		}
	}
}

// TesthourAngle_SunNeverRises tests the polar night case where the sun never rises.
// This occurs when numerator/denominator > 1, typically in polar regions during winter.
func TestHourAngle_SunNeverRises(t *testing.T) {
	// Arctic winter: high latitude with negative declination (sun south of equator)
	// This creates a situation where the sun never gets above the horizon
	latitude := 75.0     // Far north latitude
	declination := -20.0 // Sun is south of equator (winter)

	result := hourAngle(latitude, declination)

	if result != math.MaxFloat64 {
		t.Errorf("Expected math.MaxFloat64 for sun never rising, got %v", result)
	}
}

// TesthourAngle_SunNeverSets tests the midnight sun case where the sun never sets.
// This occurs when numerator/denominator < -1, typically in polar regions during summer.
func TestHourAngle_SunNeverSets(t *testing.T) {
	// Arctic summer: high latitude with positive declination (sun north of equator)
	// This creates a situation where the sun stays above the horizon all day
	latitude := 75.0    // Far north latitude
	declination := 20.0 // Sun is north of equator (summer)

	result := hourAngle(latitude, declination)

	if result != -1*math.MaxFloat64 {
		t.Errorf("Expected -1*math.MaxFloat64 for sun never setting, got %v", result)
	}
}

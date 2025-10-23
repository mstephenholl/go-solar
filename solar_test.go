package solar

import (
	"testing"
	"time"
)

var dataSunriseSunset = []struct {
	inLatitude  float64
	inLongitude float64
	inYear      int
	inMonth     time.Month
	inDay       int
	outSunrise  time.Time
	outSunset   time.Time
}{
	// 1970-01-01 - prime meridian
	{
		0, 0,
		1970, time.January, 1,
		time.Date(1970, time.January, 1, 5, 59, 54, 0, time.UTC),
		time.Date(1970, time.January, 1, 18, 7, 8, 0, time.UTC),
	},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{
		43.65, -79.38,
		2000, time.January, 1,
		time.Date(2000, time.January, 1, 12, 50, 59, 0, time.UTC),
		time.Date(2000, time.January, 1, 21, 50, 37, 0, time.UTC),
	},
	// 2004-04-01 - (52° N, 5° E)
	{
		52, 5,
		2004, time.April, 1,
		time.Date(2004, time.April, 1, 5, 13, 39, 0, time.UTC),
		time.Date(2004, time.April, 1, 18, 13, 28, 0, time.UTC),
	},
	// 2020-06-15 - Igloolik, Canada
	{
		69.3321443, -81.6781126,
		2020, time.June, 25,
		time.Time{},
		time.Time{},
	},
}

func TestSunriseSunset(t *testing.T) {
	for _, tt := range dataSunriseSunset {
		loc := NewLocation(tt.inLatitude, tt.inLongitude)
		tm := NewTime(tt.inYear, tt.inMonth, tt.inDay)
		vSunrise, vSunset, err := SunriseSunset(loc, tm)

		// For polar regions, we expect an error
		if tt.outSunrise.IsZero() && tt.outSunset.IsZero() {
			if err == nil {
				t.Fatalf("Expected error for polar region, got nil")
			}
			continue
		}

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if vSunrise != tt.outSunrise {
			t.Fatalf("%s != %s", vSunrise.String(), tt.outSunrise.String())
		}
		if vSunset != tt.outSunset {
			t.Fatalf("%s != %s", vSunset.String(), tt.outSunset.String())
		}
	}
}

// Benchmark for the Sunrise function
func BenchmarkSunrise(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Sunrise(loc, tm)
	}
}

// Benchmark for the Sunset function
func BenchmarkSunset(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Sunset(loc, tm)
	}
}

// Benchmark for the SunriseSunset function (both values)
func BenchmarkSunriseSunset(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = SunriseSunset(loc, tm)
	}
}

// Benchmark for polar regions (edge case)
func BenchmarkSunriseSunset_Polar(b *testing.B) {
	loc := NewLocation(69.3321443, -81.6781126)
	tm := NewTime(2020, time.June, 25)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = SunriseSunset(loc, tm)
	}
}

// Benchmark for different latitudes (parallel)
func BenchmarkSunriseSunset_Latitudes(b *testing.B) {
	testCases := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{"Equator", 0, 0},
		{"Tropical", 23.5, -45},
		{"MidLatitude", 45, -75},
		{"Polar", 70, -80},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			loc := NewLocation(tc.latitude, tc.longitude)
			tm := NewTime(2024, time.June, 21)
			for i := 0; i < b.N; i++ {
				_, _, _ = SunriseSunset(loc, tm)
			}
		})
	}
}

// TestSunriseSunset_PolarNight tests that ErrSunNeverRises is returned
// for polar night conditions (winter in polar regions).
func TestSunriseSunset_PolarNight(t *testing.T) {
	// Arctic location (Svalbard, Norway) in December - polar night
	loc := NewLocation(78.2232, 15.6267)
	tm := NewTime(2024, time.December, 21) // Winter solstice

	_, _, err := SunriseSunset(loc, tm)

	if err != ErrSunNeverRises {
		t.Errorf("Expected ErrSunNeverRises for polar night, got: %v", err)
	}
}

// TestSunriseSunset_MidnightSun tests that ErrSunNeverSets is returned
// for midnight sun conditions (summer in polar regions).
func TestSunriseSunset_MidnightSun(t *testing.T) {
	// Arctic location (Igloolik, Nunavut) in June - midnight sun
	loc := NewLocation(69.3321443, -81.6781126)
	tm := NewTime(2020, time.June, 25)

	_, _, err := SunriseSunset(loc, tm)

	if err != ErrSunNeverSets {
		t.Errorf("Expected ErrSunNeverSets for midnight sun, got: %v", err)
	}
}

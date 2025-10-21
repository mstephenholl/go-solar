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
		time.Date(1970, time.January, 1, 18, 0o7, 0o7, 0, time.UTC),
	},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{
		43.65, -79.38,
		2000, time.January, 1,
		time.Date(2000, time.January, 1, 12, 51, 0o0, 0, time.UTC),
		time.Date(2000, time.January, 1, 21, 50, 36, 0, time.UTC),
	},
	// 2004-04-01 - (52° N, 5° E)
	{
		52, 5,
		2004, time.April, 1,
		time.Date(2004, time.April, 1, 5, 13, 40, 0, time.UTC),
		time.Date(2004, time.April, 1, 18, 13, 27, 0, time.UTC),
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
		vSunrise, vSunset := SunriseSunset(tt.inLatitude, tt.inLongitude, tt.inYear, tt.inMonth, tt.inDay)
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
	latitude := 43.65
	longitude := -79.38
	year := 2024
	month := time.June
	day := 21

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sunrise(latitude, longitude, year, month, day)
	}
}

// Benchmark for the Sunset function
func BenchmarkSunset(b *testing.B) {
	latitude := 43.65
	longitude := -79.38
	year := 2024
	month := time.June
	day := 21

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sunset(latitude, longitude, year, month, day)
	}
}

// Benchmark for the SunriseSunset function (both values)
func BenchmarkSunriseSunset(b *testing.B) {
	latitude := 43.65
	longitude := -79.38
	year := 2024
	month := time.June
	day := 21

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SunriseSunset(latitude, longitude, year, month, day)
	}
}

// Benchmark for polar regions (edge case)
func BenchmarkSunriseSunset_Polar(b *testing.B) {
	latitude := 69.3321443
	longitude := -81.6781126
	year := 2020
	month := time.June
	day := 25

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SunriseSunset(latitude, longitude, year, month, day)
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
			for i := 0; i < b.N; i++ {
				_, _ = SunriseSunset(tc.latitude, tc.longitude, 2024, time.June, 21)
			}
		})
	}
}

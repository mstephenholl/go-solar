package solar

import (
	"testing"
	"time"
)

// Sunrise is defined to be when the Sun is 50 arc minutes below the horizon.
// This is due to atmospheric refraction and measuring the position of the top
// rather than the center of the Sun.
// https://en.wikipedia.org/wiki/Sunrise#Angle
var sunriseElevation = -50.0 / 60.0

var dataElevation = []struct {
	inLatitude  float64
	inLongitude float64
	inElevation float64
	inYear      int
	inMonth     time.Month
	inDay       int
	outFirst    time.Time
	outSecond   time.Time
}{
	// 1970-01-01 - prime meridian
	{
		0, 0, sunriseElevation,
		1970, time.January, 1,
		time.Date(1970, time.January, 1, 5, 59, 54, 0, time.UTC),
		time.Date(1970, time.January, 1, 18, 0o7, 0o7, 0, time.UTC),
	},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{
		43.65, -79.38, sunriseElevation,
		2000, time.January, 1,
		time.Date(2000, time.January, 1, 12, 51, 0o0, 0, time.UTC),
		time.Date(2000, time.January, 1, 21, 50, 36, 0, time.UTC),
	},
	// 2004-04-01 - (52° N, 5° E)
	{
		52, 5, sunriseElevation,
		2004, time.April, 1,
		time.Date(2004, time.April, 1, 5, 13, 40, 0, time.UTC),
		time.Date(2004, time.April, 1, 18, 13, 27, 0, time.UTC),
	},
	// 2020-06-15 - Igloolik, Canada
	{
		69.3321443, -81.6781126, sunriseElevation,
		2020, time.June, 25,
		time.Time{},
		time.Time{},
	},
	// 2022-08-27 - London, end of Shabbat
	{
		51.5072, -0.1276, -8.5,
		2022, time.August, 27,
		time.Date(2022, time.August, 27, 4, 11, 31, 0, time.UTC),
		time.Date(2022, time.August, 27, 19, 53, 6, 0, time.UTC),
	},
	// 2022-06-21 - London, highest point around noon
	{
		51.5072, -0.1276, 61.93,
		2022, time.June, 21,
		time.Date(2022, time.June, 21, 12, 0, 12, 0, time.UTC),
		time.Date(2022, time.June, 21, 12, 4, 16, 0, time.UTC),
	},
	// 2022-06-21 - London, too high, never reached
	{
		51.5072, -0.1276, 61.94,
		2022, time.June, 21,
		time.Time{},
		time.Time{},
	},
	// 2022-06-21 - London, too low, never reached
	{
		51.5072, -0.1276, -16,
		2022, time.June, 21,
		time.Time{},
		time.Time{},
	},
	// 2022-11-26 - New York City, end of Shabbat
	{
		40.7128, -74.006, -8.5,
		2022, time.November, 26,
		time.Date(2022, time.November, 26, 11, 11, 19, 0, time.UTC),
		time.Date(2022, time.November, 26, 22, 15, 46, 0, time.UTC),
	},
}

func TestTimeOfElevation(t *testing.T) {
	for _, tt := range dataElevation {
		vFirst, vSecond := TimeOfElevation(tt.inLatitude, tt.inLongitude, tt.inElevation, tt.inYear, tt.inMonth, tt.inDay)
		if Abs(vFirst.Unix()-tt.outFirst.Unix()) > 2 {
			t.Fatalf("%s != %s", vFirst.String(), tt.outFirst.String())
		}
		if Abs(vSecond.Unix()-tt.outSecond.Unix()) > 2 {
			t.Fatalf("%s != %s", vSecond.String(), tt.outSecond.String())
		}
	}
}

func TestElevation(t *testing.T) {
	for _, tt := range dataElevation {
		if tt.outFirst.IsZero() || tt.outSecond.IsZero() {
			continue // Not reversible from output
		}

		vFirst := Elevation(tt.inLatitude, tt.inLongitude, tt.outFirst)
		if !AlmostEqual(vFirst, tt.inElevation, 2.0) {
			t.Fatalf("%f != %f", vFirst, tt.inElevation)
		}
		vSecond := Elevation(tt.inLatitude, tt.inLongitude, tt.outSecond)
		if !AlmostEqual(vSecond, tt.inElevation, 2.0) {
			t.Fatalf("%f != %f", vSecond, tt.inElevation)
		}
	}
}
// Benchmark for the Elevation function
func BenchmarkElevation(b *testing.B) {
	latitude := 40.7128
	longitude := -74.0060
	when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Elevation(latitude, longitude, when)
	}
}

// Benchmark for the TimeOfElevation function
func BenchmarkTimeOfElevation(b *testing.B) {
	latitude := 51.5072
	longitude := -0.1276
	elevation := -8.5
	year := 2024
	month := time.June
	day := 21

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TimeOfElevation(latitude, longitude, elevation, year, month, day)
	}
}

// Benchmark for different elevation angles
func BenchmarkTimeOfElevation_Angles(b *testing.B) {
	testCases := []struct {
		name      string
		elevation float64
	}{
		{"Sunrise", -50.0 / 60.0},
		{"CivilTwilight", -6.0},
		{"NauticalTwilight", -12.0},
		{"AstronomicalTwilight", -18.0},
		{"GoldenHour", 6.0},
		{"SolarNoon", 60.0},
	}

	latitude := 40.7128
	longitude := -74.0060

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = TimeOfElevation(latitude, longitude, tc.elevation, 2024, time.June, 21)
			}
		})
	}
}

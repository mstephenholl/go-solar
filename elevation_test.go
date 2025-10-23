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
		loc := NewLocation(tt.inLatitude, tt.inLongitude)
		tm := NewTime(tt.inYear, tt.inMonth, tt.inDay)
		vFirst, vSecond := TimeOfElevation(loc, tt.inElevation, tm)
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

		loc := NewLocation(tt.inLatitude, tt.inLongitude)
		vFirst := Elevation(loc, tt.outFirst)
		if !AlmostEqual(vFirst, tt.inElevation, 2.0) {
			t.Fatalf("%f != %f", vFirst, tt.inElevation)
		}
		vSecond := Elevation(loc, tt.outSecond)
		if !AlmostEqual(vSecond, tt.inElevation, 2.0) {
			t.Fatalf("%f != %f", vSecond, tt.inElevation)
		}
	}
}

// Benchmark for the Elevation function
func BenchmarkElevation(b *testing.B) {
	loc := NewLocation(40.7128, -74.0060)
	when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Elevation(loc, when)
	}
}

// Benchmark for the TimeOfElevation function
func BenchmarkTimeOfElevation(b *testing.B) {
	loc := NewLocation(51.5072, -0.1276)
	elevation := -8.5
	tm := NewTime(2024, time.June, 21)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TimeOfElevation(loc, elevation, tm)
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

	loc := NewLocation(40.7128, -74.0060)
	tm := NewTime(2024, time.June, 21)

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = TimeOfElevation(loc, tc.elevation, tm)
			}
		})
	}
}

// TestElevationFromNMEA_RMC tests ElevationFromNMEA with RMC sentences
func TestElevationFromNMEA_RMC(t *testing.T) {
	elevation, err := ElevationFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// At Toronto (43.65°N, 79.38°W) on March 23, 1994 at 12:35:19 UTC
	// UTC 12:35 is ~7:35 AM local (EST is UTC-5), so sun should be relatively low
	// Expected elevation around 10-20 degrees in the morning
	if elevation < 0 || elevation > 30 {
		t.Errorf("elevation %.2f degrees seems unreasonable for Toronto at 7:35 AM in March", elevation)
	}
}

// TestElevationFromNMEA_GGA tests ElevationFromNMEA with GGA sentences
func TestElevationFromNMEA_GGA(t *testing.T) {
	// GGA requires external date
	elevation, err := ElevationFromNMEA(validGGA, 1994, time.March, 23)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should match RMC result (same location and time)
	elevationRMC, _ := ElevationFromNMEA(validRMC, 0, 0, 0)
	if !AlmostEqual(elevation, elevationRMC, 0.1) {
		t.Errorf("GGA elevation %.2f != RMC elevation %.2f", elevation, elevationRMC)
	}
}

// TestElevationFromNMEA_GGA_MissingDate tests that GGA fails without date
func TestElevationFromNMEA_GGA_MissingDate(t *testing.T) {
	_, err := ElevationFromNMEA(validGGA, 0, 0, 0)
	if err == nil {
		t.Error("expected error for GGA sentence without date, got nil")
	}
}

// TestElevationFromNMEA_SouthernHemisphere tests southern hemisphere
func TestElevationFromNMEA_SouthernHemisphere(t *testing.T) {
	// Sydney at 12:00:00 UTC on June 21 (winter solstice)
	elevation, err := ElevationFromNMEA(validRMCSouth, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 12:00 UTC in Sydney (UTC+10) is 22:00 local time - nighttime
	// Sun should be well below the horizon (negative elevation)
	if elevation > 0 {
		t.Errorf("elevation %.2f degrees should be negative (nighttime) for Sydney at 22:00 in winter", elevation)
	}
}

// TestElevationFromNMEA_InvalidSentences tests error handling
func TestElevationFromNMEA_InvalidSentences(t *testing.T) {
	testCases := []struct {
		name  string
		nmea  string
		year  int
		month time.Month
		day   int
	}{
		{"InvalidChecksum", "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*FF", 0, 0, 0},
		{"MissingPrefix", "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71", 0, 0, 0},
		{"UnsupportedType", "$GPGSV,3,1,12,01,45,123,45,05,12,234,38*7F", 0, 0, 0},
		{"Empty", "", 0, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ElevationFromNMEA(tc.nmea, tc.year, tc.month, tc.day)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}
}

// TestTimeOfElevationFromNMEA_RMC tests TimeOfElevationFromNMEA with RMC
func TestTimeOfElevationFromNMEA_RMC(t *testing.T) {
	// Calculate civil twilight times for Toronto on March 23, 1994
	morning, evening, err := TimeOfElevationFromNMEA(validRMC, -6.0, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got valid times
	if morning.IsZero() || evening.IsZero() {
		t.Error("expected valid times for civil twilight, got zero times")
	}

	// Morning should be before evening
	if !morning.Before(evening) {
		t.Errorf("morning (%v) should be before evening (%v)", morning, evening)
	}

	// Morning should be on March 23, 1994 in UTC
	if morning.Year() != 1994 || morning.Month() != time.March || morning.Day() != 23 {
		t.Errorf("morning time has wrong date: %v", morning)
	}

	// Evening could be on March 23 or 24 depending on timezone and calculation
	// Just verify it's reasonable (within 2 days of morning)
	daysDiff := Abs(evening.Unix()-morning.Unix()) / 86400
	if daysDiff > 2 {
		t.Errorf("evening (%v) is too far from morning (%v)", evening, morning)
	}
}

// TestTimeOfElevationFromNMEA_GGA tests TimeOfElevationFromNMEA with GGA
func TestTimeOfElevationFromNMEA_GGA(t *testing.T) {
	// Calculate sunrise/sunset for Toronto
	morning, evening, err := TimeOfElevationFromNMEA(validGGA, sunriseElevation, 1994, time.March, 23)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got valid times
	if morning.IsZero() || evening.IsZero() {
		t.Error("expected valid times for sunrise/sunset, got zero times")
	}

	// Should match RMC result
	morningRMC, eveningRMC, _ := TimeOfElevationFromNMEA(validRMC, sunriseElevation, 0, 0, 0)
	if Abs(morning.Unix()-morningRMC.Unix()) > 2 {
		t.Errorf("GGA morning %v != RMC morning %v", morning, morningRMC)
	}
	if Abs(evening.Unix()-eveningRMC.Unix()) > 2 {
		t.Errorf("GGA evening %v != RMC evening %v", evening, eveningRMC)
	}
}

// TestTimeOfElevationFromNMEA_MultipleAngles tests different elevation angles
func TestTimeOfElevationFromNMEA_MultipleAngles(t *testing.T) {
	testCases := []struct {
		name      string
		elevation float64
	}{
		{"Sunrise", -50.0 / 60.0},
		{"CivilTwilight", -6.0},
		{"NauticalTwilight", -12.0},
		{"AstronomicalTwilight", -18.0},
		{"GoldenHour", 6.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			morning, evening, err := TimeOfElevationFromNMEA(validRMC, tc.elevation, 0, 0, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// All these angles should be reachable for Toronto in March
			if morning.IsZero() || evening.IsZero() {
				t.Errorf("expected valid times for %s, got zero times", tc.name)
			}

			// Morning should be before evening
			if !morning.Before(evening) {
				t.Errorf("%s: morning should be before evening", tc.name)
			}
		})
	}
}

// TestTimeOfElevationFromNMEA_UnreachableElevation tests unreachable elevations
func TestTimeOfElevationFromNMEA_UnreachableElevation(t *testing.T) {
	// Try to find when sun is at 70 degrees (unreachable at Toronto latitude in March)
	morning, evening, err := TimeOfElevationFromNMEA(validRMC, 70.0, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return zero times for unreachable elevation
	if !morning.IsZero() || !evening.IsZero() {
		t.Errorf("expected zero times for unreachable elevation, got morning=%v evening=%v", morning, evening)
	}
}

// TestTimeOfElevationFromNMEA_GGA_MissingDate tests GGA without date
func TestTimeOfElevationFromNMEA_GGA_MissingDate(t *testing.T) {
	_, _, err := TimeOfElevationFromNMEA(validGGA, -6.0, 0, 0, 0)
	if err == nil {
		t.Error("expected error for GGA sentence without date, got nil")
	}
}

// Benchmark for ElevationFromNMEA with RMC
func BenchmarkElevationFromNMEA_RMC(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ElevationFromNMEA(validRMC, 0, 0, 0)
	}
}

// Benchmark for ElevationFromNMEA with GGA
func BenchmarkElevationFromNMEA_GGA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ElevationFromNMEA(validGGA, 1994, time.March, 23)
	}
}

// Benchmark for TimeOfElevationFromNMEA with RMC
func BenchmarkTimeOfElevationFromNMEA_RMC(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = TimeOfElevationFromNMEA(validRMC, -6.0, 0, 0, 0)
	}
}

// Benchmark for TimeOfElevationFromNMEA with GGA
func BenchmarkTimeOfElevationFromNMEA_GGA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = TimeOfElevationFromNMEA(validGGA, -6.0, 1994, time.March, 23)
	}
}

// Benchmark TimeOfElevationFromNMEA with different elevation angles
func BenchmarkTimeOfElevationFromNMEA_Angles(b *testing.B) {
	testCases := []struct {
		name      string
		elevation float64
	}{
		{"Sunrise", -50.0 / 60.0},
		{"CivilTwilight", -6.0},
		{"NauticalTwilight", -12.0},
		{"AstronomicalTwilight", -18.0},
		{"GoldenHour", 6.0},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = TimeOfElevationFromNMEA(validRMC, tc.elevation, 0, 0, 0)
			}
		})
	}
}

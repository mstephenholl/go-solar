package solar

import (
	"testing"
	"time"
)

var dataTestMeanSolarNoon = []struct {
	inLongitude float64
	inYear      int
	inMonth     time.Month
	inDay       int
	out         float64
}{
	// 1970-01-01 - prime meridian
	{0, 1970, 1, 1, 2440588},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{-79.38, 2000, 1, 1, 2451545.2205},
	// 2004-04-01 - (52° N, 5° E)
	{5, 2004, 4, 1, 2453096.98611},
}

func TestMeanSolarNoon(t *testing.T) {
	for _, tt := range dataTestMeanSolarNoon {
		v := MeanSolarNoon(tt.inLongitude, tt.inYear, tt.inMonth, tt.inDay)
		if Round(v, DefaultPlaces) != Round(tt.out, DefaultPlaces) {
			t.Fatalf("%f != %f", v, tt.out)
		}
	}
}

// TestMeanSolarNoonFromNMEA_RMC tests MeanSolarNoonFromNMEA with RMC sentences
func TestMeanSolarNoonFromNMEA_RMC(t *testing.T) {
	// Using the standard validRMC from nmea_test.go
	// Toronto (43.65° N, 79.38° W) on March 23, 1994
	solarNoon, err := MeanSolarNoonFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got a valid time
	if solarNoon.IsZero() {
		t.Error("expected valid solar noon time, got zero time")
	}

	// Solar noon should be around midday UTC (adjusted for longitude)
	// Toronto is at -79.38°W, so solar noon is ~5.3 hours after UTC noon
	// Expected around 17:18 UTC
	if solarNoon.Hour() < 16 || solarNoon.Hour() > 18 {
		t.Errorf("solar noon hour %d seems unreasonable for Toronto (expected ~17)", solarNoon.Hour())
	}

	// Date should match the NMEA sentence (March 23, 1994)
	if solarNoon.Year() != 1994 || solarNoon.Month() != time.March || solarNoon.Day() != 23 {
		t.Errorf("solar noon has wrong date: %v (expected 1994-03-23)", solarNoon)
	}
}

// TestMeanSolarNoonFromNMEA_GGA tests MeanSolarNoonFromNMEA with GGA sentences
func TestMeanSolarNoonFromNMEA_GGA(t *testing.T) {
	// GGA requires external date
	solarNoon, err := MeanSolarNoonFromNMEA(validGGA, 1994, time.March, 23)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should match RMC result (same location and date)
	solarNoonRMC, _ := MeanSolarNoonFromNMEA(validRMC, 0, 0, 0)

	// Times should be very close (within a few seconds)
	if Abs(solarNoon.Unix()-solarNoonRMC.Unix()) > 5 {
		t.Errorf("GGA solar noon %v != RMC solar noon %v", solarNoon, solarNoonRMC)
	}
}

// TestMeanSolarNoonFromNMEA_GGA_MissingDate tests that GGA fails without date
func TestMeanSolarNoonFromNMEA_GGA_MissingDate(t *testing.T) {
	_, err := MeanSolarNoonFromNMEA(validGGA, 0, 0, 0)
	if err == nil {
		t.Error("expected error for GGA sentence without date, got nil")
	}
}

// TestMeanSolarNoonFromNMEA_SouthernHemisphere tests southern hemisphere
func TestMeanSolarNoonFromNMEA_SouthernHemisphere(t *testing.T) {
	// Sydney, Australia - RMC from June 21, 2021
	solarNoon, err := MeanSolarNoonFromNMEA(validRMCSouth, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got a valid time
	if solarNoon.IsZero() {
		t.Error("expected valid solar noon time, got zero time")
	}

	// Sydney is at 151.15°E, so solar noon is ~10 hours before UTC noon
	// Expected around 02:00 UTC
	if solarNoon.Hour() < 1 || solarNoon.Hour() > 4 {
		t.Errorf("solar noon hour %d seems unreasonable for Sydney (expected ~02)", solarNoon.Hour())
	}

	// Date should be June 21, 2021
	if solarNoon.Year() != 2021 || solarNoon.Month() != time.June || solarNoon.Day() != 21 {
		t.Errorf("solar noon has wrong date: %v (expected 2021-06-21)", solarNoon)
	}
}

// TestMeanSolarNoonFromNMEA_InvalidSentences tests error handling
func TestMeanSolarNoonFromNMEA_InvalidSentences(t *testing.T) {
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
			_, err := MeanSolarNoonFromNMEA(tc.nmea, tc.year, tc.month, tc.day)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}
}

// TestMeanSolarNoonFromNMEA_ConsistencyWithDirect tests consistency with direct calculation
func TestMeanSolarNoonFromNMEA_ConsistencyWithDirect(t *testing.T) {
	// Parse NMEA to get position and date
	solarNoonFromNMEA, err := MeanSolarNoonFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Calculate directly using known Toronto coordinates and date
	// Toronto: 43.65° N, 79.38° W, March 23, 1994
	jd := MeanSolarNoon(-79.38, 1994, time.March, 23)
	solarNoonDirect := JulianDayToTime(jd)

	// Should be very close (within 1 second due to rounding in NMEA coordinates)
	if Abs(solarNoonFromNMEA.Unix()-solarNoonDirect.Unix()) > 1 {
		t.Errorf("NMEA solar noon %v != direct calculation %v", solarNoonFromNMEA, solarNoonDirect)
	}
}

// BenchmarkMeanSolarNoonFromNMEA_RMC benchmarks MeanSolarNoonFromNMEA with RMC
func BenchmarkMeanSolarNoonFromNMEA_RMC(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MeanSolarNoonFromNMEA(validRMC, 0, 0, 0)
	}
}

// BenchmarkMeanSolarNoonFromNMEA_GGA benchmarks MeanSolarNoonFromNMEA with GGA
func BenchmarkMeanSolarNoonFromNMEA_GGA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MeanSolarNoonFromNMEA(validGGA, 1994, time.March, 23)
	}
}

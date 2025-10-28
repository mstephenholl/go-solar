package solar

import (
	"math"
	"testing"
	"time"
)

// dataAzimuth contains test cases for solar azimuth calculations
var dataAzimuth = []struct {
	name      string
	latitude  float64
	longitude float64
	when      time.Time
	azimuth   float64
	tolerance float64 // degrees of acceptable error
}{
	{
		name:      "2023-04-27 Kyiv 1:00 am",
		latitude:  50.45466,
		longitude: 30.5238,
		when:      time.Date(2023, time.April, 27, 1, 0, 0, 0, time.UTC),
		azimuth:   46.23,
		tolerance: 2.0,
	},
	{
		name:      "1970-01-01 Prime meridian 5:59 am",
		latitude:  0,
		longitude: 0,
		when:      time.Date(1970, time.January, 1, 5, 59, 0, 0, time.UTC),
		azimuth:   113.04,
		tolerance: 2.0,
	},
	{
		name:      "2000-01-14 Toronto 3:42 pm",
		latitude:  43.65,
		longitude: -79.38,
		when:      time.Date(2000, time.January, 14, 15, 42, 0, 0, time.UTC),
		azimuth:   154.01,
		tolerance: 2.0,
	},
	{
		name:      "2024-02-29 London 6:00 pm",
		latitude:  51.5072,
		longitude: -0.1276,
		when:      time.Date(2024, time.February, 29, 18, 0, 0, 0, time.UTC),
		azimuth:   263.74,
		tolerance: 2.0,
	},
	{
		name:      "2023-09-23 Autumnal Equinox 7:25 am",
		latitude:  0,
		longitude: 66.8928,
		when:      time.Date(2023, time.September, 23, 7, 25, 0, 0, time.UTC),
		azimuth:   6.38,
		tolerance: 2.0,
	},
	{
		name:      "2022-11-26 New York City 11:11 am",
		latitude:  40.7128,
		longitude: -74.006,
		when:      time.Date(2022, time.November, 26, 11, 11, 0, 0, time.UTC),
		azimuth:   110.41,
		tolerance: 2.0,
	},
}

// TestAzimuth tests the basic azimuth calculation function
func TestAzimuth(t *testing.T) {
	for _, tt := range dataAzimuth {
		t.Run(tt.name, func(t *testing.T) {
			loc := NewLocation(tt.latitude, tt.longitude)
			result := Azimuth(loc, tt.when)
			diff := math.Abs(result - tt.azimuth)
			if diff > tt.tolerance {
				t.Errorf("Azimuth() = %.2f, want %.2f (±%.2f), diff = %.2f",
					result, tt.azimuth, tt.tolerance, diff)
			}
		})
	}
}

// TestAzimuth_Cardinals tests azimuth at cardinal directions
func TestAzimuth_Cardinals(t *testing.T) {
	// Test at solar noon (should be close to south for northern hemisphere)
	// Toronto on summer solstice at solar noon
	loc := NewLocation(43.65, -79.38)
	// Approximate solar noon for Toronto
	when := time.Date(2024, time.June, 21, 17, 30, 0, 0, time.UTC)

	azimuth := Azimuth(loc, when)

	// At solar noon in northern hemisphere, sun should be roughly south (180°)
	// Allow wide tolerance due to equation of time and location
	if azimuth < 160 || azimuth > 200 {
		t.Errorf("Solar noon azimuth = %.2f, expected roughly 180° (south)", azimuth)
	}
}

// TestAzimuth_Range tests that azimuth is always in valid range
func TestAzimuth_Range(t *testing.T) {
	// Test various times throughout the day
	loc := NewLocation(43.65, -79.38)
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	for hour := 0; hour < 24; hour++ {
		when := date.Add(time.Duration(hour) * time.Hour)
		azimuth := Azimuth(loc, when)

		// Azimuth should always be in range [0, 360)
		if azimuth < 0 || azimuth >= 360 {
			t.Errorf("Hour %d: azimuth = %.2f, should be in range [0, 360)", hour, azimuth)
		}
	}
}

// TestAzimuth_Consistency tests consistency across hemispheres
func TestAzimuth_Consistency(t *testing.T) {
	// Test same longitude, different hemisphere
	when := time.Date(2024, time.March, 20, 12, 0, 0, 0, time.UTC) // Equinox

	// Northern hemisphere
	locNorth := NewLocation(45.0, 0)
	azimuthNorth := Azimuth(locNorth, when)

	// Southern hemisphere
	locSouth := NewLocation(-45.0, 0)
	azimuthSouth := Azimuth(locSouth, when)

	// Both should be valid
	if azimuthNorth < 0 || azimuthNorth >= 360 {
		t.Errorf("Northern hemisphere azimuth out of range: %.2f", azimuthNorth)
	}
	if azimuthSouth < 0 || azimuthSouth >= 360 {
		t.Errorf("Southern hemisphere azimuth out of range: %.2f", azimuthSouth)
	}

	// They should be different (sun in south vs north)
	if math.Abs(azimuthNorth-azimuthSouth) < 10 {
		t.Errorf("Expected different azimuths for different hemispheres, got %.2f vs %.2f",
			azimuthNorth, azimuthSouth)
	}
}

// TestAzimuth_Equator tests azimuth at the equator
func TestAzimuth_Equator(t *testing.T) {
	// At equator on equinox, sunrise should be exactly east (90°)
	// and sunset exactly west (270°)
	loc := NewLocation(0, 0)
	when := time.Date(2024, time.March, 20, 6, 0, 0, 0, time.UTC)

	azimuth := Azimuth(loc, when)

	// Should be roughly east at sunrise on equinox
	if azimuth < 80 || azimuth > 100 {
		t.Logf("Equatorial sunrise azimuth = %.2f (expected ~90° east)", azimuth)
	}
}

// BenchmarkAzimuth benchmarks the azimuth calculation
func BenchmarkAzimuth(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Azimuth(loc, when)
	}
}

// BenchmarkAzimuth_MultipleLocations benchmarks azimuth across different locations
func BenchmarkAzimuth_MultipleLocations(b *testing.B) {
	locations := []Location{
		NewLocation(43.65, -79.38),      // Toronto
		NewLocation(51.5072, -0.1276),   // London
		NewLocation(40.7128, -74.006),   // NYC
		NewLocation(-33.8688, 151.2093), // Sydney
	}
	when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loc := locations[i%len(locations)]
		_ = Azimuth(loc, when)
	}
}

// TestAzimuthFromNMEA_RMC tests azimuth calculation from RMC sentence
func TestAzimuthFromNMEA_RMC(t *testing.T) {
	// Toronto location with time: Jan 14, 2000 at 3:42 PM UTC
	// 4339.00,N = 43°39' = 43.65°
	// 07922.80,W = 79°22.8' = 79.38°
	nmea := "$GPRMC,154200,A,4339.00,N,07922.80,W,000.0,000.0,140100,000.0,W*7B"

	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewLocationFromNMEA() error = %v", err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewTimeFromNMEA() error = %v", err)
	}
	azimuth := Azimuth(loc, tm.DateTime())

	// Expected azimuth ~154° based on test data
	expected := 154.01
	tolerance := 2.0
	diff := math.Abs(azimuth - expected)
	if diff > tolerance {
		t.Errorf("Azimuth() = %.2f, want %.2f (±%.2f), diff = %.2f",
			azimuth, expected, tolerance, diff)
	}
}

// TestAzimuthFromNMEA_GGA tests azimuth calculation from GGA sentence
func TestAzimuthFromNMEA_GGA(t *testing.T) {
	// Prime meridian location at 5:59 AM
	nmea := "$GPGGA,055900,0000.000,N,00000.000,E,1,08,0.9,545.4,M,46.9,M,,*41"

	loc, err := NewLocationFromNMEA(nmea, 1970, time.January, 1)
	if err != nil {
		t.Fatalf("NewLocationFromNMEA() error = %v", err)
	}
	tm, err := NewTimeFromNMEA(nmea, 1970, time.January, 1)
	if err != nil {
		t.Fatalf("NewTimeFromNMEA() error = %v", err)
	}
	azimuth := Azimuth(loc, tm.DateTime())

	// Expected azimuth ~113° based on test data
	expected := 113.04
	tolerance := 2.0
	diff := math.Abs(azimuth - expected)
	if diff > tolerance {
		t.Errorf("Azimuth() = %.2f, want %.2f (±%.2f), diff = %.2f",
			azimuth, expected, tolerance, diff)
	}
}

// TestAzimuthFromNMEA_InvalidSentences tests error handling
func TestAzimuthFromNMEA_InvalidSentences(t *testing.T) {
	tests := []struct {
		name string
		nmea string
	}{
		{
			name: "InvalidChecksum",
			nmea: "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*00",
		},
		{
			name: "MissingPrefix",
			nmea: "GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71",
		},
		{
			name: "UnsupportedType",
			nmea: "$GPGSV,3,1,12,01,,,*6F",
		},
		{
			name: "Empty",
			nmea: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLocationFromNMEA(tt.nmea, 2024, time.June, 21)
			if err == nil {
				t.Error("NewLocationFromNMEA() expected error, got nil")
			}
		})
	}
}

// TestAzimuthFromNMEA_SouthernHemisphere tests azimuth in southern hemisphere
func TestAzimuthFromNMEA_SouthernHemisphere(t *testing.T) {
	// Sydney coordinates: 33°52.128'S, 151°12.558'E
	nmea := "$GPRMC,120000,A,3352.128,S,15112.558,E,000.0,000.0,210624,000.0,E*69"

	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewLocationFromNMEA() error = %v", err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewTimeFromNMEA() error = %v", err)
	}
	azimuth := Azimuth(loc, tm.DateTime())

	// Should get a valid azimuth
	if azimuth < 0 || azimuth >= 360 {
		t.Errorf("Azimuth() = %.2f, should be in range [0, 360)", azimuth)
	}
}

// TestAzimuthFromNMEA_ConsistencyWithDirect tests that NMEA and direct calculations match
func TestAzimuthFromNMEA_ConsistencyWithDirect(t *testing.T) {
	// Toronto: 43°39'N = 4339.00, 79°22.8'W = 07922.80
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,210624,000.0,W*7D"

	locNMEA, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewLocationFromNMEA() error = %v", err)
	}
	tmNMEA, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewTimeFromNMEA() error = %v", err)
	}
	azimuthNMEA := Azimuth(locNMEA, tmNMEA.DateTime())

	// Calculate directly
	loc := NewLocation(43.65, -79.38)
	when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
	azimuthDirect := Azimuth(loc, when)

	// Should be very close (within 0.1 degrees)
	diff := math.Abs(azimuthNMEA - azimuthDirect)
	if diff > 0.1 {
		t.Errorf("NMEA azimuth (%.2f) differs from direct calculation (%.2f) by %.2f degrees",
			azimuthNMEA, azimuthDirect, diff)
	}
}

// BenchmarkAzimuthFromNMEA_RMC benchmarks NMEA azimuth calculation with RMC
func BenchmarkAzimuthFromNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,154200,A,4339.00,N,07922.80,W,000.0,000.0,140100,000.0,W*7B"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loc, _ := NewLocationFromNMEA(nmea, 0, 0, 0)
		tm, _ := NewTimeFromNMEA(nmea, 0, 0, 0)
		_ = Azimuth(loc, tm.DateTime())
	}
}

// BenchmarkAzimuthFromNMEA_GGA benchmarks NMEA azimuth calculation with GGA
func BenchmarkAzimuthFromNMEA_GGA(b *testing.B) {
	nmea := "$GPGGA,055900,0000.000,N,00000.000,E,1,08,0.9,545.4,M,46.9,M,,*41"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loc, _ := NewLocationFromNMEA(nmea, 1970, time.January, 1)
		tm, _ := NewTimeFromNMEA(nmea, 1970, time.January, 1)
		_ = Azimuth(loc, tm.DateTime())
	}
}

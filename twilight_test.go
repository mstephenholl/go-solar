package solar

import (
	"math"
	"testing"
	"time"
)

// TestDawn tests the basic dawn calculation function
func TestDawn(t *testing.T) {
	// Toronto on June 21, 2024
	latitude := 43.65
	longitude := -79.38

	loc := NewLocation(latitude, longitude)
	tm := NewTime(2024, time.June, 21)

	// Test civil dawn (default)
	civilDawn := Dawn(loc, tm)
	if civilDawn.IsZero() {
		t.Error("Dawn() returned zero time for civil twilight")
	}

	// Verify it's before sunrise
	sunrise, err := Sunrise(loc, tm)
	if err != nil {
		t.Fatalf("Sunrise() error = %v", err)
	}
	if !civilDawn.Before(sunrise) {
		t.Errorf("Civil dawn (%v) should be before sunrise (%v)", civilDawn, sunrise)
	}
}

// TestDusk tests the basic dusk calculation function
func TestDusk(t *testing.T) {
	// Toronto on June 21, 2024
	latitude := 43.65
	longitude := -79.38

	loc := NewLocation(latitude, longitude)
	tm := NewTime(2024, time.June, 21)

	// Test civil dusk (default)
	civilDusk := Dusk(loc, tm)
	if civilDusk.IsZero() {
		t.Error("Dusk() returned zero time for civil twilight")
	}

	// Verify it's after sunset
	sunset, err := Sunset(loc, tm)
	if err != nil {
		t.Fatalf("Sunset() error = %v", err)
	}
	if !civilDusk.After(sunset) {
		t.Errorf("Civil dusk (%v) should be after sunset (%v)", civilDusk, sunset)
	}
}

// TestDawnDusk tests the combined dawn/dusk calculation
func TestDawnDusk(t *testing.T) {
	// Toronto on June 21, 2024
	latitude := 43.65
	longitude := -79.38

	loc := NewLocation(latitude, longitude)
	tm := NewTime(2024, time.June, 21)

	// Test civil twilight (default)
	dawn, dusk := DawnDusk(loc, tm)

	if dawn.IsZero() || dusk.IsZero() {
		t.Error("DawnDusk() returned zero times")
	}

	if !dawn.Before(dusk) {
		t.Errorf("Dawn (%v) should be before dusk (%v)", dawn, dusk)
	}

	// Verify consistency with individual functions
	separateDawn := Dawn(loc, tm)
	separateDusk := Dusk(loc, tm)

	if !dawn.Equal(separateDawn) {
		t.Errorf("DawnDusk dawn (%v) != Dawn() (%v)", dawn, separateDawn)
	}
	if !dusk.Equal(separateDusk) {
		t.Errorf("DawnDusk dusk (%v) != Dusk() (%v)", dusk, separateDusk)
	}
}

// TestTwilightTypes tests all three twilight types
func TestTwilightTypes(t *testing.T) {
	// Toronto on June 21, 2024
	latitude := 43.65
	longitude := -79.38
	year := 2024
	month := time.June
	day := 21

	loc := NewLocation(latitude, longitude)
	tm := NewTime(year, month, day)

	tests := []struct {
		name          string
		twilightType  TwilightType
		expectedAngle float64
	}{
		{"Civil", Civil, CivilTwilightAngle},
		{"Nautical", Nautical, NauticalTwilightAngle},
		{"Astronomical", Astronomical, AstronomicalTwilightAngle},
	}

	var previousDawn, previousDusk time.Time

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Dawn
			dawn := Dawn(loc, tm, tt.twilightType)
			if dawn.IsZero() {
				t.Errorf("Dawn(%s) returned zero time", tt.name)
			}

			// Test Dusk
			dusk := Dusk(loc, tm, tt.twilightType)
			if dusk.IsZero() {
				t.Errorf("Dusk(%s) returned zero time", tt.name)
			}

			// Test DawnDusk
			combinedDawn, combinedDusk := DawnDusk(loc, tm, tt.twilightType)
			if !combinedDawn.Equal(dawn) || !combinedDusk.Equal(dusk) {
				t.Errorf("DawnDusk(%s) inconsistent with separate calls", tt.name)
			}

			// Verify order: astronomical < nautical < civil
			if i > 0 {
				if !dawn.Before(previousDawn) {
					t.Errorf("%s dawn (%v) should be before %s dawn (%v)",
						tt.name, dawn, tests[i-1].name, previousDawn)
				}
				if !dusk.After(previousDusk) {
					t.Errorf("%s dusk (%v) should be after %s dusk (%v)",
						tt.name, dusk, tests[i-1].name, previousDusk)
				}
			}

			previousDawn = dawn
			previousDusk = dusk

			// Verify times match TimeOfElevation
			expectedDawn, expectedDusk := TimeOfElevation(loc, tt.expectedAngle, tm)
			if !dawn.Equal(expectedDawn) {
				t.Errorf("Dawn(%s) = %v, want %v", tt.name, dawn, expectedDawn)
			}
			if !dusk.Equal(expectedDusk) {
				t.Errorf("Dusk(%s) = %v, want %v", tt.name, dusk, expectedDusk)
			}
		})
	}
}

// TestDawnDusk_PolarRegions tests twilight in polar regions
func TestDawnDusk_PolarRegions(t *testing.T) {
	// Igloolik, Nunavut during summer
	latitude := 69.3321443
	longitude := -81.6781126

	loc := NewLocation(latitude, longitude)
	tm := NewTime(2020, time.June, 25)

	// During polar day/night, some twilight types may not occur
	dawn, dusk := DawnDusk(loc, tm)

	// Both should be zero (continuous daylight) or both should be valid
	if (dawn.IsZero() && !dusk.IsZero()) || (!dawn.IsZero() && dusk.IsZero()) {
		t.Errorf("Polar region: expected both zero or both valid, got dawn=%v, dusk=%v", dawn, dusk)
	}
}

// TestDawnFromNMEA_RMC tests dawn calculation from RMC sentence
func TestDawnFromNMEA_RMC(t *testing.T) {
	// Toronto location: Jan 1, 2000
	// 4339.00,N = 43°39' = 43.65°
	// 07922.80,W = 79°22.8' = 79.38°
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,010100,000.0,W*7E"

	dawn, err := DawnFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("DawnFromNMEA() error = %v", err)
	}

	if dawn.IsZero() {
		t.Error("DawnFromNMEA() returned zero time")
	}

	// Verify it matches direct calculation
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2000, time.January, 1)
	directDawn := Dawn(loc, tm)
	diff := dawn.Sub(directDawn)
	if math.Abs(diff.Seconds()) > 1 {
		t.Errorf("NMEA dawn differs from direct calculation by %v", diff)
	}
}

// TestDuskFromNMEA_GGA tests dusk calculation from GGA sentence
func TestDuskFromNMEA_GGA(t *testing.T) {
	// Prime meridian location
	nmea := "$GPGGA,120000,0000.000,N,00000.000,E,1,08,0.9,545.4,M,46.9,M,,*4B"

	dusk, err := DuskFromNMEA(nmea, 2024, time.June, 21)
	if err != nil {
		t.Fatalf("DuskFromNMEA() error = %v", err)
	}

	if dusk.IsZero() {
		t.Error("DuskFromNMEA() returned zero time")
	}

	// Verify it matches direct calculation
	loc := NewLocation(0, 0)
	tm := NewTime(2024, time.June, 21)
	directDusk := Dusk(loc, tm)
	diff := dusk.Sub(directDusk)
	if math.Abs(diff.Seconds()) > 1 {
		t.Errorf("NMEA dusk differs from direct calculation by %v", diff)
	}
}

// TestDawnDuskFromNMEA tests combined NMEA calculation
func TestDawnDuskFromNMEA(t *testing.T) {
	// Toronto location: Jan 1, 2000
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,010100,000.0,W*7E"

	dawn, dusk, err := DawnDuskFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("DawnDuskFromNMEA() error = %v", err)
	}

	if dawn.IsZero() || dusk.IsZero() {
		t.Error("DawnDuskFromNMEA() returned zero times")
	}

	if !dawn.Before(dusk) {
		t.Errorf("Dawn (%v) should be before dusk (%v)", dawn, dusk)
	}

	// Verify consistency with individual NMEA functions
	separateDawn, _ := DawnFromNMEA(nmea, 0, 0, 0)
	separateDusk, _ := DuskFromNMEA(nmea, 0, 0, 0)

	if !dawn.Equal(separateDawn) {
		t.Errorf("DawnDuskFromNMEA dawn (%v) != DawnFromNMEA() (%v)", dawn, separateDawn)
	}
	if !dusk.Equal(separateDusk) {
		t.Errorf("DawnDuskFromNMEA dusk (%v) != DuskFromNMEA() (%v)", dusk, separateDusk)
	}
}

// TestDawnDuskFromNMEA_TwilightTypes tests NMEA functions with different twilight types
func TestDawnDuskFromNMEA_TwilightTypes(t *testing.T) {
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,010100,000.0,W*7E"

	// Test all three types
	civilDawn, civilDusk, err := DawnDuskFromNMEA(nmea, 0, 0, 0, Civil)
	if err != nil {
		t.Fatalf("Civil twilight error: %v", err)
	}

	nauticalDawn, nauticalDusk, err := DawnDuskFromNMEA(nmea, 0, 0, 0, Nautical)
	if err != nil {
		t.Fatalf("Nautical twilight error: %v", err)
	}

	astronomicalDawn, astronomicalDusk, err := DawnDuskFromNMEA(nmea, 0, 0, 0, Astronomical)
	if err != nil {
		t.Fatalf("Astronomical twilight error: %v", err)
	}

	// Verify order: astronomical < nautical < civil
	if !astronomicalDawn.Before(nauticalDawn) {
		t.Error("Astronomical dawn should be before nautical dawn")
	}
	if !nauticalDawn.Before(civilDawn) {
		t.Error("Nautical dawn should be before civil dawn")
	}

	if !civilDusk.Before(nauticalDusk) {
		t.Error("Civil dusk should be before nautical dusk")
	}
	if !nauticalDusk.Before(astronomicalDusk) {
		t.Error("Nautical dusk should be before astronomical dusk")
	}
}

// TestDawnDuskFromNMEA_InvalidSentences tests error handling
func TestDawnDuskFromNMEA_InvalidSentences(t *testing.T) {
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
			_, err := DawnFromNMEA(tt.nmea, 2024, time.June, 21)
			if err == nil {
				t.Error("DawnFromNMEA() expected error, got nil")
			}

			_, err = DuskFromNMEA(tt.nmea, 2024, time.June, 21)
			if err == nil {
				t.Error("DuskFromNMEA() expected error, got nil")
			}

			_, _, err = DawnDuskFromNMEA(tt.nmea, 2024, time.June, 21)
			if err == nil {
				t.Error("DawnDuskFromNMEA() expected error, got nil")
			}
		})
	}
}

// BenchmarkDawn benchmarks the dawn calculation
func BenchmarkDawn(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Dawn(loc, tm)
	}
}

// BenchmarkDusk benchmarks the dusk calculation
func BenchmarkDusk(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Dusk(loc, tm)
	}
}

// BenchmarkDawnDusk benchmarks the combined calculation
func BenchmarkDawnDusk(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DawnDusk(loc, tm)
	}
}

// BenchmarkDawn_AllTypes benchmarks all twilight types
func BenchmarkDawn_AllTypes(b *testing.B) {
	loc := NewLocation(43.65, -79.38)
	tm := NewTime(2024, time.June, 21)

	types := []struct {
		name string
		tt   TwilightType
	}{
		{"Civil", Civil},
		{"Nautical", Nautical},
		{"Astronomical", Astronomical},
	}

	for _, tt := range types {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Dawn(loc, tm, tt.tt)
			}
		})
	}
}

// BenchmarkDawnFromNMEA_RMC benchmarks NMEA dawn calculation
func BenchmarkDawnFromNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,010100,000.0,W*7E"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DawnFromNMEA(nmea, 0, 0, 0)
	}
}

// BenchmarkDuskFromNMEA_GGA benchmarks NMEA dusk calculation
func BenchmarkDuskFromNMEA_GGA(b *testing.B) {
	nmea := "$GPGGA,120000,0000.000,N,00000.000,E,1,08,0.9,545.4,M,46.9,M,,*4B"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DuskFromNMEA(nmea, 2024, time.June, 21)
	}
}

// BenchmarkDawnDuskFromNMEA benchmarks combined NMEA calculation
func BenchmarkDawnDuskFromNMEA(b *testing.B) {
	nmea := "$GPRMC,120000,A,4339.00,N,07922.80,W,000.0,000.0,010100,000.0,W*7E"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DawnDuskFromNMEA(nmea, 0, 0, 0)
	}
}

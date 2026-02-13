package solar

import (
	"errors"
	"math"
	"testing"
	"time"
)

// Test data - valid NMEA sentences
const (
	// Toronto coordinates: 43.6532° N, 79.3832° W
	// Date: March 23, 1994 (equinox)
	validRMC = "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"

	// Same location, different format
	validGGA = "$GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,*5C"

	// Southern hemisphere
	validRMCSouth = "$GPRMC,120000,A,3351.000,S,15115.000,E,000.0,000.0,210621,000.0,E*6B"
)

func TestSunriseFromNMEA_RMC(t *testing.T) {
	// Test with valid RMC sentence
	loc, err := NewLocationFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunrise, err := Sunrise(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunrise.IsZero() {
		t.Error("expected valid sunrise time, got zero time")
	}

	// Verify the year is in reasonable range (1994 or nearby due to UTC conversion)
	if sunrise.Year() < 1993 || sunrise.Year() > 1995 {
		t.Errorf("expected year around 1994, got %v", sunrise)
	}
}

func TestSunsetFromNMEA_RMC(t *testing.T) {
	// Test with valid RMC sentence
	loc, err := NewLocationFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunset, err := Sunset(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunset.IsZero() {
		t.Error("expected valid sunset time, got zero time")
	}

	// Verify the year is in reasonable range
	if sunset.Year() < 1993 || sunset.Year() > 1995 {
		t.Errorf("expected year around 1994, got %v", sunset)
	}
}

func TestSunriseSunsetFromNMEA_RMC(t *testing.T) {
	loc, err := NewLocationFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(validRMC, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunrise, sunset, err := SunriseSunset(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunrise.IsZero() || sunset.IsZero() {
		t.Error("expected valid sunrise/sunset times")
	}

	// Sunrise should be before sunset
	if !sunrise.Before(sunset) {
		t.Errorf("sunrise (%v) should be before sunset (%v)", sunrise, sunset)
	}

	// Both should be on the same day
	if sunrise.Year() != sunset.Year() || sunrise.Month() != sunset.Month() || sunrise.Day() != sunset.Day() {
		t.Errorf("sunrise and sunset should be on same day: sunrise=%v, sunset=%v", sunrise, sunset)
	}
}

func TestSunriseFromNMEA_GGA(t *testing.T) {
	// GGA requires external date
	loc, err := NewLocationFromNMEA(validGGA, 1994, time.March, 23)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(validGGA, 1994, time.March, 23)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunrise, err := Sunrise(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunrise.IsZero() {
		t.Error("expected valid sunrise time, got zero time")
	}

	// Verify the date matches what we provided
	if sunrise.Year() != 1994 || sunrise.Month() != time.March || sunrise.Day() != 23 {
		t.Errorf("expected date 1994-03-23, got %v", sunrise)
	}
}

func TestSunriseFromNMEA_GGA_MissingDate(t *testing.T) {
	// GGA without date should fail
	_, err := NewLocationFromNMEA(validGGA, 0, 0, 0)
	if err == nil {
		t.Fatal("expected error for GGA without date")
	}

	if !errors.Is(err, ErrInvalidDate) {
		t.Errorf("expected ErrInvalidDate, got %v", err)
	}
}

func TestSunriseSunsetFromNMEA_SouthernHemisphere(t *testing.T) {
	// Sydney area - June 21 is winter solstice
	loc, err := NewLocationFromNMEA(validRMCSouth, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(validRMCSouth, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunrise, sunset, err := SunriseSunset(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunrise.IsZero() || sunset.IsZero() {
		t.Error("expected valid sunrise/sunset times")
	}

	// Verify year is in reasonable range (2021 or nearby due to UTC conversion)
	if sunrise.Year() < 2020 || sunrise.Year() > 2022 {
		t.Errorf("expected year around 2021, got %v", sunrise)
	}
}

func TestSunriseSunsetFromNMEA_DateLineCrossing(t *testing.T) {
	// International Date Line crossing - Equator at 180° longitude
	// Date: January 1, 2020 (00:30:45 UTC)
	nmea := "$GPRMC,003045,A,0000.000,N,18000.000,E,000.0,000.0,010120,000.0,E*7F"

	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating location: %v", err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error creating time: %v", err)
	}
	sunrise, sunset, err := SunriseSunset(loc, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sunrise.IsZero() || sunset.IsZero() {
		t.Error("expected valid sunrise/sunset times")
	}

	// At the equator on the International Date Line, sunrise and sunset should occur
	// Verify year is 2020 (or nearby due to UTC conversion)
	if sunrise.Year() < 2019 || sunrise.Year() > 2021 {
		t.Errorf("expected year around 2020, got %v", sunrise)
	}

	// At longitude 180°E (+180), local time is UTC+12
	// Local sunrise at ~6:00 becomes ~18:00 UTC (previous day)
	// Local sunset at ~18:00 becomes ~06:00 UTC (current day)
	// So in UTC, "sunrise" occurs in evening and "sunset" occurs in morning
	//
	// At the equator, day and night are roughly equal (~12 hours each)
	// Allow some variance for the equation of time and calculation precision
	if sunrise.Hour() < 16 || sunrise.Hour() > 20 {
		t.Errorf("expected sunrise around 18:00 UTC at 180°E, got %02d:%02d UTC", sunrise.Hour(), sunrise.Minute())
	}

	if sunset.Hour() < 4 || sunset.Hour() > 8 {
		t.Errorf("expected sunset around 06:00 UTC at 180°E, got %02d:%02d UTC", sunset.Hour(), sunset.Minute())
	}

	// Due to date line crossing and timezone offset:
	// - Sunrise: Dec 31 2019 ~18:00 UTC (local Jan 1 2020 ~06:00)
	// - Sunset:  Jan 1  2020 ~06:00 UTC (local Jan 1 2020 ~18:00)
	// So chronologically, sunrise comes before sunset (correct)
	if !sunrise.Before(sunset) {
		t.Errorf("sunrise (%v) should be before sunset (%v)", sunrise, sunset)
	}

	// Verify they span across the calendar day boundary in UTC
	if sunrise.Day() == sunset.Day() {
		t.Errorf("expected sunrise and sunset to be on different UTC days at date line, both are on day %d", sunrise.Day())
	}
}

func TestParseNMEA_InvalidSentences(t *testing.T) {
	testCases := []struct {
		name    string
		nmea    string
		wantErr error
		year    int
		month   time.Month
		day     int
	}{
		{
			name:    "missing dollar sign",
			nmea:    "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71",
			wantErr: ErrInvalidNMEA,
		},
		{
			name:    "missing checksum",
			nmea:    "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W",
			wantErr: ErrInvalidNMEA,
		},
		{
			name:    "invalid checksum",
			nmea:    "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*FF",
			wantErr: ErrInvalidChecksum,
		},
		{
			name:    "unsupported sentence type",
			nmea:    "$GPGSV,3,1,12,01,,,42,02,,,44,03,,,42,04,,,43*7B",
			wantErr: ErrUnsupportedSentence,
		},
		{
			name:    "RMC with invalid status",
			nmea:    "$GPRMC,123519,V,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*66",
			wantErr: ErrInvalidNMEA,
		},
		{
			name:    "too few fields",
			nmea:    "$GPRMC,123519*6A",
			wantErr: ErrInvalidNMEA,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewLocationFromNMEA(tc.nmea, tc.year, tc.month, tc.day)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestParseLatitude(t *testing.T) {
	testCases := []struct {
		name    string
		latStr  string
		nsStr   string
		want    float64
		wantErr bool
	}{
		{
			name:   "Toronto North",
			latStr: "4339.192",
			nsStr:  "N",
			want:   43.6532,
		},
		{
			name:   "Sydney South",
			latStr: "3351.000",
			nsStr:  "S",
			want:   -33.85,
		},
		{
			name:   "Equator",
			latStr: "0000.000",
			nsStr:  "N",
			want:   0.0,
		},
		{
			name:    "Invalid N/S",
			latStr:  "4339.192",
			nsStr:   "X",
			wantErr: true,
		},
		{
			name:    "Empty fields",
			latStr:  "",
			nsStr:   "N",
			wantErr: true,
		},
		{
			name:    "Invalid format",
			latStr:  "43.39",
			nsStr:   "N",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLatitude(tc.latStr, tc.nsStr)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tc.want) > 0.0001 {
				t.Errorf("expected %f, got %f", tc.want, got)
			}
		})
	}
}

func TestParseLongitude(t *testing.T) {
	testCases := []struct {
		name    string
		lonStr  string
		ewStr   string
		want    float64
		wantErr bool
	}{
		{
			name:   "Toronto West",
			lonStr: "07922.992",
			ewStr:  "W",
			want:   -79.3832,
		},
		{
			name:   "Sydney East",
			lonStr: "15115.000",
			ewStr:  "E",
			want:   151.25,
		},
		{
			name:   "Prime Meridian",
			lonStr: "00000.000",
			ewStr:  "E",
			want:   0.0,
		},
		{
			name:   "International Date Line",
			lonStr: "18000.000",
			ewStr:  "E",
			want:   180.0,
		},
		{
			name:    "Invalid E/W",
			lonStr:  "07922.992",
			ewStr:   "X",
			wantErr: true,
		},
		{
			name:    "Empty fields",
			lonStr:  "",
			ewStr:   "E",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLongitude(tc.lonStr, tc.ewStr)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tc.want) > 0.0001 {
				t.Errorf("expected %f, got %f", tc.want, got)
			}
		})
	}
}

func TestParseNMEATime(t *testing.T) {
	testCases := []struct {
		name     string
		timeStr  string
		year     int
		month    time.Month
		day      int
		wantHour int
		wantMin  int
		wantSec  int
		wantErr  bool
	}{
		{
			name:     "simple time",
			timeStr:  "123519",
			year:     2024,
			month:    time.June,
			day:      21,
			wantHour: 12,
			wantMin:  35,
			wantSec:  19,
		},
		{
			name:     "time with fractional seconds",
			timeStr:  "123519.50",
			year:     2024,
			month:    time.June,
			day:      21,
			wantHour: 12,
			wantMin:  35,
			wantSec:  19,
		},
		{
			name:     "midnight",
			timeStr:  "000000",
			year:     2024,
			month:    time.June,
			day:      21,
			wantHour: 0,
			wantMin:  0,
			wantSec:  0,
		},
		{
			name:    "too short",
			timeStr: "1235",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "invalid hour",
			timeStr: "XX3519",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseNMEATime(tc.timeStr, tc.year, tc.month, tc.day)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Hour() != tc.wantHour || got.Minute() != tc.wantMin || got.Second() != tc.wantSec {
				t.Errorf("expected %02d:%02d:%02d, got %02d:%02d:%02d",
					tc.wantHour, tc.wantMin, tc.wantSec,
					got.Hour(), got.Minute(), got.Second())
			}

			if got.Year() != tc.year || got.Month() != tc.month || got.Day() != tc.day {
				t.Errorf("expected date %d-%02d-%02d, got %d-%02d-%02d",
					tc.year, tc.month, tc.day,
					got.Year(), got.Month(), got.Day())
			}
		})
	}
}

func TestValidateChecksum(t *testing.T) {
	testCases := []struct {
		name        string
		sentence    string
		checksumStr string
		wantErr     bool
	}{
		{
			name:        "valid RMC checksum",
			sentence:    "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W",
			checksumStr: "71",
			wantErr:     false,
		},
		{
			name:        "valid GGA checksum",
			sentence:    "GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,",
			checksumStr: "5C",
			wantErr:     false,
		},
		{
			name:        "invalid checksum",
			sentence:    "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W",
			checksumStr: "FF",
			wantErr:     true,
		},
		{
			name:        "invalid checksum format",
			sentence:    "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W",
			checksumStr: "XY",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateChecksum(tc.sentence, tc.checksumStr)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Example tests for documentation
func ExampleNewLocationFromNMEA() {
	// Using an RMC sentence (includes date: March 23, 1994)
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"
	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	sunrise, err := Sunrise(loc, tm)
	if err != nil {
		panic(err)
	}

	// Sunrise is returned in UTC
	_ = sunrise // Use the sunrise time
}

func ExampleNewLocationFromNMEA_gga() {
	// Using a GGA sentence (requires external date)
	nmea := "$GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,*5C"
	loc, err := NewLocationFromNMEA(nmea, 2024, time.June, 21)
	if err != nil {
		panic(err)
	}
	tm, err := NewTimeFromNMEA(nmea, 2024, time.June, 21)
	if err != nil {
		panic(err)
	}
	sunrise, err := Sunrise(loc, tm)
	if err != nil {
		panic(err)
	}

	_ = sunrise // Use the sunrise time
}

func ExampleSunset() {
	// Using an RMC sentence
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"
	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	sunset, err := Sunset(loc, tm)
	if err != nil {
		panic(err)
	}

	_ = sunset // Use the sunset time
}

func ExampleSunriseSunset() {
	// Using an RMC sentence to get both sunrise and sunset
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"
	loc, err := NewLocationFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	tm, err := NewTimeFromNMEA(nmea, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	sunrise, sunset, err := SunriseSunset(loc, tm)
	if err != nil {
		panic(err)
	}

	_, _ = sunrise, sunset // Use the times
}

// Benchmarks

func BenchmarkSunriseFromNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"

	b.ResetTimer()
	for b.Loop() {
		loc, _ := NewLocationFromNMEA(nmea, 0, 0, 0)
		tm, _ := NewTimeFromNMEA(nmea, 0, 0, 0)
		_, _ = Sunrise(loc, tm)
	}
}

func BenchmarkSunsetFromNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"

	b.ResetTimer()
	for b.Loop() {
		loc, _ := NewLocationFromNMEA(nmea, 0, 0, 0)
		tm, _ := NewTimeFromNMEA(nmea, 0, 0, 0)
		_, _ = Sunset(loc, tm)
	}
}

func BenchmarkSunriseSunsetFromNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"

	b.ResetTimer()
	for b.Loop() {
		loc, _ := NewLocationFromNMEA(nmea, 0, 0, 0)
		tm, _ := NewTimeFromNMEA(nmea, 0, 0, 0)
		_, _, _ = SunriseSunset(loc, tm)
	}
}

func BenchmarkSunriseFromNMEA_GGA(b *testing.B) {
	nmea := "$GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,*5C"

	b.ResetTimer()
	for b.Loop() {
		loc, _ := NewLocationFromNMEA(nmea, 1994, time.March, 23)
		tm, _ := NewTimeFromNMEA(nmea, 1994, time.March, 23)
		_, _ = Sunrise(loc, tm)
	}
}

func BenchmarkParseNMEA_RMC(b *testing.B) {
	nmea := "$GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W*71"

	b.ResetTimer()
	for b.Loop() {
		_, _ = parseNMEA(nmea, 0, 0, 0)
	}
}

func BenchmarkParseNMEA_GGA(b *testing.B) {
	nmea := "$GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,*5C"

	b.ResetTimer()
	for b.Loop() {
		_, _ = parseNMEA(nmea, 1994, time.March, 23)
	}
}

func BenchmarkParseLatitude(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_, _ = parseLatitude("4339.192", "N")
	}
}

func BenchmarkParseLongitude(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_, _ = parseLongitude("07922.992", "W")
	}
}

func BenchmarkValidateChecksum(b *testing.B) {
	sentence := "GPRMC,123519,A,4339.192,N,07922.992,W,022.4,084.4,230394,003.1,W"

	b.ResetTimer()
	for b.Loop() {
		_ = validateChecksum(sentence, "71")
	}
}

// Additional coverage tests for edge cases

func TestSunsetFromNMEA_Errors(t *testing.T) {
	// Invalid NMEA should return error
	_, err := NewLocationFromNMEA("invalid", 2024, time.June, 21)
	if err == nil {
		t.Error("Expected error for invalid NMEA, got nil")
	}
}

func TestSunriseSunsetFromNMEA_Errors(t *testing.T) {
	// Invalid NMEA should return error
	_, err := NewLocationFromNMEA("invalid", 2024, time.June, 21)
	if err == nil {
		t.Error("Expected error for invalid NMEA, got nil")
	}
}

func TestParseNMEA_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		nmea    string
		year    int
		month   time.Month
		day     int
		wantErr bool
	}{
		{
			name:    "Sentence type too short",
			nmea:    "$AB*00",
			wantErr: true,
		},
		{
			name:    "Empty sentence with checksum",
			nmea:    "$*00",
			wantErr: true,
		},
		{
			name:    "Whitespace around NMEA (should succeed)",
			nmea:    "  $GPGGA,123519,4339.192,N,07922.992,W,1,08,0.9,545.4,M,46.9,M,,*5C  ",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseNMEA(tt.nmea, tt.year, tt.month, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNMEA() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseGGA_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		fields  []string
		year    int
		month   time.Month
		day     int
		wantErr bool
	}{
		{
			name:    "Too few fields",
			fields:  []string{"GPGGA", "123519"},
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Missing year",
			fields:  []string{"GPGGA", "123519", "4339.192", "N", "07922.992", "W", "1"},
			year:    0,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Missing month",
			fields:  []string{"GPGGA", "123519", "4339.192", "N", "07922.992", "W", "1"},
			year:    2024,
			month:   0,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Missing day",
			fields:  []string{"GPGGA", "123519", "4339.192", "N", "07922.992", "W", "1"},
			year:    2024,
			month:   time.June,
			day:     0,
			wantErr: true,
		},
		{
			name:    "Invalid time format",
			fields:  []string{"GPGGA", "99", "4339.192", "N", "07922.992", "W", "1"},
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Invalid latitude format",
			fields:  []string{"GPGGA", "123519", "invalid", "N", "07922.992", "W", "1"},
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Invalid longitude format",
			fields:  []string{"GPGGA", "123519", "4339.192", "N", "invalid", "W", "1"},
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseGGA(tt.fields, tt.year, tt.month, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGGA() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRMC_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		fields  []string
		wantErr bool
		expYear int // Expected year for successful parses
	}{
		{
			name:    "Too few fields",
			fields:  []string{"GPRMC", "123519"},
			wantErr: true,
		},
		{
			name:    "Date string too short",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "2303"},
			wantErr: true,
		},
		{
			name:    "Invalid day in date",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "XX0394"},
			wantErr: true,
		},
		{
			name:    "Invalid month in date",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "23XX94"},
			wantErr: true,
		},
		{
			name:    "Invalid year in date",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "2303XX"},
			wantErr: true,
		},
		{
			name:    "Year >= 50 converts to 19xx",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "230395"},
			wantErr: false,
			expYear: 1995,
		},
		{
			name:    "Year < 50 converts to 20xx",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "230325"},
			wantErr: false,
			expYear: 2025,
		},
		{
			name:    "Invalid time string",
			fields:  []string{"GPRMC", "99", "A", "4339.192", "N", "07922.992", "W", "022.4", "084.4", "230394"},
			wantErr: true,
		},
		{
			name:    "Invalid latitude",
			fields:  []string{"GPRMC", "123519", "A", "invalid", "N", "07922.992", "W", "022.4", "084.4", "230394"},
			wantErr: true,
		},
		{
			name:    "Invalid longitude",
			fields:  []string{"GPRMC", "123519", "A", "4339.192", "N", "invalid", "W", "022.4", "084.4", "230394"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRMC(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRMC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.expYear > 0 {
				if result.Time.Year() != tt.expYear {
					t.Errorf("Expected year %d, got %d", tt.expYear, result.Time.Year())
				}
			}
		})
	}
}

func TestParseNMEATime_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		timeStr   string
		year      int
		month     time.Month
		day       int
		wantErr   bool
		wantNano  int // Expected nanoseconds for successful parses
		checkNano bool
	}{
		{
			name:    "Invalid minute",
			timeStr: "12XX19",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:    "Invalid second",
			timeStr: "1235XX",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
		{
			name:      "Fractional seconds - short (should pad)",
			timeStr:   "123519.12",
			year:      2024,
			month:     time.June,
			day:       21,
			wantErr:   false,
			wantNano:  120000000,
			checkNano: true,
		},
		{
			name:      "Fractional seconds - long (should truncate)",
			timeStr:   "123519.1234567890123",
			year:      2024,
			month:     time.June,
			day:       21,
			wantErr:   false,
			wantNano:  123456789,
			checkNano: true,
		},
		{
			name:    "Invalid fractional seconds",
			timeStr: "123519.XXX",
			year:    2024,
			month:   time.June,
			day:     21,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNMEATime(tt.timeStr, tt.year, tt.month, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNMEATime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkNano {
				if result.Nanosecond() != tt.wantNano {
					t.Errorf("Expected nanosecond = %d, got %d", tt.wantNano, result.Nanosecond())
				}
			}
		})
	}
}

func TestParseLatitude_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		latStr  string
		nsStr   string
		wantErr bool
	}{
		{
			name:    "Empty N/S indicator",
			latStr:  "4339.192",
			nsStr:   "",
			wantErr: true,
		},
		{
			name:    "Decimal too early",
			latStr:  "4.39192",
			nsStr:   "N",
			wantErr: true,
		},
		{
			name:    "Invalid degrees",
			latStr:  "XX39.192",
			nsStr:   "N",
			wantErr: true,
		},
		{
			name:    "Invalid minutes",
			latStr:  "43XX.192",
			nsStr:   "N",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseLatitude(tt.latStr, tt.nsStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLatitude() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseLongitude_AdditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		lonStr  string
		ewStr   string
		wantErr bool
	}{
		{
			name:    "Empty E/W indicator",
			lonStr:  "07922.992",
			ewStr:   "",
			wantErr: true,
		},
		{
			name:    "Decimal too early",
			lonStr:  "7.922992",
			ewStr:   "W",
			wantErr: true,
		},
		{
			name:    "Invalid degrees",
			lonStr:  "XXX22.992",
			ewStr:   "W",
			wantErr: true,
		},
		{
			name:    "Invalid minutes",
			lonStr:  "079XX.992",
			ewStr:   "W",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseLongitude(tt.lonStr, tt.ewStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLongitude() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

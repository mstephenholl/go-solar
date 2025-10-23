package solar

import (
	"testing"
	"time"
)

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{"Toronto", 43.65, -79.38},
		{"New York", 40.7128, -74.0060},
		{"London", 51.5074, -0.1278},
		{"Tokyo", 35.6762, 139.6503},
		{"Sydney", -33.8688, 151.2093},
		{"North Pole", 90.0, 0.0},
		{"South Pole", -90.0, 0.0},
		{"Equator Prime Meridian", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := NewLocation(tt.latitude, tt.longitude)

			if loc.Latitude() != tt.latitude {
				t.Errorf("Latitude() = %v, want %v", loc.Latitude(), tt.latitude)
			}
			if loc.Longitude() != tt.longitude {
				t.Errorf("Longitude() = %v, want %v", loc.Longitude(), tt.longitude)
			}
		})
	}
}

func TestNewLocationFromNMEA(t *testing.T) {
	tests := []struct {
		name      string
		nmea      string
		year      int
		month     time.Month
		day       int
		wantLat   float64
		wantLon   float64
		wantErr   bool
		tolerance float64
	}{
		{
			name:      "Valid GGA sentence",
			nmea:      "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47",
			year:      2025,
			month:     time.January,
			day:       15,
			wantLat:   48.1173,
			wantLon:   11.5167,
			wantErr:   false,
			tolerance: 0.0001,
		},
		{
			name:      "Valid RMC sentence",
			nmea:      "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A",
			year:      0, // Ignored for RMC
			month:     0,
			day:       0,
			wantLat:   48.1173,
			wantLon:   11.5167,
			wantErr:   false,
			tolerance: 0.0001,
		},
		{
			name:      "Toronto coordinates in GGA",
			nmea:      "$GPGGA,154200,4339.00,N,07922.80,W,1,08,0.9,545.4,M,46.9,M,,*53",
			year:      2025,
			month:     time.October,
			day:       21,
			wantLat:   43.65,
			wantLon:   -79.38,
			wantErr:   false,
			tolerance: 0.01,
		},
		{
			name:      "Southern hemisphere in RMC",
			nmea:      "$GPRMC,120000,A,3352.128,S,15112.558,E,000.0,000.0,210624,000.0,E*69",
			year:      0,
			month:     0,
			day:       0,
			wantLat:   -33.8688,
			wantLon:   151.2093,
			wantErr:   false,
			tolerance: 0.001,
		},
		{
			name:    "Invalid NMEA sentence",
			nmea:    "not a valid nmea sentence",
			year:    2025,
			month:   time.January,
			day:     15,
			wantErr: true,
		},
		{
			name:    "Empty NMEA sentence",
			nmea:    "",
			year:    2025,
			month:   time.January,
			day:     15,
			wantErr: true,
		},
		{
			name:    "Invalid checksum",
			nmea:    "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*FF",
			year:    2025,
			month:   time.January,
			day:     15,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := NewLocationFromNMEA(tt.nmea, tt.year, tt.month, tt.day)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewLocationFromNMEA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !AlmostEqual(loc.Latitude(), tt.wantLat, tt.tolerance) {
				t.Errorf("Latitude() = %v, want %v (tolerance %v)", loc.Latitude(), tt.wantLat, tt.tolerance)
			}
			if !AlmostEqual(loc.Longitude(), tt.wantLon, tt.tolerance) {
				t.Errorf("Longitude() = %v, want %v (tolerance %v)", loc.Longitude(), tt.wantLon, tt.tolerance)
			}
		})
	}
}

func TestLocationString(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		want      string
	}{
		{"Toronto", 43.65, -79.38, "43.6500°N, 79.3800°W"},
		{"Sydney", -33.8688, 151.2093, "33.8688°S, 151.2093°E"},
		{"Prime Meridian", 0, 0, "0.0000°N, 0.0000°E"},
		{"North Pole", 90, 0, "90.0000°N, 0.0000°E"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := NewLocation(tt.latitude, tt.longitude)
			got := loc.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTime(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
	}{
		{"New Year 2025", 2025, time.January, 1},
		{"Summer Solstice 2025", 2025, time.June, 21},
		{"Winter Solstice 2025", 2025, time.December, 21},
		{"Leap Day 2024", 2024, time.February, 29},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTime(tt.year, tt.month, tt.day)

			if tm.Year() != tt.year {
				t.Errorf("Year() = %v, want %v", tm.Year(), tt.year)
			}
			if tm.Month() != tt.month {
				t.Errorf("Month() = %v, want %v", tm.Month(), tt.month)
			}
			if tm.Day() != tt.day {
				t.Errorf("Day() = %v, want %v", tm.Day(), tt.day)
			}

			// Verify time is UTC and at midnight
			dt := tm.DateTime()
			if dt.Location() != time.UTC {
				t.Errorf("DateTime().Location() = %v, want UTC", dt.Location())
			}
			if dt.Hour() != 0 || dt.Minute() != 0 || dt.Second() != 0 {
				t.Errorf("DateTime() time = %02d:%02d:%02d, want 00:00:00", dt.Hour(), dt.Minute(), dt.Second())
			}
		})
	}
}

func TestNewTimeFromDateTime(t *testing.T) {
	tests := []struct {
		name string
		when time.Time
	}{
		{
			"UTC time",
			time.Date(2025, time.October, 21, 14, 30, 0, 0, time.UTC),
		},
		{
			"Local time (should convert to UTC)",
			time.Date(2025, time.October, 21, 9, 30, 0, 0, time.FixedZone("EST", -5*60*60)),
		},
		{
			"Now",
			time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTimeFromDateTime(tt.when)

			// Verify result is in UTC
			dt := tm.DateTime()
			if dt.Location() != time.UTC {
				t.Errorf("DateTime().Location() = %v, want UTC", dt.Location())
			}

			// Verify time matches (in UTC)
			expected := tt.when.UTC()
			if !dt.Equal(expected) {
				t.Errorf("DateTime() = %v, want %v", dt, expected)
			}
		})
	}
}

func TestNewTimeFromNMEA(t *testing.T) {
	tests := []struct {
		name     string
		nmea     string
		year     int
		month    time.Month
		day      int
		wantHour int
		wantMin  int
		wantSec  int
		wantYear int
		wantMon  time.Month
		wantDay  int
		wantErr  bool
	}{
		{
			name:     "GGA sentence with time",
			nmea:     "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47",
			year:     2025,
			month:    time.January,
			day:      15,
			wantHour: 12,
			wantMin:  35,
			wantSec:  19,
			wantYear: 2025,
			wantMon:  time.January,
			wantDay:  15,
			wantErr:  false,
		},
		{
			name:     "RMC sentence with date (ignores provided date)",
			nmea:     "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A",
			year:     9999, // Should be ignored
			month:    time.December,
			day:      31,
			wantHour: 12,
			wantMin:  35,
			wantSec:  19,
			wantYear: 1994, // From RMC: 230394 = March 23, 1994
			wantMon:  time.March,
			wantDay:  23,
			wantErr:  false,
		},
		{
			name:    "Invalid NMEA sentence",
			nmea:    "not valid",
			year:    2025,
			month:   time.January,
			day:     15,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := NewTimeFromNMEA(tt.nmea, tt.year, tt.month, tt.day)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTimeFromNMEA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			dt := tm.DateTime()
			if dt.Hour() != tt.wantHour {
				t.Errorf("Hour() = %v, want %v", dt.Hour(), tt.wantHour)
			}
			if dt.Minute() != tt.wantMin {
				t.Errorf("Minute() = %v, want %v", dt.Minute(), tt.wantMin)
			}
			if dt.Second() != tt.wantSec {
				t.Errorf("Second() = %v, want %v", dt.Second(), tt.wantSec)
			}
			if tm.Year() != tt.wantYear {
				t.Errorf("Year() = %v, want %v", tm.Year(), tt.wantYear)
			}
			if tm.Month() != tt.wantMon {
				t.Errorf("Month() = %v, want %v", tm.Month(), tt.wantMon)
			}
			if tm.Day() != tt.wantDay {
				t.Errorf("Day() = %v, want %v", tm.Day(), tt.wantDay)
			}
		})
	}
}

func TestTimeString(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
		want  string
	}{
		{"2025-01-01", 2025, time.January, 1, "2025-01-01"},
		{"2025-12-31", 2025, time.December, 31, "2025-12-31"},
		{"2024-02-29", 2024, time.February, 29, "2024-02-29"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTime(tt.year, tt.month, tt.day)
			got := tm.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewLocation(43.65, -79.38)
	}
}

func BenchmarkNewLocationFromNMEA(b *testing.B) {
	nmea := "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47"
	for i := 0; i < b.N; i++ {
		_, _ = NewLocationFromNMEA(nmea, 2025, time.January, 15)
	}
}

func BenchmarkNewTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewTime(2025, time.October, 21)
	}
}

func BenchmarkNewTimeFromDateTime(b *testing.B) {
	when := time.Date(2025, time.October, 21, 14, 30, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		_ = NewTimeFromDateTime(when)
	}
}

func BenchmarkNewTimeFromNMEA(b *testing.B) {
	nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A"
	for i := 0; i < b.N; i++ {
		_, _ = NewTimeFromNMEA(nmea, 0, 0, 0)
	}
}

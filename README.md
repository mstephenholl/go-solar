# go-solar ☀️

[![CI](https://github.com/mstephenholl/go-solar/actions/workflows/ci.yml/badge.svg)](https://github.com/mstephenholl/go-solar/actions/workflows/ci.yml)
[![CodeQL](https://github.com/mstephenholl/go-solar/actions/workflows/codeql.yml/badge.svg)](https://github.com/mstephenholl/go-solar/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mstephenholl/go-solar)](https://goreportcard.com/report/github.com/mstephenholl/go-solar)
[![GoDoc](https://pkg.go.dev/badge/github.com/mstephenholl/go-solar)](https://pkg.go.dev/github.com/mstephenholl/go-solar)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mstephenholl/go-solar)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/mstephenholl/go-solar)](https://github.com/mstephenholl/go-solar/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.txt)

A modern, well-tested Go package for calculating sunrise, sunset, and solar elevation at any location on Earth. Uses the [sunrise equation](https://en.wikipedia.org/wiki/Sunrise_equation#Complete_calculation_on_Earth) method with support for edge cases like polar nights and midnight sun.

## ✨ Features

- 🌅 Calculate sunrise and sunset times for any location
- 🌄 Calculate dawn and dusk with civil, nautical, and astronomical twilight
- 📐 Determine solar elevation and azimuth angles
- 🧭 Calculate solar azimuth (compass direction of the sun)
- 🛰️ Parse NMEA GPS sentences (GGA, RMC) for location-based calculations
- 🌍 Handle edge cases (polar night, midnight sun)
- 🚀 High performance with zero allocations for core functions
- ✅ 94%+ test coverage on production code
- 🔧 Generic helper functions (Go 1.18+)
- 📚 Comprehensive documentation and examples

Forked from Nathan Osman's package [`go-sunrise`](https://github.com/nathan-osman/go-sunrise) with modern enhancements.

## 📦 Installation

```bash
go get github.com/mstephenholl/go-solar
```

**Requirements:** Go 1.21 or later

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/mstephenholl/go-solar"
)

func main() {
    // Create location for Toronto
    loc := solar.NewLocation(43.65, -79.38)

    // Create time for January 1, 2000
    t := solar.NewTime(2000, time.January, 1)

    // Calculate sunrise
    sunrise, err := solar.Sunrise(loc, t)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sunrise: %s\n", sunrise.Format("15:04:05 MST"))
    // Output: Sunrise: 12:51:00 UTC

    // Calculate sunset
    sunset, err := solar.Sunset(loc, t)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sunset: %s\n", sunset.Format("15:04:05 MST"))
    // Output: Sunset: 21:50:36 UTC

    // Or get both at once (more efficient)
    sunrise, sunset, err = solar.SunriseSunset(loc, t)
    if err != nil {
        log.Fatal(err)
    }
}
```

## 📖 Usage Examples

### Creating Locations and Times

```go
import "github.com/mstephenholl/go-solar"

// Create a location from latitude and longitude
loc := solar.NewLocation(40.7128, -74.0060) // New York City

// Create a time from date components
t := solar.NewTime(2024, time.June, 21)

// Or create from a time.Time object
now := time.Now()
t := solar.NewTimeFromDateTime(now)
```

### Individual Sunrise or Sunset

```go
// Create location and time
loc := solar.NewLocation(40.7128, -74.0060)
t := solar.NewTime(2024, time.June, 21)

// Get just sunrise
sunrise, err := solar.Sunrise(loc, t)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sunrise: %s\n", sunrise.Format("15:04 MST"))

// Get just sunset
sunset, err := solar.Sunset(loc, t)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sunset: %s\n", sunset.Format("15:04 MST"))
```

### Both Sunrise and Sunset

```go
loc := solar.NewLocation(40.7128, -74.0060)
t := solar.NewTime(2024, time.June, 21)

// Calculate both at once (more efficient)
sunrise, sunset, err := solar.SunriseSunset(loc, t)
if err != nil {
    // Handle polar region cases
    if err == solar.ErrSunNeverRises {
        fmt.Println("Polar night - sun never rises on this day")
    } else if err == solar.ErrSunNeverSets {
        fmt.Println("Midnight sun - sun never sets on this day")
    }
    return
}
fmt.Printf("Sunrise: %s, Sunset: %s\n", sunrise.Format("15:04"), sunset.Format("15:04"))
```

### Solar Elevation Angle

```go
// Get sun's elevation at a specific time
loc := solar.NewLocation(40.7128, -74.0060)
when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
elevation := solar.Elevation(loc, when)
fmt.Printf("Sun elevation: %.2f degrees\n", elevation)
```

### Solar Azimuth Angle

```go
// Get sun's azimuth (compass direction) at a specific time
// Azimuth: 0° = North, 90° = East, 180° = South, 270° = West
loc := solar.NewLocation(40.7128, -74.0060)
when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
azimuth := solar.Azimuth(loc, when)
fmt.Printf("Sun azimuth: %.2f degrees\n", azimuth)
```

### Dawn and Dusk (Twilight Times)

Calculate dawn and dusk using civil, nautical, or astronomical twilight definitions:

```go
loc := solar.NewLocation(40.7128, -74.0060)
t := solar.NewTime(2024, time.June, 21)

// Calculate civil dawn and dusk (default, -6° sun angle)
dawn, dusk := solar.DawnDusk(loc, t)
fmt.Printf("Dawn: %s, Dusk: %s\n", dawn.Format("15:04"), dusk.Format("15:04"))

// Or get them individually
dawn := solar.Dawn(loc, t)
dusk := solar.Dusk(loc, t)

// Nautical twilight (-12° sun angle, for marine navigation)
dawn, dusk := solar.DawnDusk(loc, t, solar.Nautical)

// Astronomical twilight (-18° sun angle, for astronomy)
dawn, dusk := solar.DawnDusk(loc, t, solar.Astronomical)
```

**Twilight Types:**
- **Civil** (-6°): Default. Enough light for outdoor activities without artificial lighting
- **Nautical** (-12°): Horizon visible at sea for navigation, general ground outlines visible
- **Astronomical** (-18°): Sky dark enough for astronomical observations

### Custom Elevation Times

```go
loc := solar.NewLocation(40.7128, -74.0060)
t := solar.NewTime(2024, time.June, 21)

// Find when sun reaches a specific elevation (e.g., golden hour at 6°)
morning, evening := solar.TimeOfElevation(loc, 6.0, t)
fmt.Printf("Golden hour: %s to %s\n", morning.Format("15:04"), evening.Format("15:04"))
```

### Working with NMEA GPS Sentences

The package supports parsing location and time data from NMEA GPS sentences, which you can then use with any solar calculation function.

**Supported NMEA sentence types:**
- **RMC** (Recommended Minimum): Includes date, no external date needed
- **GGA** (GPS Fix Data): Requires external date parameters

**Features:**
- Automatic checksum validation
- Supports both hemispheres (N/S, E/W)
- Handles 2-digit year conversion (00-49 → 2000-2049, 50-99 → 1950-1999)
- Detailed error messages for debugging

#### Parsing NMEA Sentences

```go
// Using RMC sentence (includes date)
nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"

// Parse location from NMEA
loc, err := solar.NewLocationFromNMEA(nmea, 0, 0, 0)  // For RMC, date params ignored
if err != nil {
    log.Fatal(err)
}

// Parse time from NMEA
t, err := solar.NewTimeFromNMEA(nmea, 0, 0, 0)  // For RMC, date params ignored
if err != nil {
    log.Fatal(err)
}

// Using GGA sentence (requires external date)
nmea = "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*5C"
loc, err = solar.NewLocationFromNMEA(nmea, 2024, time.June, 21)  // Must provide date for GGA
if err != nil {
    log.Fatal(err)
}
t, err = solar.NewTimeFromNMEA(nmea, 2024, time.June, 21)
if err != nil {
    log.Fatal(err)
}
```

#### Using NMEA Data with Solar Functions

Once you've parsed the NMEA sentence into Location and Time, you can use them with any solar calculation function:

```go
nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"

// Parse NMEA once
loc, err := solar.NewLocationFromNMEA(nmea, 0, 0, 0)
if err != nil {
    log.Fatal(err)
}
t, err := solar.NewTimeFromNMEA(nmea, 0, 0, 0)
if err != nil {
    log.Fatal(err)
}

// Use with any solar function
sunrise, err := solar.Sunrise(loc, t)
sunset, err := solar.Sunset(loc, t)
noon := solar.MeanSolarNoon(loc, t)
dawn := solar.Dawn(loc, t, solar.Civil)
dusk := solar.Dusk(loc, t, solar.Nautical)
azimuth := solar.Azimuth(loc, t.DateTime())
elevation := solar.Elevation(loc, t.DateTime())

// For time-sensitive calculations (azimuth, elevation),
// use time.Now() or the GPS time from NMEA
currentAzimuth := solar.Azimuth(loc, time.Now())
```

**Benefits of this approach:**
- Parse NMEA sentence once, use with multiple calculations
- Compose any combination of Location/Time sources
- Clearer separation between parsing and calculation logic
- More flexible for complex GPS data processing workflows

### Using Generic Helpers

```go
// Generic absolute value - works with any signed type
fmt.Println(solar.Abs(-42))      // int: 42
fmt.Println(solar.Abs(-3.14))    // float64: 3.14

// Floating-point comparison with tolerance
if solar.AlmostEqual(1.0, 1.00001, 0.001) {
    fmt.Println("Values are approximately equal")
}
```

## 🧪 Development

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run benchmarks
make test-bench

# Run all CI checks
make ci
```

### Code Quality

```bash
# Format code
make fmt

# Run linters
make lint

# Security scan
make security
```

## 📊 Performance

All generic helper functions are highly optimized:

```
BenchmarkAbs-11          1000000000   0.25 ns/op   0 B/op   0 allocs/op
BenchmarkAlmostEqual-11  1000000000   0.25 ns/op   0 B/op   0 allocs/op
BenchmarkMin-11          1000000000   0.25 ns/op   0 B/op   0 allocs/op
BenchmarkMax-11          1000000000   0.25 ns/op   0 B/op   0 allocs/op
```

## 🚀 Releases

This project uses **automated releases** with [Calendar Versioning (CalVer)](https://calver.org/).

### Version Format

Releases follow the **CalVer** pattern: `YYYY.MM.MICRO`

- `YYYY` - Full year (e.g., 2025)
- `MM` - Zero-padded month (01-12)
- `MICRO` - Incrementing number for releases within the same month (0, 1, 2, ...)

**Examples:**
- `v2025.10.0` - First release in October 2025
- `v2025.10.1` - Second release in October 2025
- `v2025.11.0` - First release in November 2025

**Why CalVer?** This format provides clear, chronological versioning that makes it easy to understand when a release was created. Similar to Ubuntu's versioning scheme (e.g., 22.04, 24.04).

### Automated Release Process

Every successful push to the `master` branch automatically triggers:

1. **CI Quality Gates** (must pass before release):
   - ✅ Runs full test suite across multiple OS and Go versions
   - ✅ Runs linter checks (golangci-lint)
   - ✅ Runs security scans (gosec, govulncheck)
   - ✅ Builds the package
   - ✅ Runs benchmarks

2. **Release Creation** (only if CI passes):
   - 🏷️ Creates a CalVer git tag (YYYY.MM.MICRO)
   - 📝 Generates categorized changelog from commits
   - 🎉 Creates GitHub release with notes

**Quality First:** Releases are only created when all CI quality gates pass successfully.

### Changelog

Each release includes:
- Categorized commits (Features, Enhancements, Bug Fixes, etc.)
- Installation instructions
- Links to documentation and full changelog

### Installing a Specific Version

```bash
# Latest release
go get github.com/mstephenholl/go-solar

# Specific version (CalVer format)
go get github.com/mstephenholl/go-solar@v2025.10.0
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE.txt](LICENSE.txt) file for details.

## 🙏 Acknowledgments

- Original `go-sunrise` package by [Nathan Osman](https://github.com/nathan-osman)
- Based on the [sunrise equation](https://en.wikipedia.org/wiki/Sunrise_equation) algorithm
- Modernized with Go 1.18+ generics and comprehensive testing

## 📚 Documentation

Full documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/mstephenholl/go-solar).

## 🐛 Issues & Support

- Report bugs via [GitHub Issues](https://github.com/mstephenholl/go-solar/issues)
- For questions, check existing issues or create a new one
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines

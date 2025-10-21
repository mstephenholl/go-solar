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
- 📐 Determine solar elevation angles
- 🛰️ Parse NMEA GPS sentences (GGA, RMC) for location-based calculations
- 🌍 Handle edge cases (polar night, midnight sun)
- 🚀 High performance with zero allocations for core functions
- ✅ 100% test coverage on production code
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
    "time"

    "github.com/mstephenholl/go-solar"
)

func main() {
    // Toronto coordinates
    latitude := 43.65
    longitude := -79.38

    // Calculate sunrise for January 1, 2000
    rise := solar.Sunrise(latitude, longitude, 2000, time.January, 1)
    fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
    // Output: Sunrise: 12:51:00 UTC

    // Calculate sunset for the same day
    set := solar.Sunset(latitude, longitude, 2000, time.January, 1)
    fmt.Printf("Sunset: %s\n", set.Format("15:04:05 MST"))
    // Output: Sunset: 21:50:36 UTC

    // Or get both at once
    rise, set = solar.SunriseSunset(latitude, longitude, 2000, time.January, 1)
}
```

## 📖 Usage Examples

### Individual Sunrise or Sunset

```go
import "github.com/mstephenholl/go-solar"

latitude := 40.7128   // New York City
longitude := -74.0060

// Get just sunrise
rise := solar.Sunrise(latitude, longitude, 2024, time.June, 21)
fmt.Printf("Sunrise: %s\n", rise.Format("15:04 MST"))

// Get just sunset
set := solar.Sunset(latitude, longitude, 2024, time.June, 21)
fmt.Printf("Sunset: %s\n", set.Format("15:04 MST"))
```

### Both Sunrise and Sunset

```go
// Calculate both at once (more efficient)
rise, set := solar.SunriseSunset(latitude, longitude, 2024, time.June, 21)

// Check for special cases (polar regions)
if rise.IsZero() && set.IsZero() {
    fmt.Println("Sun does not rise or set on this day")
}
```

### Solar Elevation Angle

```go
// Get sun's elevation at a specific time
when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
elevation := solar.Elevation(latitude, longitude, when)
fmt.Printf("Sun elevation: %.2f degrees\n", elevation)
```

### Custom Elevation Times

```go
// Find when sun reaches a specific elevation (e.g., golden hour at -6°)
morning, evening := solar.TimeOfElevation(
    latitude, longitude,
    -6.0,  // Elevation angle in degrees
    2024, time.June, 21,
)
```

### Working with NMEA GPS Sentences

Calculate sunrise and sunset directly from NMEA 0183 GPS sentences (GGA or RMC):

```go
// Using RMC sentence (includes date)
nmea := "$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*71"
sunrise, err := solar.SunriseFromNMEA(nmea, 0, 0, 0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sunrise: %s\n", sunrise.Format("15:04 MST"))

// Using GGA sentence (requires external date)
nmea = "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*5C"
sunset, err := solar.SunsetFromNMEA(nmea, 2024, time.June, 21)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sunset: %s\n", sunset.Format("15:04 MST"))

// Get both sunrise and sunset from NMEA
sunrise, sunset, err := solar.SunriseSunsetFromNMEA(nmea, 2024, time.June, 21)
```

**Supported NMEA sentence types:**
- **RMC** (Recommended Minimum): Includes date, no external date needed
- **GGA** (GPS Fix Data): Requires external date parameters

**Features:**
- Automatic checksum validation
- Supports both hemispheres (N/S, E/W)
- Handles 2-digit year conversion (00-49 → 2000-2049, 50-99 → 1950-1999)
- Detailed error messages for debugging

### Using Generic Helpers

```go
// Generic absolute value - works with any signed type
fmt.Println(solar.Abs(-42))      // int: 42
fmt.Println(solar.Abs(-3.14))    // float64: 3.14

// Floating-point comparison with tolerance
if solar.AlmostEqual(1.0, 1.00001, 0.001) {
    fmt.Println("Values are approximately equal")
}

// Min/Max for any numeric type
fmt.Println(solar.Min(5, 10))        // 5
fmt.Println(solar.Max(3.14, 2.71))   // 3.14

// Clamp values to a range
fmt.Println(solar.Clamp(15, 0, 10))  // 10
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

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make ci-quick`)
5. Commit with conventional commits (`git commit -m 'feat: add amazing feature'`)
6. Push to your fork (`git push origin feature/amazing-feature`)
7. Open a Pull Request

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

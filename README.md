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
- 🌍 Handle edge cases (polar night, midnight sun)
- 🚀 High performance with zero allocations
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
    // Calculate sunrise and sunset for Toronto on January 1, 2000
    rise, set := sunrise.SunriseSunset(
        43.65, -79.38,          // Toronto coordinates (lat, long)
        2000, time.January, 1,  // Date
    )

    fmt.Printf("Sunrise: %s\n", rise.Format("15:04:05 MST"))
    fmt.Printf("Sunset:  %s\n", set.Format("15:04:05 MST"))
    // Output:
    // Sunrise: 12:51:00 UTC
    // Sunset:  21:50:36 UTC
}
```

## 📖 Usage Examples

### Basic Sunrise/Sunset Calculation

```go
import "github.com/mstephenholl/go-solar"

// Calculate for any location and date
latitude := 40.7128   // New York City
longitude := -74.0060
year, month, day := 2024, time.June, 21

rise, set := sunrise.SunriseSunset(latitude, longitude, year, month, day)

// Check for special cases (polar regions)
if rise.IsZero() && set.IsZero() {
    fmt.Println("Sun does not rise or set on this day")
}
```

### Solar Elevation Angle

```go
// Get sun's elevation at a specific time
when := time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC)
elevation := sunrise.Elevation(latitude, longitude, when)
fmt.Printf("Sun elevation: %.2f degrees\n", elevation)
```

### Custom Elevation Times

```go
// Find when sun reaches a specific elevation (e.g., golden hour at -6°)
morning, evening := sunrise.TimeOfElevation(
    latitude, longitude,
    -6.0,  // Elevation angle in degrees
    2024, time.June, 21,
)
```

### Using Generic Helpers

```go
// Generic absolute value - works with any signed type
fmt.Println(sunrise.Abs(-42))      // int: 42
fmt.Println(sunrise.Abs(-3.14))    // float64: 3.14

// Floating-point comparison with tolerance
if sunrise.AlmostEqual(1.0, 1.00001, 0.001) {
    fmt.Println("Values are approximately equal")
}

// Min/Max for any numeric type
fmt.Println(sunrise.Min(5, 10))        // 5
fmt.Println(sunrise.Max(3.14, 2.71))   // 3.14

// Clamp values to a range
fmt.Println(sunrise.Clamp(15, 0, 10))  // 10
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

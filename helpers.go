package solar

import (
	"math"
)

// Numeric represents all numeric types that support basic arithmetic operations.
// This includes all integer and floating-point types.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Signed represents all signed numeric types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Float represents all floating-point types.
type Float interface {
	~float32 | ~float64
}

// Abs returns the absolute value of x for any numeric type.
// This is a generic replacement for type-specific abs functions.
//
// Example:
//
//	Abs(-5)     // returns 5 (int)
//	Abs(-5.5)   // returns 5.5 (float64)
//	Abs(int64(-10)) // returns 10 (int64)
func Abs[T Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// AlmostEqual checks if two floating-point numbers are within a given tolerance.
// This is useful for comparing floating-point results where exact equality
// is not possible due to rounding errors.
//
// Example:
//
//	AlmostEqual(1.0, 1.00001, 0.001)  // true
//	AlmostEqual(1.0, 1.1, 0.001)      // false
func AlmostEqual[T Float](a, b, tolerance T) bool {
	return Abs(a-b) <= tolerance
}

// Min returns the minimum of two values for any ordered numeric type.
//
// Example:
//
//	Min(5, 10)      // returns 5
//	Min(3.14, 2.71) // returns 2.71
func Min[T Numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two values for any ordered numeric type.
//
// Example:
//
//	Max(5, 10)      // returns 10
//	Max(3.14, 2.71) // returns 3.14
func Max[T Numeric](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp restricts a value to be within a specified range [minVal, maxVal].
// If the value is less than minVal, it returns minVal.
// If the value is greater than maxVal, it returns maxVal.
// Otherwise, it returns the value unchanged.
//
// Example:
//
//	Clamp(5, 0, 10)   // returns 5
//	Clamp(-5, 0, 10)  // returns 0
//	Clamp(15, 0, 10)  // returns 10
func Clamp[T Numeric](value, minVal, maxVal T) T {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// DegreesToRadians converts degrees to radians for any floating-point type.
// This is a generic helper for angle conversions.
//
// Example:
//
//	DegreesToRadians(180.0)  // returns math.Pi
//	DegreesToRadians(90.0)   // returns math.Pi/2
func DegreesToRadians[T Float](degrees T) T {
	return degrees * T(math.Pi) / 180
}

// RadiansToDegrees converts radians to degrees for any floating-point type.
// This is a generic helper for angle conversions.
//
// Example:
//
//	RadiansToDegrees(math.Pi)    // returns 180.0
//	RadiansToDegrees(math.Pi/2)  // returns 90.0
func RadiansToDegrees[T Float](radians T) T {
	return radians * 180 / T(math.Pi)
}

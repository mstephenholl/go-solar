package solar

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

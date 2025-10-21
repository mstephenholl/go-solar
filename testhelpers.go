package sunrise

import (
	"testing"
)

// TestCase represents a generic table-driven test case with input and expected output.
// This can be used to simplify table-driven tests across the codebase.
type TestCase[I any, O any] struct {
	Name     string
	Input    I
	Expected O
}

// RunTestCases executes a slice of test cases using a provided test function.
// This is a generic helper for table-driven tests.
//
// Example:
//
//	tests := []TestCase[float64, float64]{
//	    {Name: "zero", Input: 0.0, Expected: 0.0},
//	    {Name: "positive", Input: 5.0, Expected: 25.0},
//	}
//	RunTestCases(t, tests, func(input float64) float64 {
//	    return input * input
//	}, AlmostEqual[float64])
func RunTestCases[I any, O comparable](
	t *testing.T,
	tests []TestCase[I, O],
	testFunc func(I) O,
	equalFunc func(O, O) bool,
) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := testFunc(tt.Input)
			if !equalFunc(got, tt.Expected) {
				t.Errorf("got %v, want %v", got, tt.Expected)
			}
		})
	}
}

// CompareFunc is a function that compares two values and returns true if they are equal.
type CompareFunc[T any] func(a, b T) bool

// EqualWithTolerance returns a comparison function that checks if two floating-point
// values are within the specified tolerance.
//
// Example:
//
//	compare := EqualWithTolerance[float64](0.001)
//	compare(1.0, 1.0001) // returns true
func EqualWithTolerance[T Float](tolerance T) CompareFunc[T] {
	return func(a, b T) bool {
		return AlmostEqual(a, b, tolerance)
	}
}

// Equal returns a comparison function for exact equality.
// This is useful for integer types or when exact comparison is needed.
func Equal[T comparable]() CompareFunc[T] {
	return func(a, b T) bool {
		return a == b
	}
}

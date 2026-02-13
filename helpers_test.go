package solar

import (
	"testing"
)

func TestAbs(t *testing.T) {
	// Test with int
	if got := Abs(-5); got != 5 {
		t.Errorf("Abs(-5) = %v, want 5", got)
	}
	if got := Abs(5); got != 5 {
		t.Errorf("Abs(5) = %v, want 5", got)
	}

	// Test with int64
	if got := Abs(int64(-10)); got != 10 {
		t.Errorf("Abs(int64(-10)) = %v, want 10", got)
	}

	// Test with float64
	if got := Abs(-5.5); got != 5.5 {
		t.Errorf("Abs(-5.5) = %v, want 5.5", got)
	}
	if got := Abs(5.5); got != 5.5 {
		t.Errorf("Abs(5.5) = %v, want 5.5", got)
	}

	// Test with float32
	if got := Abs(float32(-3.14)); got != 3.14 {
		t.Errorf("Abs(float32(-3.14)) = %v, want 3.14", got)
	}
}

func TestAlmostEqual(t *testing.T) {
	tests := []struct {
		name      string
		a         float64
		b         float64
		tolerance float64
		want      bool
	}{
		{"equal values", 1.0, 1.0, 0.001, true},
		{"within tolerance", 1.0, 1.00001, 0.001, true},
		{"outside tolerance", 1.0, 1.1, 0.001, false},
		{"negative within tolerance", -1.0, -1.00001, 0.001, true},
		{"negative outside tolerance", -1.0, -1.1, 0.001, false},
		{"zero tolerance exact", 1.0, 1.0, 0.0, true},
		{"zero tolerance different", 1.0, 1.00001, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AlmostEqual(tt.a, tt.b, tt.tolerance); got != tt.want {
				t.Errorf("AlmostEqual(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.tolerance, got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkAbs(b *testing.B) {
	for b.Loop() {
		Abs(-42.5)
	}
}

func BenchmarkAlmostEqual(b *testing.B) {
	for b.Loop() {
		AlmostEqual(1.0, 1.00001, 0.001)
	}
}

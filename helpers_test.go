package solar

import (
	"math"
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

func TestMin(t *testing.T) {
	// Test with int
	if got := Min(5, 10); got != 5 {
		t.Errorf("Min(5, 10) = %v, want 5", got)
	}
	if got := Min(10, 5); got != 5 {
		t.Errorf("Min(10, 5) = %v, want 5", got)
	}

	// Test with float64
	if got := Min(3.14, 2.71); got != 2.71 {
		t.Errorf("Min(3.14, 2.71) = %v, want 2.71", got)
	}

	// Test with negative numbers
	if got := Min(-5, -10); got != -10 {
		t.Errorf("Min(-5, -10) = %v, want -10", got)
	}
}

func TestMax(t *testing.T) {
	// Test with int
	if got := Max(5, 10); got != 10 {
		t.Errorf("Max(5, 10) = %v, want 10", got)
	}
	if got := Max(10, 5); got != 10 {
		t.Errorf("Max(10, 5) = %v, want 10", got)
	}

	// Test with float64
	if got := Max(3.14, 2.71); got != 3.14 {
		t.Errorf("Max(3.14, 2.71) = %v, want 3.14", got)
	}

	// Test with negative numbers
	if got := Max(-5, -10); got != -5 {
		t.Errorf("Max(-5, -10) = %v, want -5", got)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
		max   int
		want  int
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -5, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.value, tt.min, tt.max); got != tt.want {
				t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestDegreesToRadians(t *testing.T) {
	tests := []struct {
		name    string
		degrees float64
		want    float64
	}{
		{"180 degrees", 180.0, math.Pi},
		{"90 degrees", 90.0, math.Pi / 2},
		{"360 degrees", 360.0, 2 * math.Pi},
		{"0 degrees", 0.0, 0.0},
		{"45 degrees", 45.0, math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DegreesToRadians(tt.degrees)
			if !AlmostEqual(got, tt.want, 1e-10) {
				t.Errorf("DegreesToRadians(%v) = %v, want %v", tt.degrees, got, tt.want)
			}
		})
	}
}

func TestRadiansToDegrees(t *testing.T) {
	tests := []struct {
		name    string
		radians float64
		want    float64
	}{
		{"pi radians", math.Pi, 180.0},
		{"pi/2 radians", math.Pi / 2, 90.0},
		{"2*pi radians", 2 * math.Pi, 360.0},
		{"0 radians", 0.0, 0.0},
		{"pi/4 radians", math.Pi / 4, 45.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RadiansToDegrees(tt.radians)
			if !AlmostEqual(got, tt.want, 1e-10) {
				t.Errorf("RadiansToDegrees(%v) = %v, want %v", tt.radians, got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-42.5)
	}
}

func BenchmarkAlmostEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AlmostEqual(1.0, 1.00001, 0.001)
	}
}

func BenchmarkMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Min(5, 10)
	}
}

func BenchmarkMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(5, 10)
	}
}

func BenchmarkClamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Clamp(5, 0, 10)
	}
}

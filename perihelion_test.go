package solar

import (
	"testing"
)

var dataArgumentOfPerihelion = []struct {
	in  float64
	out float64
}{
	// 2017-01-31 - example from the wikipedia article
	{2457784.958333, 102.98437},

	// 1970-01-01 - prime meridian
	{2440588, 102.83467},
	// 2000-01-01 - Toronto (43.65° N, 79.38° W)
	{2451545.2205, 102.93005},
	// 2004-04-01 - (52° N, 5° E)
	{2453096.98611, 102.94356},
}

func TestArgumentOfPerihelion(t *testing.T) {
	for _, tt := range dataArgumentOfPerihelion {
		v := ArgumentOfPerihelion(tt.in)
		if Round(v, DefaultPlaces) != Round(tt.out, DefaultPlaces) {
			t.Fatalf("%f != %f", v, tt.out)
		}
	}
}

package sunrise

// ArgumentOfPerihelion calculates the argument of periapsis for the earth on
// the given Julian day.
func ArgumentOfPerihelion(d float64) float64 {
	return PerihelionBase + PerihelionRate*(d-J2000)/JulianCenturyDays
}

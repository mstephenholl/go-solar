package sunrise

import (
	"math"
)

const (
	// Degree provides a precise fraction for converting between degrees and
	// radians.
	Degree = math.Pi / 180

	// J2000 is the Julian date for January 1, 2000, 12:00:00 TT.
	// This is the standard epoch used in modern astronomical calculations.
	J2000 = 2451545

	// SolarMeanAnomalyBase is the mean anomaly of the sun at J2000 epoch (in degrees).
	// This represents the angular distance from perihelion at the reference epoch.
	SolarMeanAnomalyBase = 357.5291

	// SolarMeanAnomalyRate is the daily change in mean anomaly (degrees per day).
	// This accounts for Earth's orbital motion around the sun.
	SolarMeanAnomalyRate = 0.98560028

	// FullCircleDegrees represents a complete rotation in degrees.
	FullCircleDegrees = 360.0

	// HalfCircleDegrees represents a semicircle in degrees.
	HalfCircleDegrees = 180.0

	// EquationOfCenterC1 is the first coefficient in the equation of center calculation.
	// This is the dominant term in the Earth's orbital eccentricity correction.
	EquationOfCenterC1 = 1.9148

	// EquationOfCenterC2 is the second coefficient in the equation of center calculation.
	// This accounts for higher-order orbital perturbations.
	EquationOfCenterC2 = 0.0200

	// EquationOfCenterC3 is the third coefficient in the equation of center calculation.
	// This is a small correction term for orbital accuracy.
	EquationOfCenterC3 = 0.0003

	// SinDeclinationCoefficient is sin(ε) where ε is Earth's axial tilt (obliquity).
	// Earth's axial tilt is approximately 23.44°, and sin(23.44°) ≈ 0.39779
	SinDeclinationCoefficient = 0.39779

	// EquationOfTimeC1 is the first coefficient in the equation of time calculation.
	// This accounts for the eccentricity of Earth's orbit.
	EquationOfTimeC1 = 0.0053

	// EquationOfTimeC2 is the second coefficient in the equation of time calculation.
	// This accounts for the obliquity of the ecliptic.
	EquationOfTimeC2 = 0.0069

	// SunriseCorrectionAngle is the correction angle for sunrise/sunset calculations.
	// This is approximately -0.833° (-50 arc minutes), accounting for:
	// - Atmospheric refraction (34 arc minutes)
	// - Sun's angular diameter (16 arc minutes)
	SunriseCorrectionAngle = -0.01449

	// PerihelionBase is the argument of perihelion at J2000 epoch (in degrees).
	// This is the angle from the vernal equinox to perihelion.
	PerihelionBase = 102.93005

	// PerihelionRate is the rate of change of perihelion per Julian century.
	// This accounts for the precession of Earth's orbit.
	PerihelionRate = 0.3179526

	// JulianCenturyDays is the number of days in a Julian century.
	JulianCenturyDays = 36525.0

	// LongitudeDivisor is used in mean solar noon calculation.
	// Dividing longitude by 360 converts it to a fractional day offset.
	LongitudeDivisor = 360.0
)

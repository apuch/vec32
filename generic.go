package vec32

import (
	"math"
)

// Compare two floats, if they are almost equal
//
// Floats are equal, if:
//
// - they are inf of the same sign ( Inf==Inf)
//
// - if one is less then 100*FLOAT_MIN and the difference between both is less than 100*FLOAT_MIN
//
// - the relative difference ( (a-b)/(a+b)/2 ) is less than 100 EPS
//
// - NaN is unequal to NaN and everything else
//
func AlmostEqual(a, b float32) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a < 100*FLOAT_MIN || b < 100*FLOAT_MIN {
		return diff < 100*FLOAT_MIN
	}
	return 2*diff/(a+b) < 10*EPS
}

// AlmostEqual() in R3
func AlmostEqual3(a, b *Vec3) bool {
	return AlmostEqual(a.X, b.X) && AlmostEqual(a.Y, b.Y) && AlmostEqual(a.Z, b.Z)
}

// f32-fabs
func Abs(v float32) float32 {
	if v >= 0 {
		return v
	}
	return -v
}

// Return +/- Inf (see math.Inf())
func Inf(s int) float32 {
	return float32(math.Inf(s))
}

// See math.IsInf()
func IsInf(v float32, s int) bool {
	return math.IsInf(float64(v), s)
}

// Return NaN (see math.NaN())
func NaN() float32 {
	return float32(math.NaN())
}

// See math.IsNaN()
func IsNaN(v float32) bool {
	return math.IsNaN(float64(v))
}

// Return Sqrt (see math.Sqrt())
func Sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

// Max
func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// Min
func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

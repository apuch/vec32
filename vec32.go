package vec32

// Basic Vector algebra
//
// Adding, Substracting, etc. Just a basic (stupid) implementation at
// the moment

import (
	"fmt"
	"math"
)

// Abstract Vector type.
//
// This interface encapsulate common functionality of a vector. The result
// have to be a scalar as the cross product of two 3D-vectors is a 3D-vector
// again and not an abstract base type anymore.
type Vector interface {
	// Get the dimension (2, 3, ....)
	Dim() int
	// The (euklidian) length
	Length() float32
	// The (euklidian) length squared (no square root, prefere this)
	LengthSq() float32
}

// A two dimensional Vector
type Vec2 struct {
	X, Y float32
}

// A three dimensional Vector
type Vec3 struct {
	X, Y, Z float32
}

///////////////////////////////////////////////////////////////////////////

// Dimension of the vector
func (v Vec2) Dim() int { return 2 }

// The euklidian length
func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// The euklidian length squared
func (v Vec2) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

// string representation (octave style)
func (v Vec2) String() string {
	return fmt.Sprintf("[%g; %g]", v.X, v.Y)
}

// Dot product
func (v Vec2) Dot(v2 Vec2) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Add two vectors
func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v.X + v2.X, v.Y + v2.Y}
}

// Substract two vectors
func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v.X - v2.X, v.Y - v2.Y}
}

// Multiply two vectors elementwise
func (v Vec2) Mul(v2 Vec2) Vec2 {
	return Vec2{v.X * v2.X, v.Y * v2.Y}
}

// get the normalized vector ( |v| == 1 )
func (v Vec2) Normalize() Vec2 {
	n := v.Length()
	return Vec2{v.X / n, v.Y / n}
}

// Scale a vector by a factor
func (v Vec2) Scale(s float32) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

//////////////////////////////////////////////////////////////////////////

// Dimension of the vector
func (v Vec3) Dim() int { return 3 }

// The euklidian length
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// The euklidian length squared
func (v Vec3) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// string representation (octave style)
func (v Vec3) String() string {
	return fmt.Sprintf("[%g; %g; %g]", v.X, v.Y, v.Z)
}

// Dot product
func (v Vec3) Dot(v2 Vec3) float32 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Add two vectors
func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

// Substract two vectors
func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

// Multiply two vectors elementwise
func (v Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

// get the normalized vector ( |v| == 1 )
func (v Vec3) Normalize() Vec3 {
	n := v.Length()
	return Vec3{v.X / n, v.Y / n, v.Z / n}
}

// Scale a vector by a factor
func (v Vec3) Scale(s float32) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

///////////////////////////////////////////////////////////////////////////

// f32-fabs
func Abs(v float32) float32 {
	if v >= 0 {
		return v
	}
	return -v
}

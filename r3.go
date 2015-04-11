package vec32

import (
	"fmt"
)

// Create a new Vector
func NewVec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z, 0}
}

// Dimension of the vector
func (v *Vec3) Dim() int { return 3 }

// The euklidian length
func (v *Vec3) Length() float32 {
	return LengthR3(v)
}

// the euklidian length
func LengthR3(v *Vec3) float32

func lengthR3(v *Vec3) float32 {
	return Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// The euklidian length squared
func (v *Vec3) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// string representation (octave style)
func (v *Vec3) String() string {
	return fmt.Sprintf("[%g; %g; %g]", v.X, v.Y, v.Z)
}

// Dot product
func (v *Vec3) Dot(v2 *Vec3) float32 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Dot product
func DotR3(v *Vec3) float32 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Add two vectors
func (v *Vec3) Add(v2 *Vec3) *Vec3 {
	return &Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z, 0}
}

// Add two vectors (explicit)
func Add3(v1, v2, v3 *Vec3)

func add3(v1, v2, v3 *Vec3) {
	v3.X = v1.X + v2.X
	v3.Y = v1.Y + v2.Y
	v3.Z = v1.Z + v2.Z
}

// Substract two vectors (explicit)
func Sub3(v1, v2, v3 *Vec3) {
	v3.X = v1.X - v2.X
	v3.Y = v1.Y - v2.Y
	v3.Z = v1.Z - v2.Z
}

// Substract two vectors
func (v *Vec3) Sub(v2 *Vec3) *Vec3 {
	return &Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z, 0}
}

// Multiply two vectors elementwise
func (v *Vec3) Mul(v2 *Vec3) *Vec3 {
	return &Vec3{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z, 0}
}

// get the normalized vector ( |v| == 1 )
func (v *Vec3) Normalize() *Vec3 {
	n := v.Length()
	return &Vec3{v.X / n, v.Y / n, v.Z / n, 0}
}

// Scale a vector by a factor
func (v *Vec3) Scale(s float32) *Vec3 {
	return &Vec3{v.X * s, v.Y * s, v.Z * s, 0}
}

// Cross-Product
func (a *Vec3) Cross(b *Vec3) *Vec3 {
	return &Vec3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X, 0}
}

// Cross-Product (implicit)
func Cross3(a, b, c *Vec3) {
	c.X = a.Y*b.Z - a.Z*b.Y
	c.Y = a.Z*b.X - a.X*b.Z
	c.Z = a.X*b.Y - a.Y*b.X
}

// Join two boxes, adding the second to the first
func (bb *OrthoBox) Add(bb2 *OrthoBox) {
	orthoBoxAdd(bb, bb2)
}

func OrthoBoxAdd(bb1, bb2 *OrthoBox)

func orthoBoxAdd(bb1, bb2 *OrthoBox) {
	bb1.P0.X = Min(bb1.P0.X, bb2.P0.X)
	bb1.P0.Y = Min(bb1.P0.Y, bb2.P0.Y)
	bb1.P0.Z = Min(bb1.P0.Z, bb2.P0.Z)
	bb1.P1.X = Max(bb1.P1.X, bb2.P1.X)
	bb1.P1.Y = Max(bb1.P1.Y, bb2.P1.Y)
	bb1.P1.Z = Max(bb1.P1.Z, bb2.P1.Z)
}

// Equal
func (a *Vec3) IsEqual(b *Vec3) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

// for testing -- just return
func doNop()

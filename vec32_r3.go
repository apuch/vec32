package vec32

import (
	"fmt"
)

// Dimension of the vector
func (v *Vec3) Dim() int { return 3 }

// The euklidian length
func (v *Vec3) Length() float32 {
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

// Add two vectors
func (v *Vec3) Add(v2 *Vec3) *Vec3 {
	return &Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

// Add two vectors (explicit)
func Add3(v1, v2, v3 *Vec3) {
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
	return &Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

// Multiply two vectors elementwise
func (v *Vec3) Mul(v2 *Vec3) *Vec3 {
	return &Vec3{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

// get the normalized vector ( |v| == 1 )
func (v *Vec3) Normalize() *Vec3 {
	n := v.Length()
	return &Vec3{v.X / n, v.Y / n, v.Z / n}
}

// Scale a vector by a factor
func (v *Vec3) Scale(s float32) *Vec3 {
	return &Vec3{v.X * s, v.Y * s, v.Z * s}
}

// Cross-Product
func (a *Vec3) Cross(b *Vec3) *Vec3 {
	return &Vec3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

// Cross-Product (implicit)
func Cross3(a, b, c *Vec3) {
	c.X = a.Y*b.Z - a.Z*b.Y
	c.Y = a.Z*b.X - a.X*b.Z
	c.Z = a.X*b.Y - a.Y*b.X
}

// Equal
func (a *Vec3) IsEqual(b *Vec3) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

package vec32

import (
	"fmt"
)

// Dimension of the vector
func (v Vec2) Dim() int { return 2 }

// The euklidian length
func (v *Vec2) Length() float32 {
	return Sqrt(v.X*v.X + v.Y*v.Y)
}

// The euklidian length squared
func (v *Vec2) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

// string representation (octave style)
func (v *Vec2) String() string {
	return fmt.Sprintf("[%g; %g]", v.X, v.Y)
}

// Dot product
func (v *Vec2) Dot(v2 *Vec2) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Add two vectors
func (v *Vec2) Add(v2 *Vec2) *Vec2 {
	return &Vec2{v.X + v2.X, v.Y + v2.Y}
}

// Add two vectors (explicit)
func Add2(v1, v2, v3 *Vec2) {
	v3.X = v1.X + v2.X
	v3.Y = v1.Y + v2.Y
}

// Substract two vectors
func (v *Vec2) Sub(v2 *Vec2) *Vec2 {
	return &Vec2{v.X - v2.X, v.Y - v2.Y}
}

// Multiply two vectors elementwise
func (v *Vec2) Mul(v2 *Vec2) *Vec2 {
	return &Vec2{v.X * v2.X, v.Y * v2.Y}
}

// get the normalized vector ( |v| == 1 )
func (v *Vec2) Normalize() *Vec2 {
	n := v.Length()
	return &Vec2{v.X / n, v.Y / n}
}

// Scale a vector by a factor
func (v *Vec2) Scale(s float32) *Vec2 {
	return &Vec2{v.X * s, v.Y * s}
}

package vec32

import (
	"fmt"
	"math"
)

type Vector interface {
	Dim() int
	Length() float32
	LengthSq() float32
}

type Vec2 struct {
	X, Y float32
}

func (v *Vec2) Dim() int { return 2 }

func (v *Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vec2) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vec2) String() string {
	return fmt.Sprintf("[%g; %g]", v.X, v.Y)
}

func (v *Vec2) Dot(v2 *Vec2) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

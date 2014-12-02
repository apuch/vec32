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
	x, y float32
}

func (v *Vec2) Dim() int { return 2 }

func (v *Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}

func (v *Vec2) LengthSq() float32 {
	return v.x*v.x + v.y*v.y
}

func (v *Vec2) String() string {
	return fmt.Sprintf("[%g; %g]", v.x, v.y)
}

func (v *Vec2) Dot(v2 *Vec2) float32 {
	return v.x*v2.x + v.y*v2.y
}

package vec32

// Basic Vector algebra
//
// Adding, Substracting, etc. Just a basic (stupid) implementation at
// the moment

const EPS = float32(1.192e-7)
const FLOAT_MIN = 2.938735877055719e-39

// A two dimensional Vector
type Vec2 struct {
	X, Y float32
}

// A three dimensional Vector
type Vec3 struct {
	X, Y, Z float32
}

// Triangle in R3
type Triangle struct {
	P1, P2, P3 *Vec3
}

// A Ray
type Ray struct {
	P0 Vec3
	N  Vec3
}

// Box othogonal to axis
//
// P1 is always above P0
type OrthoBox struct {
	P0, P1 Vec3
}

type Intersection struct {
	t, u, v float32
}

// a generic object you can see
type Object interface {
	OrthoBox() OrthoBox
	Intersect(r *Ray, i *Intersection) float32
}

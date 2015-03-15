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

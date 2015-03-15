package vec32

// Basic Vector algebra
//
// Adding, Substracting, etc. Just a basic (stupid) implementation at
// the moment

const EPS = float32(1.192e-7)
const FLOAT_MIN = 2.938735877055719e-39

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

// Triangle in R3
type Triangle struct {
	P1, P2, P3 *Vec3
}

// A Ray
type Ray struct {
	P0 Vec3
	N  Vec3
}

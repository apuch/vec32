package vec32

import (
	"fmt"
	"io/ioutil"
)

var INF float32
var INF_NEG float32
var ORTHO_EMPTY OrthoBox

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
	X, Y, Z, pad float32
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

// Struct holding all informations of a mesh
type Mesh struct {
	// The vertices we have
	Verts []Vec3
	Tris  []Triangle
}

// Box othogonal to axis
//
// P1 is always above P0
type OrthoBox struct {
	P0, P1 Vec3
}

// Nice string representation of an orthobox
func (b *OrthoBox) String() string {
	return fmt.Sprintf("{%s->%s}", b.P0.String(), b.P1.String())
}

type Intersection struct {
	t, u, v float32
}

// a generic object you can see
type Object interface {
	OrthoBox() OrthoBox
}

func (bb *OrthoBox) Area() float32 {
	dx := bb.P1.X - bb.P0.X
	dy := bb.P1.Y - bb.P0.Y
	dz := bb.P1.Z - bb.P0.Z
	return 2 * (dx*dy + dx*dz + dy*dz)
}

func init() {
	INF = Inf(1)
	INF_NEG = Inf(-1)
	ORTHO_EMPTY = OrthoBox{NewVec3(INF, INF, INF), NewVec3(INF_NEG, INF_NEG, INF_NEG)}
	InitLogging(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
}

func (tri *Triangle) String() string {
	return fmt.Sprintf("[%s-%s-%s]", tri.P1, tri.P2, tri.P3)
}

package vec32

import (
	"testing"
)

var v1 = Vec3{1, 2, 3}
var v2 = Vec3{4, 5, 6}
var v3 = Vec3{7, 8, 9}

func TestTriangle(t *testing.T) {
	tri := Triangle{&v1, &v2, &v3}
	if !tri.P1.IsEqual(&v1) || !tri.P2.IsEqual(&v2) || !tri.P3.IsEqual(&v3) {
		t.Errorf("Initialization failed")
	}
}

func TestRay(t *testing.T) {
	r := NewRay(&v1, &v2)
	n := v2.Sub(v1).Normalize()
	if !r.P0.IsEqual(&v1) || !r.N.IsEqual(&n) {
		t.Errorf("Initialization failed")
	}
}

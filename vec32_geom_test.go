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
	n := v2.Sub(&v1).Normalize()
	if !r.P0.IsEqual(&v1) || !r.N.IsEqual(n) {
		t.Errorf("Initialization failed")
	}
}

func TestIntersect(t *testing.T) {
	tri := Triangle{&Vec3{0, 0, 0}, &Vec3{1, 0, 0}, &Vec3{0, 1, 0}}
	var cases = []struct {
		ray Ray
		t   float32
	}{
		{Ray{Vec3{0.3, 0.3, -2}, Vec3{0, 0, 1}}, 2},
		{Ray{Vec3{-0.3, 0.3, -2}, Vec3{0, 0, 1}}, Inf(1)},
		{Ray{Vec3{0.3, -0.3, -2}, Vec3{0, 0, 1}}, Inf(1)},
		{Ray{Vec3{0, 0, 2}, Vec3{0, 0, 1}}, Inf(1)},
		{Ray{Vec3{0, 0, 2}, Vec3{0, 1, 0}}, Inf(1)},
	}
	var inter Intersection
	for i, tc := range cases {
		valT := tc.ray.Intersect(&tri, &inter, Inf(1))
		if !AlmostEqual(valT, tc.t) {
			t.Errorf("failed at tc %d: expected %f, got %f", i, tc.t, valT)
		}
	}
}

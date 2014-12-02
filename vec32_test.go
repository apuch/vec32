package vec32

import (
	"math"
	"testing"
)

func TestVec2(t *testing.T) {
	v := Vec2{3, 4}
	v2 := Vec2{5, 6}
	testVec(t, v, 2, 5.0)
	testFloat(t, "X", 3, v.X)
	testFloat(t, "Y", 4, v.Y)
	testString(t, "String()", "[3; 4]", v.String())
	testFloat(t, "Dot()", 3*5+4*6, v.Dot(v2))
	testVec2(t, "Add()", Vec2{8, 10}, v.Add(v2))
	testVec2(t, "Sub()", Vec2{-2, -2}, v.Sub(v2))
	testVec2(t, "Normalize()", Vec2{3.0 / 5.0, 4.0 / 5.0}, v.Normalize())
	testVec2(t, "Scale()", Vec2{6, 8}, v.Scale(2))
	testVec2(t, "Mul()", Vec2{15, 24}, v.Mul(v2))
}

func TestNormNaN(t *testing.T) {
	v2 := Vec2{}.Normalize()
	testNaN(t, "X", v2.X)
	testNaN(t, "Y", v2.Y)
}

func testVec(t *testing.T, v Vector, dim int, length float32) {
	testFloat(t, "Dimension", float32(dim), float32(v.Dim()))
	testFloat(t, "Length", length, v.Length())
	testFloat(t, "LengthSq", length*length, v.LengthSq())
}

func testFloat(t *testing.T, name string, exp, cur float32) {
	if exp != cur {
		t.Errorf("%s is wrong - expected %f got %f", name, exp, cur)
	}
}

func testString(t *testing.T, name, exp, cur string) {
	if exp != cur {
		t.Errorf("%s is wrong - expected %s got %s", name, exp, cur)
	}
}

func testVec2(t *testing.T, name string, exp, cur Vec2) {
	d := exp.Sub(cur).LengthSq() / exp.LengthSq()
	if Abs(d) > 1e-6 {
		t.Errorf("%s is wrong - expected %s got %s", name, exp, cur)
	}
}

func testNaN(t *testing.T, name string, val float32) {
	if !math.IsNaN(float64(val)) {
		t.Errorf("%s is wrong - %g is not NaN", val)
	}
}

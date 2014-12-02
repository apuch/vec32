package vec32

import (
	"testing"
)

func TestVec2(t *testing.T) {
	v := &Vec2{3, 4}
	testVec(t, v, 2, 5.0)
	testString(t, "String()", "[3; 4]", v.String())
	testFloat(t, "Dot()", 3*5+4*6, v.Dot(&Vec2{5, 6}))
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

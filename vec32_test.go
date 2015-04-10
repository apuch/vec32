package vec32

import (
	//"io/ioutil"
	"math"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	InitLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	os.Exit(m.Run())
}

func TestVec2(t *testing.T) {
	v := Vec2{3, 4}
	v2 := Vec2{5, 6}
	testVec2Generic(t, v, 2, 5.0)
	testFloat(t, "X", 3, v.X)
	testFloat(t, "Y", 4, v.Y)
	testString(t, "String()", "[3; 4]", v.String())
	testFloat(t, "Dot()", 3*5+4*6, v.Dot(&v2))
	testVec2(t, "Add()", &Vec2{8, 10}, v.Add(&v2))
	testVec2(t, "Sub()", &Vec2{-2, -2}, v.Sub(&v2))
	testVec2(t, "Normalize()", &Vec2{3.0 / 5.0, 4.0 / 5.0}, v.Normalize())
	testVec2(t, "Scale()", &Vec2{6, 8}, v.Scale(2))
	testVec2(t, "Mul()", &Vec2{15, 24}, v.Mul(&v2))
}

func TestVec3(t *testing.T) {
	v := NewVec3(3, 4, 5)
	v2 := NewVec3(5, 6, 7)
	n := float32(math.Sqrt(3*3 + 4*4 + 5*5))
	testVec3Generic(t, v, 3, n)
	testFloat(t, "X", 3, v.X)
	testFloat(t, "Y", 4, v.Y)
	testFloat(t, "Z", 5, v.Z)
	testString(t, "String()", "[3; 4; 5]", v.String())
	testFloat(t, "Dot()", 3*5+4*6+5*7, v.Dot(&v2))
	testVec3(t, "Add()", NewVec3(8, 10, 12), *v.Add(&v2))
	testVec3(t, "Sub()", NewVec3(-2, -2, -2), *v.Sub(&v2))
	testVec3(t, "Normalize()", NewVec3(3.0/n, 4.0/n, 5.0/n), *v.Normalize())
	testVec3(t, "Scale()", NewVec3(6, 8, 10), *v.Scale(2))
	testVec3(t, "Mul()", NewVec3(15, 24, 35), *v.Mul(&v2))
}

func TestNormNaN(t *testing.T) {
	v1 := Vec2{}
	v2 := v1.Normalize()
	testNaN(t, "X", v2.X)
	testNaN(t, "Y", v2.Y)
}

func TestAlmostEqual(t *testing.T) {
	var cases = []struct {
		a, b    float32
		isEqual bool
	}{
		{0, 0, true},
		{3, 3, true},
		{3, 4, false},
		{4, 3, false},
		{Inf(+1), Inf(+1), true},
		{Inf(-1), Inf(+1), false},
		{Inf(+1), Inf(-1), false},
		{Inf(-1), Inf(-1), true},
		{0, FLOAT_MIN * 10, true},
		{FLOAT_MIN * 10, 0, true},
		{FLOAT_MIN * 80, FLOAT_MIN * 120, true},
		{FLOAT_MIN * 120, FLOAT_MIN * 80, true},
		{5, 5 * (1 + 10*EPS), true},
		{5, 5 * (1 + 200*EPS), false},
		{NaN(), NaN(), false},
		{NaN(), 0, false},
	}
	for i, tc := range cases {
		if AlmostEqual(tc.a, tc.b) != tc.isEqual {
			t.Errorf("tc %d - fail", i)
		}
	}
}

func TestEqual(t *testing.T) {
	var cases = []struct {
		a, b    Vec3
		isEqual bool
	}{
		{NewVec3(1, 2, 3), NewVec3(1, 2, 3), true},
		{NewVec3(0, 0, 0), NewVec3(0, 0, 0), true},
		{NewVec3(0, 0, 0), NewVec3(FLOAT_MIN, 0, 0), false},
		{NewVec3(0, 0, 0), NewVec3(0, FLOAT_MIN, 0), false},
		{NewVec3(0, 0, 0), NewVec3(0, 0, FLOAT_MIN), false},
	}
	for i, tc := range cases {
		if tc.a.IsEqual(&tc.b) != tc.isEqual || tc.b.IsEqual(&tc.a) != tc.isEqual {
			t.Errorf("tc %d failed", i)
		}
	}

}

func TestCross(t *testing.T) {
	c := NewVec3(5*7-6*6, 5*6-4*7, 4*6-5*5)
	testVec3(t, "crossProduct", c, *v3_1.Cross(&v3_2))
}

func TestOrthoBoxArea(t *testing.T) {
	var cases = []struct {
		bb   OrthoBox
		area float32
	}{
		{OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}, 6},
		{OrthoBox{NewVec3(1, 1, 1), NewVec3(1, 1, 1)}, 0},
		{OrthoBox{NewVec3(1, 1, 1), NewVec3(2, 2, 2)}, 6},
		{OrthoBox{NewVec3(1, 1, 1), NewVec3(5, 2, 2)}, 18},
		{OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 0)}, 2},
	}
	for i, tc := range cases {
		area := tc.bb.Area()
		if area != tc.area {
			t.Errorf("tc %d: expected area: %f current area: %f", i, tc.area, area)
		}
	}
}

var (
	v2_1 = Vec2{3, 4}
	v2_2 = Vec2{5, 6}
	v2_3 = Vec2{0, 0}
	v3_1 = NewVec3(4, 5, 6)
	v3_2 = NewVec3(5, 6, 7)
	v3_3 = NewVec3(0, 0, 0)
)

func BenchmarkLengthSq2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2_1.Dot(&v2_2)
	}
}

func BenchmarkAdd2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2_1.Add(&v2_2)
	}
}

func BenchmarkAdd2E(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add2(&v2_1, &v2_2, &v2_3)
	}
}

func BenchmarkAdd3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3_1.Add(&v3_2)
	}
}

func BenchmarkAdd3E(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add3(&v3_1, &v3_2, &v3_3)
	}
}

func testVec2Generic(t *testing.T, v Vec2, dim int, length float32) {
	testFloat(t, "Dimension", float32(dim), float32(v.Dim()))
	testFloat(t, "Length", length, v.Length())
	testFloat(t, "LengthSq", length*length, v.LengthSq())
}

func testVec3Generic(t *testing.T, v Vec3, dim int, length float32) {
	testFloat(t, "Dimension", float32(dim), float32(v.Dim()))
	testFloat(t, "Length", length, v.Length())
	testFloat(t, "LengthSq", length*length, v.LengthSq())
}

func testFloat(t *testing.T, name string, exp, cur float32) {
	if !AlmostEqual(exp, cur) {
		t.Errorf("%s is wrong - expected %f got %f", name, exp, cur)
	}
}

func testString(t *testing.T, name, exp, cur string) {
	if exp != cur {
		t.Errorf("%s is wrong - expected %s got %s", name, exp, cur)
	}
}

func testVec2(t *testing.T, name string, exp, cur *Vec2) {
	v := exp.Sub(cur)
	if !AlmostEqual(v.Length(), 0) {
		t.Errorf("%s is wrong - expected %s got %s", name, exp.String(), cur.String())
	}
}

func testVec3(t *testing.T, name string, exp, cur Vec3) {
	v := exp.Sub(&cur)
	if !AlmostEqual(v.Length(), 0) {
		t.Errorf("%s is wrong - expected %s got %s", name, exp.String(), cur.String())
	}
}

func testNaN(t *testing.T, name string, val float32) {
	if !IsNaN(val) {
		t.Errorf("%s is wrong - %g is not NaN", name, val)
	}
}

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

func TestVec3(t *testing.T) {
	v := Vec3{3, 4, 5}
	v2 := Vec3{5, 6, 7}
	n := float32(math.Sqrt(3*3 + 4*4 + 5*5))
	testVec(t, v, 3, n)
	testFloat(t, "X", 3, v.X)
	testFloat(t, "Y", 4, v.Y)
	testFloat(t, "Z", 5, v.Z)
	testString(t, "String()", "[3; 4; 5]", v.String())
	testFloat(t, "Dot()", 3*5+4*6+5*7, v.Dot(v2))
	testVec3(t, "Add()", Vec3{8, 10, 12}, v.Add(v2))
	testVec3(t, "Sub()", Vec3{-2, -2, -2}, v.Sub(v2))
	testVec3(t, "Normalize()", Vec3{3.0 / n, 4.0 / n, 5.0 / n}, v.Normalize())
	testVec3(t, "Scale()", Vec3{6, 8, 10}, v.Scale(2))
	testVec3(t, "Mul()", Vec3{15, 24, 35}, v.Mul(v2))
}

func TestNormNaN(t *testing.T) {
	v2 := Vec2{}.Normalize()
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

func TestCross(t *testing.T) {
	c := Vec3{5*7 - 6*6, 5*6 - 4*7, 4*6 - 5*5}
	testVec3(t, "crossProduct", c, v3_1.Cross(v3_2))
}

var (
	v2_1 = Vec2{3, 4}
	v2_2 = Vec2{5, 6}
	v2_3 = Vec2{0, 0}
	v3_1 = Vec3{4, 5, 6}
	v3_2 = Vec3{5, 6, 7}
	v3_3 = Vec3{0, 0, 0}
)

func BenchmarkLengthSq2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2_1.Dot(v2_2)
	}
}

func BenchmarkLengthSq2E(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Dot(&v2_1, &v2_2)
	}
}

func BenchmarkAdd2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2_1.Add(v2_2)
	}
}

func BenchmarkAdd2E(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add2(&v2_1, &v2_2, &v2_3)
	}
}

func BenchmarkAdd3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3_1.Add(v3_2)
	}
}

func BenchmarkAdd3E(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add3(&v3_1, &v3_2, &v3_3)
	}
}

func testVec(t *testing.T, v Vector, dim int, length float32) {
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

func testVec2(t *testing.T, name string, exp, cur Vec2) {
	if !AlmostEqual(exp.Sub(cur).Length(), 0) {
		t.Errorf("%s is wrong - expected %s got %s", name, exp, cur)
	}
}

func testVec3(t *testing.T, name string, exp, cur Vec3) {
	if !AlmostEqual(exp.Sub(cur).Length(), 0) {
		t.Errorf("%s is wrong - expected %s got %s", name, exp, cur)
	}
}

func testNaN(t *testing.T, name string, val float32) {
	if !IsNaN(val) {
		t.Errorf("%s is wrong - %g is not NaN", name, val)
	}
}

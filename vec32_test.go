package vec32

import (
	//"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	InitLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	os.Exit(m.Run())
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

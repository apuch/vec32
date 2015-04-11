package vec32

import (
	"math"
	"testing"
)

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

func BenchmarkR3Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3_1.Add(&v3_2)
	}
}

func BenchmarkR3AddE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add3(&v3_1, &v3_2, &v3_3)
	}
}

func BenchmarkR3AddEGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add3(&v3_1, &v3_2, &v3_3)
	}
}

func BenchmarkR3LengthSq(b *testing.B) {
	v := NewVec3(3, 4, 5)
	for i := 0; i < b.N; i++ {
		LengthR3(&v)
	}
}

func BenchmarkR3LengthSqGeneric(b *testing.B) {
	v := NewVec3(3, 4, 5)
	for i := 0; i < b.N; i++ {
		lengthR3(&v)
	}
}
func BenchmarkR3LengthSqIndirect(b *testing.B) {
	v := NewVec3(3, 4, 5)
	for i := 0; i < b.N; i++ {
		v.Length()
	}
}

func BenchmarkR3Dot(b *testing.B) {
	v := NewVec3(3, 4, 5)
	for i := 0; i < b.N; i++ {
		DotR3(&v)
	}
}

func BenchmarkR3Empty(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkR3Nop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doNop()
	}
}

func BenchmarkR3OrthoboxAdd(b *testing.B) {
	bb1 := OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}
	bb2 := OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}
	for i := 0; i < b.N; i++ {
		OrthoBoxAdd(&bb1, &bb2)
	}
}

func BenchmarkR3OrthoboxAddGeneric(b *testing.B) {
	bb1 := OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}
	bb2 := OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}
	for i := 0; i < b.N; i++ {
		orthoBoxAdd(&bb1, &bb2)
	}
}

func testVec3Generic(t *testing.T, v Vec3, dim int, length float32) {
	testFloat(t, "Dimension", float32(dim), float32(v.Dim()))
	testFloat(t, "Length", length, v.Length())
	testFloat(t, "LengthSq", length*length, v.LengthSq())
}

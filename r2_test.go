package vec32

import (
	"testing"
)

func testVec2Generic(t *testing.T, v Vec2, dim int, length float32) {
	testFloat(t, "Dimension", float32(dim), float32(v.Dim()))
	testFloat(t, "Length", length, v.Length())
	testFloat(t, "LengthSq", length*length, v.LengthSq())
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

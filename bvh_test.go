package vec32

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	opts := NewBVHDefaultOptions()
	opts.TraversalCost = INF // no split at all

	var cases = []struct {
		file string
		bb   OrthoBox
		cost float32
	}{
		{"paulbourke.net.sample1.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}, 12 * 6},
		{"two_cubes.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(4, 1, 1)}, 24 * 18},
		{"people.sc.fsu.edu.helix.ply", OrthoBox{NewVec3(-23.75, -44.6, -30.25), NewVec3(36.75, 47.65, 30.25)}, 189728000.000000},
	}
	for i, tc := range cases {
		bvh, e := buildBVH(t, i, tc.file, opts)
		if e != nil {
			continue
		}
		bb := bvh.OrthoBox()
		testBVHOrthoBox(t, 1, &tc.bb, &bb)
		if bvh.Cost() != tc.cost {
			t.Errorf("Expected a cost of %f, got %f", tc.cost, bvh.Cost())
		}
	}
}

func TestSimpleSplit(t *testing.T) {
	opts := NewBVHDefaultOptions()

	var cases = []struct {
		file string
		bb   OrthoBox
		cost float32
	}{
		// cuts of the first side
		{"paulbourke.net.sample1.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 1)}, 56},
		{"two_cubes.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(4, 1, 1)}, 112},
		{"two_cubes_y.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 4, 1)}, 112},
		{"two_cubes_2y.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 8, 1)}, 176},
		{"two_cubes_z.ply", OrthoBox{NewVec3(0, 0, 0), NewVec3(1, 1, 4)}, 112},
		{"people.sc.fsu.edu.helix.ply", OrthoBox{NewVec3(-23.75, -44.6, -30.25), NewVec3(36.75, 47.65, 30.25)}, 258309.718750},
	}
	for i, tc := range cases {
		bvh, _ := buildBVH(t, i, tc.file, opts)
		if bvh == nil {
			continue
		}
		bb := bvh.OrthoBox()
		testBVHOrthoBox(t, 1, &tc.bb, &bb)
		if bvh.Cost() != tc.cost {
			t.Errorf("tc %d: Expected a cost of %f, got %f", i, tc.cost, bvh.Cost())
		}
	}
}

func BenchmarkBVHBuilding(b *testing.B) {
	m, _ := getMesh(nil, 0, "people.sc.fsu.edu.helix.ply")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewBVHTree(m, nil)
	}
}

func testBVHOrthoBox(t *testing.T, i int, exp, cur *OrthoBox) {
	if !exp.P0.IsEqual(&cur.P0) || !exp.P1.IsEqual(&cur.P1) {
		t.Errorf("tc %d: wrong orthobox - exp: %s cur: %s", i, exp, cur)
	}
}

func buildBVH(t *testing.T, i int, filename string, opts *BVHBuildOptions) (*BVHTree, error) {
	var e error
	var m *Mesh
	if m, e = getMesh(t, i, filename); e != nil {
		return nil, e
	}
	var bvh *BVHTree
	if bvh, e = NewBVHTree(m, opts); e != nil {
		t.Errorf("tc %d: error on creating mesh: `%s`", i, e.Error())
		return nil, e
	}
	return bvh, nil
}

func getMesh(t *testing.T, i int, fileName string) (*Mesh, error) {
	var f *os.File
	var err error
	var m *Mesh
	if f, err = os.Open("test/ply/" + fileName); err != nil {
		t.Errorf("tc %d: could not open file %s: %s ", i, fileName, err.Error())
		return nil, err
	}
	if m, err = ReadPLY(f); err != nil {
		t.Errorf("tc %d: unexpected error on reading mesh: %s", i+1, err.Error())
		return nil, err
	}
	return m, nil
}

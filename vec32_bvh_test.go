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
		{"paulbourke.net.sample1.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{1, 1, 1}}, 12 * 6},
		{"two_cubes.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{4, 1, 1}}, 24 * 18},
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
		{"paulbourke.net.sample1.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{1, 1, 1}}, 64},
		{"two_cubes.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{4, 1, 1}}, 88},
		{"two_cubes_y.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{1, 4, 1}}, 128},
		{"two_cubes_2y.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{1, 8, 1}}, 192},
		{"two_cubes_z.ply", OrthoBox{Vec3{0, 0, 0}, Vec3{1, 1, 4}}, 128},
	}
	for i, tc := range cases {
		bvh, _ := buildBVH(t, i, tc.file, opts)
		if bvh == nil {
			continue
		}
		bb := bvh.OrthoBox()
		testBVHOrthoBox(t, 1, &tc.bb, &bb)
		if bvh.Cost() != tc.cost {
			t.Errorf("Expected a cost of %f, got %f", tc.cost, bvh.Cost())
		}
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

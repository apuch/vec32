package vec32

import (
	"os"
	"testing"
)

func TestBB(t *testing.T) {
	bvh, e := buildBVH(t, 0, "paulbourke.net.sample1.ply")
	if e != nil {
		return
	}
	bb := bvh.OrthoBox()
	testBVHOrthoBox(t, 1, &OrthoBox{Vec3{0, 0, 0}, Vec3{1, 1, 1}}, &bb)

}

func testBVHOrthoBox(t *testing.T, i int, exp, cur *OrthoBox) {
	if !exp.P0.IsEqual(&cur.P0) || !exp.P1.IsEqual(&cur.P1) {
		t.Errorf("tc %d: wrong orthobox - exp: %s cur: %s", i, exp, cur)
	}
}

func buildBVH(t *testing.T, i int, filename string) (*BVHTree, error) {
	var e error
	var m *Mesh
	if m, e = getMesh(t, i, filename); e != nil {
		return nil, e
	}
	var bvh *BVHTree
	if bvh, e = NewBVHTree(m, nil); e != nil {
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
	}
	return m, nil
}

package vec32

// The root structure to pass around
type BVHTree struct {
	root *bvhNode
}

// options / tweaking parameter for creating the BVH tree
type BVHBuildOptions struct {
}

type bvhNode struct {
	bb          OrthoBox
	left, right *bvhNode
	tris        []Triangle
}

type bvhBuildNode struct {
	idx int
	bb  OrthoBox
	p   Vec3
}

type bvhBuilder struct {
	nodes []bvhBuildNode
	bvh   BVHTree
	m     *Mesh
}

// create a new BVH Tree
func NewBVHTree(m *Mesh, opt *BVHBuildOptions) (*BVHTree, error) {
	bvhb := bvhBuilder{m: m}
	var e error
	if e = bvhb.createBuildNodes(); e != nil {
		return nil, e
	}
	return &bvhb.bvh, nil
}

func (bvhb *bvhBuilder) createBuildNodes() error {
	bvhb.nodes = make([]bvhBuildNode, len(bvhb.m.Tris))
	bvhb.bvh.root = &bvhNode{bb: ORTHO_EMPTY}
	for i, tri := range bvhb.m.Tris {
		tri.OrthoBox(&bvhb.nodes[i].bb)
		tri.Center(&bvhb.nodes[i].p)
		bvhb.bvh.root.bb.Add(&bvhb.nodes[i].bb)
	}
	return nil
}

// get the bounding box
func (bvh *BVHTree) OrthoBox() OrthoBox {
	return bvh.root.bb
}

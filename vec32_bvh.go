package vec32

// The root structure to pass around
type BVHTree struct {
	root *bvhNode
	Opt  *BVHBuildOptions
}

// options / tweaking parameter for creating the BVH tree
type BVHBuildOptions struct {
	TraversalCost  float32
	TrisPerNodeMin int
	TrisPerNodeMax int
}

var BVHDefaultOpt = BVHBuildOptions{
	TraversalCost:  0.0,
	TrisPerNodeMin: 1,
	TrisPerNodeMax: 1000,
}

type bvhNode struct {
	bb          OrthoBox
	left, right *bvhNode
	tris        []int
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

const BIN_COUNT = 12

// create a new BVH Tree
func NewBVHTree(m *Mesh, opt *BVHBuildOptions) (*BVHTree, error) {
	if opt == nil {
		opt = &BVHDefaultOpt
	}
	bvhb := bvhBuilder{m: m}
	bvhb.bvh.Opt = opt
	var e error
	if e = bvhb.createBuildNodes(); e != nil {
		return nil, e
	}
	return &bvhb.bvh, nil
}

func (bvhb *bvhBuilder) createBuildNodes() error {
	bvhb.nodes = make([]bvhBuildNode, len(bvhb.m.Tris))
	bvhb.bvh.root = &bvhNode{bb: ORTHO_EMPTY}
	bvhb.bvh.root.tris = make([]int, len(bvhb.m.Tris))
	for i := 0; i < len(bvhb.m.Tris); i++ {
		bvhb.bvh.root.tris[i] = i
	}
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

// get the cost of a build
//
// less is better by the way
func (bvh *BVHTree) Cost() float32 {
	return bvh.root.Cost()
}

func (n *bvhNode) Cost() float32 {
	if n.left == nil || n.right == nil {
		return n.bb.Area() * float32(len(n.tris))
	}
	return n.left.Cost() + n.right.Cost()
}

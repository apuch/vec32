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

const DEBUG_LOG = 0

// create a new BVH Tree
func NewBVHTree(m *Mesh, opt *BVHBuildOptions) (*BVHTree, error) {
	if opt == nil {
		opt = NewBVHDefaultOptions()
	}
	bvhb := bvhBuilder{m: m}
	bvhb.bvh.Opt = opt
	var e error
	if e = bvhb.createBuildNodes(); e != nil {
		return nil, e
	}
	bvhb.doSplits(bvhb.bvh.root)
	return &bvhb.bvh, nil
}

func NewBVHDefaultOptions() *BVHBuildOptions {
	return &BVHBuildOptions{
		TraversalCost:  0.0,
		TrisPerNodeMin: 1,
		TrisPerNodeMax: 999999,
	}
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

type bvhBin struct {
	bb  OrthoBox
	cnt int
}

type bvhSplit struct {
	bbLeft, bbRight OrthoBox
	dim             int
}

func (bvhb *bvhBuilder) doSplits(n *bvhNode) {
	bvhb.getSplit(n)
	if n.left != nil {
		bvhb.doSplits(n.left)
		bvhb.doSplits(n.right)
	}
}

func (bvhb *bvhBuilder) getSplit(n *bvhNode) {
	if DEBUG_LOG > 0 {
		Trace.Println("==========================================================")
		Trace.Printf("bb: %s objs: %d cost: %f", n.bb.String(), len(n.tris),
			n.bb.Area()*float32(len(n.tris)))
	}
	if DEBUG_LOG > 2 {
		for t := 0; t < len(n.tris); t++ {
			tri := bvhb.m.Tris[n.tris[t]]
			var bb OrthoBox
			tri.OrthoBox(&bb)
			Trace.Printf("tri %d: %s P0: %s", t, bb.String(),
				bvhb.nodes[n.tris[t]].p.String())
		}
	}

	if len(n.tris) < bvhb.bvh.Opt.TrisPerNodeMin {
		return
	}
	dist := n.bb.P1.Sub(&n.bb.P0)
	dimVec := getDimVec(&n.bb)
	if dimVec.LengthSq() == 0 {
		return
	}

	var bins [BIN_COUNT]bvhBin
	for i := 0; i < BIN_COUNT; i++ {
		bins[i].bb = ORTHO_EMPTY
	}

	k1 := BIN_COUNT * (1 - 10*EPS) / dist.Dot(&dimVec)
	k0 := n.bb.P0.Dot(&dimVec)

	for _, idx := range n.tris {
		bin := int(k1 * (bvhb.nodes[idx].p.Dot(&dimVec) - k0))
		bins[bin].cnt += 1
		//OrthoBoxAdd(&bins[bin].bb, &bvhb.nodes[idx].bb)
		bins[bin].bb.Add(&bvhb.nodes[idx].bb)
	}

	var costLeft, costRight [BIN_COUNT]float32
	cnt := 0
	bb := ORTHO_EMPTY
	for i := 0; i < BIN_COUNT; i++ {
		bb.Add(&bins[i].bb)
		cnt += bins[i].cnt
		if cnt > 0 {
			costLeft[i] = float32(cnt) * bb.Area()
		} else {
			costLeft[i] = 0
		}
	}

	cnt = 0
	bb = ORTHO_EMPTY
	for i := BIN_COUNT - 1; i >= 0; i-- {
		bb.Add(&bins[i].bb)
		cnt += bins[i].cnt
		if cnt > 0 {
			costRight[i] = float32(cnt) * bb.Area()
		} else {
			costRight[i] = 0
		}
	}

	if DEBUG_LOG > 1 {
		for i := 0; i < BIN_COUNT; i++ {
			Trace.Printf("Bin %d: %s * %d cost left: %f right: %f",
				i, bins[i].bb.String(), bins[i].cnt, costLeft[i], costRight[i])
		}
	}

	bestCost := INF
	bestCostIdx := 0
	for i := 0; i < BIN_COUNT-1; i++ {
		cost := costLeft[i] + costRight[i+1]
		if cost < bestCost {
			bestCost = cost
			bestCostIdx = i
		}
	}

	bestCost += bvhb.bvh.Opt.TraversalCost

	selfCost := float32(len(n.tris)) * n.bb.Area()
	if bestCost >= selfCost {
		if DEBUG_LOG > 0 {
			Trace.Println("Split isn't worth it")
		}
		return
	}

	posLeft := 0
	posRight := len(n.tris) - 1
	bbLeft := ORTHO_EMPTY
	bbRight := ORTHO_EMPTY

	for posLeft <= posRight {
		for ; posLeft < len(n.tris) && posLeft <= posRight; posLeft += 1 {
			bin := int(k1 * (bvhb.nodes[n.tris[posLeft]].p.Dot(&dimVec) - k0))
			if bin > bestCostIdx {
				break
			}
			bbLeft.Add(&bvhb.nodes[n.tris[posLeft]].bb)
		}
		for ; posRight > posLeft-1 && posRight >= 0; posRight -= 1 {
			bin := int(k1 * (bvhb.nodes[n.tris[posRight]].p.Dot(&dimVec) - k0))
			if bin <= bestCostIdx {
				break
			}
			bbRight.Add(&bvhb.nodes[n.tris[posRight]].bb)
		}
		if posLeft < posRight {
			tmp := n.tris[posLeft]
			n.tris[posLeft] = n.tris[posRight]
			n.tris[posRight] = tmp

			bbLeft.Add(&bvhb.nodes[n.tris[posLeft]].bb)
			bbRight.Add(&bvhb.nodes[n.tris[posRight]].bb)

			posLeft += 1
			posRight -= 1
		} else if posLeft == posRight {
			Error.Println("posLeft == posRight - that should not happen")
		}
	}

	n.left = &bvhNode{
		tris: n.tris[:posLeft],
		bb:   bbLeft,
	}
	n.right = &bvhNode{
		tris: n.tris[posRight+1:],
		bb:   bbRight,
	}

	if DEBUG_LOG > 0 {
		Trace.Printf("split for cost %f (%.2f%%) at %f (bin %d) "+
			"[dim: %s] - tris: %d vs. %d",
			bestCost, bestCost*100.0/selfCost, float32(bestCostIdx+1)/k1+k0, bestCostIdx,
			dimVec.String(), len(n.left.tris), len(n.right.tris))
	}
}

func getDimVec(bb *OrthoBox) Vec3 {
	dist := bb.P1.Sub(&bb.P0)
	if dist.X >= dist.Y && dist.X >= dist.Z && dist.X > 10*EPS {
		return NewVec3(1, 0, 0)
	} else if dist.Y >= dist.X && dist.Y >= dist.Z && dist.Y > 10*EPS {
		return NewVec3(0, 1, 0)
	} else if dist.Z >= dist.X && dist.Z >= dist.Y && dist.Z > 10*EPS {
		return NewVec3(0, 0, 1)
	} else {
		return NewVec3(0, 0, 0)
	}
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

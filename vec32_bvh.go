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
		TrisPerNodeMin: 0,
		TrisPerNodeMax: 0,
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
	dist := n.bb.P1.Sub(&n.bb.P0)

	var bins [BIN_COUNT]bvhBin
	for i := 0; i < BIN_COUNT; i++ {
		bins[i].bb = ORTHO_EMPTY
	}

	var k0, k1 float32

	splitX := dist.X >= dist.Y && dist.X >= dist.Z
	splitY := dist.Y >= dist.X && dist.Y >= dist.Z
	splitZ := dist.Z >= dist.X && dist.Z >= dist.Y
	if splitX {
		k1 = BIN_COUNT * (1 - 10*EPS) / dist.X
		k0 = n.bb.P0.X
		for _, idx := range n.tris {
			bin := int(k1 * (bvhb.nodes[idx].p.X - k0))
			bins[bin].cnt += 1
			bins[bin].bb.Add(&bvhb.nodes[idx].bb)
		}

	} else if splitY {
		k1 = BIN_COUNT * (1 - 10*EPS) / dist.Y
		k0 = n.bb.P0.Y
		for _, idx := range n.tris {
			bin := int(k1 * (bvhb.nodes[idx].p.Y - k0))
			bins[bin].cnt += 1
			bins[bin].bb.Add(&bvhb.nodes[idx].bb)
		}
	} else if splitZ {
		k1 = BIN_COUNT * (1 - 10*EPS) / dist.Z
		k0 = n.bb.P0.Z
		for _, idx := range n.tris {
			bin := int(k1 * (bvhb.nodes[idx].p.Z - k0))
			bins[bin].cnt += 1
			bins[bin].bb.Add(&bvhb.nodes[idx].bb)
		}
	}

	var costLeft, costRight [BIN_COUNT - 1]float32
	cnt := bins[0].cnt
	bb := bins[0].bb
	for i := 1; i < BIN_COUNT-1; i++ {
		if bins[i].cnt == 0 {
			costLeft[i] = INF
			continue
		}
		bb.Add(&bins[i].bb)
		cnt += bins[i].cnt
		costLeft[i] = float32(cnt) * bb.Area()
	}

	cnt = bins[BIN_COUNT-1].cnt
	bb = bins[BIN_COUNT-1].bb
	for i := BIN_COUNT - 2; i >= 0; i-- {
		if bins[i].cnt == 0 {
			costLeft[i] = INF
			continue
		}
		bb.Add(&bins[i].bb)
		cnt += bins[i].cnt
		costRight[i] = float32(cnt) * bb.Area()
	}

	bestCost := INF
	bestCostIdx := 0
	for i := 0; i < BIN_COUNT-1; i++ {
		cost := costLeft[i] + costRight[i]
		if cost < bestCost {
			bestCost = cost
			bestCostIdx = i
		}
	}

	bestCost += bvhb.bvh.Opt.TraversalCost

	selfCost := float32(len(n.tris)) * n.bb.Area()
	if bestCost > selfCost {
		//Trace.Printf("best cost: %f - self cost: %f -> no split", bestCost, selfCost)
		return
	}

	posLeft := 0
	posRight := len(n.tris) - 1
	if splitX {
		for posLeft <= posRight {
			for ; posLeft < len(n.tris) && posLeft <= posRight; posLeft += 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.X - k0))
				if bin > bestCostIdx {
					break
				}
			}
			for ; posRight > posLeft-1 && posRight >= 0; posRight -= 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.X - k0))
				if bin < bestCostIdx {
					break
				}
			}
			if posLeft < posRight {
				tmp := n.tris[posLeft]
				n.tris[posLeft] = n.tris[posRight]
				n.tris[posRight] = tmp
				posLeft += 1
				posRight -= 1
			}
		}
	} else if splitY {
		for posLeft <= posRight {
			for ; posLeft < len(n.tris) && posLeft <= posRight; posLeft += 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.Y - k0))
				if bin > bestCostIdx {
					break
				}
			}
			for ; posRight > posLeft-1 && posRight >= 0; posRight -= 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.Y - k0))
				if bin < bestCostIdx {
					break
				}
			}
			if posLeft < posRight {
				tmp := n.tris[posLeft]
				n.tris[posLeft] = n.tris[posRight]
				n.tris[posRight] = tmp
				posLeft += 1
				posRight -= 1
			}
		}
	} else if splitZ {
		for posLeft <= posRight {
			for ; posLeft < len(n.tris) && posLeft <= posRight; posLeft += 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.Z - k0))
				if bin > bestCostIdx {
					break
				}
			}
			for ; posRight > posLeft-1 && posRight >= 0; posRight -= 1 {
				bin := int(k1 * (bvhb.nodes[posLeft].p.Z - k0))
				if bin < bestCostIdx {
					break
				}
			}
			if posLeft < posRight {
				tmp := n.tris[posLeft]
				n.tris[posLeft] = n.tris[posRight]
				n.tris[posRight] = tmp
				posLeft += 1
				posRight -= 1
			}
		}
	}
	bbLeft := bins[0].bb
	bbRight := bins[BIN_COUNT-1].bb

	for i := 0; i <= bestCostIdx; i++ {
		bbLeft.Add(&bins[i].bb)
	}
	for i := bestCostIdx + 1; i < BIN_COUNT; i++ {
		bbRight.Add(&bins[i].bb)
	}

	n.left = &bvhNode{
		tris: n.tris[:posLeft],
		bb:   bbLeft,
	}
	n.right = &bvhNode{
		tris: n.tris[posRight+1 : len(n.tris)],
		bb:   bbRight,
	}
	//Trace.Printf("split at %f (bin %d) - tris: %d vs. %d boxes: %s vs %s",
	//	bestCost, bestCostIdx, len(n.left.tris), len(n.right.tris),
	//	bbLeft.String(), bbRight.String())
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

package layout

import (
	"math"
)

// "Fast and Simple Horizontal Coordinate Assignment" by Ulrik Brandes and Boris Kopf, 2002
// Computes horizontal coordinate in layered graph, given ordering within each layer.
// Produces result such that neighbors are close and long edges cross Layers are straight.
// Works on fully connected graphs.
// Assuming nodes do not have width.
type BrandesKopfLayersNodesHorizontalAssigner struct {
	Delta int // distance between nodes, including fake ones
}

func (s BrandesKopfLayersNodesHorizontalAssigner) NodesHorizontalCoordinates(_ Graph, g LayeredGraph) map[uint64]int {
	typeOneSegments := preprocessing(g)
	root, align := verticalAlignment(g, typeOneSegments)
	x := horizontalCompaction(g, root, align, s.Delta)
	// TODO: balancing by taking median for every node across 4 runs for each run as in algorithm
	return x
}

// Alg 1.
// Type 1 conflicts arise when a non-inner segment (normal edge) crosses an inner segment (edge between two fake nodes).
// The algorithm traverses Layers from left to right (index l) while maintaining the upper neighbors,
// v(i)_k0 and v(i)_k1, of the two closest inner Segments.
func preprocessing(g LayeredGraph) (typeOneSegments map[[2]uint64]bool) {
	typeOneSegments = map[[2]uint64]bool{}

	for i := range g.Layers() {
		if i == (len(g.Layers()) - 1) {
			continue
		}
		nextLayer := g.Layers()[i+1]

		k0 := 0
		l := 0

		for l1, v := range nextLayer {
			var upperNeighborInnerSegment uint64
			for _, u := range g.UpperNeighbors(v) {
				if g.IsInnerSegment([2]uint64{u, v}) {
					upperNeighborInnerSegment = u
					break
				}
			}

			if (l1 == (len(nextLayer) - 1)) || upperNeighborInnerSegment != 0 {
				k1 := len(g.Layers()[i]) - 1
				if upperNeighborInnerSegment != 0 {
					k1 = g.NodeYX[upperNeighborInnerSegment][1]
				}
				for l <= l1 {
					for k, u := range g.UpperNeighbors(nextLayer[l]) {
						if k < k0 || k > k1 {
							typeOneSegments[[2]uint64{u, v}] = true
						}
					}
					l += 1
				}
				k0 = k1
			}
		}
	}

	return typeOneSegments
}

// Alg 2.
// Obtain a leftmost alignment with upper neighbors.
// A maximal set of vertically aligned vertices is called a block, and we define the root of a block to be its topmost vertex.
// Blocks are stored as cyclicly linked lists, each node has reference to its lower aligned neighbor and lowest refers to topmost.
// Each node has additional reference to root of its block.
func verticalAlignment(g LayeredGraph, typeOneSegments map[[2]uint64]bool) (root map[uint64]uint64, align map[uint64]uint64) {
	root = make(map[uint64]uint64, len(g.NodeYX))
	align = make(map[uint64]uint64, len(g.NodeYX))

	for v := range g.NodeYX {
		root[v] = v
		align[v] = v
	}

	for i := range g.Layers() {
		r := 0
		for _, v := range g.Layers()[i] {
			upNeighbors := g.UpperNeighbors(v)
			if d := len(upNeighbors); d > 0 {
				for m := d / 2; m < ((d+1)/2) && (m < d); m++ {
					u := upNeighbors[m]
					if align[v] == v {
						if !typeOneSegments[[2]uint64{u, v}] && r < g.NodeYX[u][1] {
							align[u] = v
							root[v] = root[u]
							align[v] = root[v]
							r = g.NodeYX[u][1]
						}
					}
				}
			}
		}
	}

	return root, align
}

// part of Alg 3.
func placeBlock(g LayeredGraph, x map[uint64]int, root map[uint64]uint64, align map[uint64]uint64, sink map[uint64]uint64, shift map[uint64]int, delta int, v uint64) {
	if _, ok := x[v]; !ok {
		x[v] = 0
		flag := true
		for w := v; flag; flag = v != w {
			if g.NodeYX[w][1] > 0 {
				u := root[g.Layers()[g.NodeYX[w][0]][g.NodeYX[w][1]-1]]
				placeBlock(g, x, root, align, sink, shift, delta, u)
				if sink[v] == v {
					sink[v] = sink[u]
				}
				if sink[v] != sink[u] {
					if s := x[v] - x[u] - delta; s < shift[sink[u]] {
						shift[sink[u]] = s
					}
				} else {
					if s := x[u] + delta; s > x[v] {
						x[v] = s
					}
				}
			}
			w = align[w]
		}
	}
}

// Alg 3.
// All node of a block are assigned the coordinate of the root.
// Partition each block in to classes.
// Class is defined by reachable sink which has the topmost root
// Within each class, we apply a longest path layering,
// i.e. the relative coordinate of a block with respect to the defining
// sink is recursively determined to be the maximum coordinate of
// the preceding blocks in the same class, plus minimum separation.
// For each class, from top to bottom, we then compute the absolute coordinates
// of its members by placing the class with minimum separation from previously placed classes.
func horizontalCompaction(g LayeredGraph, root map[uint64]uint64, align map[uint64]uint64, delta int) (x map[uint64]int) {
	sink := map[uint64]uint64{}
	shift := map[uint64]int{}
	x = map[uint64]int{}

	for v := range g.NodeYX {
		sink[v] = v
		shift[v] = math.MaxInt
	}

	// root coordinates relative to sink
	for v := range g.NodeYX {
		if root[v] == v {
			placeBlock(g, x, root, align, sink, shift, delta, v)
		}
	}

	// absolute coordinates
	for v := range g.NodeYX {
		x[v] = x[root[v]]
		if s := shift[sink[root[v]]]; s < math.MaxInt {
			x[v] += s
		}
	}

	return x
}

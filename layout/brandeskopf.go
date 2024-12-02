package layout

import (
	"math"
	"sort"
)

// "Fast and Simple Horizontal Coordinate Assignment" by Ulrik Brandes and Boris Kopf, 2002
// Computes horizontal coordinate in layered graph, given ordering within each layer.
// Produces result such that neighbors are close and long edges cross Layers are straight.
// Works on fully connected graphs.
// Assuming nodes do not have width.
type BrandesKopfLayersNodesHorizontalAssigner struct {
	Delta int // distance between nodes, including fake ones
}

type Neighbors struct {
	Up   map[uint64][]uint64
	Down map[uint64][]uint64
}

func computeOrderedNeighbors(g LayeredGraph) Neighbors {
	n := Neighbors{
		Up:   make(map[uint64][]uint64),
		Down: make(map[uint64][]uint64),
	}

	for e := range g.Segments {
		n.Down[e[0]] = append(n.Down[e[0]], e[1])
		n.Up[e[1]] = append(n.Up[e[1]], e[0])
	}
	for _, d := range n.Down {
		sort.Slice(d, func(i, j int) bool { return g.NodeYX[d[i]][1] < g.NodeYX[d[j]][1] })
	}

	for _, d := range n.Up {
		sort.Slice(d, func(i, j int) bool { return g.NodeYX[d[i]][1] < g.NodeYX[d[j]][1] })
	}

	return n
}

type xResult struct {
	minX, maxX int
}

func (x xResult) width() int {
	return x.maxX - x.minX
}

func (s BrandesKopfLayersNodesHorizontalAssigner) NodesHorizontalCoordinates(_ Graph, g LayeredGraph) map[uint64]int {
	neighbors := computeOrderedNeighbors(g)
	typeOneSegments := preprocessing(g, neighbors)

	dirTL := TopLeft{}
	rootTL, alignTL := dirTL.verticalAlignment(g, typeOneSegments, neighbors)
	xTL := dirTL.horizontalCompaction(g, rootTL, alignTL, s.Delta)

	resTL := xResult{
		minX: math.MaxInt,
		maxX: math.MinInt,
	}
	for _, v := range xTL {
		if resTL.minX < v {
			resTL.minX = v
		}
		if resTL.maxX > v {
			resTL.maxX = v
		}
	}

	dirTR := TopRight{}
	rootTR, alignTR := dirTR.verticalAlignment(g, typeOneSegments, neighbors)
	xTR := dirTR.horizontalCompaction(g, rootTR, alignTR, s.Delta)

	resTR := xResult{
		minX: math.MaxInt,
		maxX: math.MinInt,
	}
	for _, v := range xTR {
		if resTR.minX < v {
			resTR.minX = v
		}
		if resTR.maxX > v {
			resTR.maxX = v
		}
	}

	best := resTL
	if resTR.width() < resTL.width() {
		best = resTR
	}

	shiftTL := best.minX - resTL.minX
	shiftTR := best.maxX - resTR.maxX


	x := make(map[uint64]int, len(xTL))
	for n := range g.NodeYX {
		tl := xTL[n] + shiftTL
		tr := xTR[n] + shiftTR
		x[n] = tl + tr / 2
	}

	// TODO: balancing by taking median for every node across 4 runs for each run as in algorithm

	return x
}

// Alg 1.
// Type 1 conflicts arise when a non-inner segment (normal edge) crosses an inner segment (edge between two fake nodes).
// The algorithm traverses Layers from left to right (index l) while maintaining the upper neighbors,
// v(i)_k0 and v(i)_k1, of the two closest inner Segments.
func preprocessing(g LayeredGraph, n Neighbors) (typeOneSegments map[[2]uint64]bool) {
	typeOneSegments = map[[2]uint64]bool{}

	layers := g.Layers()
	for i := range layers {
		if i == (len(layers) - 1) {
			continue
		}
		nextLayer := layers[i+1]

		k0 := 0
		l := 0

		for l1, v := range nextLayer {
			var upperNeighborInnerSegment uint64
			for _, u := range n.Up[v] {
				if g.IsInnerSegment([2]uint64{u, v}) {
					upperNeighborInnerSegment = u
					break
				}
			}

			if (l1 == (len(nextLayer) - 1)) || upperNeighborInnerSegment != 0 {
				k1 := len(layers[i]) - 1
				if upperNeighborInnerSegment != 0 {
					k1 = g.NodeYX[upperNeighborInnerSegment][1]
				}
				for l <= l1 {
					for k, u := range n.Up[nextLayer[l]] {
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

type TopLeft struct{}

// Alg 2.
// Obtain a leftmost alignment with upper neighbors.
// A maximal set of vertically aligned vertices is called a block, and we define the root of a block to be its topmost vertex.
// Blocks are stored as cyclicly linked lists, each node has reference to its lower aligned neighbor and lowest refers to topmost.
// Each node has additional reference to root of its block.
func (s TopLeft) verticalAlignment(g LayeredGraph, typeOneSegments map[[2]uint64]bool, n Neighbors) (root map[uint64]uint64, align map[uint64]uint64) {
	root = make(map[uint64]uint64, len(g.NodeYX))
	align = make(map[uint64]uint64, len(g.NodeYX))

	for v := range g.NodeYX {
		root[v] = v
		align[v] = v
	}

	layers := g.Layers()
	for i := range layers {
		r := -1
		for _, v := range layers[i] {
			upNeighbors := n.Up[v]
			if d := len(upNeighbors); d > 0 {
				for m := d / 2; m < ((d + 1) / 2); m++ {
					if align[v] == v {
						u := upNeighbors[m]
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
func (s TopLeft) placeBlock(g LayeredGraph, x map[uint64]int, root map[uint64]uint64, align map[uint64]uint64, sink map[uint64]uint64, shift map[uint64]int, delta int, v uint64, layers [][]uint64) {
	if _, ok := x[v]; !ok {
		x[v] = 0
		flag := true
		w := v
		for ; flag; flag = v != w {
			if g.NodeYX[w][1] > 0 {
				u := root[layers[g.NodeYX[w][0]][g.NodeYX[w][1]-1]]
				s.placeBlock(g, x, root, align, sink, shift, delta, u, layers)
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
		for align[w] != v {
			w = align[w]
			x[w] = x[v]
			sink[w] = sink[v]
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
func (s TopLeft) horizontalCompaction(g LayeredGraph, root map[uint64]uint64, align map[uint64]uint64, delta int) (x map[uint64]int) {
	sink := map[uint64]uint64{}
	shift := map[uint64]int{}
	x = map[uint64]int{}

	for v := range g.NodeYX {
		sink[v] = v
		shift[v] = math.MaxInt
	}

	layers := g.Layers()
	// root coordinates relative to sink
	for v := range g.NodeYX {
		if root[v] == v {
			s.placeBlock(g, x, root, align, sink, shift, delta, v, layers)
		}
	}

	// class offsets
	for i := 0; i < len(layers); i++ {
		layer := layers[i]
		vfirst := layer[0]
		if sink[vfirst] == vfirst {
			if shift[sink[vfirst]] == math.MaxInt {
				shift[sink[vfirst]] = 0
			}
			j := i
			k := 0
			for {
				v := layers[j][k]

				for align[v] != root[v] {
					v = align[v]
					j++
					if g.NodeYX[v][1] > 0 {
						u := layers[g.NodeYX[v][0]][g.NodeYX[v][1]-1]
						shifted := shift[sink[v]] + x[v] - (x[u] + delta)
						if shifted < shift[sink[u]] {
							shift[sink[u]] = shifted
						}
					}
				}
				k = g.NodeYX[v][1] + 1

				if k > len(layers[j])-1 || sink[v] != sink[layers[j][k]] {
					break
				}
			}
		}
	}

	// absolute coordinates
	for v := range g.NodeYX {
		x[v] += shift[sink[v]]
	}

	return x
}

type TopRight struct{} // RIGHT UP in paper

// Alg 2.
// Obtain a leftmost alignment with upper neighbors.
// A maximal set of vertically aligned vertices is called a block, and we define the root of a block to be its topmost vertex.
// Blocks are stored as cyclicly linked lists, each node has reference to its lower aligned neighbor and lowest refers to topmost.
// Each node has additional reference to root of its block.
func (s TopRight) verticalAlignment(g LayeredGraph, typeOneSegments map[[2]uint64]bool, n Neighbors) (root map[uint64]uint64, align map[uint64]uint64) {
	root = make(map[uint64]uint64, len(g.NodeYX))
	align = make(map[uint64]uint64, len(g.NodeYX))

	for v := range g.NodeYX {
		root[v] = v
		align[v] = v
	}

	layers := g.Layers()
	for i := range layers {
		r := math.MaxInt
		for j := len(layers[i]) - 1; j >= 0; j-- {
			v := layers[i][j]
			upNeighbors := n.Up[v]
			if d := len(upNeighbors); d > 0 {
				for m := ((d + 1) / 2) - 1; m >= d/2; m-- {
					if align[v] == v {
						u := upNeighbors[m]
						if !typeOneSegments[[2]uint64{u, v}] && r > g.NodeYX[u][1] {
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
func (s TopRight) placeBlock(g LayeredGraph, x map[uint64]int, root map[uint64]uint64, align map[uint64]uint64, sink map[uint64]uint64, shift map[uint64]int, delta int, v uint64, layers [][]uint64) {
	if _, ok := x[v]; !ok {
		x[v] = 0
		flag := true
		w := v
		for ; flag; flag = v != w {
			if g.NodeYX[w][1] < len(layers[g.NodeYX[w][0]])-1 {
				u := root[layers[g.NodeYX[w][0]][g.NodeYX[w][1]+1]]
				s.placeBlock(g, x, root, align, sink, shift, delta, u, layers)
				if sink[v] == v {
					sink[v] = sink[u]
				}
				if sink[v] != sink[u] {
					if s := x[v] + x[u] + delta; s > shift[sink[u]] {
						shift[sink[u]] = s
					}
				} else {
					if s := x[u] - delta; s < x[v] {
						x[v] = s
					}
				}
			}
			w = align[w]
		}
		for align[w] != v {
			w = align[w]
			x[w] = x[v]
			sink[w] = sink[v]
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
func (s TopRight) horizontalCompaction(g LayeredGraph, root map[uint64]uint64, align map[uint64]uint64, delta int) (x map[uint64]int) {
	sink := map[uint64]uint64{}
	shift := map[uint64]int{}
	x = map[uint64]int{}

	for v := range g.NodeYX {
		sink[v] = v
		shift[v] = math.MinInt
	}

	layers := g.Layers()
	// root coordinates relative to sink
	for v := range g.NodeYX {
		if root[v] == v {
			s.placeBlock(g, x, root, align, sink, shift, delta, v, layers)
		}
	}

	// class offsets
	for i := 0; i < len(layers); i++ {
		layer := layers[i]
		vfirst := layer[0]
		if sink[vfirst] == vfirst {
			if shift[sink[vfirst]] == math.MinInt {
				shift[sink[vfirst]] = 0
			}
			j := i
			k := len(layers[j]) - 1 // changed
			for {
				v := layers[j][k]

				for align[v] != root[v] {
					v = align[v]
					j++
					if g.NodeYX[v][1] < len(layers[j])-1 { // changed
						u := layers[g.NodeYX[v][0]][g.NodeYX[v][1]+1]     // changed
						shifted := shift[sink[v]] - x[v] + (x[u] + delta) // changed
						if shifted > shift[sink[u]] {                     // changed
							shift[sink[u]] = shifted
						}
					}
				}
				k = g.NodeYX[v][1] - 1

				if k < 0 || sink[v] != sink[layers[j][k]] {
					break
				}
			}
		}
	}

	// absolute coordinates
	for v := range g.NodeYX {
		x[v] += shift[sink[v]]
	}

	return x
}

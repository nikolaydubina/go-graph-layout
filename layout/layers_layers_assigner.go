package layout

import "fmt"

// Expects that graph g does not have cycles.
// This step creates fake nodes and splits long edges into segments.
func NewLayeredGraph(g Graph) LayeredGraph {
	nodeYX := assignLevels(g)
	edges := makeEdges(g, nodeYX)
	return LayeredGraph{
		NodeYX:   nodeYX,
		Segments: makeSegments(edges),
		Dummy:    makeDummy(edges),
		Edges:    edges,
	}
}

func maxNodeID(g Graph) uint64 {
	var maxNodeID uint64
	for e := range g.Edges {
		if e[0] > maxNodeID {
			maxNodeID = e[0]
		}
		if e[1] > maxNodeID {
			maxNodeID = e[1]
		}
	}
	return maxNodeID
}

func assignLevels(g Graph) map[uint64][2]int {
	nodeYX := make(map[uint64][2]int, len(g.Nodes))
	for _, root := range g.Roots() {
		nodeYX[root] = [2]int{0, 0}
		for que := []uint64{root}; len(que) > 0; {
			// pop
			p := que[0]
			if len(que) > 1 {
				que = que[1:]
			} else {
				que = nil
			}

			// set max depth for each child
			for e := range g.Edges {
				if parent, child := e[0], e[1]; parent == p {
					if l := nodeYX[parent][0] + 1; l > nodeYX[child][0] {
						nodeYX[child] = [2]int{l, 0}
					}
					que = append(que, child)
				}
			}
		}
	}
	return nodeYX
}

// for each long edge breaks it down to multiple segments, for short edge just adds it
func makeSegments(edges map[[2]uint64][]uint64) map[[2]uint64]bool {
	segments := map[[2]uint64]bool{}
	for e, nodes := range edges {
		switch {
		case len(nodes) == 2:
			segments[e] = true
		case len(nodes) > 2:
			for i := range nodes {
				if i == 0 {
					continue
				}
				segments[[2]uint64{nodes[i-1], nodes[i]}] = true
			}
		default:
			panic(fmt.Errorf("edge(%v) has only one node(%v) but at least 2 expected", e, nodes))
		}
	}
	return segments
}

// extracts all fake nodes for edges that are long into separate map
func makeDummy(edges map[[2]uint64][]uint64) map[uint64]bool {
	dummy := map[uint64]bool{}
	for _, nodes := range edges {
		if len(nodes) > 2 {
			for i, n := range nodes {
				if i == 0 || i == (len(nodes)-1) {
					continue
				}
				dummy[n] = true
			}
		}
	}
	return dummy
}

// makeEdges split long edges into segments and add fake nodes
// adds new fake nodes to nodeYX
func makeEdges(g Graph, nodeYX map[uint64][2]int) map[[2]uint64][]uint64 {
	edges := make(map[[2]uint64][]uint64, len(g.Edges))

	nextFakeNodeID := maxNodeID(g) + 1
	for e := range g.Edges {
		fromLayer := nodeYX[e[0]][0]
		toLayer := nodeYX[e[1]][0]

		newEdge := []uint64{}
		newEdge = append(newEdge, e[0])

		if (toLayer - fromLayer) > 1 {
			for layer := fromLayer + 1; layer < toLayer; layer++ {
				nodeYX[nextFakeNodeID] = [2]int{layer, 0}
				newEdge = append(newEdge, nextFakeNodeID)
				nextFakeNodeID++
			}
		}

		newEdge = append(newEdge, e[1])

		edges[e] = newEdge
	}

	return edges
}

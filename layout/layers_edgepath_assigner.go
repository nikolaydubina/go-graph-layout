package layout

import (
	"fmt"
)

// StraightEdgePathAssigner will check node locations for each fake/real node in path and set edge path to go through middle of it.
type StraightEdgePathAssigner struct{}

func (l StraightEdgePathAssigner) UpdateGraphLayout(g Graph, lg LayeredGraph, allNodesXY map[uint64][2]int) {
	numAssignedEdges := 0
	for e, nodes := range lg.Edges {
		if _, ok := g.Edges[e]; !ok {
			panic(fmt.Errorf("layered graph edge(%v) is not found in the original graph", e))
		}

		path := make([][2]int, len(nodes))
		for i, n := range nodes {
			xy := allNodesXY[n]
			path[i] = [2]int{
				xy[0] + (g.Nodes[n].W / 2),
				xy[1] + (g.Nodes[n].H / 2),
			}
		}

		g.Edges[e] = Edge{Path: path}
		numAssignedEdges++
	}

	if numAssignedEdges != len(g.Edges) {
		panic(fmt.Errorf("layered graph has wrong number of edges(%d) vs graph num edges (%d)", numAssignedEdges, len(g.Edges)))
	}
}

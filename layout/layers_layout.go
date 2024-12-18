package layout

type CycleRemover interface {
	RemoveCycles(g Graph)
	Restore(g Graph)
}

type NodesHorizontalCoordinatesAssigner interface {
	NodesHorizontalCoordinates(g Graph, lg LayeredGraph) map[uint64]int
}

type NodesVerticalCoordinatesAssigner interface {
	NodesVerticalCoordinates(g Graph, lg LayeredGraph) map[uint64]int
}

// Kozo Sugiyama algorithm breaks down layered graph construction in phases.
type SugiyamaLayersStrategyGraphLayout struct {
	CycleRemover                       CycleRemover
	LevelsAssigner                     func(g Graph) LayeredGraph
	OrderingAssigner                   func(g Graph, lg LayeredGraph)
	NodesHorizontalCoordinatesAssigner NodesHorizontalCoordinatesAssigner
	NodesVerticalCoordinatesAssigner   NodesVerticalCoordinatesAssigner
	EdgePathAssigner                   func(g Graph, lg LayeredGraph, allNodesXY map[uint64][2]int)
}

// UpdateGraphLayout breaks down layered graph construction in phases.
func (l SugiyamaLayersStrategyGraphLayout) UpdateGraphLayout(g Graph) {
	l.CycleRemover.RemoveCycles(g)

	lg := l.LevelsAssigner(g)
	if err := lg.Validate(); err != nil {
		panic(err)
	}

	l.OrderingAssigner(g, lg)

	nodeX := l.NodesHorizontalCoordinatesAssigner.NodesHorizontalCoordinates(g, lg)
	nodeY := l.NodesVerticalCoordinatesAssigner.NodesVerticalCoordinates(g, lg)

	// real and fake node coordinates
	allNodesXY := make(map[uint64][2]int, len(g.Nodes))
	for n := range lg.NodeYX {
		allNodesXY[n] = [2]int{nodeX[n], nodeY[n]}
	}

	// export coordinates for edges
	l.EdgePathAssigner(g, lg, allNodesXY)

	// export coordinates to real nodes
	for n, node := range g.Nodes {
		g.Nodes[n] = Node{
			XY: [2]int{nodeX[n] - node.W/2, nodeY[n] - node.H/2},
			W:  node.W,
			H:  node.H,
		}
	}

	l.CycleRemover.Restore(g)
}

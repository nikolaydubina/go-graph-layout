package layout

import (
	"math/rand"
)

// SimpleCycleRemover will keep testing for cycles, if cycle found will randomly reverse one edge in cycle.
// When restoring, will reverse previously reversed edges.
type SimpleCycleRemover struct {
	Reversed map[[2]uint64]bool
}

func NewSimpleCycleRemover() SimpleCycleRemover {
	return SimpleCycleRemover{
		Reversed: map[[2]uint64]bool{},
	}
}

func getCycleDFS(g Graph, que []uint64) []uint64 {
	if len(que) == 0 {
		return que
	}

	p := que[len(que)-1]
	for e := range g.Edges {
		if e[0] == p {
			// check if cycle
			for i, t := range que {
				if e[1] == t {
					return que[i:]
				}
			}

			// DFS deep call
			if t := getCycleDFS(g, append(que, e[1])); len(t) > 0 {
				return t
			}
		}
	}

	return nil
}

func getCycle(g Graph) []uint64 {
	for _, root := range g.Roots() {
		if t := getCycleDFS(g, []uint64{root}); len(t) > 0 {
			return t
		}
	}
	return nil
}

func reverseEdge(g Graph, e [2]uint64) {
	delete(g.Edges, e)
	g.Edges[[2]uint64{e[1], e[0]}] = Edge{}
}

func (s SimpleCycleRemover) RemoveCycles(g Graph) {
	for cycle := getCycle(g); len(cycle) > 0; cycle = getCycle(g) {
		// pick edge randomly
		i := rand.Intn(len(cycle) - 1)
		e := [2]uint64{cycle[i], cycle[i+1]}

		reverseEdge(g, e)
		s.Reversed[e] = true
	}
}

func (s SimpleCycleRemover) Restore(g Graph) {
	for e := range s.Reversed {
		reverseEdge(g, e)
		delete(s.Reversed, e)
	}
}

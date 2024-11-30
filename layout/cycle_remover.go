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

func getCycleDFS(neighbors map[uint64][]uint64, que []uint64) []uint64 {
	if len(que) == 0 {
		return que
	}

	p := que[len(que)-1]
	for _, d := range neighbors[p] {
		// check if cycle
		for i, t := range que {
			if d == t {
				return que[i:]
			}
		}

		// DFS deep call
		if t := getCycleDFS(neighbors, append(que, d)); len(t) > 0 {
			return t
		}
	}

	return nil
}

func getCycle(roots []uint64, neighbors map[uint64][]uint64) []uint64 {
	for _, root := range roots {
		if t := getCycleDFS(neighbors, []uint64{root}); len(t) > 0 {
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
	neighbors := make(map[uint64][]uint64)
	for e := range g.Edges {
		neighbors[e[0]] = append(neighbors[e[0]], e[1])
	}
	for cycle := getCycle(g.Roots(), neighbors); len(cycle) > 0; cycle = getCycle(g.Roots(), neighbors) {
		// pick edge randomly
		i := rand.Intn(len(cycle) - 1)
		e := [2]uint64{cycle[i], cycle[i+1]}

		reverseEdge(g, e)
		s.Reversed[e] = true

		neighbors[e[0]] = deleteValue(neighbors[e[0]], e[1])
		neighbors[e[1]] = append(neighbors[e[1]], e[0])
	}
}

func deleteValue(slice []uint64, value uint64) []uint64 {
	for i, n := range slice {
		if n == value {
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}

func (s SimpleCycleRemover) Restore(g Graph) {
	for e := range s.Reversed {
		reverseEdge(g, e)
		delete(s.Reversed, e)
	}
}

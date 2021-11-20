package svg

import (
	"fmt"
	"strings"
)

// Graph is rendered graph.
type Graph struct {
	ID    string
	Nodes map[uint64]Node
	Edges map[[2]uint64]Edge
}

// Render creates root svg element
func (g Graph) Render() string {
	body := []string{
		fmt.Sprintf(`<g id="%s">`, g.ID),
	}

	for _, edge := range g.Edges {
		body = append(body, edge.Render())
	}

	// draw nodes always on top of edges
	for _, node := range g.Nodes {
		body = append(body, node.Render())
	}

	body = append(body, "</g>")

	return strings.Join(body, "\n")
}

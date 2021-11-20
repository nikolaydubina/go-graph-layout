package svg

import (
	"fmt"
	"strings"
)

// Edge is polylines of straight lines going through all points.
type Edge struct {
	Path [][2]int
}

func (e Edge) Render() string {
	var points []string
	for _, point := range e.Path {
		points = append(points, fmt.Sprintf("%d,%d", point[0], point[1]))
	}
	return fmt.Sprintf(`<polyline style="fill:none;stroke-width:1;stroke:black;" points="%s"></polyline>`, strings.Join(points, " "))
}

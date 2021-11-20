package svg

import (
	"fmt"
	"strings"
)

type Renderable interface {
	Render() string
}

type SVG struct {
	ID          string
	Definitions []Renderable
	Body        Renderable
}

func (s SVG) Render() string {
	defs := make([]string, 0, len(s.Definitions))
	for _, d := range s.Definitions {
		defs = append(defs, d.Render())
	}
	return strings.Join(
		[]string{
			fmt.Sprintf(`<svg id="%s" xmlns="http://www.w3.org/2000/svg" style="width: 100%%; height: 100%%;">`, s.ID),
			`<defs>`,
			strings.Join(defs, "\n"),
			`</defs>`,
			s.Body.Render(),
			`</svg>`,
		},
		"\n",
	)
}

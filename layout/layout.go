package layout

// Layout is something that can update graph layout
type Layout interface {
	UpdateGraphLayout(g Graph)
}

// SequenceLayout applies sequence of layouts
type SequenceLayout struct {
	Layouts []Layout
}

func (s SequenceLayout) UpdateGraphLayout(g Graph) {
	for _, l := range s.Layouts {
		l.UpdateGraphLayout(g)
	}
}

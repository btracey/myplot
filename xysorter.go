package myplot

import "github.com/gonum/plot/plotter"

// Sorts by lower to higher x
type XYSorter struct {
	plotter.XYs
}

func (xy XYSorter) Len() int {
	return len(xy.XYs)
}

func (xy XYSorter) Swap(i, j int) {
	xy.XYs[i], xy.XYs[j] = xy.XYs[j], xy.XYs[i]
}

func (xy XYSorter) Less(i, j int) bool {
	return xy.XYs[i].X < xy.XYs[j].X
}

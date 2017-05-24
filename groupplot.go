package myplot

type PlotAssigner interface {
	Plot() (plotName, lineName string, pt struct{ X, Y float64 })
}

func GroupPlot(p []PlotAssigner) {

}

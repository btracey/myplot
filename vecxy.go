package myplot

type VecXY struct {
	X []float64
	Y []float64
}

func (v VecXY) Len() int {
	if len(v.X) != len(v.Y) {
		panic("length mismatch")
	}
	return len(v.X)
}

func (v VecXY) XY(i int) (x, y float64) {
	return v.X[i], v.Y[i]
}

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

type VecXYZ struct {
	X []float64
	Y []float64
	Z []float64
}

func (v VecXYZ) Len() int {
	if len(v.X) != len(v.Y) {
		panic("myplot: length mismatch")
	}
	if len(v.X) != len(v.Z) {
		panic("myplot: length mismatch")
	}
	return len(v.X)
}

func (v VecXYZ) XYZ(i int) (x, y, z float64) {
	return v.X[i], v.Y[i], v.Z[i]
}

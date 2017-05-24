package myplot

import "github.com/gonum/matrix/mat64"

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

type GridMat struct {
	x []float64
	y []float64
	z *mat64.Dense
}

func NewGridMat(x, y []float64, z *mat64.Dense) GridMat {
	return GridMat{
		x: x,
		y: y,
		z: z,
	}

}

func (g GridMat) Dims() (c, r int) {
	return len(g.x), len(g.y)
}

func (g GridMat) X(c int) float64 {
	return g.x[c]
}

func (g GridMat) Y(r int) float64 {
	return g.y[r]
}

func (g GridMat) Z(c, r int) float64 {
	return g.z.At(c, r)
}

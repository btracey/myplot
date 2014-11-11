package myplot

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"fmt"
	"math"
)

// coloredscatter implements the Plotter interface, drawing
// a bubble plot of x, y, z triples where the z value
// determines the color of the dot.
type ColoredScatter struct {
	plotter.XYZs
	plot.GlyphStyle
	Colormapper
	minZ float64
	maxZ float64
}

// NewColoredScatter returns a ColoredScatter that uses the
// default glyph style and colormap.
func NewColoredScatter(xyzs plotter.XYZer) (*ColoredScatter, error) {
	data, err := plotter.CopyXYZs(xyzs)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &ColoredScatter{}, fmt.Errorf("Must have more than zero points")
	}

	minz := data[0].Z
	maxz := data[0].Z
	for _, d := range data {
		minz = math.Min(minz, d.Z)
		maxz = math.Max(maxz, d.Z)
	}

	c := &ColoredScatter{
		XYZs:        data,
		GlyphStyle:  plotter.DefaultGlyphStyle,
		Colormapper: BlueRed(),
		//minZ: minz,
		//maxZ: maxz,
	}
	c.SetScale(minz, maxz)
	return c, nil
}

func (c *ColoredScatter) Min() float64 {
	return c.minZ
}

func (c *ColoredScatter) Max() float64 {
	return c.maxZ
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (c *ColoredScatter) DataRange() (xmin, xmax, ymin, ymax float64) {
	return plotter.XYRange(plotter.XYValues{c.XYZs})
}

func (c *ColoredScatter) SetMin(min float64) {
	c.minZ = min
	s, ok := c.Colormapper.(ScaledColormapper)
	if ok {
		s.SetScale(min, c.maxZ)
	}
}

func (c *ColoredScatter) SetMax(max float64) {
	c.maxZ = max
	s, ok := c.Colormapper.(ScaledColormapper)
	if ok {
		s.SetScale(c.minZ, max)
	}
}

func (c *ColoredScatter) SetScale(min, max float64) {
	c.maxZ = max
	c.minZ = min
	s, ok := c.Colormapper.(ScaledColormapper)
	if ok {
		s.SetScale(min, max)
	}
}

func (c *ColoredScatter) SetColormap(m Colormapper) {
	c.Colormapper = m
	s, ok := c.Colormapper.(ScaledColormapper)
	if ok {
		s.SetScale(c.minZ, c.maxZ)
	}
}

// Transform a value to an interval between 0 and 1
func BoundedNormalize(val, min, max float64) (norm float64) {
	norm = (val - min) / (max - min)
	if norm < 0 {
		norm = 0
	} else if norm > 1 {
		norm = 1
	}
	return norm
}

func (c *ColoredScatter) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	s, ok := c.Colormapper.(ScaledColormapper)
	if ok {
		s.SetScale(c.minZ, c.maxZ)
	}
	//fmt.Println("In myplot", "minz = ", c.minZ, "maxz = ", c.maxZ)
	for _, p := range c.XYZs {
		c.GlyphStyle.Color = c.Colormapper.Colormap(p.Z)
		point := plot.Point{X: trX(p.X), Y: trY(p.Y)}
		da.DrawGlyph(c.GlyphStyle, point)
	}
}

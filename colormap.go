package myplot

import "image/color"

// Maps a value to a color
type Colormapper interface {
	Colormap(float64) color.Color
}

// Single Color is a colormap which returns a single color regardless
// of the value of the input
type Uniform struct {
	Value color.Color
}

// Returns the color present in the struct
func (s *Uniform) Colormap(v float64) color.Color {
	return s.Value
}

// Maps a value in [0,1] to a color
type ScaledColormapper interface {
	Colormapper
	SetScale(min float64, max float64)
}

// Implementation of the jet colormap
type Jet struct {
	// reference: http://www.metastine.com/?p=7
	min float64
	max float64
}

func (c *Jet) SetScale(min, max float64) {
	c.min = min
	c.max = max
}

func (c *Jet) Scale() (float64, float64) {
	return c.min, c.max
}

func (c *Jet) SetMax(val float64) {
	c.max = val
}

func (c *Jet) SetMin(val float64) {
	c.min = val
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func (c *Jet) Colormap(z float64) color.Color {
	v := BoundedNormalize(z, c.min, c.max)

	fourvalue := 4 * v
	var red float64
	red1 := fourvalue - 1.5
	red2 := -fourvalue + 4.5
	if red1 < red2 {
		red = red1
	} else {
		red = red2
	}
	red = clamp(red)

	var blue float64
	blue1 := fourvalue - 0.5
	blue2 := -fourvalue + 3.5
	if blue1 < blue2 {
		blue = blue1
	} else {
		blue = blue2
	}
	blue = clamp(blue)

	var green float64
	green1 := fourvalue + 0.5
	green2 := -fourvalue + 2.5
	if green1 < green2 {
		green = green1
	} else {
		green = green2
	}
	green = clamp(green)

	red8 := uint8(255 * red)
	blue8 := uint8(255 * blue)
	green8 := uint8(255 * green)
	return color.RGBA{red8, blue8, green8, 255}
}

// Linearly maps the colors between light gray and black
type Grayscale struct {
	min      float64
	max      float64
	Inverted bool // Flip the direction of the maximum
}

func (c *Grayscale) SetScale(min, max float64) {
	c.min = min
	c.max = max
}

func (c *Grayscale) Scale() (float64, float64) {
	return c.min, c.max
}

func (c *Grayscale) SetMax(val float64) {
	c.max = val
}

func (c *Grayscale) SetMin(val float64) {
	c.min = val
}

func (c *Grayscale) Colormap(z float64) color.Color {
	v := BoundedNormalize(z, c.min, c.max)
	var val float64
	if c.Inverted {
		val = v
	} else {
		val = (1 - v)
	}
	val *= 255 * 0.9
	u8v := uint8(val)
	return color.RGBA{u8v, u8v, u8v, 255}
}

func BlueRed() *Diverging {
	d := &Diverging{}
	blue := color.RGBA{R: 59, B: 192, G: 76, A: 255}
	red := color.RGBA{R: 180, B: 38, G: 4, A: 255}
	d.SetColors(blue, red)
	//d.Low = blue
	//d.high = red
	return d
}

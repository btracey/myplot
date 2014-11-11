package myplot

import (
	"image/color"
	"math"
)

// Adopted using
// http://www.sandia.gov/~kmorel/documents/ColorMaps/ColorMapsExpanded.pdf
// and
// http://www.sandia.gov/~kmorel/documents/ColorMaps/diverging_map.m

const ref_X = 0.95047
const ref_Y = 1.000
const ref_Z = 1.08883

// Comes from matlab, so RGB need to be between 0 and 1
type RGBFloat struct {
	R, G, B, Alpha float64
}

func NewRGBFloat(c color.Color) *RGBFloat {
	r, g, b, a := c.RGBA()
	return &RGBFloat{R: float64(r) / 65535, G: float64(g) / 65535, B: float64(b) / 65535, Alpha: float64(a) / 65535}
}

func (rgb *RGBFloat) RGBA() (r, g, b, a uint32) {
	return uint32(rgb.R * 65535), uint32(rgb.G * 65535), uint32(rgb.B * 65535), uint32(rgb.Alpha * 65535)
}

// Converts RGB to XYZ
// Note from original author:
// The following performs a "gamma correction" specified by the sRGB color
// space.  sRGB is defined by a canonical definition of a display monitor and
// has been standardized by the International Electrotechnical Commission (IEC
// 61966-2-1).  The nonlinearity of the correction is designed to make the
// colors more perceptually uniform.  This color space has been adopted by
// several applications including Adobe Photoshop and Microsoft Windows color
// management.  OpenGL is agnostic on its RGB color space, but it is reasonable
// to assume it is close to this one.
func (rgb *RGBFloat) XYZ() *XYZFloat {
	r := rgb.R
	g := rgb.G
	b := rgb.B
	if r > 0.04045 {
		r = math.Pow((r+0.055)/1.055, 2.4)
	} else {
		r /= 12.92
	}

	if g > 0.04045 {
		g = math.Pow((g+0.055)/1.055, 2.4)
	} else {
		g /= 12.92
	}
	if b > 0.04045 {
		b = math.Pow((b+0.055)/1.055, 2.4)
	} else {
		b /= 12.92
	}

	//Observer. = 2 deg, Illuminant = D65
	x := r*0.4124 + g*0.3576 + b*0.1805
	y := r*0.2126 + g*0.7152 + b*0.0722
	z := r*0.0193 + g*0.1192 + b*0.9505

	return &XYZFloat{X: x, Y: y, Z: z, Alpha: rgb.Alpha}
}

func (rgb *RGBFloat) Lab() *LabFloat {
	return rgb.XYZ().Lab()
}

func (rgb *RGBFloat) MSH() *MSHFloat {
	return rgb.Lab().MSH()
}

type XYZFloat struct {
	X, Y, Z, Alpha float64
}

func NewXYZFloat(c color.Color) *XYZFloat {
	rgb := NewRGBFloat(c)
	return rgb.XYZ()
}

func (xyz *XYZFloat) RGBA() (r, g, b, a uint32) {
	return xyz.RGB().RGBA()
}

func (xyz *XYZFloat) Lab() *LabFloat {
	x := xyz.X
	y := xyz.Y
	z := xyz.Z

	var_X := x / ref_X
	var_Y := y / ref_Y
	var_Z := z / ref_Z

	if var_X > 0.008856 {
		var_X = math.Pow(var_X, 1.0/3.0)
	} else {
		var_X = (7.787 * var_X) + (16.0 / 116.0)
	}

	if var_Y > 0.008856 {
		var_Y = math.Pow(var_Y, 1.0/3)
	} else {
		var_Y = (7.787 * var_Y) + (16.0 / 116.0)
	}

	if var_Z > 0.008856 {
		var_Z = math.Pow(var_Z, 1.0/3)
	} else {
		var_Z = (7.787 * var_Z) + (16.0 / 116.0)
	}

	L := (116 * var_Y) - 16
	a := 500 * (var_X - var_Y)
	b := 200 * (var_Y - var_Z)

	return &LabFloat{L: L, A: a, B: b, Alpha: xyz.Alpha}
}

func (xyz *XYZFloat) RGB() *RGBFloat {

	x := xyz.X
	y := xyz.Y
	z := xyz.Z
	r := x*3.24063 + y*-1.53721 + z*-0.498629
	g := x*-0.968931 + y*1.87576 + z*0.0415175
	b := x*0.0557101 + y*-0.204021 + z*1.0570

	// The following performs a "gamma correction" specified by the sRGB color
	// space.  sRGB is defined by a canonical definition of a display monitor and
	// has been standardized by the International Electrotechnical Commission (IEC
	// 61966-2-1).  The nonlinearity of the correction is designed to make the
	// colors more perceptually uniform.  This color space has been adopted by
	// several applications including Adobe Photoshop and Microsoft Windows color
	// management.  OpenGL is agnostic on its RGB color space, but it is reasonable
	// to assume it is close to this one.
	if r > 0.0031308 {
		r = 1.055*math.Pow(r, (1/2.4)) - 0.055
	} else {
		r *= 12.92
	}
	if g > 0.0031308 {
		g = 1.055*math.Pow(g, (1/2.4)) - 0.055
	} else {
		g *= 12.92
	}
	if b > 0.0031308 {
		b = 1.055*math.Pow(b, (1/2.4)) - 0.055
	} else {
		b *= 12.92
	}

	// Clip colors. ideally we would do something that is perceptually closest
	// (since we can see colors outside of the display gamut), but this seems to
	// work well enough.
	maxVal := r
	if maxVal < g {
		maxVal = g
	}
	if maxVal < b {
		maxVal = b
	}
	if maxVal > 1.0 {
		r = r / maxVal
		g = g / maxVal
		b = b / maxVal
	}
	if r < 0 {
		r = 0
	}
	if g < 0 {
		g = 0
	}
	if b < 0 {
		b = 0
	}

	return &RGBFloat{R: r, G: g, B: b, Alpha: xyz.Alpha}
}

type LabFloat struct {
	L, A, B, Alpha float64
}

func NewLabFloat(c color.Color) *LabFloat {
	rgb := NewRGBFloat(c)
	return rgb.XYZ().Lab()
}

func (lab *LabFloat) RGBA() (r, g, b, a uint32) {
	return lab.RGB().RGBA()
}

func (lab *LabFloat) RGB() *RGBFloat {
	return lab.XYZ().RGB()
}

func (lab *LabFloat) MSH() *MSHFloat {
	L := lab.L
	a := lab.A
	b := lab.B
	M := math.Sqrt(L*L + a*a + b*b)
	var s, h float64
	if M > 0.001 {
		s = math.Acos(L / M)
	}
	if s > 0.001 {
		h = math.Atan2(b, a)
	}
	return &MSHFloat{M: M, S: s, H: h, Alpha: lab.Alpha}
}

func (lab *LabFloat) XYZ() *XYZFloat {
	L := lab.L
	a := lab.A
	b := lab.B

	var_Y := (L + 16) / 116
	var_X := a/500 + var_Y
	var_Z := var_Y - b/200

	if math.Pow(var_Y, 3) > 0.008856 {
		var_Y = math.Pow(var_Y, 3)
	} else {
		var_Y = (var_Y - 16.0/116.0) / 7.787
	}

	if math.Pow(var_X, 3) > 0.008856 {
		var_X = math.Pow(var_X, 3)
	} else {
		var_X = (var_X - 16.0/116.0) / 7.787
	}

	if math.Pow(var_Z, 3) > 0.008856 {
		var_Z = math.Pow(var_Z, 3)
	} else {
		var_Z = (var_Z - 16.0/116.0) / 7.787
	}

	x := ref_X * var_X
	y := ref_Y * var_Y
	z := ref_Z * var_Z
	return &XYZFloat{X: x, Y: y, Z: z, Alpha: lab.Alpha}
}

type MSHFloat struct {
	M, S, H, Alpha float64
}

func NewMSHFloat(c color.Color) *MSHFloat {
	rgb := NewRGBFloat(c)
	return rgb.MSH()
}

func (m *MSHFloat) RGBA() (r, g, b, a uint32) {
	return m.RGB().RGBA()
}

func (m *MSHFloat) RGB() *RGBFloat {
	return m.Lab().XYZ().RGB()
}

func (m *MSHFloat) Lab() *LabFloat {
	return &LabFloat{
		L:     m.M * math.Cos(m.S),
		A:     m.M * math.Sin(m.S) * math.Cos(m.H),
		B:     m.M * math.Sin(m.S) * math.Sin(m.H),
		Alpha: m.Alpha,
	}
}

type Diverging struct {
	Low  *MSHFloat // Low valued color
	High *MSHFloat // High valued color
	min  float64   // Min value of the data
	max  float64   // Max value of the data
	// Value from which the data is interpolated in either direction. Set automatically
	// as the average unless set by the user (to emphasize a certain color for example)
	midpoint float64
	midSet   bool
}

func (d *Diverging) SetColors(low, high color.Color) {
	d.Low = NewMSHFloat(low)
	d.High = NewMSHFloat(high)
}

// Set the min and max values for the scale
func (d *Diverging) SetScale(min, max float64) {
	d.min = min
	d.max = max
	if !d.midSet {
		d.midpoint = (max - min) / 2
	}
}

func (d *Diverging) SetMidpoint(mid float64) {
	d.midpoint = mid
	d.midSet = true
}

func (d *Diverging) ResetMidpoint() {
	d.midpoint = (d.max + d.min) / 2
	d.midSet = false
}

func (d *Diverging) Colormap(v float64) color.Color {
	lowTmp := &MSHFloat{M: d.Low.M, S: d.Low.S, H: d.Low.H, Alpha: d.Low.Alpha}
	highTmp := &MSHFloat{M: d.High.M, S: d.High.S, H: d.High.H, Alpha: d.High.Alpha}
	if lowTmp.S > 0.05 && highTmp.S > 0.05 {
		Mmid := math.Max(lowTmp.M, highTmp.M)
		Mmid = math.Max(88.0, Mmid)
		if v < d.midpoint {
			highTmp.M = Mmid
			highTmp.S = 0
			highTmp.H = 0
		} else {
			lowTmp.M = Mmid
			lowTmp.S = 0
			lowTmp.H = 0
		}
	}

	if lowTmp.S < 0.05 && highTmp.S > 0.05 {
		lowTmp.H = AdjustHue(highTmp, lowTmp.M)
	} else if highTmp.S < 0.05 && lowTmp.S > 0.05 {
		highTmp.H = AdjustHue(lowTmp, highTmp.M)
	}

	if v < d.midpoint {
		v = (v - d.min) / (d.midpoint - d.min)
	} else {
		v = (v - d.midpoint) / (d.max - d.midpoint)
	}

	c := &MSHFloat{M: (1-v)*lowTmp.M + v*highTmp.M,
		H:     (1-v)*lowTmp.H + v*highTmp.H,
		S:     (1-v)*lowTmp.S + v*highTmp.S,
		Alpha: (1-v)*lowTmp.Alpha + v*highTmp.Alpha,
	}
	return c
}

// Colormap that blends from low color to high color through white
// Assumes both colors are saturated (unlike other codes)

func AdjustHue(msh *MSHFloat, unsatM float64) (h float64) {
	if msh.M >= unsatM-0.1 {
		//The best we can do is hold hue constant.
		h = msh.H
	} else {
		// This equation is designed to make the perceptual change of the
		// interpolation to be close to constant.
		hueSpin := msh.S * math.Sqrt(unsatM*unsatM-msh.M*msh.M) / (msh.M * math.Sin(msh.S))

		// Spin hue away from 0 except in purple hues.
		if msh.H > -math.Pi/3 {
			h = msh.H + hueSpin
		} else {
			h = msh.H - hueSpin
		}
	}
	return h
}

// Given two angular orientations, returns the smallest angle between the two.
func AngleDiff(a, b float64) float64 {
	ca := math.Cos(a)
	sa := math.Sin(a)
	cb := math.Cos(b)
	sb := math.Sin(b)

	return math.Acos(ca*cb + sa*sb)
}

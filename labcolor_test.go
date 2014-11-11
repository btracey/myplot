package myplot

import (
	"fmt"
	"image/color"
	"math"
	"testing"
)

//fmt.Println("a")

type testcolor struct {
	color.Color
	name string
}

const FLOATTOL = 1E-5 // Not sure why this needs to be so high to pass the tests, aren't doing single precision arithmetic...

//var Blue = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
//var Red = &color.RGBA{R: 255, G: 0, B: 0, A: 255}
//var Random = &color.RGBA{R: 126, G: 234, B: 15, A: 255}

func GetColors() []*testcolor {
	vec := make([]*testcolor, 0, 3)

	c := &testcolor{Color: &color.RGBA{R: 0, G: 0, B: 255, A: 255},
		name: "PureBlue",
	}
	vec = append(vec, c)
	c = &testcolor{Color: &color.RGBA{R: 255, G: 0, B: 0, A: 255},
		name: "PureRed",
	}
	c = &testcolor{Color: &color.RGBA{R: 0, G: 255, B: 0, A: 255},
		name: "PureGreen",
	}
	vec = append(vec, c)
	c = &testcolor{Color: &color.RGBA{R: 126, G: 234, B: 15, A: 255},
		name: "Random1",
	}
	vec = append(vec, c)
	c = &testcolor{Color: &color.RGBA{R: 226, G: 34, B: 115, A: 112},
		name: "Random2",
	}
	vec = append(vec, c)
	return vec
}

// Checks if two colors are the same
func CheckColorMatch(color *testcolor, test color.Color, t *testing.T) {
	r, g, b, a := color.RGBA()
	//fmt.Println("real", r, g, b, a)
	r2, g2, b2, a2 := test.RGBA()
	//fmt.Println("pred", r2, g2, b2, a2)

	if r2 != r {
		t.Errorf("r doesn't match for %v. %v expected, %v found", color.name, r, r2)
	}
	if g2 != g {
		t.Errorf("g doesn't match for %v. %v expected, %v found", color.name, g, g2)
	}
	if b2 != b {
		t.Errorf("b doesn't match for %v. %v expected, %v found", color.name, b, b2)
	}
	if a2 != a {
		t.Errorf("a doesn't match for %v. %v expected, %v found", color.name, a, a2)
	}
}

// checks if two RGBFloats are the same
func CheckFloatMatch(a, b *RGBFloat, name string, t *testing.T) {
	if math.Abs(a.R-b.R) > FLOATTOL {
		t.Errorf("r float doesn't match for %v. %v expected, %v found", name, a.R, b.R)
	}
	if math.Abs(a.G-b.G) > FLOATTOL {
		t.Errorf("g float doesn't match for %v. %v expected, %v found", name, a.G, b.G)
	}
	if math.Abs(a.B-b.B) > FLOATTOL {
		t.Errorf("b float doesn't match for %v. %v expected, %v found", name, a.B, b.B)
	}
	if a.Alpha != b.Alpha {
		t.Errorf("a float doesn't match for %v. %v expected, %v found", name, a.Alpha, b.Alpha)
	}
}

// Test if can create an RGBFloat and convert back
func TestRGBFloatConversion(t *testing.T) {
	colors := GetColors()
	for _, color := range colors {
		rgbfloat := NewRGBFloat(color)
		CheckColorMatch(color, rgbfloat, t)
	}
}

// Test if can create a XYZ and convert back
func TestXYZFloatConversion(t *testing.T) {
	colors := GetColors()
	for _, color := range colors {
		rgbfloat := NewRGBFloat(color)
		xyzfloat := NewXYZFloat(color)
		CheckFloatMatch(rgbfloat, xyzfloat.RGB(), color.name, t)
		CheckColorMatch(color, xyzfloat, t)
	}
}

// Test if can create Lab and convert back
func TestLabFloatConversion(t *testing.T) {
	colors := GetColors()
	for _, color := range colors {
		rgbfloat := NewRGBFloat(color)
		labfloat := NewLabFloat(color)
		CheckFloatMatch(rgbfloat, labfloat.RGB(), color.name, t)
		CheckColorMatch(color, labfloat, t)
	}
}

func TestMSHFloatConversion(t *testing.T) {
	colors := GetColors()
	for _, color := range colors {
		rgbfloat := NewRGBFloat(color)
		mshfloat := NewMSHFloat(color)
		CheckFloatMatch(rgbfloat, mshfloat.RGB(), color.name, t)
		CheckColorMatch(color, mshfloat, t)
	}
}

func TestDivergent(t *testing.T) {
	d := &Diverging{}
	blue := color.RGBA{R: 59, B: 192, G: 76, A: 255}
	red := color.RGBA{R: 180, B: 38, G: 4, A: 255}
	v := 0.8
	d.SetScale(0, 1)
	c := d.Color(v)
	cfloat := NewRGBFloat(c)
	realans := &RGBFloat{R: 0.791805202746456, G: 0.672320372863664, B: 1}
	CheckFloatMatch(realans, cfloat, "0.8", t)
	v := 0.3
	c = d.Color(v)
	cfloat = NewRGBFloat(c)
	realans = &RGBFloat{R: 0.791805202746456, G: 0.672320372863664, B: 1}
	CheckFloatMatch(realans, cfloat, "0.8", t)
}

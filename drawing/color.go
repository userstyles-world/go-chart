package drawing

import (
	"strconv"
)

var (
	// ColorTransparent is a fully transparent color.
	ColorTransparent = Color{}

	// ColorWhite is white.
	ColorWhite = Color{R: 255, G: 255, B: 255, A: 255}

	// ColorBlack is black.
	ColorBlack = Color{R: 0, G: 0, B: 0, A: 255}

	// ColorRed is red.
	ColorRed = Color{R: 255, G: 0, B: 0, A: 255}

	// ColorGreen is green.
	ColorGreen = Color{R: 0, G: 255, B: 0, A: 255}

	// ColorBlue is blue.
	ColorBlue = Color{R: 0, G: 0, B: 255, A: 255}
)

func parseHex(hex string) uint8 {
	v, _ := strconv.ParseInt(hex, 16, 16)
	return uint8(v)
}

// ColorFromHex returns a color from a css hex code.
func ColorFromHex(hex string) Color {
	var c Color
	if len(hex) == 3 {
		c.R = parseHex(string(hex[0])) * 0x11
		c.G = parseHex(string(hex[1])) * 0x11
		c.B = parseHex(string(hex[2])) * 0x11
	} else {
		c.R = parseHex(hex[0:2])
		c.G = parseHex(hex[2:4])
		c.B = parseHex(hex[4:6])
	}
	c.A = 255
	return c
}

// ColorFromAlphaMixedRGBA returns the system alpha mixed rgba values.
func ColorFromAlphaMixedRGBA(r, g, b, a uint32) Color {
	fa := float64(a) / 255.0
	var c Color
	c.R = uint8(float64(r) / fa)
	c.G = uint8(float64(g) / fa)
	c.B = uint8(float64(b) / fa)
	c.A = uint8(a | (a >> 8))
	return c
}

// ColorChannelFromFloat returns a normalized byte from a given float value.
func ColorChannelFromFloat(v float64) uint8 {
	return uint8(v * 255)
}

// Color is our internal color type because color.Color is bullshit.
type Color struct {
	R, G, B, A uint8
}

// RGBA returns the color as a pre-alpha mixed color set.
func (c Color) RGBA() (r, g, b, a uint32) {
	fa := float64(c.A) / 255.0
	r = uint32(float64(uint32(c.R)) * fa)
	r |= r << 8
	g = uint32(float64(uint32(c.G)) * fa)
	g |= g << 8
	b = uint32(float64(uint32(c.B)) * fa)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// IsZero returns if the color has been set or not.
func (c Color) IsZero() bool {
	return c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0
}

// IsTransparent returns if the colors alpha channel is zero.
func (c Color) IsTransparent() bool {
	return c.A == 0
}

// WithAlpha returns a copy of the color with a given alpha.
func (c Color) WithAlpha(a uint8) Color {
	return Color{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

// Equals returns true if the color equals another.
func (c Color) Equals(other Color) bool {
	return c.R == other.R &&
		c.G == other.G &&
		c.B == other.B &&
		c.A == other.A
}

// AverageWith averages two colors.
func (c Color) AverageWith(other Color) Color {
	return Color{
		R: (c.R + other.R) >> 1,
		G: (c.G + other.G) >> 1,
		B: (c.B + other.B) >> 1,
		A: c.A,
	}
}

var cCache = make(map[Color]string)

// String returns a css string representation of the color.
func (c Color) String() string {
	if a, ok := cCache[c]; ok {
		return a
	}
	fa := float64(c.A) / float64(255)
	result := "rgba(" + strconv.FormatUint(uint64(c.R), 10) + "," + strconv.FormatUint(uint64(c.G), 10) + "," + strconv.FormatUint(uint64(c.B), 10) + "," + fastFloater(fa) + ")"
	cCache[c] = result
	return result
}

// F is between 0 and 1. 0 is fully transparent and 1 is fully opaque.
// We should return 0.1, 0.2, 0.3 etc.
func fastFloater(f float64) string {
	if f > 0.95 {
		return "1"
	} else if f < 0.05 {
		return "0"
	} else if f >= 0.85 && f < 0.95 {
		return "0.9"
	} else if f >= 0.75 && f < 0.85 {
		return "0.8"
	} else if f >= 0.65 && f < 0.75 {
		return "0.7"
	} else if f >= 0.55 && f < 0.65 {
		return "0.6"
	} else if f >= 0.45 && f < 0.55 {
		return "0.5"
	} else if f >= 0.35 && f < 0.45 {
		return "0.4"
	} else if f >= 0.25 && f < 0.35 {
		return "0.3"
	} else if f >= 0.15 && f < 0.25 {
		return "0.2"
	} else {
		return "0.1"
	}
}

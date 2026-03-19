package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents an RGBA color with uint8 components.
type Color struct {
	R, G, B, A uint8
}

// RGB creates a Color from RGB components with full opacity.
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// RGBA creates a Color from RGBA components.
func RGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Hex returns the color as "#rrggbb".
func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// HexA returns the color as "#rrggbbaa".
func (c Color) HexA() string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

// RGBAString returns "rgba(r,g,b,a)" for SVG usage.
func (c Color) RGBAString() string {
	alpha := float64(c.A) / 255.0
	return fmt.Sprintf("rgba(%d,%d,%d,%.3g)", c.R, c.G, c.B, alpha)
}

// WithAlpha returns a new Color with the given alpha (0.0-1.0).
func (c Color) WithAlpha(a float64) Color {
	a = math.Max(0, math.Min(1, a))
	return Color{R: c.R, G: c.G, B: c.B, A: uint8(a * 255)}
}

// Lighten mixes the color with white by the given amount (0.0-1.0).
func (c Color) Lighten(amount float64) Color {
	amount = math.Max(0, math.Min(1, amount))
	mix := func(v uint8, target uint8) uint8 {
		return uint8(float64(v) + (float64(target)-float64(v))*amount)
	}
	return Color{
		R: mix(c.R, 255),
		G: mix(c.G, 255),
		B: mix(c.B, 255),
		A: c.A,
	}
}

// Darken mixes the color with black by the given amount (0.0-1.0).
func (c Color) Darken(amount float64) Color {
	amount = math.Max(0, math.Min(1, amount))
	mix := func(v uint8) uint8 {
		return uint8(float64(v) * (1.0 - amount))
	}
	return Color{
		R: mix(c.R),
		G: mix(c.G),
		B: mix(c.B),
		A: c.A,
	}
}

// Parse parses a color string in "#rrggbb", "#rgb", or "#rrggbbaa" format.
// Also accepts named colors. Returns black on error.
func Parse(s string) Color {
	s = strings.TrimSpace(strings.ToLower(s))

	// Check named colors
	if c, ok := namedColors[s]; ok {
		return c
	}

	if !strings.HasPrefix(s, "#") {
		return Color{A: 255}
	}

	hex := s[1:]
	switch len(hex) {
	case 3:
		// #rgb -> #rrggbb
		r, _ := strconv.ParseUint(string([]byte{hex[0], hex[0]}), 16, 8)
		g, _ := strconv.ParseUint(string([]byte{hex[1], hex[1]}), 16, 8)
		b, _ := strconv.ParseUint(string([]byte{hex[2], hex[2]}), 16, 8)
		return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	case 6:
		r, _ := strconv.ParseUint(hex[0:2], 16, 8)
		g, _ := strconv.ParseUint(hex[2:4], 16, 8)
		b, _ := strconv.ParseUint(hex[4:6], 16, 8)
		return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	case 8:
		r, _ := strconv.ParseUint(hex[0:2], 16, 8)
		g, _ := strconv.ParseUint(hex[2:4], 16, 8)
		b, _ := strconv.ParseUint(hex[4:6], 16, 8)
		a, _ := strconv.ParseUint(hex[6:8], 16, 8)
		return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}

	return Color{A: 255}
}

// IsTransparent returns true if the color has zero alpha.
func (c Color) IsTransparent() bool {
	return c.A == 0
}

// SVGColor returns the best SVG color representation.
// If fully opaque, returns hex. If transparent, returns "none".
// Otherwise returns rgba().
func (c Color) SVGColor() string {
	if c.A == 0 {
		return "none"
	}
	if c.A == 255 {
		return c.Hex()
	}
	return c.RGBAString()
}

// Named colors map
var namedColors = map[string]Color{
	"white":       {255, 255, 255, 255},
	"black":       {0, 0, 0, 255},
	"red":         {255, 0, 0, 255},
	"green":       {0, 128, 0, 255},
	"blue":        {0, 0, 255, 255},
	"steelblue":   {70, 130, 180, 255},
	"coral":       {255, 127, 80, 255},
	"orange":      {255, 165, 0, 255},
	"purple":      {128, 0, 128, 255},
	"gray":        {128, 128, 128, 255},
	"grey":        {128, 128, 128, 255},
	"lightgray":   {211, 211, 211, 255},
	"lightgrey":   {211, 211, 211, 255},
	"darkgray":    {169, 169, 169, 255},
	"darkgrey":    {169, 169, 169, 255},
	"transparent": {0, 0, 0, 0},
}

// Predefined colors for convenience
var (
	White       = Color{255, 255, 255, 255}
	Black       = Color{0, 0, 0, 255}
	Red         = Color{255, 0, 0, 255}
	Green       = Color{0, 128, 0, 255}
	Blue        = Color{0, 0, 255, 255}
	SteelBlue   = Color{70, 130, 180, 255}
	Coral       = Color{255, 127, 80, 255}
	Orange      = Color{255, 165, 0, 255}
	Purple      = Color{128, 0, 128, 255}
	Gray        = Color{128, 128, 128, 255}
	LightGray   = Color{211, 211, 211, 255}
	DarkGray    = Color{169, 169, 169, 255}
	Transparent = Color{0, 0, 0, 0}
)

// Package colormap provides color maps for mapping scalar values to colors.
package colormap

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
)

// ColorMap maps a normalized value t ∈ [0, 1] to a Color.
type ColorMap func(t float64) color.Color

// keyframe is a color stop at position t ∈ [0, 1].
type keyframe struct {
	t   float64
	col color.Color
}

// interpolate performs piecewise linear interpolation between keyframes.
func interpolate(t float64, stops []keyframe) color.Color {
	t = math.Max(0, math.Min(1, t))
	if len(stops) == 0 {
		return color.Black
	}
	if len(stops) == 1 {
		return stops[0].col
	}
	// Find the surrounding pair
	for i := 1; i < len(stops); i++ {
		if t <= stops[i].t {
			lo := stops[i-1]
			hi := stops[i]
			span := hi.t - lo.t
			if span == 0 {
				return lo.col
			}
			f := (t - lo.t) / span
			r := uint8(float64(lo.col.R) + f*float64(int(hi.col.R)-int(lo.col.R)))
			g := uint8(float64(lo.col.G) + f*float64(int(hi.col.G)-int(lo.col.G)))
			b := uint8(float64(lo.col.B) + f*float64(int(hi.col.B)-int(lo.col.B)))
			return color.RGB(r, g, b)
		}
	}
	return stops[len(stops)-1].col
}

// Viridis is a perceptually uniform colormap from dark purple to yellow.
var Viridis ColorMap = func(t float64) color.Color {
	stops := []keyframe{
		{0.00, color.RGB(68, 1, 84)},
		{0.25, color.RGB(59, 82, 139)},
		{0.50, color.RGB(33, 145, 140)},
		{0.75, color.RGB(94, 201, 98)},
		{1.00, color.RGB(253, 231, 37)},
	}
	return interpolate(t, stops)
}

// Plasma is a perceptually uniform colormap from dark blue to yellow.
var Plasma ColorMap = func(t float64) color.Color {
	stops := []keyframe{
		{0.00, color.RGB(13, 8, 135)},
		{0.25, color.RGB(126, 3, 168)},
		{0.50, color.RGB(204, 71, 120)},
		{0.75, color.RGB(248, 149, 64)},
		{1.00, color.RGB(240, 249, 33)},
	}
	return interpolate(t, stops)
}

// RdBu is a diverging colormap from red (low) through white (zero) to blue (high).
// For diverging scales, use DivergingScale option so zero is centered.
var RdBu ColorMap = func(t float64) color.Color {
	stops := []keyframe{
		{0.00, color.RGB(178, 24, 43)},
		{0.25, color.RGB(239, 138, 98)},
		{0.50, color.RGB(247, 247, 247)},
		{0.75, color.RGB(103, 169, 207)},
		{1.00, color.RGB(33, 102, 172)},
	}
	return interpolate(t, stops)
}

// Greys is a sequential colormap from white to black.
var Greys ColorMap = func(t float64) color.Color {
	stops := []keyframe{
		{0.00, color.RGB(255, 255, 255)},
		{1.00, color.RGB(50, 50, 50)},
	}
	return interpolate(t, stops)
}

// Blues is a sequential colormap from light blue to dark blue.
var Blues ColorMap = func(t float64) color.Color {
	stops := []keyframe{
		{0.00, color.RGB(222, 235, 247)},
		{0.33, color.RGB(158, 202, 225)},
		{0.67, color.RGB(66, 146, 198)},
		{1.00, color.RGB(8, 48, 107)},
	}
	return interpolate(t, stops)
}

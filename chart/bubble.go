package chart

import (
	"math"
	"sort"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// BubbleChart renders a scatter plot where a third variable is encoded as circle radius.
type BubbleChart struct {
	xs, ys, sizes []float64
	style         BubbleStyle
	col           color.Color
}

// NewBubble creates a new BubbleChart.
func NewBubble(xs, ys, sizes []float64, style BubbleStyle) *BubbleChart {
	if style.SizeMin == 0 {
		style.SizeMin = 4
	}
	if style.SizeMax == 0 {
		style.SizeMax = 30
	}
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &BubbleChart{xs: xs, ys: ys, sizes: sizes, style: style, col: style.Color}
}

func (b *BubbleChart) Label() string          { return b.style.Label }
func (b *BubbleChart) Color() color.Color     { return b.col }
func (b *BubbleChart) SetColor(c color.Color) { b.col = c }

// DataRange is based on x/y only; bubble sizes do not affect axis ranges.
func (b *BubbleChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(b.xs) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = b.xs[0], b.xs[0]
	yMin, yMax = b.ys[0], b.ys[0]
	for i := 1; i < len(b.xs); i++ {
		if b.xs[i] < xMin {
			xMin = b.xs[i]
		}
		if b.xs[i] > xMax {
			xMax = b.xs[i]
		}
		if b.ys[i] < yMin {
			yMin = b.ys[i]
		}
		if b.ys[i] > yMax {
			yMax = b.ys[i]
		}
	}
	return
}

// Draw renders bubbles largest-first so smaller ones are never hidden.
func (b *BubbleChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	n := len(b.xs)
	if n == 0 {
		return
	}

	// Compute size range of the data
	sMin, sMax := b.sizes[0], b.sizes[0]
	for _, s := range b.sizes[1:] {
		if s < sMin {
			sMin = s
		}
		if s > sMax {
			sMax = s
		}
	}
	sSpan := sMax - sMin

	rMin := b.style.SizeMin
	rMax := b.style.SizeMax

	radii := make([]float64, n)
	for i, s := range b.sizes {
		t := 0.5
		if sSpan > 0 {
			t = (s - sMin) / sSpan
		}
		radii[i] = rMin + t*(rMax-rMin)
	}

	// Draw order: largest first (back to front)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		return radii[order[a]] > radii[order[b]]
	})

	opacity := b.style.Opacity
	if opacity == 0 || opacity > 1 {
		opacity = 1.0
	}
	fillOpacity := math.Min(opacity, 0.65)

	white := color.Color{R: 255, G: 255, B: 255, A: 180}

	for _, i := range order {
		px := xScale.Map(b.xs[i])
		py := yScale.Map(b.ys[i])
		r := radii[i]

		canvas.DrawCircle(px, py, r, render.Style{
			Fill:        b.col,
			FillOpacity: fillOpacity,
			Stroke:      white,
			StrokeWidth: 1.0,
		})
	}
}

// DrawLabels draws per-point text labels centered on each bubble. It is called
// outside the plot clip path so labels near the edges are never truncated.
func (b *BubbleChart) DrawLabels(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(b.style.Labels) == 0 {
		return
	}

	const labelFontSize = 11.0
	const labelMinRadius = 10.0

	n := len(b.xs)
	sMin, sMax := b.sizes[0], b.sizes[0]
	for _, s := range b.sizes[1:] {
		if s < sMin { sMin = s }
		if s > sMax { sMax = s }
	}
	sSpan := sMax - sMin
	rMin := b.style.SizeMin
	rMax := b.style.SizeMax

	opacity := b.style.Opacity
	if opacity == 0 || opacity > 1 {
		opacity = 1.0
	}
	fillOpacity := math.Min(opacity, 0.65)

	for i := 0; i < n; i++ {
		if i >= len(b.style.Labels) || b.style.Labels[i] == "" {
			continue
		}
		t := 0.5
		if sSpan > 0 {
			t = (b.sizes[i] - sMin) / sSpan
		}
		r := rMin + t*(rMax-rMin)
		if r < labelMinRadius {
			continue
		}

		px := xScale.Map(b.xs[i])
		py := yScale.Map(b.ys[i])

		// Choose label color based on perceived brightness of the effective bubble
		// color (raw fill blended with white at fillOpacity).
		effR := float64(b.col.R)*fillOpacity + 255*(1-fillOpacity)
		effG := float64(b.col.G)*fillOpacity + 255*(1-fillOpacity)
		effB := float64(b.col.B)*fillOpacity + 255*(1-fillOpacity)
		lum := 0.299*effR + 0.587*effG + 0.114*effB
		var labelCol color.Color
		if lum > 140 {
			labelCol = color.Color{R: 30, G: 30, B: 30, A: 220}
		} else {
			labelCol = color.Color{R: 255, G: 255, B: 255, A: 220}
		}

		canvas.DrawText(px, py, b.style.Labels[i], render.TextStyle{
			Color:    labelCol,
			FontSize: labelFontSize,
			Anchor:   "middle",
			Baseline: "middle",
			Bold:     true,
		})
	}
}

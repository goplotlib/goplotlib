package chart

import (
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// StackedAreaChart renders multiple area series accumulated bottom-up.
type StackedAreaChart struct {
	xs     []float64
	series [][]float64
	labels []string
	colors []color.Color
}

// NewStackedArea creates a new StackedAreaChart.
func NewStackedArea(xs []float64, series [][]float64, labels []string) *StackedAreaChart {
	colors := make([]color.Color, len(series))
	for i := range colors {
		colors[i] = color.Transparent
	}
	return &StackedAreaChart{xs: xs, series: series, labels: labels, colors: colors}
}

// Label returns "" — legend entries come from LegendEntries().
func (sa *StackedAreaChart) Label() string { return "" }

// Color returns Transparent — colors are assigned per series.
func (sa *StackedAreaChart) Color() color.Color { return color.Transparent }

// AssignColors assigns palette colors to unset series.
func (sa *StackedAreaChart) AssignColors(palette []color.Color, idx *int) {
	for i, c := range sa.colors {
		if c.IsTransparent() {
			sa.colors[i] = palette[*idx%len(palette)]
			*idx++
		}
	}
}

// LegendEntries implements MultiLegend — one entry per series.
func (sa *StackedAreaChart) LegendEntries() []LegendEntry {
	entries := make([]LegendEntry, len(sa.series))
	for i, lbl := range sa.labels {
		entries[i] = LegendEntry{Label: lbl, Col: sa.colors[i], IsBar: false}
	}
	return entries
}

// DataRange returns x = data range, y = 0 to max column sum.
func (sa *StackedAreaChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(sa.xs) == 0 || len(sa.series) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = sa.xs[0], sa.xs[0]
	for _, x := range sa.xs[1:] {
		if x < xMin {
			xMin = x
		}
		if x > xMax {
			xMax = x
		}
	}
	yMin, yMax = 0, 0
	for i := range sa.xs {
		colSum := 0.0
		for _, s := range sa.series {
			if i < len(s) {
				colSum += s[i]
			}
		}
		if colSum > yMax {
			yMax = colSum
		}
	}
	return
}

// Draw renders the stacked area chart.
func (sa *StackedAreaChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(sa.xs) == 0 || len(sa.series) == 0 {
		return
	}
	n := len(sa.xs)

	// cum[k][i] = sum of series 0..k-1 at point i; cum[0] = 0 (baseline)
	cum := make([][]float64, len(sa.series)+1)
	for k := range cum {
		cum[k] = make([]float64, n)
	}
	for k, s := range sa.series {
		for i := 0; i < n; i++ {
			v := 0.0
			if i < len(s) {
				v = s[i]
			}
			cum[k+1][i] = cum[k][i] + v
		}
	}

	for k, col := range sa.colors {
		// Filled polygon: top edge (cum[k+1]) then back along bottom edge (cum[k])
		var poly render.Path
		poly.MoveToP(xScale.Map(sa.xs[0]), yScale.Map(cum[k+1][0]))
		for i := 1; i < n; i++ {
			poly.LineToP(xScale.Map(sa.xs[i]), yScale.Map(cum[k+1][i]))
		}
		for i := n - 1; i >= 0; i-- {
			poly.LineToP(xScale.Map(sa.xs[i]), yScale.Map(cum[k][i]))
		}
		poly.CloseP()
		canvas.DrawPath(poly, render.Style{
			Fill:        col,
			FillOpacity: 0.75,
		})

		// Top edge line for visual separation
		var topLine render.Path
		topLine.MoveToP(xScale.Map(sa.xs[0]), yScale.Map(cum[k+1][0]))
		for i := 1; i < n; i++ {
			topLine.LineToP(xScale.Map(sa.xs[i]), yScale.Map(cum[k+1][i]))
		}
		canvas.DrawPath(topLine, render.Style{
			Stroke:      col,
			StrokeWidth: 1.5,
		})
	}
}

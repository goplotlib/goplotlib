package chart

import (
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// StackedBarChart renders multiple bar series stacked vertically within each category.
type StackedBarChart struct {
	categories []string
	series     [][]float64
	labels     []string
	colors     []color.Color
	style      StackedBarStyle
}

// NewStackedBar creates a new StackedBarChart.
func NewStackedBar(categories []string, series [][]float64, labels []string, style StackedBarStyle) *StackedBarChart {
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	colors := make([]color.Color, len(series))
	for i := range colors {
		colors[i] = color.Transparent
	}
	return &StackedBarChart{
		categories: categories,
		series:     series,
		labels:     labels,
		colors:     colors,
		style:      style,
	}
}

// Label returns "" — legend entries come from LegendEntries().
func (sb *StackedBarChart) Label() string { return "" }

// Color returns Transparent — colors are assigned per series.
func (sb *StackedBarChart) Color() color.Color { return color.Transparent }

// Categories returns the category label list.
func (sb *StackedBarChart) Categories() []string { return sb.categories }

// AssignColors assigns palette colors to series that have no explicit color.
func (sb *StackedBarChart) AssignColors(palette []color.Color, idx *int) {
	for i, c := range sb.colors {
		if c.IsTransparent() {
			sb.colors[i] = palette[*idx%len(palette)]
			*idx++
		}
	}
}

// LegendEntries implements MultiLegend — one entry per series.
func (sb *StackedBarChart) LegendEntries() []LegendEntry {
	entries := make([]LegendEntry, len(sb.series))
	for i, lbl := range sb.labels {
		entries[i] = LegendEntry{Label: lbl, Col: sb.colors[i], IsBar: true}
	}
	return entries
}

// DataRange returns x = category indices, y = min/max of column sums.
func (sb *StackedBarChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(sb.categories) == 0 || len(sb.series) == 0 {
		return 0, 1, 0, 1
	}
	n := len(sb.categories)
	xMin, xMax = 0, float64(n-1)
	yMin, yMax = 0, 0
	for i := 0; i < n; i++ {
		posSum, negSum := 0.0, 0.0
		for _, s := range sb.series {
			if i < len(s) {
				if s[i] >= 0 {
					posSum += s[i]
				} else {
					negSum += s[i]
				}
			}
		}
		if posSum > yMax {
			yMax = posSum
		}
		if negSum < yMin {
			yMin = negSum
		}
	}
	return
}

// Draw renders the stacked bar chart.
func (sb *StackedBarChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	catScale, ok := xScale.(*scale.CategoricalScale)
	if !ok {
		return
	}

	bandW := catScale.BandWidth()
	barW := bandW * 0.75

	n := len(sb.categories)
	for i := 0; i < n; i++ {
		posAcc := 0.0 // accumulated positive data value
		negAcc := 0.0 // accumulated negative data value
		barX := catScale.Map(float64(i)) - barW/2

		for k, s := range sb.series {
			if i >= len(s) {
				continue
			}
			v := s[i]
			if v == 0 {
				continue
			}

			fillColor := sb.colors[k]
			if sb.style.Opacity < 1.0 && sb.style.Opacity > 0 {
				fillColor = fillColor.WithAlpha(sb.style.Opacity)
			}

			var segY, segH float64
			if v > 0 {
				segY = yScale.Map(posAcc + v)
				segH = yScale.Map(posAcc) - segY
				posAcc += v
			} else {
				segY = yScale.Map(negAcc)
				segH = yScale.Map(negAcc+v) - segY
				negAcc += v
			}
			if segH < 1 {
				segH = 1
			}

			canvas.DrawRect(barX, segY, barW, segH, 0, 0, render.Style{
				Fill:        fillColor,
				FillOpacity: 1.0,
			})
		}
	}
}

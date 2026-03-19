package chart

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// HBarChart renders a horizontal bar chart with a categorical y-axis.
type HBarChart struct {
	categories []string
	values     []float64
	style      BarStyle
	col        color.Color
}

// NewHBar creates a new HBarChart.
func NewHBar(categories []string, values []float64, style BarStyle) *HBarChart {
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &HBarChart{categories: categories, values: values, style: style, col: style.Color}
}

func (hb *HBarChart) Label() string          { return hb.style.Label }
func (hb *HBarChart) Color() color.Color     { return hb.col }
func (hb *HBarChart) SetColor(c color.Color) { hb.col = c }
func (hb *HBarChart) Categories() []string   { return hb.categories }

// DataRange returns the value range on x and the category index range on y.
func (hb *HBarChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(hb.values) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = 0, 0
	for _, v := range hb.values {
		if v < xMin {
			xMin = v
		}
		if v > xMax {
			xMax = v
		}
	}
	yMin = 0
	yMax = float64(len(hb.values) - 1)
	return
}

// Draw renders the horizontal bar chart.
// xScale is the linear value scale; yScale must be a *scale.CategoricalScale.
func (hb *HBarChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	catScale, ok := yScale.(*scale.CategoricalScale)
	if !ok {
		return
	}

	bandH := catScale.BandWidth()
	barH := bandH * 0.75
	xZero := xScale.Map(0)

	rx := 0.0
	if !hb.style.SquareBars {
		rx = 4.0
	}

	for i, v := range hb.values {
		centerY := catScale.Map(float64(i))
		barY := centerY - barH/2

		barX := xZero
		barW := xScale.Map(v) - xZero
		if v < 0 {
			barX = xScale.Map(v)
			barW = math.Abs(barW)
		}
		if barW < 1 {
			barW = 1
		}

		fillColor := hb.col
		if hb.style.Opacity < 1.0 && hb.style.Opacity > 0 {
			fillColor = fillColor.WithAlpha(hb.style.Opacity)
		}
		canvas.DrawRect(barX, barY, barW, barH, rx, rx, render.Style{
			Fill:        fillColor,
			FillOpacity: 1.0,
		})
	}
}

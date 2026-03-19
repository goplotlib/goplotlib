package chart

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// BarChart renders a bar chart with categorical x-axis.
type BarChart struct {
	categories []string
	values     []float64
	style      BarStyle
	col        color.Color
	groupIdx   int // index within a grouped set (0 = first or only)
	groupCount int // total bars per category band (1 = ungrouped)
}

// SetGroupInfo configures this bar as part of a grouped set.
func (bc *BarChart) SetGroupInfo(idx, count int) { bc.groupIdx = idx; bc.groupCount = count }

// NewBar creates a new BarChart.
func NewBar(categories []string, values []float64, style BarStyle) *BarChart {
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &BarChart{
		categories: categories,
		values:     values,
		style:      style,
		col:        style.Color,
	}
}

// Label returns the series label.
func (bc *BarChart) Label() string { return bc.style.Label }

// Color returns the series color.
func (bc *BarChart) Color() color.Color { return bc.col }

// SetColor sets the series color.
func (bc *BarChart) SetColor(c color.Color) { bc.col = c }

// DataRange returns [xMin, xMax, yMin, yMax].
// xMin/xMax are indices for the categorical scale.
// yMin always includes 0 for the baseline.
func (bc *BarChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(bc.values) == 0 {
		return 0, 1, 0, 1
	}
	xMin = 0
	xMax = float64(len(bc.values) - 1)
	yMin = 0
	yMax = 0
	for _, v := range bc.values {
		if v < yMin {
			yMin = v
		}
		if v > yMax {
			yMax = v
		}
	}
	// Include 0 in range
	if yMin > 0 {
		yMin = 0
	}
	return
}

// Categories returns the list of category labels.
func (bc *BarChart) Categories() []string { return bc.categories }

// Draw renders the bar chart.
func (bc *BarChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	catScale, ok := xScale.(*scale.CategoricalScale)
	if !ok {
		return
	}

	bandW := catScale.BandWidth()
	barW := bandW * 0.75 // may be overridden per-bar in grouped mode

	// y=0 baseline in pixel coords
	yZero := yScale.Map(0)

	rx := 0.0
	if !bc.style.SquareBars {
		rx = 4.0
	}

	for i, v := range bc.values {
		var barX float64
		if bc.groupCount > 1 {
			subBandW := bandW / float64(bc.groupCount)
			barW = subBandW * 0.85
			barX = catScale.Map(float64(i)) - bandW/2 + float64(bc.groupIdx)*subBandW + (subBandW-barW)/2
		} else {
			barW = bandW * 0.75
			barX = catScale.Map(float64(i)) - barW/2
		}
		barY := yScale.Map(v)
		barH := math.Abs(yZero - barY)

		// For negative values, barY is below yZero
		if v < 0 {
			barY = yZero
		} else {
			barY = yScale.Map(v)
		}

		if barH < 1 {
			barH = 1
		}

		fillColor := bc.col
		if bc.style.Opacity < 1.0 && bc.style.Opacity > 0 {
			fillColor = fillColor.WithAlpha(bc.style.Opacity)
		}

		barStyle := render.Style{
			Fill:        fillColor,
			FillOpacity: 1.0,
		}

		canvas.DrawRect(barX, barY, barW, barH, rx, rx, barStyle)
	}
}

package chart

import (
	"math"
	"sort"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// boxStats holds pre-computed statistics for one group.
type boxStats struct {
	q1, median, q3 float64
	whiskerLo      float64 // lowest value within Q1 - 1.5*IQR
	whiskerHi      float64 // highest value within Q3 + 1.5*IQR
	outliers       []float64
}

// computeStats sorts the data and computes box plot statistics.
func computeStats(data []float64) boxStats {
	if len(data) == 0 {
		return boxStats{}
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	q1 := percentile(sorted, 25)
	median := percentile(sorted, 50)
	q3 := percentile(sorted, 75)
	iqr := q3 - q1
	loFence := q1 - 1.5*iqr
	hiFence := q3 + 1.5*iqr

	whiskerLo := q1
	whiskerHi := q3
	var outliers []float64
	for _, v := range sorted {
		if v < loFence || v > hiFence {
			outliers = append(outliers, v)
		} else {
			if v < whiskerLo {
				whiskerLo = v
			}
			if v > whiskerHi {
				whiskerHi = v
			}
		}
	}
	return boxStats{q1: q1, median: median, q3: q3,
		whiskerLo: whiskerLo, whiskerHi: whiskerHi, outliers: outliers}
}

// percentile returns the p-th percentile of a sorted slice using linear interpolation.
func percentile(sorted []float64, p float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sorted[0]
	}
	h := (float64(n-1)) * p / 100.0
	lo := int(math.Floor(h))
	hi := int(math.Ceil(h))
	if lo == hi {
		return sorted[lo]
	}
	return sorted[lo] + (h-float64(lo))*(sorted[hi]-sorted[lo])
}

// BoxChart renders a box plot with one box per group.
type BoxChart struct {
	groups []boxStats
	labels []string // x-axis category labels
	style  BoxStyle
	col    color.Color
}

// NewBox creates a BoxChart by pre-computing statistics for each group.
func NewBox(groups [][]float64, style BoxStyle) *BoxChart {
	stats := make([]boxStats, len(groups))
	for i, g := range groups {
		stats[i] = computeStats(g)
	}
	labels := make([]string, len(groups))
	for i := range labels {
		if i < len(style.Labels) {
			labels[i] = style.Labels[i]
		}
	}
	return &BoxChart{groups: stats, labels: labels, style: style, col: style.Color}
}

func (b *BoxChart) Label() string          { return b.style.Label }
func (b *BoxChart) Color() color.Color     { return b.col }
func (b *BoxChart) SetColor(c color.Color) { b.col = c }
func (b *BoxChart) Categories() []string   { return b.labels }

// DataRange returns x as category indices and y covering all data including outliers.
func (b *BoxChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(b.groups) == 0 {
		return 0, 1, 0, 1
	}
	xMin = 0
	xMax = float64(len(b.groups) - 1)
	yMin = math.Inf(1)
	yMax = math.Inf(-1)
	for _, g := range b.groups {
		if g.whiskerLo < yMin {
			yMin = g.whiskerLo
		}
		if g.whiskerHi > yMax {
			yMax = g.whiskerHi
		}
		for _, o := range g.outliers {
			if o < yMin {
				yMin = o
			}
			if o > yMax {
				yMax = o
			}
		}
	}
	if math.IsInf(yMin, 1) {
		yMin, yMax = 0, 1
	}
	return
}

// Draw renders the box plot. xScale must be a *scale.CategoricalScale.
func (b *BoxChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	catScale, ok := xScale.(*scale.CategoricalScale)
	if !ok {
		return
	}

	bandW := catScale.BandWidth()
	boxW := bandW * 0.5

	fillColor := b.col.WithAlpha(0.7)
	borderStyle := render.Style{
		Fill:        fillColor,
		FillOpacity: 1.0,
		Stroke:      b.col,
		StrokeWidth: 1.5,
	}
	whiskStyle := render.Style{
		Stroke:      b.col,
		StrokeWidth: 1.5,
	}
	medianStyle := render.Style{
		Stroke:      color.RGB(255, 255, 255),
		StrokeWidth: 2.0,
	}
	capStyle := render.Style{
		Stroke:      b.col,
		StrokeWidth: 1.5,
	}
	outlierStyle := render.Style{
		Fill:        b.col,
		FillOpacity: 0.7,
		Stroke:      b.col,
		StrokeWidth: 1.0,
	}

	capW := boxW * 0.4

	for i, g := range b.groups {
		cx := catScale.Map(float64(i))
		boxLeft := cx - boxW/2
		q1y := yScale.Map(g.q1)
		q3y := yScale.Map(g.q3)
		medy := yScale.Map(g.median)
		loY := yScale.Map(g.whiskerLo)
		hiY := yScale.Map(g.whiskerHi)

		boxTop := math.Min(q1y, q3y)
		boxH := math.Abs(q1y - q3y)
		if boxH < 1 {
			boxH = 1
		}

		// Whisker lines
		canvas.DrawLine(cx, q1y, cx, loY, whiskStyle)
		canvas.DrawLine(cx, q3y, cx, hiY, whiskStyle)

		// Whisker caps
		canvas.DrawLine(cx-capW/2, loY, cx+capW/2, loY, capStyle)
		canvas.DrawLine(cx-capW/2, hiY, cx+capW/2, hiY, capStyle)

		// Box (Q1–Q3)
		canvas.DrawRect(boxLeft, boxTop, boxW, boxH, 2, 2, borderStyle)

		// Median line
		canvas.DrawLine(boxLeft, medy, boxLeft+boxW, medy, medianStyle)

		// Outliers
		for _, ov := range g.outliers {
			oy := yScale.Map(ov)
			canvas.DrawCircle(cx, oy, 3, outlierStyle)
		}
	}
}

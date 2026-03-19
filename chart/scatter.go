package chart

import (
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/spline"
)

// ScatterChart renders a scatter plot.
type ScatterChart struct {
	xs, ys []float64
	style  ScatterStyle
	col    color.Color
}

// NewScatter creates a new ScatterChart.
func NewScatter(xs, ys []float64, style ScatterStyle) *ScatterChart {
	if style.MarkerSize == 0 {
		style.MarkerSize = 5
	}
	if style.MarkerShape == "" {
		style.MarkerShape = "circle"
	}
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &ScatterChart{
		xs:    xs,
		ys:    ys,
		style: style,
		col:   style.Color,
	}
}

// Label returns the series label.
func (sc *ScatterChart) Label() string { return sc.style.Label }

// Color returns the series color.
func (sc *ScatterChart) Color() color.Color { return sc.col }

// SetColor sets the series color.
func (sc *ScatterChart) SetColor(c color.Color) { sc.col = c }

// DataRange returns [xMin, xMax, yMin, yMax].
func (sc *ScatterChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(sc.xs) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = sc.xs[0], sc.xs[0]
	yMin, yMax = sc.ys[0], sc.ys[0]
	for i := 1; i < len(sc.xs); i++ {
		if sc.xs[i] < xMin {
			xMin = sc.xs[i]
		}
		if sc.xs[i] > xMax {
			xMax = sc.xs[i]
		}
		if sc.ys[i] < yMin {
			yMin = sc.ys[i]
		}
		if sc.ys[i] > yMax {
			yMax = sc.ys[i]
		}
	}
	return
}

// Draw renders the scatter chart.
func (sc *ScatterChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(sc.xs) == 0 {
		return
	}

	white := color.Color{R: 255, G: 255, B: 255, A: 200}

	markerStyle := render.Style{
		Fill:        sc.col,
		FillOpacity: 1.0,
		Stroke:      white,
		StrokeWidth: 1.2,
	}

	size := sc.style.MarkerSize
	if size == 0 {
		size = 5
	}

	pts := make([]spline.Point, len(sc.xs))
	for i := range sc.xs {
		pts[i] = spline.Point{
			X: xScale.Map(sc.xs[i]),
			Y: yScale.Map(sc.ys[i]),
		}
	}

	for _, pt := range pts {
		switch sc.style.MarkerShape {
		case "square":
			canvas.DrawRect(pt.X-size, pt.Y-size, size*2, size*2, 0, 0, markerStyle)
		case "diamond":
			var p render.Path
			p.MoveToP(pt.X, pt.Y-size*1.2)
			p.LineToP(pt.X+size, pt.Y)
			p.LineToP(pt.X, pt.Y+size*1.2)
			p.LineToP(pt.X-size, pt.Y)
			p.CloseP()
			canvas.DrawPath(p, markerStyle)
		case "triangle":
			var p render.Path
			p.MoveToP(pt.X, pt.Y-size*1.2)
			p.LineToP(pt.X+size, pt.Y+size*0.8)
			p.LineToP(pt.X-size, pt.Y+size*0.8)
			p.CloseP()
			canvas.DrawPath(p, markerStyle)
		default: // "circle"
			canvas.DrawCircle(pt.X, pt.Y, size, markerStyle)
		}
	}
}

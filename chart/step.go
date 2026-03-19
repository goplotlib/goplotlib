package chart

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// StepChart renders a staircase line chart.
type StepChart struct {
	xs, ys []float64
	style  StepStyle
	col    color.Color
}

// NewStep creates a new StepChart.
func NewStep(xs, ys []float64, style StepStyle) *StepChart {
	if style.LineWidth == 0 {
		style.LineWidth = 2.5
	}
	if style.Mode == "" {
		style.Mode = "post"
	}
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &StepChart{xs: xs, ys: ys, style: style, col: style.Color}
}

func (sc *StepChart) Label() string          { return sc.style.Label }
func (sc *StepChart) Color() color.Color     { return sc.col }
func (sc *StepChart) SetColor(c color.Color) { sc.col = c }

func (sc *StepChart) DataRange() (xMin, xMax, yMin, yMax float64) {
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

// buildStepPath constructs the staircase render.Path for the given pixel points.
// mode is "pre", "post" (default), or "mid".
func buildStepPath(pts [][2]float64, mode string) render.Path {
	var p render.Path
	if len(pts) == 0 {
		return p
	}
	p.MoveToP(pts[0][0], pts[0][1])
	for i := 1; i < len(pts); i++ {
		x0, y0 := pts[i-1][0], pts[i-1][1]
		x1, y1 := pts[i][0], pts[i][1]
		switch mode {
		case "pre":
			p.LineToP(x0, y1)
			p.LineToP(x1, y1)
		case "mid":
			mid := (x0 + x1) / 2
			p.LineToP(mid, y0)
			p.LineToP(mid, y1)
			p.LineToP(x1, y1)
		default: // "post"
			p.LineToP(x1, y0)
			p.LineToP(x1, y1)
		}
	}
	return p
}

// Draw renders the step chart.
func (sc *StepChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(sc.xs) < 2 {
		return
	}

	pts := make([][2]float64, len(sc.xs))
	for i := range sc.xs {
		pts[i] = [2]float64{xScale.Map(sc.xs[i]), yScale.Map(sc.ys[i])}
	}

	mode := sc.style.Mode
	if mode == "" {
		mode = "post"
	}

	lineWidth := sc.style.LineWidth
	if lineWidth == 0 {
		lineWidth = 2.5
	}

	if sc.style.Fill {
		domain := yScale.Domain()
		baseline := yScale.Map(math.Max(domain[0], math.Min(0, domain[1])))
		fillPath := buildStepPath(pts, mode)
		fillPath.LineToP(pts[len(pts)-1][0], baseline)
		fillPath.LineToP(pts[0][0], baseline)
		fillPath.CloseP()
		canvas.DrawPath(fillPath, render.Style{
			Fill:        sc.col.WithAlpha(0.2),
			FillOpacity: 1.0,
		})
	}

	linePath := buildStepPath(pts, mode)
	canvas.DrawPath(linePath, render.Style{
		Stroke:      sc.col,
		StrokeWidth: lineWidth,
		LineCap:     "round",
		LineJoin:    "miter",
		Dash:        sc.style.Dash,
	})
}

package chart

import (
	"fmt"
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/spline"
)

// LineChart renders a line (or smooth spline) chart.
type LineChart struct {
	xs, ys []float64
	style  LineStyle
	col    color.Color
	idx    int // index used for unique gradient IDs
}

// NewLine creates a new LineChart.
func NewLine(xs, ys []float64, style LineStyle) *LineChart {
	if style.LineWidth == 0 {
		style.LineWidth = 2.5
	}
	if style.MarkerSize == 0 {
		style.MarkerSize = 5
	}
	if style.MarkerShape == "" {
		style.MarkerShape = "circle"
	}
	if style.Fill && style.FillOpacity == 0 {
		style.FillOpacity = 0.15
	}
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}
	return &LineChart{xs: xs, ys: ys, style: style, col: style.Color}
}

// SetIndex sets the chart index for unique ID generation.
func (lc *LineChart) SetIndex(i int) {
	lc.idx = i
}

// SetColor sets the series color.
func (lc *LineChart) SetColor(c color.Color) {
	lc.col = c
}

// Label returns the series label.
func (lc *LineChart) Label() string { return lc.style.Label }

// Color returns the series color.
func (lc *LineChart) Color() color.Color { return lc.col }

// DataRange returns [xMin, xMax, yMin, yMax].
func (lc *LineChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(lc.xs) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = lc.xs[0], lc.xs[0]
	yMin, yMax = lc.ys[0], lc.ys[0]
	for i := 1; i < len(lc.xs); i++ {
		if lc.xs[i] < xMin {
			xMin = lc.xs[i]
		}
		if lc.xs[i] > xMax {
			xMax = lc.xs[i]
		}
		if lc.ys[i] < yMin {
			yMin = lc.ys[i]
		}
		if lc.ys[i] > yMax {
			yMax = lc.ys[i]
		}
	}
	return
}

// Draw renders the line chart.
func (lc *LineChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(lc.xs) < 2 {
		return
	}

	// Map data points to pixel coordinates
	pts := make([]spline.Point, len(lc.xs))
	for i := range lc.xs {
		pts[i] = spline.Point{
			X: xScale.Map(lc.xs[i]),
			Y: yScale.Map(lc.ys[i]),
		}
	}

	lineWidth := lc.style.LineWidth
	if lineWidth == 0 {
		lineWidth = 2.5
	}

	lineStyle := render.Style{
		Stroke:      lc.col,
		StrokeWidth: lineWidth,
		LineCap:     "round",
		LineJoin:    "round",
		Dash:        lc.style.Dash,
	}

	var linePath render.Path

	if lc.style.Smooth && len(pts) >= 3 {
		// Use Catmull-Rom splines
		segments := spline.CatmullRomToBezier(pts)
		if len(segments) > 0 {
			seg0 := segments[0]
			linePath.MoveToP(seg0[0].X, seg0[0].Y)
			for _, seg := range segments {
				linePath.CubicToP(seg[1].X, seg[1].Y, seg[2].X, seg[2].Y, seg[3].X, seg[3].Y)
			}
		}
	} else {
		// Straight lines
		linePath.MoveToP(pts[0].X, pts[0].Y)
		for _, pt := range pts[1:] {
			linePath.LineToP(pt.X, pt.Y)
		}
	}

	// Fill area under the line
	if lc.style.Fill {
		gradID := fmt.Sprintf("linefill-%d", lc.idx)
		yBottom := yScale.Map(math.Min(0, lc.ys[0]))
		// Find the baseline (y=0 or bottom of chart)
		domain := yScale.Domain()
		baseline := yScale.Map(math.Max(domain[0], math.Min(0, domain[1])))

		// Define vertical gradient
		canvas.DefineLinearGradient(gradID,
			0, pts[0].Y,   // x1,y1 — top (use first point's y roughly)
			0, baseline,   // x2,y2 — bottom
			[]render.GradStop{
				{Offset: 0, Color: lc.col, Opacity: 0.35},
				{Offset: 1, Color: lc.col, Opacity: 0.02},
			},
		)

		_ = yBottom

		// Build fill path: follow the line then close via baseline
		var fillPath render.Path
		if lc.style.Smooth && len(pts) >= 3 {
			segments := spline.CatmullRomToBezier(pts)
			if len(segments) > 0 {
				seg0 := segments[0]
				fillPath.MoveToP(seg0[0].X, seg0[0].Y)
				for _, seg := range segments {
					fillPath.CubicToP(seg[1].X, seg[1].Y, seg[2].X, seg[2].Y, seg[3].X, seg[3].Y)
				}
			}
		} else {
			fillPath.MoveToP(pts[0].X, pts[0].Y)
			for _, pt := range pts[1:] {
				fillPath.LineToP(pt.X, pt.Y)
			}
		}
		// Close via bottom
		fillPath.LineToP(pts[len(pts)-1].X, baseline)
		fillPath.LineToP(pts[0].X, baseline)
		fillPath.CloseP()

		fillStyle := render.Style{
			GradientID:  gradID,
			FillOpacity: 1.0,
		}
		canvas.DrawPath(fillPath, fillStyle)
	}

	canvas.DrawPath(linePath, lineStyle)

	// Draw markers
	markerSize := lc.style.MarkerSize
	if markerSize > 0 {
		drawMarkers(canvas, pts, lc.col, markerSize, lc.style.MarkerShape)
	}
}

// drawMarkers draws marker shapes at each data point.
func drawMarkers(canvas render.Canvas, pts []spline.Point, col color.Color, size float64, shape string) {
	white := color.Color{R: 255, G: 255, B: 255, A: 255}
	markerFill := render.Style{
		Fill:        col,
		FillOpacity: 1.0,
		Stroke:      white,
		StrokeWidth: 1.0,
	}

	for _, pt := range pts {
		switch shape {
		case "square":
			canvas.DrawRect(pt.X-size, pt.Y-size, size*2, size*2, 0, 0, markerFill)
		case "diamond":
			var p render.Path
			p.MoveToP(pt.X, pt.Y-size*1.2)
			p.LineToP(pt.X+size, pt.Y)
			p.LineToP(pt.X, pt.Y+size*1.2)
			p.LineToP(pt.X-size, pt.Y)
			p.CloseP()
			canvas.DrawPath(p, markerFill)
		case "triangle":
			var p render.Path
			p.MoveToP(pt.X, pt.Y-size*1.2)
			p.LineToP(pt.X+size, pt.Y+size*0.8)
			p.LineToP(pt.X-size, pt.Y+size*0.8)
			p.CloseP()
			canvas.DrawPath(p, markerFill)
		default: // "circle"
			canvas.DrawCircle(pt.X, pt.Y, size, markerFill)
		}
	}
}

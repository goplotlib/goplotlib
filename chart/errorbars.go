package chart

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

const errorBarCapWidth = 6.0

// ErrorBarsChart renders vertical or horizontal error bars over data points.
// It implements Chart so that its range expands the axis bounds correctly.
type ErrorBarsChart struct {
	xs, ys     []float64
	lo, hi     []float64 // error magnitudes (positive values); bar extends y-lo[i] to y+hi[i]
	horizontal bool      // true → x error bars
	col        color.Color
}

// NewErrorBars creates vertical error bars. lo[i] and hi[i] are the downward
// and upward error magnitudes respectively (e.g. standard deviation).
func NewErrorBars(xs, ys, lo, hi []float64, col color.Color) *ErrorBarsChart {
	return &ErrorBarsChart{
		xs:  xs,
		ys:  ys,
		lo:  lo,
		hi:  hi,
		col: col,
	}
}

// NewErrorBarsX creates horizontal error bars with symmetric magnitude xErr.
func NewErrorBarsX(xs, ys, xErr []float64, col color.Color) *ErrorBarsChart {
	return &ErrorBarsChart{
		xs:         xs,
		ys:         ys,
		lo:         xErr,
		hi:         xErr,
		horizontal: true,
		col:        col,
	}
}

func (e *ErrorBarsChart) Label() string          { return "" } // error bars don't appear in legend
func (e *ErrorBarsChart) Color() color.Color     { return e.col }
func (e *ErrorBarsChart) SetColor(c color.Color) { e.col = c }

// DataRange expands to include the full extent of every error bar.
func (e *ErrorBarsChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(e.xs) == 0 {
		return 0, 1, 0, 1
	}
	xMin, xMax = e.xs[0], e.xs[0]
	yMin, yMax = e.ys[0], e.ys[0]
	for i, x := range e.xs {
		y := safeAt(e.ys, i)
		lo := math.Abs(safeAt(e.lo, i))
		hi := math.Abs(safeAt(e.hi, i))
		if e.horizontal {
			if x-lo < xMin {
				xMin = x - lo
			}
			if x+hi > xMax {
				xMax = x + hi
			}
			if y < yMin {
				yMin = y
			}
			if y > yMax {
				yMax = y
			}
		} else {
			if x < xMin {
				xMin = x
			}
			if x > xMax {
				xMax = x
			}
			if y-lo < yMin {
				yMin = y - lo
			}
			if y+hi > yMax {
				yMax = y + hi
			}
		}
	}
	return
}

// Draw renders the error bars using the provided scales.
func (e *ErrorBarsChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	style := render.Style{
		Stroke:      e.col,
		StrokeWidth: 1.5,
	}
	capW := errorBarCapWidth

	for i, x := range e.xs {
		if i >= len(e.ys) {
			break
		}
		y := e.ys[i]
		lo := math.Abs(safeAt(e.lo, i))
		hi := math.Abs(safeAt(e.hi, i))
		px := xScale.Map(x)
		py := yScale.Map(y)

		if e.horizontal {
			pxLo := xScale.Map(x - lo)
			pxHi := xScale.Map(x + hi)
			canvas.DrawLine(pxLo, py, pxHi, py, style)
			canvas.DrawLine(pxLo, py-capW/2, pxLo, py+capW/2, style)
			canvas.DrawLine(pxHi, py-capW/2, pxHi, py+capW/2, style)
		} else {
			pyLo := yScale.Map(y - lo)
			pyHi := yScale.Map(y + hi)
			canvas.DrawLine(px, pyLo, px, pyHi, style)
			canvas.DrawLine(px-capW/2, pyLo, px+capW/2, pyLo, style)
			canvas.DrawLine(px-capW/2, pyHi, px+capW/2, pyHi, style)
		}
	}
}

// safeAt returns slice[i] if in-bounds, else 0.
func safeAt(s []float64, i int) float64 {
	if i < len(s) {
		return s[i]
	}
	return 0
}

package chart

import (
	"fmt"
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/theme"
)

// PieChart renders a pie or donut chart. It does not use x/y scales.
type PieChart struct {
	labels []string
	values []float64
	colors []color.Color
	style  PieStyle
}

// NewPie creates a new PieChart.
func NewPie(labels []string, values []float64, style PieStyle) *PieChart {
	return &PieChart{labels: labels, values: values, style: style}
}

func (p *PieChart) Label() string        { return "" }
func (p *PieChart) Color() color.Color   { return color.Transparent }

// AssignColors sets per-segment colors from the palette.
func (p *PieChart) AssignColors(palette []color.Color, idx *int) {
	p.colors = make([]color.Color, len(p.values))
	for i := range p.values {
		p.colors[i] = palette[*idx%len(palette)]
		(*idx)++
	}
}

// DataRange is unused for pie charts (no axis scales).
func (p *PieChart) DataRange() (xMin, xMax, yMin, yMax float64) { return 0, 1, 0, 1 }

// Draw satisfies the Chart interface but is a no-op — use DrawInBounds instead.
func (p *PieChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {}

// DrawInBounds renders the pie/donut chart centered at (cx, cy) with the given outer radius.
func (p *PieChart) DrawInBounds(canvas render.Canvas, cx, cy, outerR float64, t theme.Theme) {
	if len(p.values) == 0 {
		return
	}

	// Sum values
	total := 0.0
	for _, v := range p.values {
		if v > 0 {
			total += v
		}
	}
	if total == 0 {
		return
	}

	innerR := outerR * p.style.DonutRadius

	// Build segments: start/end angles, colors
	type segment struct {
		start, end, mid float64
		col             color.Color
		label           string
		pct             float64
	}

	segments := make([]segment, 0, len(p.values))
	angle := -math.Pi / 2 // start at 12 o'clock
	for i, v := range p.values {
		if v <= 0 {
			continue
		}
		sweep := 2 * math.Pi * v / total
		mid := angle + sweep/2
		col := color.RGB(100, 100, 200) // fallback
		if i < len(p.colors) {
			col = p.colors[i]
		}
		lbl := ""
		if i < len(p.labels) {
			lbl = p.labels[i]
		}
		segments = append(segments, segment{
			start: angle,
			end:   angle + sweep,
			mid:   mid,
			col:   col,
			label: lbl,
			pct:   v / total,
		})
		angle += sweep
	}

	// Draw segments
	white := color.RGB(255, 255, 255)
	for i, seg := range segments {
		scx, scy := cx, cy
		// Explode: offset center outward along mid angle
		if p.style.ExplodeOffset > 0 && i == p.style.ExplodeIdx {
			scx += math.Cos(seg.mid) * p.style.ExplodeOffset
			scy += math.Sin(seg.mid) * p.style.ExplodeOffset
		}

		var path render.Path
		if innerR > 0 {
			// Donut segment
			ox1 := scx + outerR*math.Cos(seg.start)
			oy1 := scy + outerR*math.Sin(seg.start)
			ix2 := scx + innerR*math.Cos(seg.end)
			iy2 := scy + innerR*math.Sin(seg.end)

			path.MoveToP(ox1, oy1)
			appendArc(&path, scx, scy, outerR, seg.start, seg.end)
			path.LineToP(ix2, iy2)
			appendArc(&path, scx, scy, innerR, seg.end, seg.start)
			path.CloseP()
		} else {
			// Solid pie segment
			path.MoveToP(scx, scy)
			path.LineToP(scx+outerR*math.Cos(seg.start), scy+outerR*math.Sin(seg.start))
			appendArc(&path, scx, scy, outerR, seg.start, seg.end)
			path.CloseP()
		}

		canvas.DrawPath(path, render.Style{
			Fill:        seg.col,
			FillOpacity: 1.0,
			Stroke:      white,
			StrokeWidth: 1.5,
		})
	}

	// Draw labels (skip segments < 5%)
	labelR := outerR + 18
	lineR := outerR + 6
	labelStyle := render.TextStyle{
		Color:      t.TextColor,
		FontSize:   t.TickFontSize,
		FontFamily: t.FontFamily,
		Baseline:   "middle",
	}
	lineStyle := render.Style{
		Stroke:      t.TextColor,
		StrokeWidth: 0.8,
	}

	for _, seg := range segments {
		if seg.pct < 0.05 {
			continue
		}
		lx := cx + labelR*math.Cos(seg.mid)
		ly := cy + labelR*math.Sin(seg.mid)
		lx2 := cx + lineR*math.Cos(seg.mid)
		ly2 := cy + lineR*math.Sin(seg.mid)

		canvas.DrawLine(lx2, ly2, lx, ly, lineStyle)

		anchor := "start"
		if math.Cos(seg.mid) < -0.1 {
			anchor = "end"
		}
		labelStyle.Anchor = anchor
		text := fmt.Sprintf("%s %.0f%%", seg.label, seg.pct*100)
		canvas.DrawText(lx+math.Copysign(4, math.Cos(seg.mid)), ly, text, labelStyle)
	}
}

// appendArc appends a cubic-bezier-approximated arc to the path.
// Sweeps from startAngle to endAngle around (cx, cy) with given radius.
func appendArc(path *render.Path, cx, cy, r, startAngle, endAngle float64) {
	sweep := endAngle - startAngle
	// Split into segments of at most π/2 for accuracy
	nSegs := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	if nSegs < 1 {
		nSegs = 1
	}
	step := sweep / float64(nSegs)
	for i := 0; i < nSegs; i++ {
		a0 := startAngle + float64(i)*step
		a1 := a0 + step
		k := 4.0 / 3.0 * math.Tan(step/4)
		x0 := cx + r*math.Cos(a0)
		y0 := cy + r*math.Sin(a0)
		x3 := cx + r*math.Cos(a1)
		y3 := cy + r*math.Sin(a1)
		cp1x := x0 - r*k*math.Sin(a0)
		cp1y := y0 + r*k*math.Cos(a0)
		cp2x := x3 + r*k*math.Sin(a1)
		cp2y := y3 - r*k*math.Cos(a1)
		path.CubicToP(cp1x, cp1y, cp2x, cp2y, x3, y3)
	}
}

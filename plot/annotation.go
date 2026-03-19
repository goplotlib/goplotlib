package plot

import (
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/render/svg"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/theme"
)

// AnnotationOption configures an annotation.
type AnnotationOption func(*annotationStyle)

type arrowDirection int

const (
	arrowNone  arrowDirection = iota
	arrowDown
	arrowUp
	arrowLeft
	arrowRight
)

const annotArrowHead = 6.0  // arrowhead size in px
const annotOffset    = 24.0 // distance from data point to label in px

// annotationStyle holds visual properties shared by all annotation types.
// Zero value is valid; theme-dependent defaults are resolved in each draw method.
type annotationStyle struct {
	label     string
	lineColor color.Color // Transparent = use theme.TextColor
	fillColor color.Color // Transparent = derive from lineColor
	lineWidth float64
	dash      []float64
	opacity   float64
	arrowDir  arrowDirection
}

func newAnnotStyle() annotationStyle {
	return annotationStyle{
		lineColor: color.Transparent,
		fillColor: color.Transparent,
		lineWidth: 1.5,
		opacity:   1.0,
	}
}

func applyAnnotOpts(s *annotationStyle, opts []AnnotationOption) {
	for _, o := range opts {
		o(s)
	}
}

// Label sets the text label shown alongside a line or text annotation.
func Label(text string) AnnotationOption { return func(a *annotationStyle) { a.label = text } }

// Dash sets the stroke-dasharray for a line annotation, e.g. Dash(6, 3).
func Dash(d ...float64) AnnotationOption { return func(a *annotationStyle) { a.dash = d } }

// AnnotColor sets the stroke/text color for a line or text annotation.
func AnnotColor(c color.Color) AnnotationOption { return func(a *annotationStyle) { a.lineColor = c } }

// FillColor sets the fill color for a span annotation.
// Accepts a CSS hex string; an 8-digit hex embeds alpha, e.g. "#ff000020".
func FillColor(hex string) AnnotationOption {
	c := color.Parse(hex)
	return func(a *annotationStyle) { a.fillColor = c }
}

// Opacity sets the overall opacity for an annotation (0.0–1.0).
func Opacity(o float64) AnnotationOption { return func(a *annotationStyle) { a.opacity = o } }

// ArrowDown places the arrowhead pointing downward toward the annotated point.
func ArrowDown() AnnotationOption { return func(a *annotationStyle) { a.arrowDir = arrowDown } }

// ArrowUp places the arrowhead pointing upward toward the annotated point.
func ArrowUp() AnnotationOption { return func(a *annotationStyle) { a.arrowDir = arrowUp } }

// ArrowLeft places the arrowhead pointing leftward toward the annotated point.
func ArrowLeft() AnnotationOption { return func(a *annotationStyle) { a.arrowDir = arrowLeft } }

// ArrowRight places the arrowhead pointing rightward toward the annotated point.
func ArrowRight() AnnotationOption { return func(a *annotationStyle) { a.arrowDir = arrowRight } }

// annotatable is the common draw interface for all annotation types.
type annotatable interface {
	draw(c *svg.Canvas, xScale, yScale scale.Scale,
		plotX, plotY, plotW, plotH float64, t theme.Theme, clipID string)
}

// resolveLineColor returns the style's line color, falling back to the theme text color.
func (s *annotationStyle) resolveLineColor(t theme.Theme) color.Color {
	if s.lineColor.IsTransparent() {
		return t.TextColor
	}
	return s.lineColor
}

// --- hLineAnnotation ---

type hLineAnnotation struct {
	y     float64
	style annotationStyle
}

func (h *hLineAnnotation) draw(c *svg.Canvas, _, yScale scale.Scale,
	plotX, _, plotW, _ float64, t theme.Theme, clipID string) {
	py := yScale.Map(h.y)
	s := &h.style
	col := s.resolveLineColor(t)

	c.BeginClipGroup(clipID)
	c.DrawLine(plotX, py, plotX+plotW, py, render.Style{
		Stroke:        col,
		StrokeWidth:   s.lineWidth,
		Dash:          s.dash,
		StrokeOpacity: s.opacity,
	})
	c.EndGroup()

	if s.label != "" {
		c.DrawText(plotX+plotW-4, py-4, s.label, render.TextStyle{
			Color:      col,
			FontSize:   t.TickFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "end",
			Baseline:   "auto",
		})
	}
}

// --- vLineAnnotation ---

type vLineAnnotation struct {
	x     float64
	style annotationStyle
}

func (v *vLineAnnotation) draw(c *svg.Canvas, xScale, _ scale.Scale,
	_, plotY, _, plotH float64, t theme.Theme, clipID string) {
	px := xScale.Map(v.x)
	s := &v.style
	col := s.resolveLineColor(t)

	c.BeginClipGroup(clipID)
	c.DrawLine(px, plotY, px, plotY+plotH, render.Style{
		Stroke:        col,
		StrokeWidth:   s.lineWidth,
		Dash:          s.dash,
		StrokeOpacity: s.opacity,
	})
	c.EndGroup()

	if s.label != "" {
		c.DrawText(px+4, plotY+4, s.label, render.TextStyle{
			Color:      col,
			FontSize:   t.TickFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "start",
			Baseline:   "hanging",
		})
	}
}

// --- hSpanAnnotation ---

type hSpanAnnotation struct {
	y1, y2 float64
	style  annotationStyle
}

func (h *hSpanAnnotation) draw(c *svg.Canvas, _, yScale scale.Scale,
	plotX, _, plotW, _ float64, t theme.Theme, clipID string) {
	py1 := yScale.Map(h.y1)
	py2 := yScale.Map(h.y2)
	if py1 > py2 {
		py1, py2 = py2, py1
	}
	s := &h.style
	col := s.resolveLineColor(t)
	fill := s.fillColor
	if fill.IsTransparent() {
		fill = col.WithAlpha(0.15)
	}

	c.BeginClipGroup(clipID)
	c.DrawRect(plotX, py1, plotW, py2-py1, 0, 0, render.Style{
		Fill:        fill,
		FillOpacity: s.opacity,
	})
	c.EndGroup()
}

// --- vSpanAnnotation ---

type vSpanAnnotation struct {
	x1, x2 float64
	style  annotationStyle
}

func (v *vSpanAnnotation) draw(c *svg.Canvas, xScale, _ scale.Scale,
	_, plotY, _, plotH float64, t theme.Theme, clipID string) {
	px1 := xScale.Map(v.x1)
	px2 := xScale.Map(v.x2)
	if px1 > px2 {
		px1, px2 = px2, px1
	}
	s := &v.style
	col := s.resolveLineColor(t)
	fill := s.fillColor
	if fill.IsTransparent() {
		fill = col.WithAlpha(0.15)
	}

	c.BeginClipGroup(clipID)
	c.DrawRect(px1, plotY, px2-px1, plotH, 0, 0, render.Style{
		Fill:        fill,
		FillOpacity: s.opacity,
	})
	c.EndGroup()
}

// --- textAnnotation ---

type textAnnotation struct {
	x, y  float64
	text  string
	style annotationStyle
}

func (ta *textAnnotation) draw(c *svg.Canvas, xScale, yScale scale.Scale,
	_, _, _, _ float64, t theme.Theme, _ string) {
	px := xScale.Map(ta.x)
	py := yScale.Map(ta.y)
	s := &ta.style
	col := s.resolveLineColor(t)

	// Determine label position from arrow direction.
	var lx, ly float64
	var anchor, baseline string
	switch s.arrowDir {
	case arrowDown: // arrowhead at bottom → label above
		lx, ly = px, py-annotOffset
		anchor, baseline = "middle", "auto"
	case arrowUp: // arrowhead at top → label below
		lx, ly = px, py+annotOffset
		anchor, baseline = "middle", "hanging"
	case arrowLeft: // arrowhead at left → label right
		lx, ly = px+annotOffset, py
		anchor, baseline = "start", "middle"
	case arrowRight: // arrowhead at right → label left
		lx, ly = px-annotOffset, py
		anchor, baseline = "end", "middle"
	default: // no arrow — label above-right of point
		lx, ly = px+8, py-8
		anchor, baseline = "start", "auto"
	}

	if s.arrowDir != arrowNone {
		// Shaft from label to data point.
		c.DrawLine(lx, ly, px, py, render.Style{
			Stroke:      col,
			StrokeWidth: 1.2,
		})
		// Arrowhead triangle at the data point.
		ah := annotArrowHead
		var head render.Path
		switch s.arrowDir {
		case arrowDown:
			head.MoveToP(px, py)
			head.LineToP(px-ah/2, py-ah)
			head.LineToP(px+ah/2, py-ah)
		case arrowUp:
			head.MoveToP(px, py)
			head.LineToP(px-ah/2, py+ah)
			head.LineToP(px+ah/2, py+ah)
		case arrowLeft:
			head.MoveToP(px, py)
			head.LineToP(px+ah, py-ah/2)
			head.LineToP(px+ah, py+ah/2)
		case arrowRight:
			head.MoveToP(px, py)
			head.LineToP(px-ah, py-ah/2)
			head.LineToP(px-ah, py+ah/2)
		}
		head.CloseP()
		c.DrawPath(head, render.Style{Fill: col, FillOpacity: 1.0})
	}

	c.DrawText(lx, ly, ta.text, render.TextStyle{
		Color:      col,
		FontSize:   t.TickFontSize,
		FontFamily: t.FontFamily,
		Anchor:     anchor,
		Baseline:   baseline,
	})
}

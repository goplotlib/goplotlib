package render

import "github.com/goplotlib/goplotlib/color"

// Op is a path operation type.
type Op int

const (
	MoveTo  Op = iota
	LineTo
	CubicTo // 3 sets of (x,y): cp1, cp2, end — uses 6 args
	Close
)

// PathCmd is a single path command.
type PathCmd struct {
	Op   Op
	Args [6]float64 // MoveTo/LineTo: x,y; CubicTo: cp1x,cp1y,cp2x,cp2y,ex,ey; Close: unused
}

// Path is a vector path composed of path commands.
type Path struct {
	Cmds []PathCmd
}

// MoveToP appends a MoveTo command.
func (p *Path) MoveToP(x, y float64) {
	p.Cmds = append(p.Cmds, PathCmd{Op: MoveTo, Args: [6]float64{x, y}})
}

// LineToP appends a LineTo command.
func (p *Path) LineToP(x, y float64) {
	p.Cmds = append(p.Cmds, PathCmd{Op: LineTo, Args: [6]float64{x, y}})
}

// CubicToP appends a CubicTo command.
func (p *Path) CubicToP(cp1x, cp1y, cp2x, cp2y, ex, ey float64) {
	p.Cmds = append(p.Cmds, PathCmd{Op: CubicTo, Args: [6]float64{cp1x, cp1y, cp2x, cp2y, ex, ey}})
}

// CloseP appends a Close command.
func (p *Path) CloseP() {
	p.Cmds = append(p.Cmds, PathCmd{Op: Close})
}

// Style defines the drawing style for a shape.
type Style struct {
	Stroke        color.Color
	StrokeWidth   float64
	StrokeOpacity float64 // 0-1, default 1
	Fill          color.Color
	FillOpacity   float64 // 0-1, default 0 (no fill)
	LineCap       string  // "round", "butt", "square"
	LineJoin      string  // "round", "miter", "bevel"
	Dash          []float64
	FilterID      string // optional SVG filter reference
	GradientID    string // optional SVG gradient fill reference
}

// TextStyle defines text drawing style.
type TextStyle struct {
	Color      color.Color
	FontSize   float64
	FontFamily string
	Anchor     string  // "start", "middle", "end"
	Baseline   string  // "auto", "middle", "hanging", "text-top"
	Bold       bool
	Italic     bool
	Opacity    float64
}

// GradStop is a gradient stop.
type GradStop struct {
	Offset  float64 // 0-1
	Color   color.Color
	Opacity float64 // 0-1
}

// Canvas is the drawing interface that all rendering backends implement.
type Canvas interface {
	// Drawing primitives
	DrawPath(p Path, s Style)
	DrawRect(x, y, w, h, rx, ry float64, s Style)
	DrawCircle(cx, cy, r float64, s Style)
	DrawText(x, y float64, text string, s TextStyle)
	DrawLine(x1, y1, x2, y2 float64, s Style)

	// Clipping
	BeginClip(id string, x, y, w, h float64)
	BeginClipGroup(clipID string)
	EndGroup()

	// Definitions (gradients, filters)
	DefineLinearGradient(id string, x1, y1, x2, y2 float64, stops []GradStop)
	DefineDropShadow(id string, dx, dy, stdDev float64, c color.Color)

	// Output
	Bytes() []byte
}

// Package svg implements the render.Canvas interface, producing SVG output.
package svg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
)

// Canvas is an SVG rendering backend.
type Canvas struct {
	width, height float64
	defs          strings.Builder
	body          strings.Builder
	groupDepth    int
}

// New creates a new SVG Canvas with the given dimensions.
func New(width, height float64) *Canvas {
	return &Canvas{
		width:  width,
		height: height,
	}
}

// ff formats a float64 compactly, capping at 2 decimal places.
func ff(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	// Trim trailing zeros after decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// xmlEscape escapes special XML characters in text.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// styleAttr builds SVG style attributes for a shape.
func styleAttr(s render.Style) string {
	var parts []string

	// Stroke
	strokeOpacity := s.StrokeOpacity
	if strokeOpacity == 0 && s.Stroke.A > 0 {
		strokeOpacity = 1.0
	}
	if strokeOpacity > 0 && !s.Stroke.IsTransparent() {
		parts = append(parts, fmt.Sprintf(`stroke="%s"`, s.Stroke.SVGColor()))
		parts = append(parts, fmt.Sprintf(`stroke-width="%s"`, ff(s.StrokeWidth)))
		if strokeOpacity < 1.0 {
			parts = append(parts, fmt.Sprintf(`stroke-opacity="%s"`, ff(strokeOpacity)))
		}
	} else {
		parts = append(parts, `stroke="none"`)
	}

	// Fill
	if s.GradientID != "" {
		parts = append(parts, fmt.Sprintf(`fill="url(#%s)"`, s.GradientID))
		if s.FillOpacity > 0 && s.FillOpacity < 1.0 {
			parts = append(parts, fmt.Sprintf(`fill-opacity="%s"`, ff(s.FillOpacity)))
		}
	} else if s.FillOpacity > 0 && !s.Fill.IsTransparent() {
		fillColor := s.Fill
		// If the fill color itself has alpha, use it directly
		if fillColor.A < 255 {
			parts = append(parts, fmt.Sprintf(`fill="%s"`, fillColor.SVGColor()))
		} else {
			parts = append(parts, fmt.Sprintf(`fill="%s"`, fillColor.Hex()))
			if s.FillOpacity < 1.0 {
				parts = append(parts, fmt.Sprintf(`fill-opacity="%s"`, ff(s.FillOpacity)))
			}
		}
	} else {
		parts = append(parts, `fill="none"`)
	}

	// Line cap and join
	if s.LineCap != "" {
		parts = append(parts, fmt.Sprintf(`stroke-linecap="%s"`, s.LineCap))
	}
	if s.LineJoin != "" {
		parts = append(parts, fmt.Sprintf(`stroke-linejoin="%s"`, s.LineJoin))
	}

	// Dash array
	if len(s.Dash) > 0 {
		dashParts := make([]string, len(s.Dash))
		for i, d := range s.Dash {
			dashParts[i] = ff(d)
		}
		parts = append(parts, fmt.Sprintf(`stroke-dasharray="%s"`, strings.Join(dashParts, ",")))
	}

	// Filter
	if s.FilterID != "" {
		parts = append(parts, fmt.Sprintf(`filter="url(#%s)"`, s.FilterID))
	}

	return strings.Join(parts, " ")
}

// pathData converts a Path to an SVG path d attribute value.
func pathData(p render.Path) string {
	var sb strings.Builder
	for i, cmd := range p.Cmds {
		if i > 0 {
			sb.WriteByte(' ')
		}
		switch cmd.Op {
		case render.MoveTo:
			fmt.Fprintf(&sb, "M%s %s", ff(cmd.Args[0]), ff(cmd.Args[1]))
		case render.LineTo:
			fmt.Fprintf(&sb, "L%s %s", ff(cmd.Args[0]), ff(cmd.Args[1]))
		case render.CubicTo:
			fmt.Fprintf(&sb, "C%s %s %s %s %s %s",
				ff(cmd.Args[0]), ff(cmd.Args[1]),
				ff(cmd.Args[2]), ff(cmd.Args[3]),
				ff(cmd.Args[4]), ff(cmd.Args[5]))
		case render.Close:
			sb.WriteByte('Z')
		}
	}
	return sb.String()
}

// DrawPath renders a path to the SVG body.
func (c *Canvas) DrawPath(p render.Path, s render.Style) {
	d := pathData(p)
	if d == "" {
		return
	}
	fmt.Fprintf(&c.body, `<path d="%s" %s/>`, d, styleAttr(s))
	c.body.WriteByte('\n')
}

// DrawRect renders a rectangle to the SVG body.
func (c *Canvas) DrawRect(x, y, w, h, rx, ry float64, s render.Style) {
	rxAttr := ""
	if rx > 0 || ry > 0 {
		rxAttr = fmt.Sprintf(` rx="%s" ry="%s"`, ff(rx), ff(ry))
	}
	fmt.Fprintf(&c.body, `<rect x="%s" y="%s" width="%s" height="%s"%s %s/>`,
		ff(x), ff(y), ff(w), ff(h), rxAttr, styleAttr(s))
	c.body.WriteByte('\n')
}

// DrawCircle renders a circle to the SVG body.
func (c *Canvas) DrawCircle(cx, cy, r float64, s render.Style) {
	fmt.Fprintf(&c.body, `<circle cx="%s" cy="%s" r="%s" %s/>`,
		ff(cx), ff(cy), ff(r), styleAttr(s))
	c.body.WriteByte('\n')
}

// DrawText renders text to the SVG body.
func (c *Canvas) DrawText(x, y float64, text string, s render.TextStyle) {
	var attrs []string

	attrs = append(attrs, fmt.Sprintf(`x="%s" y="%s"`, ff(x), ff(y)))

	if s.FontSize > 0 {
		attrs = append(attrs, fmt.Sprintf(`font-size="%s"`, ff(s.FontSize)))
	}
	if s.FontFamily != "" {
		attrs = append(attrs, fmt.Sprintf(`font-family="%s"`, s.FontFamily))
	}
	if s.Bold {
		attrs = append(attrs, `font-weight="bold"`)
	}
	if s.Italic {
		attrs = append(attrs, `font-style="italic"`)
	}
	if s.Anchor != "" {
		attrs = append(attrs, fmt.Sprintf(`text-anchor="%s"`, s.Anchor))
	}
	if s.Baseline != "" {
		attrs = append(attrs, fmt.Sprintf(`dominant-baseline="%s"`, s.Baseline))
	}
	if !s.Color.IsTransparent() {
		attrs = append(attrs, fmt.Sprintf(`fill="%s"`, s.Color.SVGColor()))
	}
	if s.Opacity > 0 && s.Opacity < 1.0 {
		attrs = append(attrs, fmt.Sprintf(`opacity="%s"`, ff(s.Opacity)))
	}

	fmt.Fprintf(&c.body, `<text %s>%s</text>`, strings.Join(attrs, " "), xmlEscape(text))
	c.body.WriteByte('\n')
}

// DrawLine renders a line to the SVG body.
func (c *Canvas) DrawLine(x1, y1, x2, y2 float64, s render.Style) {
	fmt.Fprintf(&c.body, `<line x1="%s" y1="%s" x2="%s" y2="%s" %s/>`,
		ff(x1), ff(y1), ff(x2), ff(y2), styleAttr(s))
	c.body.WriteByte('\n')
}

// BeginClip defines a rectangular clip path in the SVG defs.
func (c *Canvas) BeginClip(id string, x, y, w, h float64) {
	fmt.Fprintf(&c.defs, `<clipPath id="%s"><rect x="%s" y="%s" width="%s" height="%s"/></clipPath>`,
		id, ff(x), ff(y), ff(w), ff(h))
	c.defs.WriteByte('\n')
}

// BeginClipGroup starts a group with a clip path applied.
func (c *Canvas) BeginClipGroup(clipID string) {
	fmt.Fprintf(&c.body, `<g clip-path="url(#%s)">`, clipID)
	c.body.WriteByte('\n')
	c.groupDepth++
}

// EndGroup closes a group element.
func (c *Canvas) EndGroup() {
	c.body.WriteString(`</g>`)
	c.body.WriteByte('\n')
	if c.groupDepth > 0 {
		c.groupDepth--
	}
}

// DefineLinearGradient adds a linear gradient definition to the SVG defs.
func (c *Canvas) DefineLinearGradient(id string, x1, y1, x2, y2 float64, stops []render.GradStop) {
	fmt.Fprintf(&c.defs, `<linearGradient id="%s" x1="%s" y1="%s" x2="%s" y2="%s" gradientUnits="userSpaceOnUse">`,
		id, ff(x1), ff(y1), ff(x2), ff(y2))
	c.defs.WriteByte('\n')
	for _, stop := range stops {
		fmt.Fprintf(&c.defs, `  <stop offset="%s" stop-color="%s" stop-opacity="%s"/>`,
			ff(stop.Offset), stop.Color.Hex(), ff(stop.Opacity))
		c.defs.WriteByte('\n')
	}
	c.defs.WriteString(`</linearGradient>`)
	c.defs.WriteByte('\n')
}

// DefineDropShadow adds a drop shadow filter definition to the SVG defs.
func (c *Canvas) DefineDropShadow(id string, dx, dy, stdDev float64, col color.Color) {
	fmt.Fprintf(&c.defs, `<filter id="%s" x="-20%%" y="-20%%" width="140%%" height="140%%">`, id)
	c.defs.WriteByte('\n')
	fmt.Fprintf(&c.defs, `  <feDropShadow dx="%s" dy="%s" stdDeviation="%s" flood-color="%s" flood-opacity="0.3"/>`,
		ff(dx), ff(dy), ff(stdDev), col.Hex())
	c.defs.WriteByte('\n')
	c.defs.WriteString(`</filter>`)
	c.defs.WriteByte('\n')
}

// DrawTextRotated draws text rotated by the given degrees around the point (cx, cy).
func (c *Canvas) DrawTextRotated(cx, cy float64, text string, s render.TextStyle, degrees float64) {
	var attrs []string

	attrs = append(attrs, fmt.Sprintf(`x="%s" y="%s"`, ff(cx), ff(cy)))
	attrs = append(attrs, fmt.Sprintf(`transform="rotate(%s,%s,%s)"`, ff(degrees), ff(cx), ff(cy)))

	if s.FontSize > 0 {
		attrs = append(attrs, fmt.Sprintf(`font-size="%s"`, ff(s.FontSize)))
	}
	if s.FontFamily != "" {
		attrs = append(attrs, fmt.Sprintf(`font-family="%s"`, s.FontFamily))
	}
	if s.Bold {
		attrs = append(attrs, `font-weight="bold"`)
	}
	if s.Italic {
		attrs = append(attrs, `font-style="italic"`)
	}
	if s.Anchor != "" {
		attrs = append(attrs, fmt.Sprintf(`text-anchor="%s"`, s.Anchor))
	}
	if s.Baseline != "" {
		attrs = append(attrs, fmt.Sprintf(`dominant-baseline="%s"`, s.Baseline))
	}
	if !s.Color.IsTransparent() {
		attrs = append(attrs, fmt.Sprintf(`fill="%s"`, s.Color.SVGColor()))
	}

	fmt.Fprintf(&c.body, `<text %s>%s</text>`, strings.Join(attrs, " "), xmlEscape(text))
	c.body.WriteByte('\n')
}

// Bytes returns the complete SVG document as bytes.
func (c *Canvas) Bytes() []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s">`,
		ff(c.width), ff(c.height))
	sb.WriteByte('\n')

	defsContent := c.defs.String()
	if defsContent != "" {
		sb.WriteString("<defs>\n")
		sb.WriteString(defsContent)
		sb.WriteString("</defs>\n")
	}

	sb.WriteString(c.body.String())
	sb.WriteString(`</svg>`)
	sb.WriteByte('\n')

	return []byte(sb.String())
}

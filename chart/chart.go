package chart

import (
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/colormap"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// Chart is the interface all chart types implement.
type Chart interface {
	Draw(canvas render.Canvas, xScale, yScale scale.Scale)
	DataRange() (xMin, xMax, yMin, yMax float64)
	Label() string
	Color() color.Color
}

// LineStyle configures a Line or Area chart.
// Zero values are sensible: no label, auto color, 2.5px line, markers size 5, no fill.
type LineStyle struct {
	Label       string
	Color       color.Color // zero = auto from palette
	Opacity     float64     // 0 → 1.0
	Smooth      bool
	Fill        bool
	FillOpacity float64   // 0 → 0.15 when Fill=true
	LineWidth   float64   // 0 → 2.5
	MarkerSize  float64   // 0 → 5
	MarkerShape string    // "" → "circle"
	Dash        []float64
}

// ScatterStyle configures a Scatter chart.
type ScatterStyle struct {
	Label       string
	Color       color.Color
	Opacity     float64 // 0 → 1.0
	MarkerSize  float64 // 0 → 5
	MarkerShape string  // "" → "circle"
}

// BarStyle configures a Bar or HBar chart.
type BarStyle struct {
	Label      string
	Color      color.Color
	Opacity    float64 // 0 → 1.0
	SquareBars bool    // set true to disable the default rounded corners
}

// StepStyle configures a Step chart.
type StepStyle struct {
	Label     string
	Color     color.Color
	Opacity   float64   // 0 → 1.0
	Mode      string    // "pre", "post" (default), "mid"
	Fill      bool
	LineWidth float64   // 0 → 2.5
	Dash      []float64
}

// StackedBarStyle configures a StackedBar chart.
type StackedBarStyle struct {
	Opacity float64 // 0 → 1.0
}

// HistogramStyle configures a Histogram chart.
type HistogramStyle struct {
	Label      string
	Color      color.Color
	Opacity    float64   // 0 → 1.0
	Bins       int       // 0 → auto (Sturges' rule)
	BinEdges   []float64
	Normalize  bool
	Cumulative bool
}

// BubbleStyle configures a Bubble chart.
type BubbleStyle struct {
	Label   string
	Color   color.Color
	Opacity float64  // 0 → 1.0
	SizeMin float64  // 0 → 4
	SizeMax float64  // 0 → 30
	Labels  []string // per-point labels drawn inside each bubble; skipped when bubble is too small to fit
}

// PieStyle configures a Pie or Donut chart.
type PieStyle struct {
	DonutRadius   float64 // 0 = solid pie; e.g. 0.55 = donut
	ExplodeIdx    int
	ExplodeOffset float64
}

// HeatmapStyle configures a Heatmap chart.
type HeatmapStyle struct {
	RowLabels      []string
	ColLabels      []string
	ColorMap       colormap.ColorMap // nil → Viridis
	CellLabels     bool
	DivergingScale bool
}

// BoxStyle configures a Box plot.
type BoxStyle struct {
	Label   string
	Color   color.Color
	Labels  []string // per-group x-axis category labels
}

// LegendEntry represents one item in the legend.
type LegendEntry struct {
	Label string
	Col   color.Color
	IsBar bool
}

// MultiLegend is implemented by chart types that contribute multiple legend entries.
type MultiLegend interface {
	LegendEntries() []LegendEntry
}

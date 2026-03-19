---
title: API Reference
weight: 5
description: Complete reference for all Figure options, Axes methods, and chart style structs.
---

## Figure options

Passed to `plot.New(opts...)` when creating a Figure.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithWidth(w)` | `float64` | `900` | Figure width in pixels |
| `WithHeight(h)` | `float64` | `550` | Figure height in pixels |
| `WithTheme(t)` | `theme.Theme` | `theme.Light` | Visual theme |
| `WithTitle(s)` | `string` | `""` | Figure-level title drawn above all axes |
| `WithSpacing(px)` | `float64` | `10` | Gap in pixels between subplot cells |

**Example:**

```go
fig := plot.New(
	plot.WithWidth(1200),
	plot.WithHeight(600),
	plot.WithTheme(theme.Dark),
	plot.WithTitle("Annual Overview"),
)
```

**Subplot layout:**

```go
// Shorthand — returns a [rows][cols]*Axes grid
axes := fig.SubPlots(2, 2)
axes[0][0].Line(xs, ys, chart.LineStyle{})
axes[0][1].Bar(cats, vals, chart.BarStyle{})

// Explicit placement
fig.SetLayout(1, 2)
ax1 := fig.AddAxes(plot.At(0, 0))
ax2 := fig.AddAxes(plot.At(0, 1))
```

---

## Axes methods

Returned by `fig.AddAxes()`. All methods return `*Axes` for chaining.

### Chart series

| Method | Description |
|--------|-------------|
| `Line(xs, ys, chart.LineStyle)` | Add a line series |
| `Area(xs, ys, chart.LineStyle)` | Add a filled area series (shorthand for Line with fill enabled) |
| `Step(xs, ys, chart.StepStyle)` | Add a step (staircase) series |
| `Bar(cats, vals, chart.BarStyle)` | Add a bar series with categorical x-axis; multiple calls produce grouped bars |
| `HBar(cats, vals, chart.BarStyle)` | Add a horizontal bar series |
| `StackedBar(cats, series, labels, chart.StackedBarStyle)` | Add a stacked bar chart |
| `StackedArea(xs, series, labels)` | Add a stacked area chart |
| `Scatter(xs, ys, chart.ScatterStyle)` | Add a scatter series |
| `Bubble(xs, ys, sizes, chart.BubbleStyle)` | Add a bubble chart (scatter with radius encoding) |
| `Histogram(values, chart.HistogramStyle)` | Bin values and render as adjacent bars |
| `Heatmap(matrix, chart.HeatmapStyle)` | Render a scalar matrix as a colour grid |
| `Pie(values, labels, chart.PieStyle)` | Add a pie or donut chart |
| `Box(groups, chart.BoxStyle)` | Add a box plot; each inner slice is one group |
| `ErrorBars(xs, ys, lo, hi, col)` | Overlay vertical error bars (symmetric or asymmetric) |
| `ErrorBarsX(xs, ys, xErr, col)` | Overlay horizontal error bars |

### Axis configuration

| Method | Description |
|--------|-------------|
| `SetTitle(s)` | Set the axes title (drawn inside the plot area, above the data) |
| `SetXLabel(s)` | Set the x-axis label |
| `SetYLabel(s)` | Set the y-axis label (rendered rotated 90°) |
| `SetXTicks(n)` | Approximate number of x-axis tick marks (default: 6) |
| `SetYTicks(n)` | Approximate number of y-axis tick marks (default: 6) |
| `SetXLim(min, max)` | Fix the visible x-axis range; pass `plot.Auto` to keep a bound auto-ranged |
| `SetYLim(min, max)` | Fix the visible y-axis range; pass `plot.Auto` to keep a bound auto-ranged |
| `SetXTickFormat(f)` | Custom formatter `func(float64) string` for x-axis tick labels |
| `SetYTickFormat(f)` | Custom formatter `func(float64) string` for y-axis tick labels |
| `SetXScale(f)` | Override x-axis scale, e.g. `scale.Log(10)` |
| `SetYScale(f)` | Override y-axis scale, e.g. `scale.Log(10)` |

### Secondary y-axis

A second independent y-axis can be added on the right side. Call `.OnY2()` immediately
after any series method to assign that series to the secondary axis:

```go
ax.Bar(months, revenue, chart.BarStyle{Label: "Revenue ($M)"})
ax.Line(months, growth, chart.LineStyle{Label: "Growth (%)"}).OnY2()
ax.SetYLabel("Revenue ($M)").SetY2Label("Growth (%)")
```

| Method | Description |
|--------|-------------|
| `OnY2()` | Assign the most recently added series to the secondary (right) y-axis |
| `SetY2Label(s)` | Set the secondary y-axis label |
| `SetY2Lim(min, max)` | Fix the secondary y-axis range; pass `plot.Auto` for auto-ranging |
| `SetY2Ticks(n)` | Approximate number of secondary y-axis tick marks (default: 6) |

### Legend

| Method | Description |
|--------|-------------|
| `Legend(pos)` | Set legend position (see positions below) |
| `NoLegend()` | Hide the legend entirely |

**Legend positions:**

```go
plot.LegendTopRight      // default — inside, top-right corner
plot.LegendTopLeft       // inside, top-left corner
plot.LegendBottomRight   // inside, bottom-right corner
plot.LegendBottomLeft    // inside, bottom-left corner
plot.LegendOutsideRight  // outside the plot, vertical list on the right
plot.LegendOutsideBottom // outside the plot, horizontal row below
plot.LegendNone          // same as NoLegend()
```

### Annotations

| Method | Description |
|--------|-------------|
| `HLine(y, opts...)` | Draw a horizontal reference line at y in data coordinates |
| `VLine(x, opts...)` | Draw a vertical reference line at x in data coordinates |
| `HSpan(y1, y2, opts...)` | Draw a horizontal shaded band between y1 and y2 |
| `VSpan(x1, x2, opts...)` | Draw a vertical shaded band between x1 and x2 |
| `Annotate(x, y, text, opts...)` | Place a text label at data point (x, y) with an optional arrow |

**Example:**

```go
ax := fig.AddAxes()
ax.SetTitle("Monthly Trend").
   SetXLabel("Month").
   SetYLabel("Value").
   Legend(plot.LegendBottomRight).
   SetXTicks(12).
   SetYTicks(8)
```

**Axis limits:**

```go
ax.SetXLim(0, 100)           // fix both bounds
ax.SetYLim(plot.Auto, 100)   // fix only upper bound
```

**Tick format examples:**

```go
ax.SetYTickFormat(func(v float64) string { return fmt.Sprintf("$%.0f", v) })
ax.SetYTickFormat(func(v float64) string { return fmt.Sprintf("%.0f%%", v*100) })
ax.SetXTickFormat(func(v float64) string {
    return time.Unix(int64(v), 0).Format("Jan 2006")
})
```

**Log scale:**

```go
import "github.com/goplotlib/goplotlib/scale"

ax.SetYScale(scale.Log(10))  // base-10 log y-axis
ax.SetXScale(scale.Log(2))   // base-2 log x-axis
```

---

## Annotation options

Passed to `HLine`, `VLine`, `HSpan`, `VSpan`, and `Annotate` as variadic `plot.AnnotationOption` arguments.

| Option | Description |
|--------|-------------|
| `Label(s)` | Text label shown near a line annotation |
| `Dash(d...)` | Stroke-dasharray pattern, e.g. `Dash(6, 3)` |
| `AnnotColor(c)` | Stroke/text color for a line or text annotation |
| `FillColor(hex)` | Fill color for a span; accepts CSS hex with optional alpha, e.g. `"#ffff0040"` |
| `Opacity(o)` | Overall opacity 0.0–1.0 |
| `ArrowDown()` | Arrow pointing downward to the annotated point |
| `ArrowUp()` | Arrow pointing upward to the annotated point |
| `ArrowLeft()` | Arrow pointing leftward to the annotated point |
| `ArrowRight()` | Arrow pointing rightward to the annotated point |

**Examples:**

```go
ax.HLine(0.5, plot.Label("threshold"), plot.Dash(6, 3))
ax.VLine(math.Pi, plot.Label("π"), plot.AnnotColor(color.Purple))
ax.HSpan(-0.5, 0.5, plot.FillColor("#4C72B0"), plot.Opacity(0.1))
ax.Annotate(peakX, peakY, "local max", plot.ArrowDown())
```

---

## Chart style structs

Each chart method accepts a typed style struct. Zero values select sensible defaults.

### LineStyle

Used by `Line` and `Area`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall series opacity, 0.0–1.0 |
| `Smooth` | `bool` | `false` | Catmull-Rom spline smoothing |
| `Fill` | `bool` | `false` | Gradient area fill under the line |
| `FillOpacity` | `float64` | `0.15` | Fill area opacity, 0.0–1.0 |
| `LineWidth` | `float64` | `2.5` | Stroke width in pixels |
| `MarkerSize` | `float64` | `0` | Marker radius in pixels; `0` disables markers |
| `MarkerShape` | `string` | `"circle"` | `"circle"`, `"square"`, `"diamond"`, `"triangle"` |
| `Dash` | `[]float64` | none | SVG stroke-dasharray values |

### ScatterStyle

Used by `Scatter`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall marker opacity, 0.0–1.0 |
| `MarkerSize` | `float64` | `5` | Marker radius in pixels |
| `MarkerShape` | `string` | `"circle"` | `"circle"`, `"square"`, `"diamond"`, `"triangle"` |

### BarStyle

Used by `Bar` and `HBar`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall bar opacity, 0.0–1.0 |
| `SquareBars` | `bool` | `false` | Set `true` to disable the default rounded corners |

### StepStyle

Used by `Step`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall series opacity, 0.0–1.0 |
| `Mode` | `string` | `"post"` | Step interpolation: `"pre"`, `"post"`, `"mid"` |
| `Fill` | `bool` | `false` | Fill area under the step line |
| `LineWidth` | `float64` | `2.5` | Stroke width in pixels |
| `Dash` | `[]float64` | none | SVG stroke-dasharray pattern |

### StackedBarStyle

Used by `StackedBar`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Opacity` | `float64` | `1.0` | Overall series opacity, 0.0–1.0 |

### HistogramStyle

Used by `Histogram`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall bar opacity, 0.0–1.0 |
| `Bins` | `int` | auto | Number of equal-width bins (0 = Sturges' rule) |
| `BinEdges` | `[]float64` | — | Explicit bin edges (overrides `Bins`) |
| `Normalize` | `bool` | `false` | Normalise to density (area = 1) |
| `Cumulative` | `bool` | `false` | Show cumulative count or density |

### BubbleStyle

Used by `Bubble`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Series label for the legend |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Opacity` | `float64` | `1.0` | Overall series opacity (fill capped at 0.65) |
| `SizeMin` | `float64` | `4` | Minimum pixel radius for bubbles |
| `SizeMax` | `float64` | `30` | Maximum pixel radius for bubbles |

### PieStyle

Used by `Pie`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `DonutRadius` | `float64` | `0` | Inner radius as fraction of outer; `0` = solid pie |
| `ExplodeIdx` | `int` | `0` | Index of segment to offset outward |
| `ExplodeOffset` | `float64` | `0` | Outward offset in pixels for the exploded segment |

### HeatmapStyle

Used by `Heatmap`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `RowLabels` | `[]string` | index | Labels for each row (y-axis) |
| `ColLabels` | `[]string` | index | Labels for each column (x-axis) |
| `ColorMap` | `colormap.ColorMap` | `Viridis` | Colour map to use |
| `CellLabels` | `bool` | `false` | Print the numeric value inside each cell |
| `DivergingScale` | `bool` | `false` | Centre normalisation at zero (for diverging maps) |

### BoxStyle

Used by `Box`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Label` | `string` | `""` | Legend label for the series |
| `Color` | `color.Color` | palette | Override the automatic palette color |
| `Labels` | `[]string` | index | Per-group x-axis category labels |

---

## Scale factories

Import `"github.com/goplotlib/goplotlib/scale"` and pass the result to `SetXScale` or `SetYScale`:

| Factory | Description |
|---------|-------------|
| `scale.Log(base)` | Logarithmic scale; domain must be strictly positive |

```go
ax.SetYScale(scale.Log(10))   // base-10 log y-axis
ax.SetXScale(scale.Log(2))    // base-2 log x-axis
```

---

## Colourmap constants

Import `"github.com/goplotlib/goplotlib/colormap"`:

| Constant | Type | Description |
|----------|------|-------------|
| `colormap.Viridis` | `ColorMap` | Perceptually uniform: blue → green → yellow |
| `colormap.Plasma` | `ColorMap` | Blue → purple → orange |
| `colormap.Greys` | `ColorMap` | White → black |
| `colormap.Blues` | `ColorMap` | White → dark blue |
| `colormap.RdBu` | `ColorMap` | Red → white → blue (diverging) |

---

## Color helpers

```go
import "github.com/goplotlib/goplotlib/color"

c1 := color.Parse("#4C72B0")          // hex string
c2 := color.RGB(70, 130, 180)         // RGB uint8
c3 := color.RGBA(70, 130, 180, 200)   // RGBA uint8
c4 := color.SteelBlue                 // predefined constant
c5 := c1.WithAlpha(0.5)               // adjust opacity
c6 := c1.Lighten(0.3)                 // mix with white
c7 := c1.Darken(0.2)                  // mix with black
```

**Predefined constants:** `White`, `Black`, `Red`, `Green`, `Blue`, `SteelBlue`, `Coral`,
`Orange`, `Purple`, `Gray`, `LightGray`, `DarkGray`, `Transparent`.

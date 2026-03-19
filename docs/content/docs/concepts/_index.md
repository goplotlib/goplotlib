---
title: Core Concepts
weight: 2
description: Understand the Figure / Axes / Chart composition model.
---

## The composition model

goplotlib uses three types arranged in a simple hierarchy:

```
Figure (size, theme)
  └── Axes (plot area, scales, labels)
        ├── LineChart (xs, ys, style)
        ├── BarChart  (categories, values, style)
        └── ScatterChart (xs, ys, style)
```

### Figure

`plot.Figure` is the root container. It defines:

- **Width and height** in pixels (`WithWidth`, `WithHeight`). Default: 900×550.
- **Theme** — the visual style applied to every Axes it contains (`WithTheme`).
  Default: `theme.Light`.
- **Figure-level title** shown centered above all Axes (`WithTitle`).

You create a Figure with `plot.New(opts...)` and render it with `fig.SVG()` (returns
`[]byte`) or `fig.SaveSVG(path)`.

### Axes

`plot.Axes` is a single rectangular plot area. One Axes occupies the full Figure.
An Axes is responsible for:

- Computing linear (or categorical) data scales from the union of all series' ranges.
- Drawing the grid, spines, tick marks, and tick labels.
- Rendering each series clipped to the plot rectangle.
- Drawing the legend when multiple labeled series are present.

Add an Axes by calling `fig.AddAxes()`. Configure it with method chaining:

```go
ax := fig.AddAxes()
ax.SetTitle("My chart").
   SetXLabel("time").
   SetYLabel("value").
   Legend(plot.LegendOutsideRight)
```

#### Legend positioning

The legend appears automatically when multiple labeled series are present. Six positions
are available — four inside the plot area and two outside it:

| Position | Behaviour |
|----------|-----------|
| `LegendTopRight` (default) | Inside, top-right corner |
| `LegendTopLeft` | Inside, top-left corner |
| `LegendBottomRight` | Inside, bottom-right corner |
| `LegendBottomLeft` | Inside, bottom-left corner |
| `LegendOutsideRight` | Vertical list to the right of the plot — never overlaps data |
| `LegendOutsideBottom` | Horizontal row below the plot — never overlaps data |

Use `ax.NoLegend()` (or `plot.LegendNone`) to hide it entirely.

### Chart (series)

A chart is a single data series. You never construct chart types directly — instead you call
the corresponding `Axes` method, which creates the series and registers it:

| Method | Series type |
|--------|-------------|
| `ax.Line(xs, ys, chart.LineStyle{...})` | `chart.LineChart` |
| `ax.Bar(categories, values, chart.BarStyle{...})` | `chart.BarChart` |
| `ax.Scatter(xs, ys, chart.ScatterStyle{...})` | `chart.ScatterChart` |

All three methods return the `*Axes` so you can chain multiple calls:

```go
ax.Line(xs, ys1, chart.LineStyle{Label: "series A"}).
   Line(xs, ys2, chart.LineStyle{Label: "series B"})
```

### Colors and themes

When a series is added without an explicit `Color` field in its style struct, the Figure
automatically assigns the next color from the active theme's palette. The `Classic`
palette is used by `Light`, `Dark`, and `Minimal` themes. The `FiveThirtyEight` theme has
its own palette.

You can override any series color explicitly:

```go
ax.Line(xs, ys, chart.LineStyle{Color: color.Parse("#e63946")})
```

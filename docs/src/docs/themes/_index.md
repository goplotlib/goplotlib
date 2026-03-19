---
title: Themes
weight: 4
description: Four built-in visual themes — Light, Dark, Minimal, and FiveThirtyEight.
---

goplotlib ships with four carefully designed themes. Pass a theme to `plot.WithTheme(t)` when
creating a Figure:

```go
import (
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

fig := plot.New(
	plot.WithWidth(860),
	plot.WithHeight(420),
	plot.WithTheme(theme.Dark),
)
```

---

## Light (default)

Clean white background with a light grey plot area. Left and bottom spines are shown.
Uses the **Classic** palette — saturated but not garish. Good for reports, notebooks,
and light-mode UIs.

```go run:theme_light.svg
package main

import (
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

func main() {
	// @nodoc
	n := 60
	xs, ys1, ys2 := make([]float64, n), make([]float64, n), make([]float64, n)
	for i := range xs {
		t := float64(i) / float64(n-1) * 6
		xs[i] = t
		ys1[i] = math.Sin(t)*1.5 + math.Sin(2*t)*0.4
		ys2[i] = math.Cos(t)*1.2 + math.Sin(3*t)*0.3
	}
	// @doc
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(350), plot.WithTheme(theme.Light))
	ax := fig.AddAxes()
	ax.Line(xs, ys1, chart.LineStyle{Label: "series A", Smooth: true, Fill: true, LineWidth: 2.5})
	ax.Line(xs, ys2, chart.LineStyle{Label: "series B", Smooth: true, LineWidth: 2.5})
	ax.SetTitle("Theme Showcase").SetXLabel("x").SetYLabel("y")
	// @nodoc
	fig.SaveSVG("docs/static/img/theme_light.svg")
	// @doc
}
```

![Light theme](/img/theme_light.svg)

---

## Dark

Deep navy background, matching plot area. Grid lines are subtle dark-blue stripes.
Text and tick labels use a soft lavender tone. Pairs well with dashboards and dark-mode
applications. Also uses the **Classic** palette.

```go run:theme_dark.svg
package main

import (
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

func main() {
	// @nodoc
	n := 60
	xs, ys1, ys2 := make([]float64, n), make([]float64, n), make([]float64, n)
	for i := range xs {
		t := float64(i) / float64(n-1) * 6
		xs[i] = t
		ys1[i] = math.Sin(t)*1.5 + math.Sin(2*t)*0.4
		ys2[i] = math.Cos(t)*1.2 + math.Sin(3*t)*0.3
	}
	// @doc
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(350), plot.WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Line(xs, ys1, chart.LineStyle{Label: "series A", Smooth: true, Fill: true, LineWidth: 2.5})
	ax.Line(xs, ys2, chart.LineStyle{Label: "series B", Smooth: true, LineWidth: 2.5})
	ax.SetTitle("Theme Showcase").SetXLabel("x").SetYLabel("y")
	// @nodoc
	fig.SaveSVG("docs/static/img/theme_dark.svg")
	// @doc
}
```

![Dark theme](/img/theme_dark.svg)

---

## Minimal

White background and plot area, no left spine — just a bottom spine. Grid lines are dashed
and very light. Produces clean, presentation-style charts with minimal visual noise.
Uses the **Classic** palette.

```go run:theme_minimal.svg
package main

import (
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

func main() {
	// @nodoc
	n := 60
	xs, ys1, ys2 := make([]float64, n), make([]float64, n), make([]float64, n)
	for i := range xs {
		t := float64(i) / float64(n-1) * 6
		xs[i] = t
		ys1[i] = math.Sin(t)*1.5 + math.Sin(2*t)*0.4
		ys2[i] = math.Cos(t)*1.2 + math.Sin(3*t)*0.3
	}
	// @doc
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(350), plot.WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Line(xs, ys1, chart.LineStyle{Label: "series A", Smooth: true, Fill: true, LineWidth: 2.5})
	ax.Line(xs, ys2, chart.LineStyle{Label: "series B", Smooth: true, LineWidth: 2.5})
	ax.SetTitle("Theme Showcase").SetXLabel("x").SetYLabel("y")
	// @nodoc
	fig.SaveSVG("docs/static/img/theme_minimal.svg")
	// @doc
}
```

![Minimal theme](/img/theme_minimal.svg)

---

## FiveThirtyEight

Inspired by the data journalism style of FiveThirtyEight. Grey background, bold white grid
lines (no spines), larger title font. Uses the **FTEPalette** — bright, high-contrast colours
that read well at small sizes.

```go run:theme_538.svg
package main

import (
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

func main() {
	// @nodoc
	n := 60
	xs, ys1, ys2 := make([]float64, n), make([]float64, n), make([]float64, n)
	for i := range xs {
		t := float64(i) / float64(n-1) * 6
		xs[i] = t
		ys1[i] = math.Sin(t)*1.5 + math.Sin(2*t)*0.4
		ys2[i] = math.Cos(t)*1.2 + math.Sin(3*t)*0.3
	}
	// @doc
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(350), plot.WithTheme(theme.FiveThirtyEight))
	ax := fig.AddAxes()
	ax.Line(xs, ys1, chart.LineStyle{Label: "series A", Smooth: true, Fill: true, LineWidth: 2.5})
	ax.Line(xs, ys2, chart.LineStyle{Label: "series B", Smooth: true, LineWidth: 2.5})
	ax.SetTitle("Theme Showcase").SetXLabel("x").SetYLabel("y")
	// @nodoc
	fig.SaveSVG("docs/static/img/theme_538.svg")
	// @doc
}
```

![FiveThirtyEight theme](/img/theme_538.svg)

---

## Palette colors

| Palette | Colors |
|---------|--------|
| Classic (Light, Dark, Minimal) | `#4C72B0` `#DD8452` `#55A868` `#C44E52` `#8172B3` `#937860` `#DA8BC3` `#8C8C8C` `#CCB974` `#64B5CD` |
| FTEPalette (FiveThirtyEight) | `#008fd5` `#fc4f30` `#e5ae38` `#6d904f` `#8b8b8b` `#810f7c` |

You can override any series color by setting the `Color` field in the chart's style struct:

```go
ax.Line(xs, ys, chart.LineStyle{Color: color.Parse("#hex")})
```

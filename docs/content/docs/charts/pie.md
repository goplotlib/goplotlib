---
title: Pie & Donut Chart
weight: 12
description: Show part-to-whole proportions as sectors of a circle, with optional donut cutout.
---

## When to use

Pie charts work best when:

- You have a small number of categories (≤ 6) and the relative proportions are what matters.
- You want to highlight one dominant segment.

Donut charts are the same chart with a circular cutout — they reduce the visual weight of the centre and are a common alternative to the solid pie.

Use `ax.Pie(labels, values, chart.PieStyle{...})`.

---

## Basic example

```go
package main

import (
	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
)

func main() {
	labels := []string{"Direct", "Organic Search", "Referral", "Social", "Email"}
	values := []float64{35, 28, 18, 12, 7}

	fig := plot.New(plot.WithWidth(700), plot.WithHeight(420))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{})
	ax.SetTitle("Traffic Sources")
}
```

![Pie chart](/img/pie_chart.svg)

---

## Donut variant

Set `DonutRadius` to the inner hole radius as a proportion of the outer radius (0 = solid pie, 0.55 = typical donut).

```go
package main

import (
	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
)

func main() {
	labels := []string{"Direct", "Organic Search", "Referral", "Social", "Email"}
	values := []float64{35, 28, 18, 12, 7}

	fig := plot.New(plot.WithWidth(700), plot.WithHeight(420))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.55})
	ax.SetTitle("Traffic Sources (Donut)")
}
```

![Donut chart](/img/donut_chart.svg)

---

## Exploding a segment

Set `ExplodeIdx` and `ExplodeOffset` to offset a segment outward along its midpoint angle, drawing attention to it.

```go
ax.Pie(labels, values, chart.PieStyle{
	DonutRadius:   0.55,
	ExplodeIdx:    0,  // explode the first segment
	ExplodeOffset: 12, // by 12 px
})
```

---

## Notes

- Segment labels are drawn outside the circle for segments that are ≥ 5 % of the total; smaller segments are unlabelled.
- Colors are assigned automatically from the active theme's palette.
- Pie charts occupy the full cell area — no axes, grid lines, or spines are drawn.

---

## Style reference

`chart.PieStyle` fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `DonutRadius` | `float64` | `0` | Inner radius as fraction of outer; `0` = solid pie |
| `ExplodeIdx` | `int` | `0` | Index of segment to offset outward |
| `ExplodeOffset` | `float64` | `0` | Outward offset in pixels for the exploded segment |

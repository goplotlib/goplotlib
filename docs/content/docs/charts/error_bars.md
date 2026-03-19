---
title: Error Bars
weight: 14
description: Overlay ±uncertainty indicators on scatter or line charts to show standard deviation, confidence intervals, or measurement error.
---

## When to use

Error bars are an overlay added on top of an existing chart series. They are useful for:

- Scientific plots showing measurement uncertainty or standard deviation.
- Confidence intervals around model predictions.
- Range indicators on any point-based chart.

Error bars implement the same `Chart` interface as other series, so they automatically expand the axis bounds to fit the full bar extent.

---

## Symmetric error bars

Vaccine antibody titre measured at seven time points post-vaccination. Passing a single error slice draws symmetric ±σ bars.

```go
package main

import (
	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
)

func main() {
	// Days after second dose
	days  := []float64{7, 14, 21, 28, 42, 60, 90}
	// Mean antibody titre (AU/mL)
	titre := []float64{12.3, 48.7, 142.5, 215.8, 189.4, 143.2, 98.6}
	// Standard deviation (symmetric ±)
	sigma := []float64{3.1, 12.4, 28.9, 41.2, 35.7, 29.8, 22.1}

	fig := plot.New(plot.WithWidth(860), plot.WithHeight(420))
	ax := fig.AddAxes()
	ax.Scatter(days, titre, chart.ScatterStyle{Label: "mean titre", MarkerSize: 6})
	ax.ErrorBars(days, titre, sigma)
	ax.SetTitle("Antibody Response Post-Vaccination").
		SetXLabel("Days after dose").
		SetYLabel("Titre (AU/mL) ± σ")
}
```

![Error bars](/img/error_bars.svg)

---

## Asymmetric error bars

Pass a second slice for the upper bound when the uncertainty is not symmetric — for example, a 95 % confidence interval from a log-normal distribution:

```go
// A/B test: conversion rate by variant with asymmetric 95% CI
variants := []float64{1, 2, 3, 4}
rates    := []float64{3.2, 4.1, 3.8, 4.7}
ciLow    := []float64{0.4, 0.5, 0.4, 0.6}
ciHigh   := []float64{0.6, 0.7, 0.5, 0.8}

ax.Scatter(variants, rates, chart.ScatterStyle{Label: "conversion %", MarkerSize: 7})
ax.ErrorBars(variants, rates, ciLow, ciHigh)
```

The bar extends from `y − ciLow[i]` to `y + ciHigh[i]`.

---

## Horizontal error bars

```go
ax.ErrorBarsX(xs, ys, xErr) // symmetric ±xErr
```

---

## Signature reference

```go
// Vertical error bars. Pass one slice for symmetric ±yLow bars,
// or two slices for asymmetric lower/upper bounds.
func (a *Axes) ErrorBars(xs, ys, yLow []float64, yHigh ...[]float64) *Axes

// Horizontal error bars (symmetric).
func (a *Axes) ErrorBarsX(xs, ys, xErr []float64) *Axes
```

Error bars automatically inherit the color of the immediately preceding chart series.

# goplotlib

A zero-dependency Go library for generating beautiful SVG charts.

**[Documentation](https://goplotlib.github.io/goplotlib)** · **[API Reference](https://pkg.go.dev/github.com/goplotlib/goplotlib)**

## Chart types

Line · Area · Step · Bar · HBar · Grouped Bar · Stacked Bar · Stacked Area · Scatter · Bubble · Pie · Donut · Histogram · Heatmap · Box · Error Bars

## Quick start

```bash
go get github.com/goplotlib/goplotlib
```

```go
package main

import (
    "math"

    "github.com/goplotlib/goplotlib/chart"
    "github.com/goplotlib/goplotlib/plot"
)

func main() {
    xs := make([]float64, 100)
    ys := make([]float64, 100)
    for i := range xs {
        xs[i] = float64(i) * 0.1
        ys[i] = math.Sin(xs[i])
    }

    fig := plot.New(plot.WithWidth(800), plot.WithHeight(400))
    ax := fig.AddAxes()
    ax.Line(xs, ys, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true})
    ax.SetTitle("Sine wave").SetXLabel("x").SetYLabel("y")
    fig.SaveSVG("chart.svg")
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for build, test, and docs instructions.

## License

MIT

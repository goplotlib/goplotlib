package main

import (
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/plot"
	"github.com/goplotlib/goplotlib/theme"
)

func main() {
	// Example 1: Line chart with smoothing
	{
		xs := make([]float64, 100)
		sin := make([]float64, 100)
		cos := make([]float64, 100)
		for i := range xs {
			xs[i] = float64(i) * 0.1
			sin[i] = math.Sin(xs[i])
			cos[i] = math.Cos(xs[i])
		}

		fig := plot.New(plot.WithWidth(900), plot.WithHeight(500))
		ax := fig.AddAxes()
		ax.Line(xs, sin, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true})
		ax.Line(xs, cos, chart.LineStyle{Label: "cos(x)", Smooth: true, Dash: []float64{8, 4}})
		ax.SetTitle("Trigonometric Functions").SetXLabel("x").SetYLabel("y")
		fig.SaveSVG("/tmp/line_chart.svg")
	}

	// Example 2: Bar chart
	{
		categories := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		values := []float64{42, 58, 71, 65, 83, 91}

		fig := plot.New(plot.WithWidth(800), plot.WithHeight(500), plot.WithTheme(theme.FiveThirtyEight))
		ax := fig.AddAxes()
		ax.Bar(categories, values, chart.BarStyle{Label: "Revenue"})
		ax.SetTitle("Monthly Revenue").SetYLabel("USD (thousands)")
		fig.SaveSVG("/tmp/bar_chart.svg")
	}

	// Example 3: Scatter plot
	{
		n := 80
		xs := make([]float64, n)
		ys := make([]float64, n)
		xs2 := make([]float64, n)
		ys2 := make([]float64, n)
		for i := range xs {
			t := float64(i) / float64(n)
			xs[i] = t*10 + pseudoRand(i, 0)*0.5
			ys[i] = 2*t + pseudoRand(i, 1)*0.8
			xs2[i] = t*10 + pseudoRand(i, 2)*0.5
			ys2[i] = 5 - 2*t + pseudoRand(i, 3)*0.8
		}

		fig := plot.New(plot.WithWidth(900), plot.WithHeight(550), plot.WithTheme(theme.Dark))
		ax := fig.AddAxes()
		ax.Scatter(xs, ys, chart.ScatterStyle{Label: "Group A", MarkerSize: 7})
		ax.Scatter(xs2, ys2, chart.ScatterStyle{Label: "Group B", MarkerSize: 7})
		ax.SetTitle("Scatter Plot").SetXLabel("X").SetYLabel("Y")
		fig.SaveSVG("/tmp/scatter_chart.svg")
	}
}

// pseudoRand returns a deterministic pseudo-random value in [-1, 1].
func pseudoRand(i, seed int) float64 {
	v := float64((i*1103515245 + seed*12345) & 0x7fffffff)
	return v/float64(0x3fffffff) - 1.0
}

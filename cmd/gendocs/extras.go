// extras.go generates SVGs that are not produced by running doc page snippets:
// showcase images for the homepage gallery, theme previews, and a handful of
// chart variants (donut, grouped bar, negative bar, log scale, line markers).
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/colormap"
	"github.com/goplotlib/goplotlib/plot"
)

// generateExtras is called by main after processing docs/src/.
func generateExtras() {
	// Showcase images for the homepage gallery.
	genShowcaseLine()
	genShowcaseBar()
	genShowcaseScatter()
	genShowcaseStep()
	genShowcaseHBar()
	genShowcaseStackedBar()
	genShowcaseStackedArea()
	genShowcaseHistogram()
	genShowcaseBubble()
	genShowcaseHeatmap()
	genShowcasePie()
	genShowcaseBox()
	genShowcaseErrorBars()

}

// save writes a figure SVG to docs/static/img/<name> and reports progress.
func save(fig *plot.Figure, name string) {
	path := imgDir + "/" + name
	if err := fig.SaveSVG(path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Println("wrote", path)
}

// --- Shared data helpers ---

func sinCosData() (xs, sinY, cosY []float64) {
	n := 80
	xs = make([]float64, n)
	sinY = make([]float64, n)
	cosY = make([]float64, n)
	for i := range xs {
		xs[i] = float64(i) * 4 * math.Pi / float64(n-1)
		sinY[i] = math.Sin(xs[i])
		cosY[i] = math.Cos(xs[i])
	}
	return
}

// --- Chart variant extras ---

// --- Showcase images ---

func genShowcaseLine() {
	months := make([]float64, 24)
	for i := range months {
		months[i] = float64(i + 1)
	}
	prices := []float64{
		3970, 4080, 4110, 4170, 4200, 4450,
		4590, 4510, 4290, 4200, 4380, 4770,
		4850, 5090, 5220, 5010, 5300, 5460,
		5540, 5650, 5620, 5700, 5870, 5960,
	}
	ma := make([]float64, 24)
	for i := range ma {
		if i < 2 {
			ma[i] = prices[i]
		} else {
			ma[i] = (prices[i] + prices[i-1] + prices[i-2]) / 3
		}
	}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Line(months, prices, chart.LineStyle{Label: "S&P 500", Smooth: true, LineWidth: 2})
	ax.Line(months, ma, chart.LineStyle{Label: "3M avg", Smooth: true, Dash: []float64{6, 3}, LineWidth: 1.5})
	ax.SetTitle("S&P 500 — Monthly Close (2023–2024)").SetYLabel("Price (USD)")
	save(fig, "showcase_line.svg")
}

func genShowcaseBar() {
	langs := []string{"JavaScript", "Python", "TypeScript", "Java", "C#", "C++", "Go", "Rust"}
	pct := []float64{63.6, 49.3, 38.5, 30.6, 27.6, 23.5, 13.5, 13.1}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Bar(langs, pct, chart.BarStyle{Label: "% of respondents"})
	ax.SetTitle("Most Used Languages — Stack Overflow 2023").SetYLabel("Usage (%)")
	save(fig, "showcase_bar.svg")
}

func genShowcaseScatter() {
	sqft1 := []float64{850, 920, 1050, 1100, 1200, 1280, 1350, 1400, 1480, 1550, 1600, 1650, 1700, 1780, 1850, 1920, 2000, 2100, 2200, 2350}
	price1 := []float64{210, 228, 251, 265, 288, 305, 318, 332, 351, 368, 378, 392, 405, 421, 438, 455, 472, 495, 518, 545}
	sqft2 := []float64{900, 980, 1080, 1150, 1250, 1320, 1400, 1480, 1560, 1640, 1720, 1800, 1900, 2000, 2100, 2250, 2400, 2550, 2700, 2900}
	price2 := []float64{285, 312, 345, 368, 402, 428, 455, 481, 510, 538, 565, 595, 632, 668, 705, 751, 798, 848, 901, 960}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Scatter(sqft1, price1, chart.ScatterStyle{Label: "Riverside", MarkerSize: 5})
	ax.Scatter(sqft2, price2, chart.ScatterStyle{Label: "Oakwood", MarkerSize: 5, MarkerShape: "diamond"})
	ax.SetTitle("House Size vs Sale Price").SetXLabel("Size (sq ft)").SetYLabel("Price ($k)")
	save(fig, "showcase_scatter.svg")
}

func genShowcaseStep() {
	meetingMonth := []float64{1, 3, 5, 6, 7, 9, 10, 12, 14, 15, 17, 18, 19, 21, 22, 23, 24}
	rate := []float64{0.25, 0.5, 1.0, 1.75, 2.5, 3.25, 4.0, 4.5, 4.75, 5.0, 5.25, 5.5, 5.5, 5.5, 5.25, 5.0, 4.75}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Step(meetingMonth, rate, chart.StepStyle{Label: "Fed Funds Rate", Fill: true, LineWidth: 2})
	ax.SetTitle("Federal Funds Rate — 2022 to 2024").SetXLabel("Month").SetYLabel("Rate (%)")
	save(fig, "showcase_step.svg")
}

func genShowcaseHBar() {
	frameworks := []string{"Node.js", "React", "jQuery", "Express", "Angular", "Vue.js", "Next.js", "Django", "FastAPI", "Spring Boot"}
	usage := []float64{42.7, 40.6, 21.9, 22.8, 17.1, 15.6, 17.9, 12.0, 12.1, 10.5}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.HBar(frameworks, usage, chart.BarStyle{Label: "% of respondents"})
	ax.SetTitle("Top Web Frameworks — Stack Overflow 2023").SetXLabel("Usage (%)")
	save(fig, "showcase_hbar.svg")
}

func genShowcaseStackedBar() {
	years := []string{"2020", "2021", "2022", "2023"}
	coal := []float64{9400, 9820, 10200, 9800}
	gas := []float64{6300, 6500, 6700, 6600}
	nuclear := []float64{2700, 2800, 2700, 2800}
	renewables := []float64{7200, 7900, 8600, 9500}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.StackedBar(years, [][]float64{coal, gas, nuclear, renewables}, []string{"Coal", "Gas", "Nuclear", "Renewables"}, chart.StackedBarStyle{})
	ax.SetTitle("Global Electricity Generation by Source").SetYLabel("TWh")
	save(fig, "showcase_stacked_bar.svg")
}

func genShowcaseStackedArea() {
	n := 12
	months := make([]float64, n)
	for i := range months {
		months[i] = float64(i + 1)
	}
	mobile := []float64{420, 435, 448, 462, 478, 491, 505, 518, 530, 542, 558, 572}
	desktop := []float64{280, 275, 272, 270, 265, 263, 261, 258, 255, 252, 250, 248}
	tablet := []float64{95, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.StackedArea(months, [][]float64{mobile, desktop, tablet}, []string{"Mobile", "Desktop", "Tablet"})
	ax.SetTitle("Monthly Active Users by Device (2024)").SetXLabel("Month").SetYLabel("Users (M)")
	save(fig, "showcase_stacked_area.svg")
}

func genShowcaseHistogram() {
	salaries := []float64{
		72, 75, 78, 80, 82, 83, 84, 85, 86, 87, 88, 89, 90, 90, 91, 92, 93, 94, 95, 95, 95, 96, 97, 98, 99, 100,
		100, 101, 102, 103, 104, 105, 106, 107, 108, 110, 112, 114, 116, 118, 120,
		125, 128, 130, 132, 135, 138, 140, 142, 145,
		150, 152, 155, 157, 160, 162, 165, 168, 170, 172, 175, 178, 180, 182, 185,
		190, 192, 195, 198, 200, 205, 210, 215, 220, 230, 240, 250,
		88, 90, 92, 94, 96, 98, 100, 102, 104, 106,
		155, 160, 165, 170, 175, 180, 185, 190,
	}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Histogram(salaries, chart.HistogramStyle{Label: "engineers", Bins: 18})
	ax.SetTitle("Software Engineer Salary Distribution").SetXLabel("Base Salary ($k)").SetYLabel("Count")
	save(fig, "showcase_histogram.svg")
}

func genShowcaseBubble() {
	gdp := []float64{65, 48, 42, 38, 55, 12, 8, 18, 72, 52, 35, 5}
	life := []float64{79, 83, 81, 77, 82, 68, 63, 72, 82, 81, 75, 64}
	pop := []float64{330, 125, 83, 67, 25, 1400, 215, 145, 8, 38, 210, 120}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Bubble(gdp, life, pop, chart.BubbleStyle{Label: "countries", SizeMin: 8, SizeMax: 45})
	ax.SetTitle("GDP per Capita vs Life Expectancy vs Population").SetXLabel("GDP per Capita ($k)").SetYLabel("Life Expectancy (years)")
	save(fig, "showcase_bubble.svg")
}

func genShowcaseHeatmap() {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	hours := []string{"6am", "8am", "10am", "12pm", "2pm", "4pm", "6pm", "8pm"}
	matrix := [][]float64{
		{1.2, 8.5, 3.1, 4.2, 3.8, 5.1, 9.2, 4.1},
		{1.3, 8.8, 3.2, 4.3, 3.9, 5.2, 9.5, 4.0},
		{1.1, 8.6, 3.0, 4.4, 4.0, 5.3, 9.3, 3.9},
		{1.2, 8.7, 3.1, 4.5, 4.1, 5.4, 9.4, 4.2},
		{1.4, 8.3, 3.3, 4.6, 4.2, 5.8, 10.2, 5.1},
		{2.1, 3.5, 6.8, 8.2, 7.9, 7.1, 6.5, 3.8},
		{1.8, 2.9, 6.2, 7.8, 7.5, 6.8, 5.9, 3.2},
	}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{
		RowLabels: days,
		ColLabels: hours,
		ColorMap:  colormap.Viridis,
	})
	ax.SetTitle("Bike Share Rides — Day × Hour (hundreds)")
	save(fig, "showcase_heatmap.svg")
}

func genShowcasePie() {
	sources := []string{"Coal", "Natural Gas", "Hydro", "Nuclear", "Wind", "Solar", "Other"}
	pct := []float64{36, 23, 15, 10, 7, 5, 4}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Pie(sources, pct, chart.PieStyle{DonutRadius: 0.55})
	ax.SetTitle("Global Electricity Generation Mix (2023)")
	save(fig, "showcase_pie.svg")
}

func genShowcaseBox() {
	auth := []float64{12, 14, 15, 16, 18, 19, 20, 21, 22, 23, 24, 25, 26, 28, 30, 32, 38, 55, 72}
	search := []float64{45, 52, 58, 63, 68, 72, 75, 78, 82, 86, 90, 95, 100, 108, 120, 145, 180, 220, 380}
	checkout := []float64{180, 195, 210, 225, 235, 245, 255, 265, 275, 285, 295, 310, 325, 345, 370, 420, 490, 580, 720}
	feed := []float64{22, 28, 32, 36, 40, 44, 48, 52, 56, 60, 65, 70, 78, 88, 102, 130, 165}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Box([][]float64{auth, search, checkout, feed}, chart.BoxStyle{Labels: []string{"Auth", "Search", "Checkout", "Feed"}})
	ax.SetTitle("API Response Times by Service").SetYLabel("Latency (ms)")
	save(fig, "showcase_box.svg")
}

func genShowcaseErrorBars() {
	variants := []float64{1, 2, 3, 4, 5}
	convRate := []float64{3.2, 4.1, 3.8, 5.2, 4.7}
	ciLow := []float64{0.4, 0.5, 0.4, 0.6, 0.5}
	ciHigh := []float64{0.4, 0.5, 0.4, 0.6, 0.5}
	fig := plot.New(plot.WithWidth(700), plot.WithHeight(380))
	ax := fig.AddAxes()
	ax.Scatter(variants, convRate, chart.ScatterStyle{Label: "variant", MarkerSize: 7})
	ax.ErrorBars(variants, convRate, ciLow, ciHigh)
	ax.HLine(3.2, plot.Label("baseline"), plot.Dash(6, 3))
	ax.SetTitle("A/B Test — Conversion Rate with 95% CI").SetXLabel("Variant").SetYLabel("Conversion (%)")
	save(fig, "showcase_error_bars.svg")
}


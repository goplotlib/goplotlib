package plot

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/colormap"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/theme"
)

var update    = flag.Bool("update", false, "regenerate all golden files")
var updateNew = flag.Bool("update-new", false, "create golden files for new tests only; fail if an existing file changes")

func checkGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	golden := filepath.Join("testdata", "golden", name+".svg")
	got = normalizeSVG(got)

	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(golden, got, 0644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}

	if *updateNew {
		if _, err := os.Stat(golden); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(golden), 0755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}
			if err := os.WriteFile(golden, got, 0644); err != nil {
				t.Fatalf("write golden: %v", err)
			}
			return
		}
		// File exists — fall through to the normal comparison so regressions still fail.
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("golden file missing — run: go test ./plot/ -update-new")
	}
	if !bytes.Equal(got, want) {
		_ = os.WriteFile(golden+".actual", got, 0644)
		t.Errorf("output differs from golden\n  visual review: go run ./cmd/testview > /tmp/review.html && open /tmp/review.html\n  accept changes: go test ./plot/ -update")
	} else {
		// Clean up stale .actual file from a previous failed run.
		_ = os.Remove(golden + ".actual")
	}
}

var floatRe = regexp.MustCompile(`-?\d+\.\d+`)

// normalizeSVG rounds all floating-point numbers to one decimal place,
// preventing test noise from sub-pixel differences that have no visual effect.
func normalizeSVG(svg []byte) []byte {
	return floatRe.ReplaceAllFunc(svg, func(b []byte) []byte {
		v, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return b
		}
		return []byte(fmt.Sprintf("%.1f", math.Round(v*10)/10))
	})
}

// --- fixtures ---

var (
	xs    = makeXs(60)
	sinYs = makeSin(xs)
	cosYs = makeCos(xs)
)

func makeXs(n int) []float64 {
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = float64(i) * (4 * math.Pi / float64(n-1))
	}
	return xs
}

func makeSin(xs []float64) []float64 {
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = math.Sin(x)
	}
	return ys
}

func makeCos(xs []float64) []float64 {
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = math.Cos(x)
	}
	return ys
}

// --- golden tests ---

func TestGoldenLineStraight(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	checkGolden(t, "line_straight", fig.SVG())
}

func TestGoldenLineSmoothedFill(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true, Fill: true})
	checkGolden(t, "line_smooth_fill", fig.SVG())
}

func TestGoldenLineDashed(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Dash: []float64{8, 4}})
	checkGolden(t, "line_dashed", fig.SVG())
}

func TestGoldenLineMarkers(t *testing.T) {
	short := xs[:12]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(short, makeSin(short), chart.LineStyle{Label: "circle", MarkerSize: 5, MarkerShape: "circle"})
	ax.Line(short, makeCos(short), chart.LineStyle{Label: "diamond", MarkerSize: 5, MarkerShape: "diamond"})
	checkGolden(t, "line_markers", fig.SVG())
}

func TestGoldenBarLight(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	checkGolden(t, "bar_light", fig.SVG())
}

func TestGoldenBarFiveThirtyEight(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.FiveThirtyEight))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	checkGolden(t, "bar_538", fig.SVG())
}

func TestGoldenScatterDark(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys1 := []float64{2.1, 3.8, 3.2, 5.1, 4.9, 6.2, 7.1, 6.8}
	xs2 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys2 := []float64{5.0, 4.2, 6.1, 3.8, 7.2, 5.5, 4.1, 8.0}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Scatter(xs1, ys1, chart.ScatterStyle{Label: "Group A", MarkerSize: 7})
	ax.Scatter(xs2, ys2, chart.ScatterStyle{Label: "Group B", MarkerSize: 7, MarkerShape: "diamond"})
	checkGolden(t, "scatter_dark", fig.SVG())
}

func TestGoldenAxisLabelsAndTitle(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.SetTitle("Trigonometric Function")
	ax.SetXLabel("x (radians)")
	ax.SetYLabel("amplitude")
	checkGolden(t, "axis_labels_title", fig.SVG())
}

func TestGoldenLegendTopRight(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendTopRight)
	checkGolden(t, "legend_top_right", fig.SVG())
}

func TestGoldenLegendTopLeft(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendTopLeft)
	checkGolden(t, "legend_top_left", fig.SVG())
}

func TestGoldenLegendBottomRight(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendBottomRight)
	checkGolden(t, "legend_bottom_right", fig.SVG())
}

func TestGoldenLegendBottomLeft(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendBottomLeft)
	checkGolden(t, "legend_bottom_left", fig.SVG())
}

func TestGoldenLegendOutsideRight(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendOutsideRight)
	checkGolden(t, "legend_outside_right", fig.SVG())
}

func TestGoldenLegendOutsideBottom(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	ax.Legend(LegendOutsideBottom)
	checkGolden(t, "legend_outside_bottom", fig.SVG())
}

func TestGoldenMinimalTheme(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	checkGolden(t, "theme_minimal", fig.SVG())
}

func TestGoldenSingleSeriesNoLegend(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	// Single series with no label — legend should not appear
	ax.Line(xs, sinYs, chart.LineStyle{Color: color.Parse("#e63946"), Smooth: true})
	checkGolden(t, "single_series_no_legend", fig.SVG())
}

func TestGoldenFigureTitle(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTitle("Trigonometric Functions"))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	checkGolden(t, "figure_title", fig.SVG())
}

func TestGoldenFigureTitleAndAxesTitle(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTitle("Annual Report"))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.SetTitle("Waveform")
	ax.SetXLabel("x (radians)")
	ax.SetYLabel("amplitude")
	checkGolden(t, "figure_and_axes_title", fig.SVG())
}

func TestGoldenAxisLimits(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	// Zoom into the first cycle and clamp y to ±1
	ax.SetXLim(0, 6.3)
	ax.SetYLim(-1.0, 1.0)
	checkGolden(t, "axis_limits", fig.SVG())
}

func TestGoldenAxisLimitsOneBound(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	// Fix only the upper y bound; lower remains auto (0 for bar)
	ax.SetYLim(Auto, 100)
	checkGolden(t, "axis_limits_one_bound", fig.SVG())
}

func TestGoldenCustomTickFormat(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{1200, 1850, 2400, 2100, 3100, 3800}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	ax.SetYTickFormat(func(v float64) string { return fmt.Sprintf("$%.0f", v) })
	checkGolden(t, "custom_tick_format_currency", fig.SVG())
}

func TestGoldenCustomTickFormatPercent(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.SetYLim(-1.0, 1.0)
	ax.SetYTickFormat(func(v float64) string { return fmt.Sprintf("%.0f%%", v*100) })
	checkGolden(t, "custom_tick_format_percent", fig.SVG())
}

func TestGoldenAnnotationHLine(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.HLine(0.5, Label("threshold"), Dash(6, 3))
	ax.HLine(0.0, Label("zero"))
	checkGolden(t, "annotation_hline", fig.SVG())
}

func TestGoldenAnnotationVLine(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.VLine(math.Pi, Label("π"))
	ax.VLine(3*math.Pi, Label("3π"), Dash(4, 4))
	checkGolden(t, "annotation_vline", fig.SVG())
}

func TestGoldenAnnotationSpans(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	ax.HSpan(-0.5, 0.5, FillColor("#4C72B0"), Opacity(0.1))
	ax.VSpan(math.Pi, 2*math.Pi, FillColor("#DD8452"), Opacity(0.1))
	checkGolden(t, "annotation_spans", fig.SVG())
}

func TestGoldenAnnotationText(t *testing.T) {
	peakX := xs[7] // xs[7] ≈ π/2, the actual sin peak
	peakY := sinYs[7]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Annotate(peakX, peakY, "local max", ArrowDown())
	ax.HLine(0.0, Dash(4, 4))
	checkGolden(t, "annotation_text", fig.SVG())
}

// --- new chart type golden tests ---

func TestGoldenHBar(t *testing.T) {
	cats := []string{"North America", "Europe", "Asia Pacific", "Latin America", "Middle East"}
	vals := []float64{84, 71, 63, 42, 29}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.HBar(cats, vals, chart.BarStyle{Label: "revenue"})
	ax.SetXLabel("Revenue ($M)")
	checkGolden(t, "hbar", fig.SVG())
}

func TestGoldenHBarNegative(t *testing.T) {
	cats := []string{"Q1", "Q2", "Q3", "Q4"}
	vals := []float64{12, -5, 18, -3}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.HBar(cats, vals, chart.BarStyle{Label: "profit/loss"})
	checkGolden(t, "hbar_negative", fig.SVG())
}

func TestGoldenArea(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Area(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.Area(xs, cosYs, chart.LineStyle{Label: "cos(x)"})
	checkGolden(t, "area", fig.SVG())
}

func TestGoldenStep(t *testing.T) {
	shortXs := xs[:20]
	shortSin := sinYs[:20]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "post (default)"})
	checkGolden(t, "step_post", fig.SVG())
}

func TestGoldenStepPre(t *testing.T) {
	shortXs := xs[:20]
	shortSin := sinYs[:20]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "pre", Mode: "pre"})
	checkGolden(t, "step_pre", fig.SVG())
}

func TestGoldenStepFill(t *testing.T) {
	shortXs := xs[:20]
	shortSin := sinYs[:20]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "inventory", Fill: true})
	checkGolden(t, "step_fill", fig.SVG())
}

func TestGoldenGroupedBar(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	q1 := []float64{42, 58, 71, 65, 83, 91}
	q2 := []float64{55, 62, 68, 78, 75, 95}
	q3 := []float64{38, 51, 60, 70, 80, 88}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, q1, chart.BarStyle{Label: "Q1"})
	ax.Bar(cats, q2, chart.BarStyle{Label: "Q2"})
	ax.Bar(cats, q3, chart.BarStyle{Label: "Q3"})
	checkGolden(t, "grouped_bar", fig.SVG())
}

func TestGoldenStackedBar(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	series := [][]float64{
		{20, 28, 35, 30, 40, 45},
		{15, 18, 22, 25, 28, 32},
		{7, 12, 14, 10, 15, 14},
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.StackedBar(cats, series, []string{"Rent", "Salaries", "Marketing"}, chart.StackedBarStyle{})
	checkGolden(t, "stacked_bar", fig.SVG())
}

func TestGoldenSubplot1x2(t *testing.T) {
	fig := New(WithWidth(1000), WithHeight(420))
	axes := fig.SubPlots(1, 2)
	axes[0][0].Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	axes[0][0].SetTitle("Sine")
	axes[0][1].Bar(
		[]string{"Jan", "Feb", "Mar", "Apr", "May"},
		[]float64{42, 58, 71, 65, 83},
		chart.BarStyle{Label: "revenue"},
	)
	axes[0][1].SetTitle("Revenue")
	checkGolden(t, "subplot_1x2", fig.SVG())
}

func TestGoldenSubplot2x2(t *testing.T) {
	fig := New(WithWidth(1000), WithHeight(700), WithTitle("Dashboard"))
	axes := fig.SubPlots(2, 2)

	axes[0][0].Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	axes[0][0].SetTitle("Sine")

	axes[0][1].Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	axes[0][1].SetTitle("Cosine")

	cats := []string{"A", "B", "C", "D", "E"}
	vals := []float64{42, 58, 71, 65, 83}
	axes[1][0].Bar(cats, vals, chart.BarStyle{Label: "metric"})
	axes[1][0].SetTitle("Bar")

	xs1 := []float64{1, 2, 3, 4, 5, 6, 7}
	ys1 := []float64{2.1, 3.8, 3.2, 5.1, 4.9, 6.2, 7.1}
	axes[1][1].Scatter(xs1, ys1, chart.ScatterStyle{Label: "data", MarkerSize: 6})
	axes[1][1].SetTitle("Scatter")

	checkGolden(t, "subplot_2x2", fig.SVG())
}

func TestGoldenSubplotAddAxesAt(t *testing.T) {
	fig := New(WithWidth(1000), WithHeight(420))
	fig.SetLayout(1, 3)

	ax0 := fig.AddAxes(At(0, 0))
	ax0.Line(xs, sinYs, chart.LineStyle{Smooth: true})
	ax0.SetTitle("sin")

	ax1 := fig.AddAxes(At(0, 1))
	ax1.Line(xs, cosYs, chart.LineStyle{Smooth: true})
	ax1.SetTitle("cos")

	ax2 := fig.AddAxes(At(0, 2))
	ax2.Step(xs[:20], sinYs[:20], chart.StepStyle{Label: "step"})
	ax2.SetTitle("step")

	checkGolden(t, "subplot_1x3", fig.SVG())
}

func TestGoldenHeatmapViridis(t *testing.T) {
	// 5x5 correlation-like matrix
	matrix := [][]float64{
		{1.00, 0.85, 0.42, 0.10, -0.20},
		{0.85, 1.00, 0.60, 0.30, 0.05},
		{0.42, 0.60, 1.00, 0.75, 0.50},
		{0.10, 0.30, 0.75, 1.00, 0.80},
		{-0.20, 0.05, 0.50, 0.80, 1.00},
	}
	labels := []string{"A", "B", "C", "D", "E"}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{
		RowLabels:  labels,
		ColLabels:  labels,
		ColorMap:   colormap.Viridis,
		CellLabels: true,
	})
	ax.SetTitle("Correlation Matrix")
	ax.NoLegend()
	checkGolden(t, "heatmap_viridis", fig.SVG())
}

func TestGoldenHeatmapRdBu(t *testing.T) {
	// Diverging data centered around zero
	matrix := [][]float64{
		{2.5, 1.2, -0.5, -1.8},
		{1.0, 0.3, -0.3, -1.0},
		{-0.8, -0.2, 0.4, 1.5},
		{-2.0, -1.2, 0.8, 2.2},
	}
	rows := []string{"W1", "W2", "W3", "W4"}
	cols := []string{"Mon", "Tue", "Wed", "Thu"}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{
		RowLabels:      rows,
		ColLabels:      cols,
		ColorMap:       colormap.RdBu,
		DivergingScale: true,
		CellLabels:     true,
	})
	ax.SetTitle("Diverging Heatmap")
	ax.NoLegend()
	checkGolden(t, "heatmap_rdbu_diverging", fig.SVG())
}

func TestGoldenHeatmapBlues(t *testing.T) {
	// Activity grid — no cell labels
	matrix := [][]float64{
		{0, 3, 7, 12, 8, 5, 1},
		{2, 5, 10, 18, 14, 9, 3},
		{1, 4, 8, 15, 11, 6, 2},
	}
	rows := []string{"Low", "Med", "High"}
	cols := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	fig := New(WithWidth(700), WithHeight(360))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{
		RowLabels: rows,
		ColLabels: cols,
		ColorMap:  colormap.Blues,
	})
	ax.SetTitle("Weekly Activity")
	ax.NoLegend()
	checkGolden(t, "heatmap_blues_activity", fig.SVG())
}

func TestGoldenStackedArea(t *testing.T) {
	series := [][]float64{
		makeSin(xs),
		makeCos(xs),
	}
	// Shift to positive values for stacking
	for i := range series[0] {
		series[0][i] = series[0][i] + 1.2
		series[1][i] = series[1][i] + 1.2
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.StackedArea(xs, series, []string{"Product A", "Product B"})
	checkGolden(t, "stacked_area", fig.SVG())
}

func TestGoldenBoxPlot(t *testing.T) {
	// Deterministic groups with varied spread and outliers
	groupA := []float64{2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 15}        // outlier at 15
	groupB := []float64{1, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12}        // tighter spread
	groupC := []float64{0, 1, 3, 5, 6, 6, 7, 8, 9, 18, 20}         // two high outliers
	fig := New(WithWidth(700), WithHeight(500))
	ax := fig.AddAxes()
	ax.Box([][]float64{groupA, groupB, groupC},
		chart.BoxStyle{Labels: []string{"Control", "Treatment A", "Treatment B"}},
	)
	ax.SetYLabel("Score")
	checkGolden(t, "box_plot", fig.SVG())
}

func TestGoldenBoxPlotMulti(t *testing.T) {
	// Two box series side by side (multiple Box calls share the same categories)
	pre  := []float64{5, 6, 7, 8, 8, 9, 10, 11, 12}
	post := []float64{8, 9, 10, 11, 11, 12, 13, 14, 20}
	fig := New(WithWidth(500), WithHeight(450))
	ax := fig.AddAxes()
	ax.Box([][]float64{pre}, chart.BoxStyle{Labels: []string{"Before"}})
	ax.Box([][]float64{post}, chart.BoxStyle{Labels: []string{"After"}})
	ax.SetTitle("Pre/Post Comparison")
	checkGolden(t, "box_plot_multi", fig.SVG())
}

func TestGoldenPie(t *testing.T) {
	labels := []string{"Chrome", "Safari", "Firefox", "Edge", "Other"}
	values := []float64{65, 19, 4, 4, 8}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{})
	ax.SetTitle("Browser Market Share")
	ax.NoLegend()
	checkGolden(t, "pie", fig.SVG())
}

func TestGoldenDonut(t *testing.T) {
	labels := []string{"Chrome", "Safari", "Firefox", "Edge", "Other"}
	values := []float64{65, 19, 4, 4, 8}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.55})
	ax.SetTitle("Browser Market Share")
	ax.NoLegend()
	checkGolden(t, "donut", fig.SVG())
}

func TestGoldenDonutExplode(t *testing.T) {
	labels := []string{"Q1", "Q2", "Q3", "Q4"}
	values := []float64{28, 35, 20, 17}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.5, ExplodeIdx: 1, ExplodeOffset: 14})
	ax.SetTitle("Quarterly Revenue")
	ax.NoLegend()
	checkGolden(t, "donut_explode", fig.SVG())
}

func TestGoldenBubble(t *testing.T) {
	xs := []float64{10, 20, 30, 40, 55, 65, 75, 85, 50, 35}
	ys := []float64{40, 60, 30, 80, 50, 70, 20, 55, 65, 45}
	sizes := []float64{100, 250, 80, 400, 150, 300, 60, 200, 350, 120}
	fig := New(WithWidth(700), WithHeight(500))
	ax := fig.AddAxes()
	ax.Bubble(xs, ys, sizes, chart.BubbleStyle{Label: "countries"})
	ax.SetXLabel("GDP per capita")
	ax.SetYLabel("Life expectancy")
	checkGolden(t, "bubble", fig.SVG())
}

func TestGoldenBubbleMultiSeries(t *testing.T) {
	xs1 := []float64{10, 25, 40, 60, 75}
	ys1 := []float64{30, 55, 45, 70, 60}
	sz1 := []float64{50, 200, 100, 350, 180}

	xs2 := []float64{15, 30, 50, 65, 80}
	ys2 := []float64{50, 40, 65, 35, 75}
	sz2 := []float64{80, 150, 250, 120, 300}

	fig := New(WithWidth(700), WithHeight(500))
	ax := fig.AddAxes()
	ax.Bubble(xs1, ys1, sz1, chart.BubbleStyle{Label: "Group A"})
	ax.Bubble(xs2, ys2, sz2, chart.BubbleStyle{Label: "Group B"})
	checkGolden(t, "bubble_multi_series", fig.SVG())
}

func TestGoldenHistogramAuto(t *testing.T) {
	// Normally-distributed-ish values via deterministic pseudo-random
	vals := make([]float64, 200)
	for i := range vals {
		// Box-Muller using fixed seed-like pattern
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Label: "response time"})
	ax.SetXLabel("ms")
	ax.SetYLabel("count")
	checkGolden(t, "histogram_auto", fig.SVG())
}

func TestGoldenHistogramBinEdges(t *testing.T) {
	vals := make([]float64, 100)
	for i := range vals {
		vals[i] = float64(i%50) * 2.0
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{
		BinEdges: []float64{0, 10, 25, 50, 75, 100},
		Label:    "custom bins",
	})
	checkGolden(t, "histogram_bin_edges", fig.SVG())
}

func TestGoldenHistogramNormalized(t *testing.T) {
	vals := make([]float64, 200)
	for i := range vals {
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Bins: 20, Normalize: true, Label: "density"})
	ax.SetYLabel("density")
	checkGolden(t, "histogram_normalized", fig.SVG())
}

func TestGoldenHistogramCumulative(t *testing.T) {
	vals := make([]float64, 200)
	for i := range vals {
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Bins: 20, Cumulative: true, Label: "cumulative"})
	checkGolden(t, "histogram_cumulative", fig.SVG())
}

// --- Line ---

func TestGoldenLineColorOverride(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin", Color: color.Parse("#e74c3c"), Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos", Color: color.Parse("#2980b9"), Smooth: true})
	checkGolden(t, "line_color_override", fig.SVG())
}

func TestGoldenLineMarkersAllShapes(t *testing.T) {
	short := xs[:10]
	fig := New(WithWidth(800), WithHeight(500))
	ax := fig.AddAxes()
	offset := func(ys []float64, d float64) []float64 {
		out := make([]float64, len(ys))
		for i, v := range ys {
			out[i] = v + d
		}
		return out
	}
	ax.Line(short, makeSin(short), chart.LineStyle{Label: "circle", MarkerSize: 6, MarkerShape: "circle"})
	ax.Line(short, offset(makeSin(short), 0.5), chart.LineStyle{Label: "square", MarkerSize: 6, MarkerShape: "square"})
	ax.Line(short, offset(makeSin(short), 1.0), chart.LineStyle{Label: "diamond", MarkerSize: 6, MarkerShape: "diamond"})
	ax.Line(short, offset(makeSin(short), 1.5), chart.LineStyle{Label: "triangle", MarkerSize: 6, MarkerShape: "triangle"})
	checkGolden(t, "line_marker_shapes", fig.SVG())
}

func TestGoldenLineFillOpacity(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin", Fill: true, FillOpacity: 0.4})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos", Fill: true, FillOpacity: 0.1})
	checkGolden(t, "line_fill_opacity", fig.SVG())
}

func TestGoldenLineOpacity(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin", LineWidth: 4, Opacity: 0.4})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos", LineWidth: 4, Opacity: 1.0})
	checkGolden(t, "line_opacity", fig.SVG())
}

func TestGoldenLineDarkTheme(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	checkGolden(t, "line_dark_theme", fig.SVG())
}

func TestGoldenLineLineWidth(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "thin (1px)", LineWidth: 1})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "thick (5px)", LineWidth: 5})
	checkGolden(t, "line_linewidth", fig.SVG())
}

// --- Bar ---

func TestGoldenBarNegative(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, -18, 71, -30, 83, -12}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "profit/loss"})
	checkGolden(t, "bar_negative", fig.SVG())
}

func TestGoldenBarNoRounded(t *testing.T) {
	cats := []string{"A", "B", "C", "D", "E"}
	vals := []float64{30, 55, 45, 70, 60}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "sales", SquareBars: true})
	checkGolden(t, "bar_no_rounded", fig.SVG())
}

func TestGoldenBarColorOverride(t *testing.T) {
	cats := []string{"A", "B", "C", "D"}
	vals := []float64{40, 65, 50, 80}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue", Color: color.Parse("#e67e22"), Opacity: 0.85})
	checkGolden(t, "bar_color_override", fig.SVG())
}

func TestGoldenBarDarkTheme(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	checkGolden(t, "bar_dark_theme", fig.SVG())
}

// --- Scatter ---

func TestGoldenScatterLight(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys1 := []float64{2.1, 3.8, 3.2, 5.1, 4.9, 6.2, 7.1, 6.8}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(xs1, ys1, chart.ScatterStyle{Label: "Group A", MarkerSize: 8})
	checkGolden(t, "scatter_light", fig.SVG())
}

func TestGoldenScatterAllShapes(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7}
	fig := New(WithWidth(800), WithHeight(500))
	ax := fig.AddAxes()
	ax.Scatter(xs1, []float64{1, 1.2, 0.9, 1.1, 1.0, 0.8, 1.2}, chart.ScatterStyle{Label: "circle", MarkerSize: 8, MarkerShape: "circle"})
	ax.Scatter(xs1, []float64{2, 2.2, 1.9, 2.1, 2.0, 1.8, 2.2}, chart.ScatterStyle{Label: "square", MarkerSize: 8, MarkerShape: "square"})
	ax.Scatter(xs1, []float64{3, 3.2, 2.9, 3.1, 3.0, 2.8, 3.2}, chart.ScatterStyle{Label: "diamond", MarkerSize: 8, MarkerShape: "diamond"})
	ax.Scatter(xs1, []float64{4, 4.2, 3.9, 4.1, 4.0, 3.8, 4.2}, chart.ScatterStyle{Label: "triangle", MarkerSize: 8, MarkerShape: "triangle"})
	checkGolden(t, "scatter_all_shapes", fig.SVG())
}

func TestGoldenScatterOpacity(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ys1 := make([]float64, 10)
	ys2 := make([]float64, 10)
	for i := range xs1 {
		ys1[i] = xs1[i]*0.5 + 1
		ys2[i] = xs1[i]*0.3 + 2
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(xs1, ys1, chart.ScatterStyle{Label: "A", MarkerSize: 14, Opacity: 0.4})
	ax.Scatter(xs1, ys2, chart.ScatterStyle{Label: "B", MarkerSize: 14, Opacity: 0.4})
	checkGolden(t, "scatter_opacity", fig.SVG())
}

// --- Step ---

func TestGoldenStepMid(t *testing.T) {
	shortXs := xs[:20]
	shortSin := sinYs[:20]
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "mid", Mode: "mid"})
	checkGolden(t, "step_mid", fig.SVG())
}

func TestGoldenStepAllModes(t *testing.T) {
	shortXs := xs[:15]
	shortSin := sinYs[:15]
	offset := func(d float64) []float64 {
		out := make([]float64, len(shortSin))
		for i, v := range shortSin {
			out[i] = v + d
		}
		return out
	}
	fig := New(WithWidth(800), WithHeight(500))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "post", Mode: "post"})
	ax.Step(shortXs, offset(0.4), chart.StepStyle{Label: "pre", Mode: "pre"})
	ax.Step(shortXs, offset(0.8), chart.StepStyle{Label: "mid", Mode: "mid"})
	checkGolden(t, "step_all_modes", fig.SVG())
}

// --- HBar ---

func TestGoldenHBarColor(t *testing.T) {
	cats := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}
	vals := []float64{55, 82, 38, 71, 60}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.HBar(cats, vals, chart.BarStyle{Color: color.Parse("#27ae60"), Label: "score"})
	checkGolden(t, "hbar_color", fig.SVG())
}

func TestGoldenHBarDarkTheme(t *testing.T) {
	cats := []string{"Alpha", "Beta", "Gamma", "Delta"}
	vals := []float64{55, 82, 38, 71}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.HBar(cats, vals, chart.BarStyle{Label: "score"})
	checkGolden(t, "hbar_dark_theme", fig.SVG())
}

// --- Area ---

func TestGoldenAreaFillOpacity(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Area(xs, sinYs, chart.LineStyle{Label: "sin", FillOpacity: 0.5})
	ax.Area(xs, cosYs, chart.LineStyle{Label: "cos", FillOpacity: 0.5})
	checkGolden(t, "area_fill_opacity", fig.SVG())
}

// --- Heatmap ---

func TestGoldenHeatmapPlasma(t *testing.T) {
	matrix := [][]float64{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		{11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20},
	}
	fig := New(WithWidth(600), WithHeight(450))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{ColorMap: colormap.Plasma, CellLabels: true})
	ax.SetTitle("Plasma Colormap")
	ax.NoLegend()
	checkGolden(t, "heatmap_plasma", fig.SVG())
}

func TestGoldenHeatmapGreys(t *testing.T) {
	matrix := [][]float64{
		{10, 20, 30},
		{40, 50, 60},
		{70, 80, 90},
	}
	// No row/col labels — should auto-generate index labels
	fig := New(WithWidth(500), WithHeight(420))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{ColorMap: colormap.Greys})
	ax.SetTitle("Greys — auto index labels")
	ax.NoLegend()
	checkGolden(t, "heatmap_greys_auto_labels", fig.SVG())
}

// --- Histogram ---

func TestGoldenHistogramExplicitBins(t *testing.T) {
	vals := make([]float64, 200)
	for i := range vals {
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Bins: 8, Label: "coarse binning"})
	ax.SetYLabel("count")
	checkGolden(t, "histogram_explicit_bins", fig.SVG())
}

func TestGoldenHistogramColor(t *testing.T) {
	vals := make([]float64, 150)
	for i := range vals {
		u1 := float64(i*7+1) / 150.0
		u2 := float64(i*13+1) / 150.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*10 + 40
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Color: color.Parse("#8e44ad"), Label: "latency"})
	checkGolden(t, "histogram_color", fig.SVG())
}

func TestGoldenHistogramDarkTheme(t *testing.T) {
	vals := make([]float64, 200)
	for i := range vals {
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Label: "response time"})
	checkGolden(t, "histogram_dark_theme", fig.SVG())
}

// --- Bubble ---

func TestGoldenBubbleSizeRange(t *testing.T) {
	xs1 := []float64{10, 20, 30, 40, 55, 65, 75, 85, 50, 35}
	ys1 := []float64{40, 60, 30, 80, 50, 70, 20, 55, 65, 45}
	sizes := []float64{100, 250, 80, 400, 150, 300, 60, 200, 350, 120}
	fig := New(WithWidth(700), WithHeight(500))
	ax := fig.AddAxes()
	ax.Bubble(xs1, ys1, sizes, chart.BubbleStyle{Label: "data", SizeMin: 3, SizeMax: 20})
	checkGolden(t, "bubble_size_range", fig.SVG())
}

func TestGoldenBubbleLabels(t *testing.T) {
	gdp := []float64{65, 48, 42, 38, 55, 12, 8, 18, 72, 52, 35, 5}
	life := []float64{79, 83, 81, 77, 82, 68, 63, 72, 82, 81, 75, 64}
	pop := []float64{330, 125, 83, 67, 25, 1400, 215, 145, 8, 38, 210, 120}
	labels := []string{"USA", "Japan", "Germany", "France", "UK", "China", "Brazil", "Russia", "Switz.", "Spain", "Mexico", "Ethiopia"}
	fig := New(WithWidth(800), WithHeight(500))
	ax := fig.AddAxes()
	ax.Bubble(gdp, life, pop, chart.BubbleStyle{SizeMin: 8, SizeMax: 50, Labels: labels})
	ax.SetTitle("GDP per Capita vs Life Expectancy vs Population")
	ax.SetXLabel("GDP per Capita ($k)").SetYLabel("Life Expectancy (years)")
	checkGolden(t, "bubble_labels", fig.SVG())
}

func TestGoldenBubbleOpacity(t *testing.T) {
	xs1 := []float64{10, 25, 40, 55, 70}
	ys1 := []float64{30, 50, 40, 60, 45}
	sizes := []float64{200, 400, 150, 350, 250}
	fig := New(WithWidth(700), WithHeight(500))
	ax := fig.AddAxes()
	ax.Bubble(xs1, ys1, sizes, chart.BubbleStyle{Label: "overlapping", Opacity: 0.4})
	ax.Bubble(xs1, []float64{50, 30, 60, 40, 55}, sizes, chart.BubbleStyle{Label: "series B", Opacity: 0.4})
	checkGolden(t, "bubble_opacity", fig.SVG())
}

// --- Pie / Donut ---

func TestGoldenPieExplode(t *testing.T) {
	labels := []string{"Chrome", "Safari", "Firefox", "Edge", "Other"}
	values := []float64{65, 19, 4, 4, 8}
	fig := New(WithWidth(600), WithHeight(500))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{ExplodeIdx: 0, ExplodeOffset: 16})
	ax.SetTitle("Exploded Pie")
	ax.NoLegend()
	checkGolden(t, "pie_explode", fig.SVG())
}

func TestGoldenPieDarkTheme(t *testing.T) {
	labels := []string{"A", "B", "C", "D"}
	values := []float64{40, 30, 20, 10}
	fig := New(WithWidth(600), WithHeight(500), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.5})
	ax.SetTitle("Dark Donut")
	ax.NoLegend()
	checkGolden(t, "pie_dark_theme", fig.SVG())
}

// --- Box ---

func TestGoldenBoxSingle(t *testing.T) {
	data := []float64{3, 5, 6, 7, 8, 8, 9, 10, 11, 12, 20}
	fig := New(WithWidth(400), WithHeight(450))
	ax := fig.AddAxes()
	ax.Box([][]float64{data}, chart.BoxStyle{Labels: []string{"Group"}})
	ax.SetYLabel("measurement")
	checkGolden(t, "box_single", fig.SVG())
}

func TestGoldenBoxDarkTheme(t *testing.T) {
	groupA := []float64{2, 4, 5, 6, 7, 8, 9, 10, 15}
	groupB := []float64{5, 6, 7, 8, 8, 9, 10, 11, 12}
	groupC := []float64{1, 3, 5, 6, 7, 8, 9, 18, 20}
	fig := New(WithWidth(700), WithHeight(500), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Box([][]float64{groupA, groupB, groupC},
		chart.BoxStyle{Labels: []string{"Control", "Treatment A", "Treatment B"}},
	)
	checkGolden(t, "box_dark_theme", fig.SVG())
}

// --- Annotations ---

func TestGoldenAnnotationArrows(t *testing.T) {
	peakX := xs[7]   // xs[7] ≈ π/2, sin peak ≈ 1
	troughX := xs[22] // xs[22] ≈ 3π/2, sin trough ≈ -1
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Annotate(peakX, sinYs[7], "peak", ArrowDown())
	ax.Annotate(troughX, sinYs[22], "trough", ArrowUp())
	ax.Annotate(xs[15], sinYs[15], "zero crossing", ArrowLeft())
	ax.Annotate(xs[3], sinYs[3], "rising", ArrowRight())
	checkGolden(t, "annotation_all_arrows", fig.SVG())
}

func TestGoldenAnnotationColor(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.HLine(0.5, Label("danger"), AnnotColor(color.Parse("#e74c3c")), Dash(6, 3))
	ax.HLine(-0.5, Label("safe"), AnnotColor(color.Parse("#27ae60")))
	ax.VLine(math.Pi, Label("π"), AnnotColor(color.Parse("#8e44ad")), Dash(4, 4))
	checkGolden(t, "annotation_color", fig.SVG())
}

// --- Axis features ---

func TestGoldenAxisTickCount(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)"})
	ax.SetXTicks(12)
	ax.SetYTicks(10)
	checkGolden(t, "axis_tick_count", fig.SVG())
}

func TestGoldenAxisXLimAutoMin(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.SetXLim(Auto, 2*math.Pi) // fix upper bound, auto lower
	checkGolden(t, "axis_xlim_auto_min", fig.SVG())
}

func TestGoldenAxisTickFormatDate(t *testing.T) {
	// Simulate monthly timestamps
	start := float64(1704067200) // 2024-01-01 UTC
	tXs := make([]float64, 12)
	tYs := []float64{100, 115, 108, 130, 125, 140, 135, 150, 145, 160, 155, 170}
	for i := range tXs {
		tXs[i] = start + float64(i)*30*24*3600
	}
	fig := New(WithWidth(900), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(tXs, tYs, chart.LineStyle{Label: "revenue", Smooth: true, MarkerSize: 5})
	ax.SetXTickFormat(func(v float64) string {
		months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
			"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		idx := int((v-start)/(30*24*3600)) % 12
		if idx < 0 {
			idx = 0
		}
		return months[idx]
	})
	ax.SetXTicks(12)
	ax.SetYLabel("Revenue ($k)")
	checkGolden(t, "axis_tick_format_date", fig.SVG())
}

// --- Themes ---

func TestGoldenThemeDark(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(900), WithHeight(500), WithTheme(theme.Dark), WithTitle("Dark Theme Dashboard"))
	axes := fig.SubPlots(1, 2)
	axes[0][0].Line(xs, sinYs, chart.LineStyle{Label: "sin", Smooth: true})
	axes[0][0].Line(xs, cosYs, chart.LineStyle{Label: "cos", Smooth: true})
	axes[0][0].SetTitle("Lines")
	axes[0][1].Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	axes[0][1].SetTitle("Bars")
	checkGolden(t, "theme_dark_dashboard", fig.SVG())
}

func TestGoldenThemeMinimalLine(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	checkGolden(t, "theme_minimal_line", fig.SVG())
}

// --- Stacked Bar ---

func TestGoldenStackedBarDarkTheme(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	series := [][]float64{
		{20, 28, 35, 30, 40, 45},
		{15, 18, 22, 25, 28, 32},
		{7, 12, 14, 10, 15, 14},
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.StackedBar(cats, series, []string{"Product A", "Product B", "Product C"}, chart.StackedBarStyle{})
	checkGolden(t, "stacked_bar_dark", fig.SVG())
}

// --- Log Scale ---

func TestGoldenLogScaleY(t *testing.T) {
	// Latency data spanning 4 decades.
	lxs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lys := []float64{0.5, 2, 8, 30, 90, 300, 800, 2500, 7000, 20000}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(lxs, lys, chart.LineStyle{Label: "latency (µs)", MarkerSize: 5})
	ax.SetYScale(scale.Log(10))
	ax.SetYLabel("Latency (µs)")
	ax.SetXLabel("Request batch size")
	checkGolden(t, "log_scale_y", fig.SVG())
}

func TestGoldenLogScaleXY(t *testing.T) {
	// Power-law: y = x^2, both axes log.
	n := 10
	lxs := make([]float64, n)
	lys := make([]float64, n)
	for i := range lxs {
		lxs[i] = math.Pow(2, float64(i))   // 1, 2, 4, 8, ..., 512
		lys[i] = lxs[i] * lxs[i]           // 1, 4, 16, 64, ..., 262144
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(lxs, lys, chart.ScatterStyle{Label: "y = x²", MarkerSize: 6})
	ax.SetXScale(scale.Log(2))
	ax.SetYScale(scale.Log(2))
	ax.SetXLabel("x (log₂)")
	ax.SetYLabel("y (log₂)")
	checkGolden(t, "log_scale_xy", fig.SVG())
}

// --- Error Bars ---

func TestGoldenErrorBarsSymmetric(t *testing.T) {
	ebXs := []float64{1, 2, 3, 4, 5, 6, 7}
	ebYs := []float64{2.1, 3.8, 3.2, 5.1, 4.6, 6.3, 5.8}
	yErr := []float64{0.4, 0.6, 0.5, 0.7, 0.5, 0.8, 0.6}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(ebXs, ebYs, chart.ScatterStyle{Label: "measurements", MarkerSize: 6})
	ax.ErrorBars(ebXs, ebYs, yErr)
	ax.SetXLabel("Sample")
	ax.SetYLabel("Value")
	checkGolden(t, "error_bars_symmetric", fig.SVG())
}

func TestGoldenErrorBarsAsymmetric(t *testing.T) {
	ebXs := []float64{1, 2, 3, 4, 5}
	ebYs := []float64{10, 25, 18, 32, 22}
	yLow := []float64{2, 5, 3, 7, 4}
	yHigh := []float64{4, 8, 5, 10, 6}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(ebXs, ebYs, chart.LineStyle{Label: "model", MarkerSize: 5})
	ax.ErrorBars(ebXs, ebYs, yLow, yHigh)
	checkGolden(t, "error_bars_asymmetric", fig.SVG())
}

func TestGoldenErrorBarsX(t *testing.T) {
	ebXs := []float64{5, 12, 8, 18, 14}
	ebYs := []float64{1, 2, 3, 4, 5}
	xErr := []float64{1.5, 2.0, 1.0, 2.5, 1.8}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(ebXs, ebYs, chart.ScatterStyle{Label: "data", MarkerSize: 6})
	ax.ErrorBarsX(ebXs, ebYs, xErr)
	ax.SetXLabel("Measurement")
	ax.SetYLabel("Trial")
	checkGolden(t, "error_bars_x", fig.SVG())
}

func TestGoldenErrorBarsWithLine(t *testing.T) {
	// Error bars on a line chart with explicit color.
	ebXs := []float64{0, 1, 2, 3, 4, 5}
	ebYs := []float64{1.0, 2.2, 4.1, 7.9, 14.8, 28.5}
	yErr := []float64{0.2, 0.4, 0.6, 1.0, 1.5, 2.5}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(ebXs, ebYs, chart.LineStyle{Label: "signal", Smooth: false, MarkerSize: 5})
	ax.ErrorBars(ebXs, ebYs, yErr)
	checkGolden(t, "error_bars_with_line", fig.SVG())
}

func TestGoldenY2Axis(t *testing.T) {
	quarters := []string{"Q1 '22", "Q2 '22", "Q3 '22", "Q4 '22", "Q1 '23", "Q2 '23", "Q3 '23", "Q4 '23"}
	revenue := []float64{210, 235, 228, 260, 275, 290, 282, 315}
	growth := make([]float64, len(revenue)-1)
	for i := range growth {
		growth[i] = (revenue[i+1] - revenue[i]) / revenue[i] * 100
	}
	// growth[0] corresponds to Q2 (0-based index 1), so x positions start at 1
	gxs := make([]float64, len(growth))
	for i := range gxs {
		gxs[i] = float64(i + 1)
	}
	qxs := make([]float64, len(quarters))
	for i := range qxs {
		qxs[i] = float64(i + 1)
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(quarters, revenue, chart.BarStyle{Label: "Revenue ($M)"})
	ax.Line(gxs, growth, chart.LineStyle{Label: "QoQ Growth (%)", Smooth: true, LineWidth: 2, MarkerSize: 5}).OnY2()
	ax.SetYLabel("Revenue ($M)").SetY2Label("QoQ Growth (%)")
	ax.SetTitle("Quarterly Revenue and Growth")
	checkGolden(t, "y2_axis", fig.SVG())
}

func TestGoldenY2AxisTwoLines(t *testing.T) {
	n := 50
	xs := make([]float64, n)
	temp := make([]float64, n)
	pressure := make([]float64, n)
	for i := range xs {
		xs[i] = float64(i)
		temp[i] = 20 + math.Sin(float64(i)*0.2)*8
		pressure[i] = 1013 + math.Cos(float64(i)*0.15)*12
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, temp, chart.LineStyle{Label: "Temperature (°C)", Smooth: true, LineWidth: 2})
	ax.Line(xs, pressure, chart.LineStyle{Label: "Pressure (hPa)", Smooth: true, LineWidth: 2}).OnY2()
	ax.SetYLabel("Temperature (°C)").SetY2Label("Pressure (hPa)")
	ax.SetTitle("Temperature vs Pressure")
	checkGolden(t, "y2_axis_two_lines", fig.SVG())
}

// --- Theme: FiveThirtyEight variants ---

func TestGoldenTheme538Line(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.FiveThirtyEight))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true, LineWidth: 2.5})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true, LineWidth: 2.5})
	ax.SetTitle("FiveThirtyEight — Line")
	checkGolden(t, "theme_538_line", fig.SVG())
}

func TestGoldenTheme538Scatter(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys1 := []float64{2.1, 3.8, 3.2, 5.1, 4.9, 6.2, 7.1, 6.8}
	xs2 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys2 := []float64{5.0, 4.2, 6.1, 3.8, 7.2, 5.5, 4.1, 8.0}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.FiveThirtyEight))
	ax := fig.AddAxes()
	ax.Scatter(xs1, ys1, chart.ScatterStyle{Label: "Group A", MarkerSize: 7})
	ax.Scatter(xs2, ys2, chart.ScatterStyle{Label: "Group B", MarkerSize: 7, MarkerShape: "diamond"})
	ax.SetTitle("FiveThirtyEight — Scatter")
	checkGolden(t, "theme_538_scatter", fig.SVG())
}

func TestGoldenTheme538StackedBar(t *testing.T) {
	cats := []string{"2020", "2021", "2022", "2023"}
	series := [][]float64{
		{30, 38, 42, 48},
		{20, 24, 28, 32},
		{10, 13, 16, 19},
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.FiveThirtyEight))
	ax := fig.AddAxes()
	ax.StackedBar(cats, series, []string{"Product A", "Product B", "Product C"}, chart.StackedBarStyle{})
	ax.SetTitle("FiveThirtyEight — Stacked Bar")
	checkGolden(t, "theme_538_stacked_bar", fig.SVG())
}

// --- Theme: Minimal variants ---

func TestGoldenThemeMinimalBar(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	vals := []float64{42, 58, 71, 65, 83, 91}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	ax.SetTitle("Minimal — Bar")
	checkGolden(t, "theme_minimal_bar", fig.SVG())
}

func TestGoldenThemeMinimalScatter(t *testing.T) {
	xs1 := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	ys1 := []float64{2.1, 3.8, 3.2, 5.1, 4.9, 6.2, 7.1, 6.8}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Scatter(xs1, ys1, chart.ScatterStyle{Label: "measurements", MarkerSize: 7})
	ax.SetTitle("Minimal — Scatter")
	checkGolden(t, "theme_minimal_scatter", fig.SVG())
}

// --- Custom theme ---

func TestGoldenCustomTheme(t *testing.T) {
	myTheme := theme.Theme{
		Background:     color.Parse("#fdf6e3"),
		PlotBackground: color.Parse("#eee8d5"),
		GridColor:      color.Parse("#cdc9b8"),
		GridWidth:      0.8,
		AxisColor:      color.Parse("#93a1a1"),
		AxisWidth:      1.2,
		TextColor:      color.Parse("#657b83"),
		TitleColor:     color.Parse("#073642"),
		FontFamily:     "system-ui, sans-serif",
		FontSize:       12,
		TitleFontSize:  16,
		LabelFontSize:  12,
		TickFontSize:   11,
		LineWidth:      2.5,
		SpineLeft:      true,
		SpineBottom:    true,
		Padding:        theme.Padding{Top: 45, Right: 30, Bottom: 60, Left: 70},
		Palette:        []color.Color{color.Parse("#268bd2"), color.Parse("#cb4b16"), color.Parse("#2aa198"), color.Parse("#859900")},
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(myTheme))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true, Fill: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	ax.SetTitle("Custom Solarized Theme")
	checkGolden(t, "custom_theme", fig.SVG())
}

// --- Legend variants ---

func TestGoldenLegendNoneExplicit(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	ax.NoLegend()
	checkGolden(t, "legend_none_explicit", fig.SVG())
}

func TestGoldenLegendManyEntries(t *testing.T) {
	short := xs[:20]
	offsets := []float64{0, 0.4, 0.8, 1.2, 1.6, 2.0}
	labels := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}
	fig := New(WithWidth(900), WithHeight(500))
	ax := fig.AddAxes()
	for i, off := range offsets {
		ys := make([]float64, len(short))
		for j, v := range makeSin(short) {
			ys[j] = v + off
		}
		ax.Line(short, ys, chart.LineStyle{Label: labels[i]})
	}
	checkGolden(t, "legend_many_entries", fig.SVG())
}

func TestGoldenLegendOutsideRightManyEntries(t *testing.T) {
	short := xs[:20]
	offsets := []float64{0, 0.5, 1.0, 1.5, 2.0}
	labels := []string{"Series Alpha", "Series Beta", "Series Gamma", "Series Delta", "Series Epsilon"}
	fig := New(WithWidth(900), WithHeight(450))
	ax := fig.AddAxes()
	for i, off := range offsets {
		ys := make([]float64, len(short))
		for j, v := range makeSin(short) {
			ys[j] = v + off
		}
		ax.Line(short, ys, chart.LineStyle{Label: labels[i]})
	}
	ax.Legend(LegendOutsideRight)
	checkGolden(t, "legend_outside_right_many_entries", fig.SVG())
}

// --- Y2 axis variants ---

func TestGoldenY2AxisDarkTheme(t *testing.T) {
	n := 50
	y2xs := make([]float64, n)
	temp := make([]float64, n)
	pressure := make([]float64, n)
	for i := range y2xs {
		y2xs[i] = float64(i)
		temp[i] = 20 + math.Sin(float64(i)*0.2)*8
		pressure[i] = 1013 + math.Cos(float64(i)*0.15)*12
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Line(y2xs, temp, chart.LineStyle{Label: "Temperature (°C)", Smooth: true, LineWidth: 2})
	ax.Line(y2xs, pressure, chart.LineStyle{Label: "Pressure (hPa)", Smooth: true, LineWidth: 2}).OnY2()
	ax.SetYLabel("Temperature (°C)").SetY2Label("Pressure (hPa)")
	ax.SetTitle("Y2 Axis — Dark Theme")
	checkGolden(t, "y2_axis_dark_theme", fig.SVG())
}

func TestGoldenY2AxisScatter(t *testing.T) {
	n := 12
	y2xs := make([]float64, n)
	sales := make([]float64, n)
	satisfaction := make([]float64, n)
	for i := range y2xs {
		y2xs[i] = float64(i + 1)
		sales[i] = 80 + float64(i)*12
		satisfaction[i] = 3.5 + math.Sin(float64(i)*0.5)*0.8
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(y2xs, sales, chart.LineStyle{Label: "Sales ($k)", Smooth: true, LineWidth: 2})
	ax.Scatter(y2xs, satisfaction, chart.ScatterStyle{Label: "Satisfaction (1–5)", MarkerSize: 7}).OnY2()
	ax.SetYLabel("Sales ($k)").SetY2Label("Satisfaction")
	ax.SetTitle("Sales vs Customer Satisfaction")
	checkGolden(t, "y2_axis_scatter", fig.SVG())
}

func TestGoldenY2AxisLimTicks(t *testing.T) {
	n := 50
	y2xs := make([]float64, n)
	temp := make([]float64, n)
	pressure := make([]float64, n)
	for i := range y2xs {
		y2xs[i] = float64(i)
		temp[i] = 20 + math.Sin(float64(i)*0.2)*8
		pressure[i] = 1013 + math.Cos(float64(i)*0.15)*12
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Line(y2xs, temp, chart.LineStyle{Label: "Temperature (°C)", Smooth: true})
	ax.Line(y2xs, pressure, chart.LineStyle{Label: "Pressure (hPa)", Smooth: true}).OnY2()
	ax.SetY2Lim(1000, 1030)
	ax.SetY2Ticks(4)
	ax.SetYLabel("Temperature (°C)").SetY2Label("Pressure (hPa)")
	checkGolden(t, "y2_axis_lim_ticks", fig.SVG())
}

// --- Missing dark/minimal theme tests for chart types ---

func TestGoldenAreaDarkTheme(t *testing.T) {
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Area(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.Area(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	ax.SetTitle("Area — Dark Theme")
	checkGolden(t, "area_dark_theme", fig.SVG())
}

func TestGoldenStepDarkTheme(t *testing.T) {
	shortXs := xs[:20]
	shortSin := sinYs[:20]
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Step(shortXs, shortSin, chart.StepStyle{Label: "signal", Fill: true})
	ax.SetTitle("Step — Dark Theme")
	checkGolden(t, "step_dark_theme", fig.SVG())
}

func TestGoldenBubbleDarkTheme(t *testing.T) {
	bxs := []float64{10, 20, 30, 40, 55, 65, 75, 85, 50, 35}
	bys := []float64{40, 60, 30, 80, 50, 70, 20, 55, 65, 45}
	sizes := []float64{100, 250, 80, 400, 150, 300, 60, 200, 350, 120}
	fig := New(WithWidth(700), WithHeight(500), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Bubble(bxs, bys, sizes, chart.BubbleStyle{Label: "countries", Opacity: 0.75})
	ax.SetTitle("Bubble — Dark Theme")
	checkGolden(t, "bubble_dark_theme", fig.SVG())
}

func TestGoldenHeatmapDarkTheme(t *testing.T) {
	matrix := [][]float64{
		{1.00, 0.85, 0.42, 0.10, -0.20},
		{0.85, 1.00, 0.60, 0.30, 0.05},
		{0.42, 0.60, 1.00, 0.75, 0.50},
		{0.10, 0.30, 0.75, 1.00, 0.80},
		{-0.20, 0.05, 0.50, 0.80, 1.00},
	}
	labels := []string{"A", "B", "C", "D", "E"}
	fig := New(WithWidth(600), WithHeight(500), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Heatmap(matrix, chart.HeatmapStyle{
		RowLabels:  labels,
		ColLabels:  labels,
		ColorMap:   colormap.Viridis,
		CellLabels: true,
	})
	ax.SetTitle("Heatmap — Dark Theme")
	ax.NoLegend()
	checkGolden(t, "heatmap_dark_theme", fig.SVG())
}

func TestGoldenStackedAreaDarkTheme(t *testing.T) {
	series := [][]float64{
		makeSin(xs),
		makeCos(xs),
	}
	for i := range series[0] {
		series[0][i] += 1.2
		series[1][i] += 1.2
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.StackedArea(xs, series, []string{"Product A", "Product B"})
	ax.SetTitle("Stacked Area — Dark Theme")
	checkGolden(t, "stacked_area_dark_theme", fig.SVG())
}

func TestGoldenErrorBarsDarkTheme(t *testing.T) {
	ebXs := []float64{1, 2, 3, 4, 5, 6, 7}
	ebYs := []float64{2.1, 3.8, 3.2, 5.1, 4.6, 6.3, 5.8}
	yErr := []float64{0.4, 0.6, 0.5, 0.7, 0.5, 0.8, 0.6}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Scatter(ebXs, ebYs, chart.ScatterStyle{Label: "measurements", MarkerSize: 6})
	ax.ErrorBars(ebXs, ebYs, yErr)
	ax.SetTitle("Error Bars — Dark Theme")
	checkGolden(t, "error_bars_dark_theme", fig.SVG())
}

// --- Mixed chart types on same axes ---

func TestGoldenMixedLineScatter(t *testing.T) {
	trendXs := xs[:30]
	trendYs := make([]float64, 30)
	scatterXs := []float64{0.5, 1.2, 2.1, 2.8, 3.5, 4.2, 5.0, 5.8, 6.5, 7.2}
	scatterYs := []float64{0.08, 0.22, 0.45, 0.65, 0.75, 0.88, 0.95, 0.92, 0.82, 0.70}
	for i, x := range trendXs {
		trendYs[i] = math.Sin(x)
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(scatterXs, scatterYs, chart.ScatterStyle{Label: "observed", MarkerSize: 6})
	ax.Line(trendXs, trendYs, chart.LineStyle{Label: "sin(x) fit", Smooth: true, Dash: []float64{6, 3}})
	ax.SetTitle("Scatter + Line Overlay")
	checkGolden(t, "mixed_line_scatter", fig.SVG())
}

func TestGoldenMixedBarLine(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	monthly := []float64{120, 145, 132, 168, 155, 180}
	cumXs := []float64{0, 1, 2, 3, 4, 5}
	cumYs := make([]float64, 6)
	cumYs[0] = monthly[0]
	for i := 1; i < 6; i++ {
		cumYs[i] = cumYs[i-1] + monthly[i]
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, monthly, chart.BarStyle{Label: "Monthly"})
	ax.Line(cumXs, cumYs, chart.LineStyle{Label: "Cumulative", LineWidth: 2.5, MarkerSize: 5}).OnY2()
	ax.SetYLabel("Monthly Sales").SetY2Label("Cumulative")
	ax.SetTitle("Monthly vs Cumulative Sales")
	checkGolden(t, "mixed_bar_line_y2", fig.SVG())
}

// --- Subplot variants ---

func TestGoldenSubplot2x1(t *testing.T) {
	cats := []string{"Q1", "Q2", "Q3", "Q4"}
	vals := []float64{110, 145, 132, 168}
	fig := New(WithWidth(700), WithHeight(700))
	axes := fig.SubPlots(2, 1)
	axes[0][0].Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	axes[0][0].Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	axes[0][0].SetTitle("Waveforms")
	axes[1][0].Bar(cats, vals, chart.BarStyle{Label: "revenue"})
	axes[1][0].SetTitle("Quarterly Revenue")
	checkGolden(t, "subplot_2x1", fig.SVG())
}

func TestGoldenSubplotCustomSpacing(t *testing.T) {
	fig := New(WithWidth(1000), WithHeight(420), WithSpacing(40))
	axes := fig.SubPlots(1, 2)
	axes[0][0].Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	axes[0][0].SetTitle("Sine")
	axes[0][1].Line(xs, cosYs, chart.LineStyle{Label: "cos(x)", Smooth: true})
	axes[0][1].SetTitle("Cosine")
	checkGolden(t, "subplot_custom_spacing", fig.SVG())
}

// --- Pie variants ---

func TestGoldenPieManySlices(t *testing.T) {
	labels := []string{"Direct", "Organic", "Social", "Email", "Referral", "Paid Search", "Display", "Video", "Other"}
	values := []float64{28, 22, 14, 12, 9, 7, 4, 2, 2}
	fig := New(WithWidth(650), WithHeight(520))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.5})
	ax.SetTitle("Traffic Sources (9 channels)")
	ax.NoLegend()
	checkGolden(t, "pie_many_slices", fig.SVG())
}

func TestGoldenPieMinimalTheme(t *testing.T) {
	labels := []string{"Chrome", "Safari", "Firefox", "Edge", "Other"}
	values := []float64{65, 19, 4, 4, 8}
	fig := New(WithWidth(600), WithHeight(500), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Pie(labels, values, chart.PieStyle{DonutRadius: 0.5})
	ax.SetTitle("Browser Share — Minimal")
	ax.NoLegend()
	checkGolden(t, "pie_minimal_theme", fig.SVG())
}

// --- Box variants ---

func TestGoldenBoxManyGroups(t *testing.T) {
	g1 := []float64{2, 4, 5, 6, 7, 8, 9, 10, 15}
	g2 := []float64{5, 6, 7, 8, 8, 9, 10, 11, 12}
	g3 := []float64{1, 3, 5, 6, 7, 8, 9, 18, 20}
	g4 := []float64{4, 6, 7, 8, 9, 10, 11, 12, 16}
	g5 := []float64{3, 5, 6, 7, 8, 9, 10, 13, 14}
	fig := New(WithWidth(900), WithHeight(500))
	ax := fig.AddAxes()
	ax.Box([][]float64{g1, g2, g3, g4, g5},
		chart.BoxStyle{Labels: []string{"Control", "Treat A", "Treat B", "Treat C", "Treat D"}},
	)
	ax.SetYLabel("Score")
	ax.SetTitle("Five-Group Box Plot")
	checkGolden(t, "box_many_groups", fig.SVG())
}

func TestGoldenBoxMinimalTheme(t *testing.T) {
	groupA := []float64{2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 15}
	groupB := []float64{1, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12}
	groupC := []float64{0, 1, 3, 5, 6, 6, 7, 8, 9, 18, 20}
	fig := New(WithWidth(700), WithHeight(500), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Box([][]float64{groupA, groupB, groupC},
		chart.BoxStyle{Labels: []string{"Control", "Treatment A", "Treatment B"}},
	)
	ax.SetYLabel("Score")
	ax.SetTitle("Box Plot — Minimal Theme")
	checkGolden(t, "box_minimal_theme", fig.SVG())
}

// --- Histogram variants ---

func TestGoldenHistogramBimodal(t *testing.T) {
	vals := make([]float64, 0, 160)
	for i := 0; i < 80; i++ {
		v := 30.0 + math.Sin(float64(i)*0.3)*8 + float64(i%5)*0.5
		vals = append(vals, v)
	}
	for i := 0; i < 80; i++ {
		v := 70.0 + math.Cos(float64(i)*0.3)*8 + float64(i%5)*0.5
		vals = append(vals, v)
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Bins: 20, Label: "bimodal"})
	ax.SetTitle("Bimodal Distribution")
	ax.SetXLabel("Value").SetYLabel("Count")
	checkGolden(t, "histogram_bimodal", fig.SVG())
}

func TestGoldenHistogramMinimalTheme(t *testing.T) {
	vals := make([]float64, 200)
	for i := range vals {
		u1 := float64(i*7+1) / 200.0
		u2 := float64(i*13+1) / 200.0
		vals[i] = math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)*15 + 50
	}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.Histogram(vals, chart.HistogramStyle{Label: "response time", Bins: 16})
	ax.SetTitle("Histogram — Minimal Theme")
	checkGolden(t, "histogram_minimal_theme", fig.SVG())
}

// --- Grouped bar variants ---

func TestGoldenGroupedBar2Series(t *testing.T) {
	cats := []string{"North", "South", "East", "West"}
	a := []float64{48, 62, 55, 71}
	b := []float64{55, 58, 70, 65}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Bar(cats, a, chart.BarStyle{Label: "2023"})
	ax.Bar(cats, b, chart.BarStyle{Label: "2024"})
	ax.SetYLabel("Revenue ($M)")
	ax.SetTitle("Grouped Bar — 2 Series")
	checkGolden(t, "grouped_bar_2series", fig.SVG())
}

func TestGoldenGroupedBarDarkTheme(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar", "Apr"}
	q1 := []float64{42, 58, 71, 65}
	q2 := []float64{55, 62, 68, 78}
	q3 := []float64{38, 51, 60, 70}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Bar(cats, q1, chart.BarStyle{Label: "Q1"})
	ax.Bar(cats, q2, chart.BarStyle{Label: "Q2"})
	ax.Bar(cats, q3, chart.BarStyle{Label: "Q3"})
	ax.SetTitle("Grouped Bar — Dark Theme")
	checkGolden(t, "grouped_bar_dark_theme", fig.SVG())
}

// --- Stacked bar variants ---

func TestGoldenStackedBarManySeries(t *testing.T) {
	cats := []string{"2020", "2021", "2022", "2023"}
	series := [][]float64{
		{20, 25, 28, 32},
		{15, 18, 20, 24},
		{10, 12, 14, 16},
		{8, 9, 11, 13},
	}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.StackedBar(cats, series, []string{"APAC", "EMEA", "Americas", "Other"}, chart.StackedBarStyle{})
	ax.SetYLabel("Revenue ($M)")
	ax.SetTitle("Stacked Bar — 4 Series")
	checkGolden(t, "stacked_bar_many_series", fig.SVG())
}

// --- Stacked area variants ---

func TestGoldenStackedAreaThreeSeries(t *testing.T) {
	n := 12
	months := make([]float64, n)
	for i := range months {
		months[i] = float64(i + 1)
	}
	mobile := []float64{42, 44, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55}
	desktop := []float64{28, 27, 27, 26, 26, 25, 25, 24, 24, 23, 23, 22}
	tablet := []float64{10, 10, 10, 11, 11, 11, 11, 12, 12, 12, 12, 13}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.StackedArea(months, [][]float64{mobile, desktop, tablet}, []string{"Mobile", "Desktop", "Tablet"})
	ax.SetXLabel("Month").SetYLabel("Users (M)")
	ax.SetTitle("Stacked Area — 3 Series")
	checkGolden(t, "stacked_area_three_series", fig.SVG())
}

// --- HBar variants ---

func TestGoldenHBarMinimalTheme(t *testing.T) {
	cats := []string{"Node.js", "React", "jQuery", "Express", "Vue.js"}
	vals := []float64{42.7, 40.6, 21.9, 22.8, 15.6}
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Minimal))
	ax := fig.AddAxes()
	ax.HBar(cats, vals, chart.BarStyle{Label: "usage %"})
	ax.SetXLabel("Usage (%)").SetTitle("HBar — Minimal Theme")
	checkGolden(t, "hbar_minimal_theme", fig.SVG())
}

// --- Annotation on dark theme ---

func TestGoldenAnnotationDarkTheme(t *testing.T) {
	peakX := xs[7]
	fig := New(WithWidth(800), WithHeight(450), WithTheme(theme.Dark))
	ax := fig.AddAxes()
	ax.Line(xs, sinYs, chart.LineStyle{Label: "sin(x)", Smooth: true})
	ax.HLine(0.5, Label("threshold"), Dash(6, 3))
	ax.VLine(math.Pi, Label("π"))
	ax.Annotate(peakX, sinYs[7], "peak", ArrowDown())
	ax.SetTitle("Annotations — Dark Theme")
	checkGolden(t, "annotation_dark_theme", fig.SVG())
}

// --- Error bars both axes ---

func TestGoldenErrorBarsBothAxes(t *testing.T) {
	ebXs := []float64{1, 2, 3, 4, 5}
	ebYs := []float64{2.1, 3.8, 3.2, 5.1, 4.6}
	xErr := []float64{0.3, 0.4, 0.2, 0.5, 0.3}
	yErr := []float64{0.4, 0.6, 0.5, 0.7, 0.5}
	fig := New(WithWidth(800), WithHeight(450))
	ax := fig.AddAxes()
	ax.Scatter(ebXs, ebYs, chart.ScatterStyle{Label: "data", MarkerSize: 7})
	ax.ErrorBars(ebXs, ebYs, yErr)
	ax.ErrorBarsX(ebXs, ebYs, xErr)
	ax.SetTitle("Error Bars — Both Axes")
	checkGolden(t, "error_bars_both_axes", fig.SVG())
}

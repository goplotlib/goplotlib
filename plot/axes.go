package plot

import (
	"fmt"
	"math"

	"github.com/goplotlib/goplotlib/chart"
	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/render/svg"
	"github.com/goplotlib/goplotlib/scale"
	"github.com/goplotlib/goplotlib/theme"
)

// LegendPosition specifies where to draw the legend.
type LegendPosition int

const (
	LegendOutsideRight  LegendPosition = iota // default: outside, right side — vertical list
	LegendOutsideBottom                        // outside, below plot — horizontal row
	LegendTopRight                             // inside, top-right corner
	LegendTopLeft                              // inside, top-left corner
	LegendBottomRight                          // inside, bottom-right corner
	LegendBottomLeft                           // inside, bottom-left corner
	LegendNone
)

// Axes holds a set of charts and layout configuration for one plot area.
type Axes struct {
	charts      []chart.Chart
	annotations []annotatable
	title       string
	xLabel      string
	yLabel      string
	legendPos   LegendPosition
	showLegend  bool
	colorIdx    int
	xTickCount  int
	yTickCount  int
	// explicit axis limits; nil means auto-range from data
	xLimMin *float64
	xLimMax *float64
	yLimMin *float64
	yLimMax *float64
	// custom tick formatters; nil means use scale default
	xTickFormat func(float64) string
	yTickFormat func(float64) string
	// custom scale factories; nil means use the default (linear or categorical)
	xScaleFactory scale.Factory
	yScaleFactory scale.Factory
	// grid position within the figure layout (0-indexed)
	row, col int
	// secondary y-axis
	chartsY2    []bool   // parallel to charts; true = draw on right y-axis
	y2Label     string
	y2LimMin    *float64
	y2LimMax    *float64
	y2TickCount int
}

func newAxes() *Axes {
	return &Axes{
		legendPos:   LegendTopRight,
		showLegend:  true,
		xTickCount:  6,
		yTickCount:  6,
		y2TickCount: 6,
	}
}

// addChart appends ch to the chart list and initialises its Y2 slot to false.
func (a *Axes) addChart(ch chart.Chart) *Axes {
	a.charts = append(a.charts, ch)
	a.chartsY2 = append(a.chartsY2, false)
	return a
}

// Line adds a line chart to the axes.
func (a *Axes) Line(xs, ys []float64, style chart.LineStyle) *Axes {
	return a.addChart(chart.NewLine(xs, ys, style))
}

// Area adds a filled area chart to the axes (equivalent to Line with fill enabled).
func (a *Axes) Area(xs, ys []float64, style chart.LineStyle) *Axes {
	style.Fill = true
	return a.Line(xs, ys, style)
}

// Step adds a step (staircase) chart to the axes.
func (a *Axes) Step(xs, ys []float64, style chart.StepStyle) *Axes {
	return a.addChart(chart.NewStep(xs, ys, style))
}

// Bar adds a bar chart to the axes.
// When multiple Bar calls share the same categories, they are rendered as a grouped bar chart.
func (a *Axes) Bar(categories []string, values []float64, style chart.BarStyle) *Axes {
	return a.addChart(chart.NewBar(categories, values, style))
}

// HBar adds a horizontal bar chart to the axes.
func (a *Axes) HBar(categories []string, values []float64, style chart.BarStyle) *Axes {
	return a.addChart(chart.NewHBar(categories, values, style))
}

// StackedBar adds a stacked bar chart to the axes.
func (a *Axes) StackedBar(categories []string, series [][]float64, labels []string, style chart.StackedBarStyle) *Axes {
	return a.addChart(chart.NewStackedBar(categories, series, labels, style))
}

// StackedArea adds a stacked area chart to the axes.
func (a *Axes) StackedArea(xs []float64, series [][]float64, labels []string) *Axes {
	return a.addChart(chart.NewStackedArea(xs, series, labels))
}

// Box adds a box plot to the axes. Each inner slice is one group.
func (a *Axes) Box(groups [][]float64, style chart.BoxStyle) *Axes {
	return a.addChart(chart.NewBox(groups, style))
}

// Pie adds a pie or donut chart to the axes.
func (a *Axes) Pie(labels []string, values []float64, style chart.PieStyle) *Axes {
	return a.addChart(chart.NewPie(labels, values, style))
}

// Bubble adds a bubble chart to the axes.
// sizes encodes a third variable as circle radius; xs, ys, and sizes must have equal length.
func (a *Axes) Bubble(xs, ys, sizes []float64, style chart.BubbleStyle) *Axes {
	return a.addChart(chart.NewBubble(xs, ys, sizes, style))
}

// Histogram bins values and renders them as adjacent bars.
func (a *Axes) Histogram(values []float64, style chart.HistogramStyle) *Axes {
	return a.addChart(chart.NewHistogram(values, style))
}

// Heatmap adds a heatmap chart to the axes.
// matrix[row][col] is the scalar value for that cell.
func (a *Axes) Heatmap(matrix [][]float64, style chart.HeatmapStyle) *Axes {
	return a.addChart(chart.NewHeatmap(matrix, style))
}

// Scatter adds a scatter chart to the axes.
func (a *Axes) Scatter(xs, ys []float64, style chart.ScatterStyle) *Axes {
	return a.addChart(chart.NewScatter(xs, ys, style))
}

// SetTitle sets the axes title.
func (a *Axes) SetTitle(title string) *Axes { a.title = title; return a }

// SetXLabel sets the x-axis label.
func (a *Axes) SetXLabel(label string) *Axes { a.xLabel = label; return a }

// SetYLabel sets the y-axis label.
func (a *Axes) SetYLabel(label string) *Axes { a.yLabel = label; return a }

// Legend sets the legend position.
func (a *Axes) Legend(pos LegendPosition) *Axes { a.legendPos = pos; a.showLegend = true; return a }

// NoLegend disables the legend.
func (a *Axes) NoLegend() *Axes { a.showLegend = false; return a }

// SetXTicks sets the approximate number of x ticks.
func (a *Axes) SetXTicks(n int) *Axes { a.xTickCount = n; return a }

// SetYTicks sets the approximate number of y ticks.
func (a *Axes) SetYTicks(n int) *Axes { a.yTickCount = n; return a }

// SetXTickFormat sets a custom formatter for x-axis tick labels.
// The function receives the raw data value and returns the display string.
func (a *Axes) SetXTickFormat(f func(float64) string) *Axes { a.xTickFormat = f; return a }

// SetYTickFormat sets a custom formatter for y-axis tick labels.
// The function receives the raw data value and returns the display string.
func (a *Axes) SetYTickFormat(f func(float64) string) *Axes { a.yTickFormat = f; return a }

// ErrorBars adds vertical error bars to the axes.
// With 3 slices: symmetric bars of ±yErr. With 4 slices: asymmetric, yLow down and yHigh up.
func (a *Axes) ErrorBars(xs, ys, yLow []float64, yHigh ...[]float64) *Axes {
	hi := yLow
	if len(yHigh) > 0 {
		hi = yHigh[0]
	}
	return a.addChart(chart.NewErrorBars(xs, ys, yLow, hi, color.Transparent))
}

// ErrorBarsX adds symmetric horizontal error bars to the axes.
func (a *Axes) ErrorBarsX(xs, ys, xErr []float64) *Axes {
	return a.addChart(chart.NewErrorBarsX(xs, ys, xErr, color.Transparent))
}

// SetXScale sets a custom scale factory for the x-axis (e.g. scale.Log(10)).
func (a *Axes) SetXScale(f scale.Factory) *Axes { a.xScaleFactory = f; return a }

// SetYScale sets a custom scale factory for the y-axis (e.g. scale.Log(10)).
func (a *Axes) SetYScale(f scale.Factory) *Axes { a.yScaleFactory = f; return a }

// OnY2 assigns the most recently added series to the secondary (right) y-axis.
// Call immediately after the series method: ax.Line(...).OnY2()
func (a *Axes) OnY2() *Axes {
	if len(a.chartsY2) > 0 {
		a.chartsY2[len(a.chartsY2)-1] = true
	}
	return a
}

// SetY2Label sets the label for the secondary y-axis.
func (a *Axes) SetY2Label(label string) *Axes { a.y2Label = label; return a }

// SetY2Lim fixes the visible secondary y-axis range.
func (a *Axes) SetY2Lim(min, max float64) *Axes {
	if !math.IsNaN(min) {
		v := min
		a.y2LimMin = &v
	}
	if !math.IsNaN(max) {
		v := max
		a.y2LimMax = &v
	}
	return a
}

// SetY2Ticks sets the approximate number of ticks on the secondary y-axis.
func (a *Axes) SetY2Ticks(n int) *Axes { a.y2TickCount = n; return a }

// HLine draws a horizontal reference line at the given y value.
func (a *Axes) HLine(y float64, opts ...AnnotationOption) *Axes {
	s := newAnnotStyle()
	applyAnnotOpts(&s, opts)
	a.annotations = append(a.annotations, &hLineAnnotation{y: y, style: s})
	return a
}

// VLine draws a vertical reference line at the given x value.
func (a *Axes) VLine(x float64, opts ...AnnotationOption) *Axes {
	s := newAnnotStyle()
	applyAnnotOpts(&s, opts)
	a.annotations = append(a.annotations, &vLineAnnotation{x: x, style: s})
	return a
}

// HSpan draws a horizontal shaded band between y1 and y2.
func (a *Axes) HSpan(y1, y2 float64, opts ...AnnotationOption) *Axes {
	s := newAnnotStyle()
	applyAnnotOpts(&s, opts)
	a.annotations = append(a.annotations, &hSpanAnnotation{y1: y1, y2: y2, style: s})
	return a
}

// VSpan draws a vertical shaded band between x1 and x2.
func (a *Axes) VSpan(x1, x2 float64, opts ...AnnotationOption) *Axes {
	s := newAnnotStyle()
	applyAnnotOpts(&s, opts)
	a.annotations = append(a.annotations, &vSpanAnnotation{x1: x1, x2: x2, style: s})
	return a
}

// Annotate places a text label at the given data coordinates with an optional arrow.
func (a *Axes) Annotate(x, y float64, text string, opts ...AnnotationOption) *Axes {
	s := newAnnotStyle()
	applyAnnotOpts(&s, opts)
	a.annotations = append(a.annotations, &textAnnotation{x: x, y: y, text: text, style: s})
	return a
}

// SetXLim fixes the visible x-axis range. Pass Auto for either bound to keep it auto-ranged.
func (a *Axes) SetXLim(min, max float64) *Axes {
	if !math.IsNaN(min) {
		v := min
		a.xLimMin = &v
	}
	if !math.IsNaN(max) {
		v := max
		a.xLimMax = &v
	}
	return a
}

// SetYLim fixes the visible y-axis range. Pass Auto for either bound to keep it auto-ranged.
func (a *Axes) SetYLim(min, max float64) *Axes {
	if !math.IsNaN(min) {
		v := min
		a.yLimMin = &v
	}
	if !math.IsNaN(max) {
		v := max
		a.yLimMax = &v
	}
	return a
}

// assignColors assigns palette colors to charts that don't have an explicit color.
func (a *Axes) assignColors(palette []color.Color) {
	idx := 0
	lastColor := palette[0] // fallback for error bars with no preceding series
	for _, ch := range a.charts {
		switch c := ch.(type) {
		case *chart.StackedBarChart:
			c.AssignColors(palette, &idx)
		case *chart.StackedAreaChart:
			c.AssignColors(palette, &idx)
		case *chart.PieChart:
			c.AssignColors(palette, &idx)
		case *chart.ErrorBarsChart:
			// Error bars inherit the color of the preceding series; don't advance palette index.
			if c.Color().IsTransparent() {
				c.SetColor(lastColor)
			}
		default:
			if ch.Color().IsTransparent() {
				col := palette[idx%len(palette)]
				idx++
				lastColor = col
				switch c := ch.(type) {
				case *chart.LineChart:
					c.SetColor(col)
				case *chart.BarChart:
					c.SetColor(col)
				case *chart.HBarChart:
					c.SetColor(col)
				case *chart.StepChart:
					c.SetColor(col)
				case *chart.ScatterChart:
					c.SetColor(col)
				case *chart.HistogramChart:
					c.SetColor(col)
				case *chart.BubbleChart:
					c.SetColor(col)
				case *chart.BoxChart:
					c.SetColor(col)
				}
			} else {
				lastColor = ch.Color()
			}
		}
	}
}

// render draws the axes content within the given cell bounds (in figure-level pixel coords).
// axesIdx is the 0-based index of this axes within the figure, used to generate a stable clip ID.
func (a *Axes) render(canvas *svg.Canvas, cellX, cellY, cellW, cellH float64, t theme.Theme, axesIdx int) {
	// Pie charts have their own layout — no axes, grid, or spines.
	for _, ch := range a.charts {
		if pc, ok := ch.(*chart.PieChart); ok {
			a.renderPie(canvas, pc, cellX, cellY, cellW, cellH, t)
			return
		}
	}

	pad := t.Padding

	// Plot area bounds within the cell
	plotX := cellX + pad.Left
	plotY := cellY + pad.Top
	plotW := cellW - pad.Left - pad.Right
	plotH := cellH - pad.Top - pad.Bottom

	// Reserve space for outside legends.
	if a.showLegend && a.hasLegend() {
		switch a.legendPos {
		case LegendOutsideRight:
			lw, _ := a.legendBoxSize()
			plotW -= lw + legendOutGap
		case LegendOutsideBottom:
			_, lh := a.legendBottomBoxSize()
			plotH -= lh + legendOutGap
		}
	}

	if plotW <= 0 || plotH <= 0 {
		return
	}

	// Detect chart types and compute merged data ranges
	var dataXMin, dataXMax, dataYMin, dataYMax float64
	hasData := false
	var barChart *chart.BarChart
	var hBarChart *chart.HBarChart
	var stackedBarChart *chart.StackedBarChart
	var heatmapChart *chart.HeatmapChart
	var histogramChart *chart.HistogramChart
	var boxChart *chart.BoxChart
	var categories []string    // x-axis category labels (bar / stacked bar)
	var hBarCategories []string // y-axis category labels (hbar)
	var barCharts []*chart.BarChart // all BarCharts — for grouped bar detection

	for i, ch := range a.charts {
		if i < len(a.chartsY2) && a.chartsY2[i] {
			continue
		}
		xMin, xMax, yMin, yMax := ch.DataRange()
		switch c := ch.(type) {
		case *chart.BarChart:
			if barChart == nil {
				barChart = c
				categories = c.Categories()
			}
			barCharts = append(barCharts, c)
		case *chart.HBarChart:
			if hBarChart == nil {
				hBarChart = c
				hBarCategories = c.Categories()
			}
		case *chart.StackedBarChart:
			if stackedBarChart == nil {
				stackedBarChart = c
				categories = c.Categories()
			}
		case *chart.HeatmapChart:
			if heatmapChart == nil {
				heatmapChart = c
			}
		case *chart.HistogramChart:
			if histogramChart == nil {
				histogramChart = c
			}
		case *chart.BoxChart:
			if boxChart == nil {
				boxChart = c
				categories = c.Categories()
			}
		}
		if !hasData {
			dataXMin, dataXMax = xMin, xMax
			dataYMin, dataYMax = yMin, yMax
			hasData = true
		} else {
			if xMin < dataXMin {
				dataXMin = xMin
			}
			if xMax > dataXMax {
				dataXMax = xMax
			}
			if yMin < dataYMin {
				dataYMin = yMin
			}
			if yMax > dataYMax {
				dataYMax = yMax
			}
		}
	}

	// Separate Y2 range
	var y2DataMin, y2DataMax float64
	hasY2Data := false
	for i, ch := range a.charts {
		if i >= len(a.chartsY2) || !a.chartsY2[i] {
			continue
		}
		_, _, yMin, yMax := ch.DataRange()
		if !hasY2Data {
			y2DataMin, y2DataMax = yMin, yMax
			hasY2Data = true
		} else {
			if yMin < y2DataMin {
				y2DataMin = yMin
			}
			if yMax > y2DataMax {
				y2DataMax = yMax
			}
		}
	}

	// Reserve right-side space for Y2 ticks and label
	const y2RightPad = 55.0 // space for right-side tick labels and label
	hasY2 := hasY2Data
	if hasY2 {
		plotW -= y2RightPad
	}

	if !hasData {
		dataXMin, dataXMax = 0, 1
		dataYMin, dataYMax = 0, 1
	}

	// Apply grouped bar layout when multiple Bar series share the same categories
	if len(barCharts) > 1 {
		for i, bc := range barCharts {
			bc.SetGroupInfo(i, len(barCharts))
		}
	} else if len(barCharts) == 1 {
		barCharts[0].SetGroupInfo(0, 1)
	}

	// Always include y=0 for vertical bar/stacked bar/histogram (only when y range is auto and not log scale)
	if (barChart != nil || stackedBarChart != nil || histogramChart != nil) && a.yLimMin == nil && a.yScaleFactory == nil {
		if dataYMin > 0 {
			dataYMin = 0
		}
	}
	// Always include x=0 for horizontal bar (values extend from zero)
	if hBarChart != nil && a.xLimMin == nil {
		if dataXMin > 0 {
			dataXMin = 0
		}
	}

	// Add 5% padding to value axes (skip y for HBar since y is categorical)
	if hBarChart == nil {
		ySpan := dataYMax - dataYMin
		if ySpan == 0 {
			ySpan = 1
		}
		if a.yLimMax == nil {
			dataYMax += ySpan * 0.05
		}
		if a.yLimMin == nil && dataYMin < 0 {
			dataYMin -= ySpan * 0.05
		}
	}
	if hBarChart != nil {
		xSpan := dataXMax - dataXMin
		if xSpan == 0 {
			xSpan = 1
		}
		if a.xLimMax == nil {
			dataXMax += xSpan * 0.05
		}
		if a.xLimMin == nil && dataXMin < 0 {
			dataXMin -= xSpan * 0.05
		}
	}

	// Override with explicit limits
	if a.xLimMin != nil {
		dataXMin = *a.xLimMin
	}
	if a.xLimMax != nil {
		dataXMax = *a.xLimMax
	}
	if a.yLimMin != nil {
		dataYMin = *a.yLimMin
	}
	if a.yLimMax != nil {
		dataYMax = *a.yLimMax
	}

	// Ensure xMin != xMax for linear x scales
	if hBarChart == nil && barChart == nil && stackedBarChart == nil {
		if dataXMax == dataXMin && a.xLimMin == nil && a.xLimMax == nil {
			dataXMin -= 0.5
			dataXMax += 0.5
		}
	}

	// Build scales
	var xScale scale.Scale
	if heatmapChart != nil {
		colLabels := heatmapChart.ColLabels()
		if len(colLabels) == 0 {
			colLabels = makeIndexLabels(heatmapChart.NumCols())
		}
		xScale = scale.NewCategorical(colLabels, plotX, plotX+plotW, 0)
	} else if (barChart != nil || stackedBarChart != nil || boxChart != nil) && len(categories) > 0 {
		xScale = scale.NewCategorical(categories, plotX, plotX+plotW, 0.2)
	} else if a.xScaleFactory != nil {
		xScale = a.xScaleFactory(dataXMin, dataXMax, plotX, plotX+plotW)
	} else {
		xScale = scale.NewLinear(dataXMin, dataXMax, plotX, plotX+plotW)
	}

	var yScale scale.Scale
	if heatmapChart != nil {
		rowLabels := heatmapChart.RowLabels()
		if len(rowLabels) == 0 {
			rowLabels = makeIndexLabels(heatmapChart.NumRows())
		}
		yScale = scale.NewCategorical(rowLabels, plotY, plotY+plotH, 0)
	} else if hBarChart != nil && len(hBarCategories) > 0 {
		// HBar: categorical y-axis, categories top-to-bottom
		yScale = scale.NewCategorical(hBarCategories, plotY, plotY+plotH, 0.2)
	} else if a.yScaleFactory != nil {
		// Inverted: data min → bottom, data max → top
		yScale = a.yScaleFactory(dataYMin, dataYMax, plotY+plotH, plotY)
	} else {
		// Default: inverted linear (data min → bottom, data max → top)
		yScale = scale.NewLinear(dataYMin, dataYMax, plotY+plotH, plotY)
	}

	// Secondary y-axis scale
	var y2Scale scale.Scale
	if hasY2 {
		// Add 5% padding to Y2 range
		y2Span := y2DataMax - y2DataMin
		if y2Span == 0 {
			y2Span = 1
		}
		if a.y2LimMax == nil {
			y2DataMax += y2Span * 0.05
		}
		if a.y2LimMin == nil && y2DataMin < 0 {
			y2DataMin -= y2Span * 0.05
		}
		if a.y2LimMin != nil {
			y2DataMin = *a.y2LimMin
		}
		if a.y2LimMax != nil {
			y2DataMax = *a.y2LimMax
		}
		y2Scale = scale.NewLinear(y2DataMin, y2DataMax, plotY+plotH, plotY)
	}

	// --- Draw plot area background ---
	canvas.DrawRect(plotX, plotY, plotW, plotH, 0, 0, render.Style{
		Fill:        t.PlotBackground,
		FillOpacity: 1.0,
	})

	// --- Clip path for plot area ---
	clipID := fmt.Sprintf("axes-clip-%d", axesIdx)
	canvas.BeginClip(clipID, plotX, plotY, plotW, plotH)

	// --- Grid lines ---
	yTicks := yScale.Ticks(a.yTickCount)
	xTicks := xScale.Ticks(a.xTickCount)

	// Apply custom tick formatters if set
	if a.xTickFormat != nil {
		for i := range xTicks {
			xTicks[i].Label = a.xTickFormat(xTicks[i].Value)
		}
	}
	if a.yTickFormat != nil {
		for i := range yTicks {
			yTicks[i].Label = a.yTickFormat(yTicks[i].Value)
		}
	}

	gridStyle := render.Style{
		Stroke:      t.GridColor,
		StrokeWidth: t.GridWidth,
		Dash:        t.GridDash,
	}

	// Horizontal grid lines — skip for HBar and Heatmap (categorical axes, cells provide separation)
	if hBarChart == nil && heatmapChart == nil {
		canvas.BeginClipGroup(clipID)
		for _, tick := range yTicks {
			canvas.DrawLine(plotX, tick.Pos, plotX+plotW, tick.Pos, gridStyle)
		}
		canvas.EndGroup()
	}

	// Vertical grid lines — skip for vertical bar / stacked bar / box / heatmap (x is categorical)
	if barChart == nil && stackedBarChart == nil && heatmapChart == nil && boxChart == nil {
		canvas.BeginClipGroup(clipID)
		for _, tick := range xTicks {
			canvas.DrawLine(tick.Pos, plotY, tick.Pos, plotY+plotH, gridStyle)
		}
		canvas.EndGroup()
	}

	// --- Set index on line charts for unique gradient IDs ---
	for i, ch := range a.charts {
		if lc, ok := ch.(*chart.LineChart); ok {
			lc.SetIndex(i)
		}
	}

	// --- Draw chart data (clipped) ---
	canvas.BeginClipGroup(clipID)
	for i, ch := range a.charts {
		ys := yScale
		if i < len(a.chartsY2) && a.chartsY2[i] && y2Scale != nil {
			ys = y2Scale
		}
		ch.Draw(canvas, xScale, ys)
	}
	canvas.EndGroup()

	// --- Bubble labels (drawn outside clip so they are never truncated at edges) ---
	for i, ch := range a.charts {
		if bc, ok := ch.(*chart.BubbleChart); ok {
			ys := yScale
			if i < len(a.chartsY2) && a.chartsY2[i] && y2Scale != nil {
				ys = y2Scale
			}
			bc.DrawLabels(canvas, xScale, ys)
		}
	}

	// --- Annotations (on top of data, below axes labels) ---
	for _, ann := range a.annotations {
		ann.draw(canvas, xScale, yScale, plotX, plotY, plotW, plotH, t, clipID)
	}

	// --- Axis spines ---
	spineStyle := render.Style{
		Stroke:      t.AxisColor,
		StrokeWidth: t.AxisWidth,
	}
	if t.SpineBottom {
		canvas.DrawLine(plotX, plotY+plotH, plotX+plotW, plotY+plotH, spineStyle)
	}
	if t.SpineLeft {
		canvas.DrawLine(plotX, plotY, plotX, plotY+plotH, spineStyle)
	}
	if t.SpineTop {
		canvas.DrawLine(plotX, plotY, plotX+plotW, plotY, spineStyle)
	}
	if t.SpineRight {
		canvas.DrawLine(plotX+plotW, plotY, plotX+plotW, plotY+plotH, spineStyle)
	}

	// --- X axis ticks and labels ---
	tickStyle := render.TextStyle{
		Color:      t.TextColor,
		FontSize:   t.TickFontSize,
		FontFamily: t.FontFamily,
		Anchor:     "middle",
		Baseline:   "hanging",
	}
	tickLineStyle := render.Style{
		Stroke:      t.AxisColor,
		StrokeWidth: 1.0,
	}
	for _, tick := range xTicks {
		// Tick mark
		canvas.DrawLine(tick.Pos, plotY+plotH, tick.Pos, plotY+plotH+5, tickLineStyle)
		// Label
		canvas.DrawText(tick.Pos, plotY+plotH+8, tick.Label, tickStyle)
	}

	// --- Y axis ticks and labels ---
	yTickLabelStyle := render.TextStyle{
		Color:      t.TextColor,
		FontSize:   t.TickFontSize,
		FontFamily: t.FontFamily,
		Anchor:     "end",
		Baseline:   "middle",
	}
	for _, tick := range yTicks {
		// Tick mark
		canvas.DrawLine(plotX-5, tick.Pos, plotX, tick.Pos, tickLineStyle)
		// Label
		canvas.DrawText(plotX-8, tick.Pos, tick.Label, yTickLabelStyle)
	}

	// --- Y2 axis ticks and labels (right side) ---
	if hasY2 && y2Scale != nil {
		y2Ticks := y2Scale.Ticks(a.y2TickCount)
		y2TickLabelStyle := render.TextStyle{
			Color:      t.TextColor,
			FontSize:   t.TickFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "start",
			Baseline:   "middle",
		}
		rightEdge := plotX + plotW
		for _, tick := range y2Ticks {
			canvas.DrawLine(rightEdge, tick.Pos, rightEdge+5, tick.Pos, tickLineStyle)
			canvas.DrawText(rightEdge+8, tick.Pos, tick.Label, y2TickLabelStyle)
		}
		// Y2 spine (right edge of plot)
		canvas.DrawLine(rightEdge, plotY, rightEdge, plotY+plotH, spineStyle)
	}

	// --- X axis label ---
	if a.xLabel != "" {
		labelStyle := render.TextStyle{
			Color:      t.TextColor,
			FontSize:   t.LabelFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "middle",
			Baseline:   "hanging",
		}
		canvas.DrawText(plotX+plotW/2, plotY+plotH+30, a.xLabel, labelStyle)
	}

	// --- Y axis label (rotated) ---
	if a.yLabel != "" {
		cx := plotX - 48
		cy := plotY + plotH/2
		// We use DrawText with a transform via a raw SVG approach.
		// The svg.Canvas doesn't expose raw SVG, so we use a workaround:
		// We write via a TextStyle and add rotation separately.
		// Since Canvas.DrawText doesn't support transform, we add a helper.
		drawRotatedText(canvas, cx, cy, a.yLabel, render.TextStyle{
			Color:      t.TextColor,
			FontSize:   t.LabelFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "middle",
			Baseline:   "middle",
		})
	}

	// --- Y2 axis label (rotated, right side) ---
	if a.y2Label != "" && hasY2 {
		cx := plotX + plotW + y2RightPad - 14
		cy := plotY + plotH/2
		drawRotatedText(canvas, cx, cy, a.y2Label, render.TextStyle{
			Color:      t.TextColor,
			FontSize:   t.LabelFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "middle",
			Baseline:   "middle",
		})
	}

	// --- Axes title ---
	if a.title != "" {
		titleStyle := render.TextStyle{
			Color:      t.TitleColor,
			FontSize:   t.TitleFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "middle",
			Baseline:   "auto",
			Bold:       true,
		}
		canvas.DrawText(plotX+plotW/2, plotY-15, a.title, titleStyle)
	}

	// --- Legend ---
	if a.showLegend && a.hasLegend() {
		a.drawLegend(canvas, plotX, plotY, plotW, plotH, t)
	}
}

// hasLegend returns true if there are legend entries to show.
func (a *Axes) hasLegend() bool {
	for _, ch := range a.charts {
		if _, ok := ch.(chart.MultiLegend); ok {
			return true
		}
	}
	if len(a.charts) <= 1 {
		for _, ch := range a.charts {
			if ch.Label() != "" {
				return true
			}
		}
		return false
	}
	for _, ch := range a.charts {
		if ch.Label() != "" {
			return true
		}
	}
	return len(a.charts) > 1
}

const (
	legendItemH    = 22.0
	legendMarkerW  = 20.0
	legendPadding  = 10.0
	legendFontSize = 12.0
	legendOutGap   = 12.0 // gap between plot right edge and outside legend
)

// legendBoxSize estimates the pixel dimensions of the legend box for the vertical (right) layout.
func (a *Axes) legendBoxSize() (w, h float64) {
	maxLabelLen := 0
	nEntries := 0
	for _, ch := range a.charts {
		if ml, ok := ch.(chart.MultiLegend); ok {
			for _, e := range ml.LegendEntries() {
				if len(e.Label) > maxLabelLen {
					maxLabelLen = len(e.Label)
				}
				nEntries++
			}
			continue
		}
		lbl := ch.Label()
		if lbl == "" {
			lbl = "Series"
		}
		if len(lbl) > maxLabelLen {
			maxLabelLen = len(lbl)
		}
		nEntries++
	}
	charWidth := legendFontSize * 0.55
	w = legendMarkerW + legendPadding*2 + float64(maxLabelLen)*charWidth + legendPadding
	h = float64(nEntries)*legendItemH + legendPadding*2
	return
}

// legendBottomBoxSize estimates the pixel dimensions of the legend box for the horizontal (bottom) layout.
func (a *Axes) legendBottomBoxSize() (w, h float64) {
	charWidth := legendFontSize * 0.55
	entryGap := 16.0 // gap between entries
	totalEntryW := 0.0
	for _, ch := range a.charts {
		if ml, ok := ch.(chart.MultiLegend); ok {
			for _, e := range ml.LegendEntries() {
				totalEntryW += legendMarkerW + 4 + float64(len(e.Label))*charWidth + entryGap
			}
			continue
		}
		lbl := ch.Label()
		if lbl == "" {
			lbl = "Series"
		}
		totalEntryW += legendMarkerW + 4 + float64(len(lbl))*charWidth + entryGap
	}
	w = totalEntryW + legendPadding*2 - entryGap // remove trailing gap, add padding
	h = legendItemH + legendPadding*2
	return
}

// legendEntry holds data for a single legend item.
type legendEntry struct {
	label string
	col   color.Color
	isBar bool
}

// drawLegend draws the legend box.
func (a *Axes) drawLegend(canvas *svg.Canvas, plotX, plotY, plotW, plotH float64, t theme.Theme) {
	const boxRadius = 4.0

	// Collect legend entries
	var entries []legendEntry
	for _, ch := range a.charts {
		if ml, ok := ch.(chart.MultiLegend); ok {
			for _, e := range ml.LegendEntries() {
				entries = append(entries, legendEntry{e.Label, e.Col, e.IsBar})
			}
			continue
		}
		lbl := ch.Label()
		if lbl == "" {
			lbl = "Series"
		}
		_, isVertBar := ch.(*chart.BarChart)
		_, isHBar := ch.(*chart.HBarChart)
		entries = append(entries, legendEntry{lbl, ch.Color(), isVertBar || isHBar})
	}
	if len(entries) == 0 {
		return
	}

	// Position and draw — outside-bottom uses a horizontal layout; all others use vertical.
	if a.legendPos == LegendOutsideBottom {
		a.drawLegendBottom(canvas, entries, plotX, plotY, plotW, plotH, t, boxRadius)
		return
	}

	boxW, boxH := a.legendBoxSize()

	// Position
	var boxX, boxY float64
	switch a.legendPos {
	case LegendOutsideRight:
		boxX = plotX + plotW + legendOutGap
		boxY = plotY + (plotH-boxH)/2
	case LegendTopLeft:
		boxX = plotX + 10
		boxY = plotY + 10
	case LegendBottomRight:
		boxX = plotX + plotW - boxW - 10
		boxY = plotY + plotH - boxH - 10
	case LegendBottomLeft:
		boxX = plotX + 10
		boxY = plotY + plotH - boxH - 10
	default: // LegendTopRight
		boxX = plotX + plotW - boxW - 10
		boxY = plotY + 10
	}

	// Legend background
	bgColor := t.Background.WithAlpha(0.9)
	borderColor := t.GridColor
	canvas.DrawRect(boxX, boxY, boxW, boxH, boxRadius, boxRadius, render.Style{
		Fill:        bgColor,
		FillOpacity: 1.0,
		Stroke:      borderColor,
		StrokeWidth: 0.8,
	})

	// Legend items
	textStyle := render.TextStyle{
		Color:      t.TextColor,
		FontSize:   legendFontSize,
		FontFamily: t.FontFamily,
		Anchor:     "start",
		Baseline:   "middle",
	}

	for i, e := range entries {
		itemY := boxY + legendPadding + float64(i)*legendItemH + legendItemH/2
		markerX := boxX + legendPadding
		labelX := markerX + legendMarkerW + 4

		if e.isBar {
			// Bar marker: small filled rectangle
			canvas.DrawRect(markerX, itemY-5, legendMarkerW, 10, 2, 2, render.Style{
				Fill:        e.col,
				FillOpacity: 1.0,
			})
		} else {
			// Line marker: short horizontal line
			canvas.DrawLine(markerX, itemY, markerX+legendMarkerW, itemY, render.Style{
				Stroke:      e.col,
				StrokeWidth: 2.5,
				LineCap:     "round",
			})
			// Small circle on line
			canvas.DrawCircle(markerX+legendMarkerW/2, itemY, 3, render.Style{
				Fill:        e.col,
				FillOpacity: 1.0,
			})
		}

		canvas.DrawText(labelX, itemY, e.label, textStyle)
	}
}

// drawLegendBottom draws a horizontal legend row below the plot area.
func (a *Axes) drawLegendBottom(canvas *svg.Canvas, entries []legendEntry, plotX, plotY, plotW, plotH float64, t theme.Theme, boxRadius float64) {
	charWidth := legendFontSize * 0.55
	entryGap := 16.0

	boxW, boxH := a.legendBottomBoxSize()
	// Center the box horizontally below the plot
	boxX := plotX + (plotW-boxW)/2
	boxY := plotY + plotH + legendOutGap

	// Legend background
	bgColor := t.Background.WithAlpha(0.9)
	borderColor := t.GridColor
	canvas.DrawRect(boxX, boxY, boxW, boxH, boxRadius, boxRadius, render.Style{
		Fill:        bgColor,
		FillOpacity: 1.0,
		Stroke:      borderColor,
		StrokeWidth: 0.8,
	})

	textStyle := render.TextStyle{
		Color:      t.TextColor,
		FontSize:   legendFontSize,
		FontFamily: t.FontFamily,
		Anchor:     "start",
		Baseline:   "middle",
	}

	itemY := boxY + boxH/2
	curX := boxX + legendPadding

	for _, e := range entries {
		markerX := curX
		labelX := markerX + legendMarkerW + 4
		labelW := float64(len(e.label)) * charWidth

		if e.isBar {
			canvas.DrawRect(markerX, itemY-5, legendMarkerW, 10, 2, 2, render.Style{
				Fill:        e.col,
				FillOpacity: 1.0,
			})
		} else {
			canvas.DrawLine(markerX, itemY, markerX+legendMarkerW, itemY, render.Style{
				Stroke:      e.col,
				StrokeWidth: 2.5,
				LineCap:     "round",
			})
			canvas.DrawCircle(markerX+legendMarkerW/2, itemY, 3, render.Style{
				Fill:        e.col,
				FillOpacity: 1.0,
			})
		}

		canvas.DrawText(labelX, itemY, e.label, textStyle)
		curX = labelX + labelW + entryGap
	}
}

// drawRotatedText draws text rotated -90 degrees around (cx, cy).
// This uses the svg.Canvas WriteRaw method, or falls back to a workaround.
func drawRotatedText(canvas *svg.Canvas, cx, cy float64, text string, s render.TextStyle) {
	canvas.DrawTextRotated(cx, cy, text, s, -90)
}

// renderPie draws a pie/donut chart centered in the cell with no axes or grid.
func (a *Axes) renderPie(canvas *svg.Canvas, pc *chart.PieChart, cellX, cellY, cellW, cellH float64, t theme.Theme) {
	// Reserve space at the top for an optional title
	topOffset := 0.0
	if a.title != "" {
		topOffset = 28.0
	}

	cx := cellX + cellW/2
	cy := cellY + topOffset + (cellH-topOffset)/2
	outerR := math.Min(cellW, cellH-topOffset)/2 - 30 // 30px margin for labels
	if outerR < 10 {
		return
	}

	pc.DrawInBounds(canvas, cx, cy, outerR, t)

	if a.title != "" {
		canvas.DrawText(cx, cellY+14, a.title, render.TextStyle{
			Color:      t.TitleColor,
			FontSize:   t.TitleFontSize,
			FontFamily: t.FontFamily,
			Anchor:     "middle",
			Baseline:   "middle",
			Bold:       true,
		})
	}
}

// makeIndexLabels generates ["0", "1", ..., "n-1"] for unlabeled heatmap axes.
func makeIndexLabels(n int) []string {
	labels := make([]string, n)
	for i := range labels {
		labels[i] = fmt.Sprintf("%d", i)
	}
	return labels
}

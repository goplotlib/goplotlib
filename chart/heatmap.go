package chart

import (
	"fmt"
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/colormap"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// HeatmapChart renders a color-encoded matrix.
// matrix[row][col] is the value at that cell.
type HeatmapChart struct {
	matrix [][]float64
	style  HeatmapStyle
}

// NewHeatmap creates a new HeatmapChart.
func NewHeatmap(matrix [][]float64, style HeatmapStyle) *HeatmapChart {
	if style.ColorMap == nil {
		style.ColorMap = colormap.Viridis
	}
	return &HeatmapChart{matrix: matrix, style: style}
}

func (h *HeatmapChart) Label() string      { return "" }
func (h *HeatmapChart) Color() color.Color { return color.Transparent }

// RowLabels returns the row label strings (may be nil).
func (h *HeatmapChart) RowLabels() []string { return h.style.RowLabels }

// ColLabels returns the column label strings (may be nil).
func (h *HeatmapChart) ColLabels() []string { return h.style.ColLabels }

// NumRows returns the number of rows in the matrix.
func (h *HeatmapChart) NumRows() int { return len(h.matrix) }

// NumCols returns the number of columns in the matrix.
func (h *HeatmapChart) NumCols() int { return h.numCols() }

// DataRange returns dummy ranges — the heatmap uses categorical scales for both axes.
func (h *HeatmapChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	cols := h.numCols()
	rows := len(h.matrix)
	if cols == 0 || rows == 0 {
		return 0, 1, 0, 1
	}
	return 0, float64(cols - 1), 0, float64(rows - 1)
}

func (h *HeatmapChart) numCols() int {
	if len(h.matrix) == 0 {
		return 0
	}
	return len(h.matrix[0])
}

// valueRange computes the min and max values in the matrix.
func (h *HeatmapChart) valueRange() (vMin, vMax float64) {
	vMin = math.Inf(1)
	vMax = math.Inf(-1)
	for _, row := range h.matrix {
		for _, v := range row {
			if v < vMin {
				vMin = v
			}
			if v > vMax {
				vMax = v
			}
		}
	}
	if math.IsInf(vMin, 1) {
		vMin, vMax = 0, 1
	}
	return
}

// normalizeVal maps a value to [0, 1] based on the configured scale type.
func (h *HeatmapChart) normalizeVal(v, vMin, vMax float64) float64 {
	if h.style.DivergingScale {
		// Center at zero: map [-absMax, +absMax] → [0, 1]
		absMax := math.Max(math.Abs(vMin), math.Abs(vMax))
		if absMax == 0 {
			return 0.5
		}
		return (v/absMax + 1) / 2
	}
	span := vMax - vMin
	if span == 0 {
		return 0.5
	}
	return (v - vMin) / span
}

// luminance computes the relative luminance of a color (0 = black, 1 = white).
func luminance(c color.Color) float64 {
	linearize := func(v uint8) float64 {
		s := float64(v) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	r := linearize(c.R)
	g := linearize(c.G)
	b := linearize(c.B)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// Draw renders the heatmap onto the canvas.
// xScale and yScale must be *scale.CategoricalScale instances.
func (h *HeatmapChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	xCat, xOK := xScale.(*scale.CategoricalScale)
	yCat, yOK := yScale.(*scale.CategoricalScale)
	if !xOK || !yOK {
		return
	}
	if len(h.matrix) == 0 || h.numCols() == 0 {
		return
	}

	vMin, vMax := h.valueRange()
	cellW := xCat.BandWidth()
	cellH := yCat.BandWidth()
	cm := h.style.ColorMap

	for rowIdx, row := range h.matrix {
		centerY := yCat.Map(float64(rowIdx))
		cellTop := centerY - cellH/2
		for colIdx, v := range row {
			centerX := xCat.Map(float64(colIdx))
			cellLeft := centerX - cellW/2

			t := h.normalizeVal(v, vMin, vMax)
			cellColor := cm(t)

			canvas.DrawRect(cellLeft, cellTop, cellW, cellH, 0, 0, render.Style{
				Fill:        cellColor,
				FillOpacity: 1.0,
			})

			if h.style.CellLabels {
				lum := luminance(cellColor)
				var textColor color.Color
				if lum > 0.35 {
					textColor = color.RGB(30, 30, 30)
				} else {
					textColor = color.RGB(240, 240, 240)
				}
				label := formatCellValue(v)
				canvas.DrawText(centerX, centerY, label, render.TextStyle{
					Color:      textColor,
					FontSize:   math.Max(8, math.Min(cellW, cellH)*0.28),
					FontFamily: "system-ui, sans-serif",
					Anchor:     "middle",
					Baseline:   "middle",
				})
			}
		}
	}
}

// formatCellValue formats a cell value for display.
func formatCellValue(v float64) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e6 {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

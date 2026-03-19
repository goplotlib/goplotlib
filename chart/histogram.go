package chart

import (
	"math"

	"github.com/goplotlib/goplotlib/color"
	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/scale"
)

// HistogramChart aggregates continuous values into bins and renders them as adjacent bars.
type HistogramChart struct {
	edges  []float64 // bin edges, len = numBins + 1
	counts []float64 // bar heights (counts, densities, or cumulative)
	style  HistogramStyle
	col    color.Color
}

// NewHistogram creates a HistogramChart by binning values at construction time.
func NewHistogram(values []float64, style HistogramStyle) *HistogramChart {
	if style.Opacity == 0 {
		style.Opacity = 1.0
	}

	edges := computeEdges(values, style)
	counts := bin(values, edges)

	if style.Normalize {
		counts = normalize(counts, edges)
	}
	if style.Cumulative {
		counts = cumulate(counts)
	}

	return &HistogramChart{edges: edges, counts: counts, style: style, col: style.Color}
}

func (h *HistogramChart) Label() string          { return h.style.Label }
func (h *HistogramChart) Color() color.Color     { return h.col }
func (h *HistogramChart) SetColor(c color.Color) { h.col = c }

// DataRange returns the full x extent of the bin edges and y from 0 to max bar height.
func (h *HistogramChart) DataRange() (xMin, xMax, yMin, yMax float64) {
	if len(h.edges) < 2 {
		return 0, 1, 0, 1
	}
	xMin = h.edges[0]
	xMax = h.edges[len(h.edges)-1]
	yMin = 0
	yMax = 0
	for _, c := range h.counts {
		if c > yMax {
			yMax = c
		}
	}
	return
}

// Draw renders each bin as a filled rectangle spanning the full bin width.
func (h *HistogramChart) Draw(canvas render.Canvas, xScale, yScale scale.Scale) {
	if len(h.counts) == 0 {
		return
	}

	yZero := yScale.Map(0)
	fillColor := h.col
	if h.style.Opacity < 1.0 && h.style.Opacity > 0 {
		fillColor = fillColor.WithAlpha(h.style.Opacity)
	}

	for i, count := range h.counts {
		x0 := xScale.Map(h.edges[i])
		x1 := xScale.Map(h.edges[i+1])
		barW := x1 - x0
		if barW < 1 {
			barW = 1
		}
		barY := yScale.Map(count)
		barH := math.Abs(yZero - barY)
		if barH < 1 {
			barH = 1
		}

		canvas.DrawRect(x0, barY, barW, barH, 0, 0, render.Style{
			Fill:        fillColor,
			FillOpacity: 1.0,
		})
	}
}

// computeEdges returns bin edges from style or Sturges' rule.
func computeEdges(values []float64, s HistogramStyle) []float64 {
	if len(s.BinEdges) >= 2 {
		return s.BinEdges
	}

	if len(values) == 0 {
		return []float64{0, 1}
	}

	vMin, vMax := values[0], values[0]
	for _, v := range values[1:] {
		if v < vMin {
			vMin = v
		}
		if v > vMax {
			vMax = v
		}
	}

	k := s.Bins
	if k <= 0 {
		// Sturges' rule
		k = int(math.Ceil(math.Log2(float64(len(values))) + 1))
		if k < 1 {
			k = 1
		}
	}

	if vMin == vMax {
		vMin -= 0.5
		vMax += 0.5
	}

	edges := make([]float64, k+1)
	step := (vMax - vMin) / float64(k)
	for i := range edges {
		edges[i] = vMin + float64(i)*step
	}
	// Ensure the last edge is exactly vMax to avoid floating-point drift
	edges[k] = vMax
	return edges
}

// bin counts how many values fall into each bin.
func bin(values []float64, edges []float64) []float64 {
	k := len(edges) - 1
	counts := make([]float64, k)
	for _, v := range values {
		if v < edges[0] || v > edges[k] {
			continue
		}
		// Binary search for the bin index
		lo, hi := 0, k-1
		for lo < hi {
			mid := (lo + hi) / 2
			if v < edges[mid+1] {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		counts[lo]++
	}
	return counts
}

// normalize converts counts to density so that sum(count * binWidth) = 1.
func normalize(counts []float64, edges []float64) []float64 {
	total := 0.0
	for i, c := range counts {
		total += c * (edges[i+1] - edges[i])
	}
	if total == 0 {
		return counts
	}
	out := make([]float64, len(counts))
	for i, c := range counts {
		out[i] = c / total
	}
	return out
}

// cumulate converts counts to cumulative values.
func cumulate(counts []float64) []float64 {
	out := make([]float64, len(counts))
	sum := 0.0
	for i, c := range counts {
		sum += c
		out[i] = sum
	}
	return out
}

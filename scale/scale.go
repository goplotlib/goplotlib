package scale

import (
	"fmt"
	"math"
	"strconv"
)

// Scale maps data values to pixel coordinates.
type Scale interface {
	Map(value float64) float64
	Domain() [2]float64
	Ticks(n int) []Tick
}

// Factory builds a Scale from a data domain and pixel range.
// Pass the result of Log(base) to Axes.SetXScale / Axes.SetYScale.
type Factory func(domainMin, domainMax, rangeMin, rangeMax float64) Scale

// Log returns a Factory that creates a LogScale with the given base.
// The data domain must be strictly positive.
func Log(base float64) Factory {
	return func(domainMin, domainMax, rangeMin, rangeMax float64) Scale {
		return NewLog(base, domainMin, domainMax, rangeMin, rangeMax)
	}
}

// Tick represents a tick mark on an axis.
type Tick struct {
	Value float64
	Label string
	Pos   float64 // pixel position
}

// LinearScale maps a linear data domain to a pixel range.
type LinearScale struct {
	DomainMin, DomainMax float64
	RangeMin, RangeMax   float64
}

// NewLinear creates a new LinearScale.
func NewLinear(domainMin, domainMax, rangeMin, rangeMax float64) *LinearScale {
	return &LinearScale{
		DomainMin: domainMin,
		DomainMax: domainMax,
		RangeMin:  rangeMin,
		RangeMax:  rangeMax,
	}
}

// Map linearly interpolates a value from domain to range.
func (s *LinearScale) Map(v float64) float64 {
	domainSpan := s.DomainMax - s.DomainMin
	if domainSpan == 0 {
		return (s.RangeMin + s.RangeMax) / 2
	}
	t := (v - s.DomainMin) / domainSpan
	return s.RangeMin + t*(s.RangeMax-s.RangeMin)
}

// Domain returns the [min, max] data domain.
func (s *LinearScale) Domain() [2]float64 {
	return [2]float64{s.DomainMin, s.DomainMax}
}

// Ticks generates approximately n nicely-spaced ticks using the "nice numbers" algorithm.
func (s *LinearScale) Ticks(n int) []Tick {
	if n <= 0 {
		n = 5
	}

	dataRange := s.DomainMax - s.DomainMin
	if dataRange == 0 {
		v := s.DomainMin
		return []Tick{{Value: v, Label: formatTickLabel(v, 1), Pos: s.Map(v)}}
	}

	magnitude := math.Pow(10, math.Floor(math.Log10(dataRange/float64(n))))

	// Try candidate step sizes
	candidates := []float64{1, 2, 2.5, 5, 10}
	step := candidates[0] * magnitude
	targetStep := dataRange / float64(n)
	minDiff := math.Abs(step - targetStep)

	for _, c := range candidates[1:] {
		candidate := c * magnitude
		diff := math.Abs(candidate - targetStep)
		if diff < minDiff {
			minDiff = diff
			step = candidate
		}
	}

	if step == 0 {
		step = 1
	}

	first := math.Ceil(s.DomainMin/step) * step
	var ticks []Tick
	for v := first; v <= s.DomainMax+step*1e-9; v += step {
		// Snap to zero if very close
		if math.Abs(v) < step*1e-9 {
			v = 0
		}
		if v > s.DomainMax+step*1e-9 {
			break
		}
		ticks = append(ticks, Tick{
			Value: v,
			Label: formatTickLabel(v, step),
			Pos:   s.Map(v),
		})
	}

	return ticks
}

// formatTickLabel formats a tick value as a string.
func formatTickLabel(v, step float64) string {
	if step >= 1 {
		// Handle large numbers
		abs := math.Abs(v)
		if abs >= 1e9 {
			return fmt.Sprintf("%.1fB", v/1e9)
		}
		if abs >= 1e6 {
			return fmt.Sprintf("%.1fM", v/1e6)
		}
		if abs >= 1e3 {
			return fmt.Sprintf("%.1fk", v/1e3)
		}
		return strconv.FormatFloat(v, 'f', 0, 64)
	}

	// Determine decimal places based on step size
	decimals := int(math.Ceil(-math.Log10(step)))
	decimals = min(max(decimals, 0), 10)
	return strconv.FormatFloat(v, 'f', decimals, 64)
}

// CategoricalScale maps category indices to pixel positions.
type CategoricalScale struct {
	Categories []string
	RangeMin   float64
	RangeMax   float64
	Padding    float64 // fraction of band width as padding, e.g. 0.2
}

// NewCategorical creates a new CategoricalScale.
func NewCategorical(categories []string, rangeMin, rangeMax float64, padding float64) *CategoricalScale {
	return &CategoricalScale{
		Categories: categories,
		RangeMin:   rangeMin,
		RangeMax:   rangeMax,
		Padding:    padding,
	}
}

// BandWidth returns the width of each band (including padding).
func (s *CategoricalScale) BandWidth() float64 {
	n := len(s.Categories)
	if n == 0 {
		return 0
	}
	totalWidth := s.RangeMax - s.RangeMin
	return totalWidth / float64(n)
}

// Map returns the center pixel position of the band at the given index.
func (s *CategoricalScale) Map(index float64) float64 {
	n := len(s.Categories)
	if n == 0 {
		return s.RangeMin
	}
	bandW := s.BandWidth()
	// Center of band at index
	return s.RangeMin + (index+0.5)*bandW
}

// Domain returns [0, len(categories)-1].
func (s *CategoricalScale) Domain() [2]float64 {
	n := len(s.Categories)
	if n == 0 {
		return [2]float64{0, 0}
	}
	return [2]float64{0, float64(n - 1)}
}

// Ticks returns one tick per category, positioned at the center of each band.
func (s *CategoricalScale) Ticks(n int) []Tick {
	ticks := make([]Tick, len(s.Categories))
	for i, cat := range s.Categories {
		ticks[i] = Tick{
			Value: float64(i),
			Label: cat,
			Pos:   s.Map(float64(i)),
		}
	}
	return ticks
}

package scale

import (
	"math"
	"sort"
	"strconv"
)

// LogScale maps a strictly-positive data domain to a pixel range using logarithm.
type LogScale struct {
	Base                 float64
	DomainMin, DomainMax float64
	RangeMin, RangeMax   float64
	logBase              float64 // cached math.Log(Base)
}

// NewLog creates a LogScale. base must be > 1 and the domain must be strictly positive.
func NewLog(base, domainMin, domainMax, rangeMin, rangeMax float64) *LogScale {
	if base <= 1 {
		base = 10
	}
	return &LogScale{
		Base:      base,
		DomainMin: domainMin,
		DomainMax: domainMax,
		RangeMin:  rangeMin,
		RangeMax:  rangeMax,
		logBase:   math.Log(base),
	}
}

func (s *LogScale) logOf(v float64) float64 {
	return math.Log(v) / s.logBase
}

// Map converts a data value to a pixel coordinate via logarithmic interpolation.
// Values ≤ 0 are clamped to DomainMin.
func (s *LogScale) Map(v float64) float64 {
	if v <= 0 {
		v = s.DomainMin
	}
	logMin := s.logOf(s.DomainMin)
	logMax := s.logOf(s.DomainMax)
	span := logMax - logMin
	if span == 0 {
		return (s.RangeMin + s.RangeMax) / 2
	}
	t := (s.logOf(v) - logMin) / span
	return s.RangeMin + t*(s.RangeMax-s.RangeMin)
}

// Domain returns [DomainMin, DomainMax].
func (s *LogScale) Domain() [2]float64 {
	return [2]float64{s.DomainMin, s.DomainMax}
}

// Ticks generates tick marks at powers of the base. When the range spans less
// than two decades, intermediate 1/2/5 multiples are added instead.
func (s *LogScale) Ticks(n int) []Tick {
	if n <= 0 {
		n = 6
	}
	logMin := s.logOf(s.DomainMin)
	logMax := s.logOf(s.DomainMax)
	firstExp := int(math.Ceil(logMin - 1e-9))
	lastExp := int(math.Floor(logMax + 1e-9))
	numPowers := lastExp - firstExp + 1

	if numPowers >= 2 {
		// Multi-decade: one tick per power of base, subsampled if too many.
		step := 1
		if numPowers > n {
			step = (numPowers + n - 1) / n
		}
		var ticks []Tick
		for e := firstExp; e <= lastExp; e += step {
			v := math.Pow(s.Base, float64(e))
			if v >= s.DomainMin*(1-1e-9) && v <= s.DomainMax*(1+1e-9) {
				ticks = append(ticks, Tick{Value: v, Label: formatLogLabel(v), Pos: s.Map(v)})
			}
		}
		return ticks
	}

	// Sub-decade: ticks at 1×, 2×, 5× multiples of surrounding powers.
	startExp := int(math.Floor(logMin)) - 1
	endExp := int(math.Ceil(logMax)) + 1
	seen := map[int64]bool{}
	var ticks []Tick
	for e := startExp; e <= endExp; e++ {
		baseV := math.Pow(s.Base, float64(e))
		for _, mult := range []float64{1, 2, 5} {
			v := baseV * mult
			key := int64(math.Round(v * 1e6))
			if seen[key] {
				continue
			}
			if v >= s.DomainMin*(1-1e-9) && v <= s.DomainMax*(1+1e-9) {
				seen[key] = true
				ticks = append(ticks, Tick{Value: v, Label: formatLogLabel(v), Pos: s.Map(v)})
			}
		}
	}
	sort.Slice(ticks, func(i, j int) bool { return ticks[i].Value < ticks[j].Value })
	return ticks
}

// formatLogLabel formats a log tick value with SI suffixes for large values.
func formatLogLabel(v float64) string {
	abs := math.Abs(v)
	switch {
	case abs >= 1e9:
		return strconv.FormatFloat(v/1e9, 'g', 4, 64) + "B"
	case abs >= 1e6:
		return strconv.FormatFloat(v/1e6, 'g', 4, 64) + "M"
	case abs >= 1e3:
		return strconv.FormatFloat(v/1e3, 'g', 4, 64) + "k"
	default:
		return strconv.FormatFloat(v, 'g', 4, 64)
	}
}

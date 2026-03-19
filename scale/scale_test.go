package scale

import (
	"math"
	"testing"
)

const eps = 1e-9

func approxEq(a, b float64) bool {
	return math.Abs(a-b) <= eps
}

// --- LinearScale ---

func TestLinearScaleMapBoundaries(t *testing.T) {
	s := NewLinear(0, 10, 0, 100)

	if got := s.Map(0); !approxEq(got, 0) {
		t.Errorf("Map(0) = %v, want 0", got)
	}
	if got := s.Map(10); !approxEq(got, 100) {
		t.Errorf("Map(10) = %v, want 100", got)
	}
	if got := s.Map(5); !approxEq(got, 50) {
		t.Errorf("Map(5) = %v, want 50", got)
	}
}

func TestLinearScaleMapMidpoint(t *testing.T) {
	s := NewLinear(0, 100, 200, 400)
	if got := s.Map(50); !approxEq(got, 300) {
		t.Errorf("Map(50) = %v, want 300", got)
	}
}

func TestLinearScaleMapInvertedRange(t *testing.T) {
	// Inverted y-axis: data max maps to top (small pixel y), data min maps to bottom (large pixel y)
	s := NewLinear(0, 10, 500, 0)
	if got := s.Map(0); !approxEq(got, 500) {
		t.Errorf("Map(0) = %v, want 500", got)
	}
	if got := s.Map(10); !approxEq(got, 0) {
		t.Errorf("Map(10) = %v, want 0", got)
	}
	if got := s.Map(5); !approxEq(got, 250) {
		t.Errorf("Map(5) = %v, want 250", got)
	}
}

func TestLinearScaleMapOutsideDomain(t *testing.T) {
	// Map is a linear extrapolation — no clamping
	s := NewLinear(0, 10, 0, 100)
	if got := s.Map(-5); !approxEq(got, -50) {
		t.Errorf("Map(-5) = %v, want -50", got)
	}
	if got := s.Map(20); !approxEq(got, 200) {
		t.Errorf("Map(20) = %v, want 200", got)
	}
}

func TestLinearScaleMapZeroDomain(t *testing.T) {
	// Zero-span domain returns midpoint of range
	s := NewLinear(5, 5, 0, 100)
	if got := s.Map(5); !approxEq(got, 50) {
		t.Errorf("Map with zero domain = %v, want midpoint 50", got)
	}
}

func TestLinearScaleMapNegativeDomain(t *testing.T) {
	s := NewLinear(-10, 10, 0, 200)
	if got := s.Map(0); !approxEq(got, 100) {
		t.Errorf("Map(0) = %v, want 100", got)
	}
	if got := s.Map(-10); !approxEq(got, 0) {
		t.Errorf("Map(-10) = %v, want 0", got)
	}
}

func TestLinearScaleDomain(t *testing.T) {
	s := NewLinear(2, 8, 0, 100)
	d := s.Domain()
	if d[0] != 2 || d[1] != 8 {
		t.Errorf("Domain() = %v, want [2 8]", d)
	}
}

// TestLinearScaleTicksNiceNumbers verifies ticks land on round values.
func TestLinearScaleTicksNiceNumbers(t *testing.T) {
	s := NewLinear(0.3, 9.7, 0, 500)
	ticks := s.Ticks(5)

	if len(ticks) == 0 {
		t.Fatal("Ticks() returned no ticks")
	}

	// All tick values should be round integers, not 0.3, 2.18, etc.
	for _, tick := range ticks {
		if tick.Value != math.Round(tick.Value) {
			t.Errorf("tick value %v is not a round integer", tick.Value)
		}
	}

	// First tick should be >= domain min
	if ticks[0].Value < 0.3 {
		t.Errorf("first tick %v is below domain min 0.3", ticks[0].Value)
	}
	// Last tick should be <= domain max
	if ticks[len(ticks)-1].Value > 9.7 {
		t.Errorf("last tick %v is above domain max 9.7", ticks[len(ticks)-1].Value)
	}
}

func TestLinearScaleTicksCrossZero(t *testing.T) {
	s := NewLinear(-5, 5, 0, 400)
	ticks := s.Ticks(5)

	hasZero := false
	for _, tick := range ticks {
		if tick.Value == 0 {
			hasZero = true
		}
	}
	if !hasZero {
		t.Error("expected a tick at 0 when domain crosses zero")
	}
}

func TestLinearScaleTicksPositions(t *testing.T) {
	s := NewLinear(0, 10, 0, 100)
	ticks := s.Ticks(5)

	for _, tick := range ticks {
		expected := s.Map(tick.Value)
		if !approxEq(tick.Pos, expected) {
			t.Errorf("tick at value %v has Pos=%v, want %v", tick.Value, tick.Pos, expected)
		}
	}
}

func TestLinearScaleTicksDecimalRange(t *testing.T) {
	s := NewLinear(0, 1, 0, 500)
	ticks := s.Ticks(5)

	if len(ticks) == 0 {
		t.Fatal("Ticks() returned no ticks")
	}

	// Labels should have decimal places, not be empty or integer-only
	for _, tick := range ticks {
		if tick.Label == "" {
			t.Errorf("tick at %v has empty label", tick.Value)
		}
	}
}

func TestLinearScaleTicksLargeNumbers(t *testing.T) {
	s := NewLinear(0, 5_000_000, 0, 500)
	ticks := s.Ticks(5)

	// Labels should use "k" or "M" suffix for readability
	for _, tick := range ticks {
		if tick.Value >= 1_000_000 {
			found := false
			for _, ch := range tick.Label {
				if ch == 'M' {
					found = true
				}
			}
			if !found {
				t.Errorf("tick at %v label %q should contain 'M' suffix", tick.Value, tick.Label)
			}
		}
	}
}

// --- CategoricalScale ---

func TestCategoricalScaleMap(t *testing.T) {
	cats := []string{"A", "B", "C", "D"}
	s := NewCategorical(cats, 0, 400, 0.2)

	bandW := s.BandWidth() // 400/4 = 100

	// Each category should map to the center of its band
	for i, cat := range cats {
		_ = cat
		got := s.Map(float64(i))
		want := float64(i)*bandW + bandW/2
		if !approxEq(got, want) {
			t.Errorf("Map(%d) = %v, want %v", i, got, want)
		}
	}
}

func TestCategoricalScaleBandWidth(t *testing.T) {
	s := NewCategorical([]string{"X", "Y", "Z"}, 0, 300, 0.2)
	if got := s.BandWidth(); !approxEq(got, 100) {
		t.Errorf("BandWidth() = %v, want 100", got)
	}
}

func TestCategoricalScaleBandWidthEmpty(t *testing.T) {
	s := NewCategorical(nil, 0, 300, 0.2)
	if got := s.BandWidth(); got != 0 {
		t.Errorf("BandWidth() with no categories = %v, want 0", got)
	}
}

func TestCategoricalScaleTicks(t *testing.T) {
	cats := []string{"Jan", "Feb", "Mar"}
	s := NewCategorical(cats, 0, 300, 0.2)
	ticks := s.Ticks(0) // n is ignored for categorical

	if len(ticks) != len(cats) {
		t.Fatalf("Ticks() returned %d ticks, want %d", len(ticks), len(cats))
	}
	for i, tick := range ticks {
		if tick.Label != cats[i] {
			t.Errorf("ticks[%d].Label = %q, want %q", i, tick.Label, cats[i])
		}
		if !approxEq(tick.Pos, s.Map(float64(i))) {
			t.Errorf("ticks[%d].Pos = %v, want %v", i, tick.Pos, s.Map(float64(i)))
		}
	}
}

func TestCategoricalScaleDomain(t *testing.T) {
	s := NewCategorical([]string{"A", "B", "C", "D", "E"}, 0, 500, 0)
	d := s.Domain()
	if d[0] != 0 || d[1] != 4 {
		t.Errorf("Domain() = %v, want [0 4]", d)
	}
}

func TestCategoricalScaleDomainEmpty(t *testing.T) {
	s := NewCategorical(nil, 0, 500, 0)
	d := s.Domain()
	if d[0] != 0 || d[1] != 0 {
		t.Errorf("Domain() with empty categories = %v, want [0 0]", d)
	}
}

// --- LogScale ---

func TestLogScaleMapBoundaries(t *testing.T) {
	s := NewLog(10, 1, 1000, 0, 300)

	if got := s.Map(1); !approxEq(got, 0) {
		t.Errorf("Map(1) = %v, want 0", got)
	}
	if got := s.Map(1000); !approxEq(got, 300) {
		t.Errorf("Map(1000) = %v, want 300", got)
	}
	// log10(10) = 1, midpoint of [0,3] in log space → pixel 100
	if got := s.Map(10); !approxEq(got, 100) {
		t.Errorf("Map(10) = %v, want 100", got)
	}
}

func TestLogScaleMapInverted(t *testing.T) {
	// Inverted range: data min → bottom (large pixel), data max → top (small pixel)
	s := NewLog(10, 1, 100, 400, 0)

	if got := s.Map(1); !approxEq(got, 400) {
		t.Errorf("Map(1) = %v, want 400", got)
	}
	if got := s.Map(100); !approxEq(got, 0) {
		t.Errorf("Map(100) = %v, want 0", got)
	}
	if got := s.Map(10); !approxEq(got, 200) {
		t.Errorf("Map(10) = %v, want 200 (midpoint)", got)
	}
}

func TestLogScaleMapNonPositiveClamp(t *testing.T) {
	s := NewLog(10, 1, 100, 0, 200)
	// Values ≤ 0 should clamp to domain min, not NaN or panic.
	got := s.Map(0)
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Errorf("Map(0) = %v, expected finite value (clamped to domain min)", got)
	}
}

func TestLogScaleDomain(t *testing.T) {
	s := NewLog(10, 1, 10000, 0, 500)
	d := s.Domain()
	if d[0] != 1 || d[1] != 10000 {
		t.Errorf("Domain() = %v, want [1 10000]", d)
	}
}

func TestLogScaleTicksMultiDecade(t *testing.T) {
	// [1, 10000] spans 4 decades → expect ticks at 1, 10, 100, 1000, 10000
	s := NewLog(10, 1, 10000, 0, 400)
	ticks := s.Ticks(6)

	if len(ticks) == 0 {
		t.Fatal("Ticks() returned no ticks")
	}
	// All tick values should be powers of 10
	for _, tick := range ticks {
		log := math.Log10(tick.Value)
		if math.Abs(log-math.Round(log)) > 1e-6 {
			t.Errorf("tick value %v is not a power of 10", tick.Value)
		}
		if tick.Pos != s.Map(tick.Value) {
			t.Errorf("tick at %v has wrong Pos %v", tick.Value, tick.Pos)
		}
	}
}

func TestLogScaleTicksSubDecade(t *testing.T) {
	// [1, 9] is sub-decade → should get intermediate ticks (1, 2, 5 style)
	s := NewLog(10, 1, 9, 0, 300)
	ticks := s.Ticks(6)

	if len(ticks) < 2 {
		t.Fatalf("Ticks() returned %d ticks, want at least 2", len(ticks))
	}
	// All tick positions should be within range
	for _, tick := range ticks {
		if tick.Pos < 0 || tick.Pos > 300 {
			t.Errorf("tick at %v has Pos %v outside [0, 300]", tick.Value, tick.Pos)
		}
	}
}

func TestLogScaleTicksLabels(t *testing.T) {
	s := NewLog(10, 1, 1e9, 0, 500)
	ticks := s.Ticks(10)

	// Labels for large values should use SI suffixes
	for _, tick := range ticks {
		if tick.Label == "" {
			t.Errorf("tick at %v has empty label", tick.Value)
		}
		if tick.Value >= 1e9 && tick.Label[len(tick.Label)-1] != 'B' {
			t.Errorf("tick at %v label %q should end with B", tick.Value, tick.Label)
		}
		if tick.Value >= 1e6 && tick.Value < 1e9 && tick.Label[len(tick.Label)-1] != 'M' {
			t.Errorf("tick at %v label %q should end with M", tick.Value, tick.Label)
		}
	}
}

func TestLogScaleBase2(t *testing.T) {
	s := NewLog(2, 1, 1024, 0, 100)
	ticks := s.Ticks(11)

	if len(ticks) == 0 {
		t.Fatal("Ticks() returned no ticks for base-2 scale")
	}
	// All ticks should be powers of 2
	for _, tick := range ticks {
		log2 := math.Log2(tick.Value)
		if math.Abs(log2-math.Round(log2)) > 1e-6 {
			t.Errorf("tick value %v is not a power of 2", tick.Value)
		}
	}
}

func TestLogFactory(t *testing.T) {
	f := Log(10)
	s := f(1, 100, 0, 200)
	ls, ok := s.(*LogScale)
	if !ok {
		t.Fatal("Log factory did not return *LogScale")
	}
	if ls.Base != 10 {
		t.Errorf("Base = %v, want 10", ls.Base)
	}
	if got := s.Map(10); !approxEq(got, 100) {
		t.Errorf("Map(10) = %v, want 100 (midpoint of [1,100] in log space)", got)
	}
}

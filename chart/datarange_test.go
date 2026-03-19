package chart

import (
	"testing"
)

// --- LineChart.DataRange ---

func TestLineChartDataRangeNormal(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5}
	ys := []float64{3, 1, 4, 1, 5}
	c := NewLine(xs, ys, LineStyle{})

	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin != 1 || xMax != 5 {
		t.Errorf("x range = [%v, %v], want [1, 5]", xMin, xMax)
	}
	if yMin != 1 || yMax != 5 {
		t.Errorf("y range = [%v, %v], want [1, 5]", yMin, yMax)
	}
}

func TestLineChartDataRangeNegativeValues(t *testing.T) {
	xs := []float64{-5, 0, 5}
	ys := []float64{-3, 0, 3}
	c := NewLine(xs, ys, LineStyle{})

	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin != -5 || xMax != 5 {
		t.Errorf("x range = [%v, %v], want [-5, 5]", xMin, xMax)
	}
	if yMin != -3 || yMax != 3 {
		t.Errorf("y range = [%v, %v], want [-3, 3]", yMin, yMax)
	}
}

func TestLineChartDataRangeSinglePoint(t *testing.T) {
	xs := []float64{7}
	ys := []float64{4}
	c := NewLine(xs, ys, LineStyle{})

	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin != 7 || xMax != 7 {
		t.Errorf("x range = [%v, %v], want [7, 7]", xMin, xMax)
	}
	if yMin != 4 || yMax != 4 {
		t.Errorf("y range = [%v, %v], want [4, 4]", yMin, yMax)
	}
}

func TestLineChartDataRangeEmpty(t *testing.T) {
	c := NewLine(nil, nil, LineStyle{})
	xMin, xMax, yMin, yMax := c.DataRange()

	// Empty data returns a safe default range
	if xMin > xMax {
		t.Errorf("empty x range [%v, %v] is inverted", xMin, xMax)
	}
	if yMin > yMax {
		t.Errorf("empty y range [%v, %v] is inverted", yMin, yMax)
	}
}

// --- BarChart.DataRange ---

func TestBarChartDataRangeIncludesZero(t *testing.T) {
	cats := []string{"A", "B", "C"}

	// All positive values: yMin must still be 0
	c := NewBar(cats, []float64{5, 10, 3}, BarStyle{})
	_, _, yMin, _ := c.DataRange()
	if yMin > 0 {
		t.Errorf("BarChart yMin = %v for positive values, want <= 0", yMin)
	}

	// All negative values: yMax must still be 0
	c2 := NewBar(cats, []float64{-5, -10, -3}, BarStyle{})
	_, _, _, yMax := c2.DataRange()
	if yMax < 0 {
		t.Errorf("BarChart yMax = %v for negative values, want >= 0", yMax)
	}
}

func TestBarChartDataRangeXIndices(t *testing.T) {
	cats := []string{"A", "B", "C", "D"}
	c := NewBar(cats, []float64{1, 2, 3, 4}, BarStyle{})
	xMin, xMax, _, _ := c.DataRange()

	if xMin != 0 {
		t.Errorf("xMin = %v, want 0", xMin)
	}
	if xMax != float64(len(cats)-1) {
		t.Errorf("xMax = %v, want %v", xMax, float64(len(cats)-1))
	}
}

func TestBarChartDataRangeMixed(t *testing.T) {
	cats := []string{"A", "B", "C"}
	c := NewBar(cats, []float64{-3, 5, -1}, BarStyle{})
	_, _, yMin, yMax := c.DataRange()

	if yMin > -3 {
		t.Errorf("yMin = %v, want <= -3", yMin)
	}
	if yMax < 5 {
		t.Errorf("yMax = %v, want >= 5", yMax)
	}
}

func TestBarChartDataRangeEmpty(t *testing.T) {
	c := NewBar(nil, nil, BarStyle{})
	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin > xMax {
		t.Errorf("empty x range [%v, %v] is inverted", xMin, xMax)
	}
	if yMin > yMax {
		t.Errorf("empty y range [%v, %v] is inverted", yMin, yMax)
	}
}

// --- ScatterChart.DataRange ---

func TestScatterChartDataRangeNormal(t *testing.T) {
	xs := []float64{2, 8, 5, 1, 9}
	ys := []float64{4, 2, 7, 3, 1}
	c := NewScatter(xs, ys, ScatterStyle{})

	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin != 1 || xMax != 9 {
		t.Errorf("x range = [%v, %v], want [1, 9]", xMin, xMax)
	}
	if yMin != 1 || yMax != 7 {
		t.Errorf("y range = [%v, %v], want [1, 7]", yMin, yMax)
	}
}

func TestScatterChartDataRangeNegative(t *testing.T) {
	xs := []float64{-10, 0, 10}
	ys := []float64{-5, 5, -5}
	c := NewScatter(xs, ys, ScatterStyle{})

	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin != -10 || xMax != 10 {
		t.Errorf("x range = [%v, %v], want [-10, 10]", xMin, xMax)
	}
	if yMin != -5 || yMax != 5 {
		t.Errorf("y range = [%v, %v], want [-5, 5]", yMin, yMax)
	}
}

func TestScatterChartDataRangeEmpty(t *testing.T) {
	c := NewScatter(nil, nil, ScatterStyle{})
	xMin, xMax, yMin, yMax := c.DataRange()

	if xMin > xMax {
		t.Errorf("empty x range [%v, %v] is inverted", xMin, xMax)
	}
	if yMin > yMax {
		t.Errorf("empty y range [%v, %v] is inverted", yMin, yMax)
	}
}

package spline

import (
	"math"
	"testing"
)

const eps = 1e-9

func approxEq(a, b float64) bool {
	return math.Abs(a-b) <= eps
}

func pointEq(a, b Point) bool {
	return approxEq(a.X, b.X) && approxEq(a.Y, b.Y)
}

// TestCatmullRomSegmentCount verifies the number of segments returned.
func TestCatmullRomSegmentCount(t *testing.T) {
	tests := []struct {
		pts  []Point
		want int
	}{
		{nil, 0},
		{[]Point{{0, 0}}, 0},
		{[]Point{{0, 0}, {1, 1}}, 1},
		{[]Point{{0, 0}, {1, 1}, {2, 0}}, 2},
		{[]Point{{0, 0}, {1, 1}, {2, 0}, {3, 1}}, 3},
	}
	for _, tt := range tests {
		got := CatmullRomToBezier(tt.pts)
		if len(got) != tt.want {
			t.Errorf("CatmullRomToBezier(%d pts) returned %d segments, want %d",
				len(tt.pts), len(got), tt.want)
		}
	}
}

// TestCatmullRomSegmentEndpoints verifies that segments chain correctly:
// the end of segment[i] equals the start of segment[i+1].
func TestCatmullRomSegmentEndpoints(t *testing.T) {
	pts := []Point{{0, 0}, {1, 2}, {3, 1}, {4, 3}, {6, 0}}
	segs := CatmullRomToBezier(pts)

	for i := 0; i < len(segs)-1; i++ {
		end := segs[i][3]
		next := segs[i+1][0]
		if !pointEq(end, next) {
			t.Errorf("seg[%d].end=%v != seg[%d].start=%v", i, end, i+1, next)
		}
	}
}

// TestCatmullRomPassesThroughPoints verifies that segment start/end equal the input points.
func TestCatmullRomPassesThroughPoints(t *testing.T) {
	pts := []Point{{0, 0}, {1, 2}, {3, 1}, {4, 3}}
	segs := CatmullRomToBezier(pts)

	for i, seg := range segs {
		if !pointEq(seg[0], pts[i]) {
			t.Errorf("seg[%d].start=%v, want pts[%d]=%v", i, seg[0], i, pts[i])
		}
		if !pointEq(seg[3], pts[i+1]) {
			t.Errorf("seg[%d].end=%v, want pts[%d]=%v", i, seg[3], i+1, pts[i+1])
		}
	}
}

// TestCatmullRomCollinear verifies that collinear points produce straight-line Bézier segments:
// the control points should lie on the line between start and end.
func TestCatmullRomCollinear(t *testing.T) {
	// All points on y = x
	pts := []Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}}
	segs := CatmullRomToBezier(pts)

	for i, seg := range segs {
		// For a straight line, cp1 and cp2 should also be on y=x
		if !approxEq(seg[1].X, seg[1].Y) {
			t.Errorf("seg[%d] cp1=%v is not on y=x line", i, seg[1])
		}
		if !approxEq(seg[2].X, seg[2].Y) {
			t.Errorf("seg[%d] cp2=%v is not on y=x line", i, seg[2])
		}
	}
}

// TestCatmullRomFirstSegmentPhantomPoint verifies that for the first segment,
// the phantom prev point is pts[0] itself, which means cp1 = pts[0] + (pts[1]-pts[0])/6.
func TestCatmullRomFirstSegmentCP1(t *testing.T) {
	pts := []Point{{0, 0}, {6, 0}, {12, 0}}
	segs := CatmullRomToBezier(pts)

	// First segment: prev = pts[0] (phantom), so cp1 = pts[0] + (pts[1]-pts[0])/6
	// = {0,0} + {6,0}/6 = {1, 0}
	wantCP1 := Point{1, 0}
	if !pointEq(segs[0][1], wantCP1) {
		t.Errorf("first segment cp1 = %v, want %v", segs[0][1], wantCP1)
	}
}

// TestCatmullRomLastSegmentPhantomPoint verifies that for the last segment,
// the phantom next point is pts[n-1] itself, which means cp2 = pts[n-1] - (pts[n-1]-pts[n-2])/6.
func TestCatmullRomLastSegmentCP2(t *testing.T) {
	pts := []Point{{0, 0}, {6, 0}, {12, 0}}
	segs := CatmullRomToBezier(pts)

	// Last segment: next = pts[2] (phantom), so cp2 = pts[2] - (pts[2]-pts[1])/6
	// = {12,0} - {6,0}/6 = {11, 0}
	wantCP2 := Point{11, 0}
	if !pointEq(segs[1][2], wantCP2) {
		t.Errorf("last segment cp2 = %v, want %v", segs[1][2], wantCP2)
	}
}

// TestCatmullRomTwoPoints verifies the minimum case (one segment, no interior points).
func TestCatmullRomTwoPoints(t *testing.T) {
	pts := []Point{{0, 0}, {10, 5}}
	segs := CatmullRomToBezier(pts)

	if len(segs) != 1 {
		t.Fatalf("expected 1 segment for 2 points, got %d", len(segs))
	}

	seg := segs[0]
	if !pointEq(seg[0], pts[0]) {
		t.Errorf("start = %v, want %v", seg[0], pts[0])
	}
	if !pointEq(seg[3], pts[1]) {
		t.Errorf("end = %v, want %v", seg[3], pts[1])
	}

	// With phantom points at both ends, both phantom = the only other point,
	// so cp1 = pts[0] + (pts[1]-pts[0])/6 and cp2 = pts[1] - (pts[1]-pts[0])/6
	wantCP1 := Point{pts[0].X + (pts[1].X-pts[0].X)/6, pts[0].Y + (pts[1].Y-pts[0].Y)/6}
	wantCP2 := Point{pts[1].X - (pts[1].X-pts[0].X)/6, pts[1].Y - (pts[1].Y-pts[0].Y)/6}
	if !pointEq(seg[1], wantCP1) {
		t.Errorf("cp1 = %v, want %v", seg[1], wantCP1)
	}
	if !pointEq(seg[2], wantCP2) {
		t.Errorf("cp2 = %v, want %v", seg[2], wantCP2)
	}
}

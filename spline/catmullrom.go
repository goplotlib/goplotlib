package spline

// Point is a 2D point.
type Point struct{ X, Y float64 }

// CatmullRomToBezier converts a sequence of points into cubic Bézier
// control points using Catmull-Rom parameterization.
// Returns a slice of segments, each segment being [start, cp1, cp2, end].
// Returns len(pts)-1 segments.
func CatmullRomToBezier(pts []Point) [][4]Point {
	n := len(pts)
	if n < 2 {
		return nil
	}

	segments := make([][4]Point, n-1)

	for i := 0; i < n-1; i++ {
		prev := pts[max0(i-1)]
		p0 := pts[i]
		p1 := pts[i+1]
		next := pts[min0(n-1, i+2)]

		// Catmull-Rom control points:
		// CP1 = P[i] + (P[i+1] - P_prev) / 6
		// CP2 = P[i+1] - (P_next - P[i]) / 6
		cp1 := Point{
			X: p0.X + (p1.X-prev.X)/6.0,
			Y: p0.Y + (p1.Y-prev.Y)/6.0,
		}
		cp2 := Point{
			X: p1.X - (next.X-p0.X)/6.0,
			Y: p1.Y - (next.Y-p0.Y)/6.0,
		}

		segments[i] = [4]Point{p0, cp1, cp2, p1}
	}

	return segments
}

func max0(a int) int {
	if a < 0 {
		return 0
	}
	return a
}

func min0(a, b int) int {
	if a < b {
		return a
	}
	return b
}

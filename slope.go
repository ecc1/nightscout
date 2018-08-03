package nightscout

// Points is the interface satisfied by a set of points.
type Points interface {
	Len() int
	X(int) float64
	Y(int) float64
}

// Line represents the linear equation y = Slope*x + Intercept.
type Line struct {
	Slope     float64
	Intercept float64
}

// FindLine performs simple linear regression on a set of points.
func FindLine(points Points) Line {
	n := points.Len()
	if n < 2 {
		panic("FindLine requires at least 2 points")
	}
	xSum, xSqSum, ySum, ySqSum, xySum := 0.0, 0.0, 0.0, 0.0, 0.0
	for i := 0; i < n; i++ {
		x := points.X(i)
		xSum += x
		xSqSum += x * x
		y := points.Y(i)
		ySum += y
		ySqSum += y * y
		xySum += x * y
	}
	s := float64(n)
	xBar := xSum / s
	yBar := ySum / s
	slope := (xySum - xSum*yBar) / (xSqSum - xSum*xBar)
	intercept := yBar - slope*xBar
	return Line{
		Slope:     slope,
		Intercept: intercept,
	}
}

// Eval evaluates a linear function at the given value.
func (l Line) Eval(x float64) float64 {
	return l.Slope*x + l.Intercept
}

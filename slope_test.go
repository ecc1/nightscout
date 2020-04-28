package nightscout

import (
	"math"
	"testing"
)

type IntPoints struct {
	xs []int
	ys []int
}

func (p IntPoints) Len() int {
	return len(p.xs)
}

func (p IntPoints) X(i int) float64 {
	return float64(p.xs[i])
}

func (p IntPoints) Y(i int) float64 {
	return float64(p.ys[i])
}

var (
	points0 = IntPoints{
		[]int{0, 10},
		[]int{50, 100},
	}

	points1 = IntPoints{
		[]int{147, 150, 152, 155, 157, 160, 163, 165, 168, 170, 173, 175, 178, 180, 183},
		[]int{5221, 5312, 5448, 5584, 5720, 5857, 5993, 6129, 6311, 6447, 6628, 6810, 6992, 7219, 7446},
	}

	points2 = IntPoints{
		[]int{17, 16, 28, 56, 13, 22, 13, 11, 32, 15, 52, 46, 58, 30},
		[]int{37, 39, 67, 95, 34, 56, 37, 27, 55, 29, 107, 76, 118, 41},
	}
)

const tolerance = 1e-12

func closeEnough(x, y float64) bool {
	return math.Abs(x-y) <= tolerance
}

func TestFindLine(t *testing.T) {
	cases := []struct {
		points IntPoints
		line   Line
	}{
		{points0, Line{5.0, 50.0}},
		{points1, Line{61.272186542110, -3906.195591884244}},
		{points2, Line{1.669862317066, 9.644736594278}},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			line := FindLine(c.points)
			if !closeEnough(line.Slope, c.line.Slope) {
				t.Errorf("FindLine: got Slope %v, want %v", line.Slope, c.line.Slope)
			}
			if !closeEnough(line.Intercept, c.line.Intercept) {
				t.Errorf("FindLine: got Intercept %v, want %v", line.Intercept, c.line.Intercept)
			}
		})
	}
}

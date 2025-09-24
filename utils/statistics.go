package utils

import (
	"math"
	"sort"

	"github.com/sajari/regression"
	"gonum.org/v1/plot/plotter"
)

func RemoveOutliers(data []plotter.XY, k float64) []plotter.XY {
	if len(data) == 0 {
		return nil
	}

	// Extract the JD values to analyze outliers on.
	jdValues := make([]float64, len(data))
	for i, dp := range data {
		jdValues[i] = dp.Y
	}

	// Sort JD values
	sort.Float64s(jdValues)

	// Calculate Q1 (25th percentile) and Q3 (75th percentile)
	q1 := percentile(jdValues, 25)
	q3 := percentile(jdValues, 75)
	iqr := q3 - q1

	lowerBound := q1 - k*iqr
	upperBound := q3 + k*iqr

	// Filter the dataset
	var filtered []plotter.XY
	for _, dp := range data {
		if dp.Y >= lowerBound && dp.Y <= upperBound {
			filtered = append(filtered, dp)
		}
	}

	return filtered
}

// Helper function to compute the percentile of a sorted slice
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	pos := (p / 100) * float64(len(sorted)-1)
	lower := int(pos)
	upper := lower + 1
	weight := pos - float64(lower)
	if upper >= len(sorted) {
		return sorted[lower]
	}
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

func PredictDeg(reg *regression.Regression, deg int, value float64) (prediction float64, err error) {
	x := make([]float64, deg)
	for d := 1; d <= deg; d++ {
		x[d-1] = math.Pow(value, float64(d))
	}

	prediction, err = reg.Predict(x)
	return
}

func FindRange(points []plotter.XY, dim int) (float64, float64) {
	lmin, lmax := math.MaxFloat64, -math.MaxFloat64
	for _, p := range points {
		val := p.X
		if dim == 1 {
			val = p.Y
		}
		if val < lmin {
			lmin = val
		}
		if val > lmax {
			lmax = val
		}
	}
	return lmin, lmax
}

// distance calculates the Euclidean distance between two points
func distance(p1, p2 plotter.XY) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Made by Claude 3.7 Sonnet

package utils

import (
	"fmt"
	"math"

	"github.com/sajari/regression"
	"gonum.org/v1/plot/plotter"
)

// KMeans performs k-means clustering on the data points
func KMeans(points []plotter.XY, k int, maxIterations int) [][]plotter.XY {
	// guessed initial centroids; to improve
	centroids := []plotter.XY{
		{16, 18.5}, // Lower curve
		{18, 14},   // Upper curve
	}

	// Initialize clusters
	clusters := make([][]plotter.XY, k)
	for i := range clusters {
		clusters[i] = []plotter.XY{}
	}

	// Iterate until convergence or maximum iterations
	for iter := 0; iter < maxIterations; iter++ {
		// Reset clusters
		for i := range clusters {
			clusters[i] = []plotter.XY{}
		}

		// Assign points to nearest centroid
		for _, point := range points {
			minDist := math.MaxFloat64
			clusterIdx := 0
			for i, centroid := range centroids {
				dist := distance(point, centroid)
				if dist < minDist {
					minDist = dist
					clusterIdx = i
				}
			}
			clusters[clusterIdx] = append(clusters[clusterIdx], point)
		}

		// Calculate new centroids
		oldCentroids := make([]plotter.XY, k)
		copy(oldCentroids, centroids)

		for i := range centroids {
			if len(clusters[i]) == 0 {
				continue // Skip empty clusters
			}

			sumX, sumY := 0.0, 0.0
			for _, point := range clusters[i] {
				sumX += point.X
				sumY += point.Y
			}
			centroids[i] = plotter.XY{
				X: sumX / float64(len(clusters[i])),
				Y: sumY / float64(len(clusters[i])),
			}
		}

		// Check for convergence
		converged := false
		for i := range centroids {
			if distance(centroids[i], oldCentroids[i]) > 0.001 {
				converged = true
				break
			}
		}
		if converged {
			break
		}
	}

	return clusters
}

// FitModels tries different regression models and returns the best one
func FitModels(points []plotter.XY) (*regression.Regression, float64) {
	// Try different polynomial degrees
	bestModel := &regression.Regression{}
	bestR2 := -1.0

	// Try polynomial regression with different degrees
	for degree := 1; degree <= 4; degree++ {
		r := &regression.Regression{}
		r.SetObserved("Jump Distance")

		// Add feature names for each polynomial term
		for i := 1; i <= degree; i++ {
			name := fmt.Sprintf("x^%d", i)
			r.SetVar(i-1, name)
		}

		// Add data points
		for _, p := range points {
			terms := make([]float64, degree)
			for i := 1; i <= degree; i++ {
				terms[i-1] = math.Pow(p.X, float64(i))
			}
			r.Train(regression.DataPoint(p.Y, terms))
		}

		// Fit the model
		err := r.Run()
		if err != nil {
			continue
		}

		// Get R-squared value
		r2 := r.R2
		if r2 > bestR2 {
			bestR2 = r2
			bestModel = r
		}
	}

	return bestModel, bestR2
}

// EvaluateModel predicts y values for a range of x values using the given model
func EvaluateModel(model *regression.Regression, minX, maxX float64, points int) ([]Point, error) {
	results := make([]Point, points)
	step := (maxX - minX) / float64(points-1)

	degree := len(model.GetCoeffs())
	for i := 0; i < points; i++ {
		x := minX + float64(i)*step
		terms := make([]float64, degree)
		for j := 1; j <= degree; j++ {
			terms[j-1] = math.Pow(x, float64(j))
		}
		y, err := model.Predict(terms)
		if err != nil {
			return nil, err
		}
		results[i] = Point{X: x, Y: y}
	}

	return results, nil
}

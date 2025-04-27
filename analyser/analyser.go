package analyser

import (
	"github.com/motzel/go-bsor/bsor"
	"math"
	"replayAnalyzer/utils"
)

func AverageSwingDistance(frames []*bsor.PositionAndRotation) float64 {
	if len(frames) < 2 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(frames); i++ {
		dx := float64(frames[i].Position.X - frames[i-1].Position.X)
		dy := float64(frames[i].Position.Y - frames[i-1].Position.Y)
		dz := float64(frames[i].Position.Z - frames[i-1].Position.Z)
		total += math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
	return total / float64(len(frames)-1)
}
func EstimateControllerType(frames []*bsor.PositionAndRotation) string {
	if len(frames) == 0 {
		return "Unknown"
	}
	avgY := 0.0
	for _, f := range frames {
		avgY += float64(f.Position.Y)
	}
	avgY /= float64(len(frames))

	if avgY > 1.0 {
		return "Valve Index" // High controller profile
	} else if avgY > 0.7 {
		return "Quest 2"
	}
	return "WMR/Other"
}
func EstimateGripStyle(leftFrames, rightFrames []*bsor.PositionAndRotation) string {
	// Very rough estimation based on hand distance
	if len(leftFrames) == 0 || len(rightFrames) == 0 {
		return "Unknown"
	}
	totalDistance := 0.0
	count := int(math.Min(float64(len(leftFrames)), float64(len(rightFrames))))
	for i := 0; i < count; i++ {
		l := leftFrames[i].Position
		r := rightFrames[i].Position
		dx := float64(l.X - r.X)
		dy := float64(l.Y - r.Y)
		dz := float64(l.Z - r.Z)
		totalDistance += math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
	avg := totalDistance / float64(count)
	if avg < 0.4 {
		return "Close Grip (Claw)"
	} else if avg < 0.6 {
		return "Neutral Grip"
	}
	return "Wide Grip"
}
func AnalyzeControllers(frames []bsor.Frame) utils.AnalysisResult {
	var left []*bsor.PositionAndRotation
	var right []*bsor.PositionAndRotation

	for _, frame := range frames {
		left = append(left, &frame.LeftHand)
		right = append(right, &frame.RightHand)
	}

	leftAvg := AverageSwingDistance(left)
	rightAvg := AverageSwingDistance(right)

	dominant := "Equal"
	if leftAvg > rightAvg*1.1 {
		dominant = "Left"
	} else if rightAvg > leftAvg*1.1 {
		dominant = "Right"
	}
	return utils.AnalysisResult{
		DominantHand:         dominant,
		AvgSwingDistanceL:    leftAvg,
		AvgSwingDistanceR:    rightAvg,
		SwingIntensityDiff:   math.Abs(leftAvg - rightAvg),
		EstimatedControllerL: EstimateControllerType(left),
		EstimatedControllerR: EstimateControllerType(right),
		EstimatedGripStyle:   EstimateGripStyle(left, right),
	}
}

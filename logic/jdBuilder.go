package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/motzel/go-bsor/bsor"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log/slog"
	"os"
	"replayAnalyzer/storage"
	"replayAnalyzer/utils"
	"sort"
)

func GenerateJDConfig(playerId string) error {
	slog.Info("Loading player's replays...")

	replays, err := storage.GetReplays(playerId)
	if err != nil {
		return err
	}

	slog.Info("Training jd prediction model...")

	var points plotter.XYs

	for _, fileName := range *replays {
		point, err := buildJDPoint(playerId, fileName)
		if err != nil {
			return err
		}
		points = append(points, *point)
	}

	if len(points) < 50 {
		slog.Info("WARNING: Please note that the reliability of the model increases with more training data and a wider range of maps. " +
			"Consider fetching more replays with different njs values.")
	}

	points = utils.RemoveOutliers(points, 1.5)

	// Grouping
	clusters := make([]utils.Cluster, 0)
	pointClusters := utils.KMeans(points, 2, 300)
	if len(pointClusters) > 3 {
		return errors.New("too many grouped pairs, this is odd")
	}

	for _, clusterPoints := range pointClusters {
		if len(clusterPoints) < 3 {
			continue
		}

		model, r2 := utils.FitModels(clusterPoints)
		fmt.Printf("Cluster with %d points - R²: %.4f\n", len(clusterPoints), r2)
		fmt.Printf("Model formula: %s", model.Formula)
		fmt.Println()

		clusters = append(clusters, utils.Cluster{
			Points: clusterPoints,
			Model:  model,
		})
	}
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Model.R2 > clusters[j].Model.R2
	})

	if len(clusters) > 2 {
		clusters = clusters[:2]
	}

	p := plot.New()
	p.Title.Text = "[NJS - JD] Cluster Regression Analysis"
	p.X.Label.Text = "Note Jump Speed"
	p.Y.Label.Text = "Jump Distance"

	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 0, 255, 255},   // Blue
		{0, 255, 0, 255},   // Green
		{255, 0, 255, 255}, // Purple
		{255, 165, 0, 255}, // Orange
	}

	for i, cluster := range clusters {
		// Scatter points for this cluster
		pts := make(plotter.XYs, len(cluster.Points))
		for j, p := range cluster.Points {
			pts[j].X = p.X
			pts[j].Y = p.Y
		}

		s, err := plotter.NewScatter(pts)
		if err != nil {
			panic(err)
		}
		s.GlyphStyle.Color = colors[i%len(colors)]
		s.GlyphStyle.Radius = vg.Points(3)
		s.GlyphStyle.Shape = draw.CircleGlyph{}
		p.Add(s)

		// Create regression curve for this cluster
		minX, maxX := utils.FindRange(cluster.Points, 0)
		curve, err := utils.EvaluateModel(cluster.Model, minX, maxX, 100)
		if err != nil {
			panic(err)
		}

		line := make(plotter.XYs, len(curve))
		for j, pt := range curve {
			line[j].X = pt.X
			line[j].Y = pt.Y
		}

		l, err := plotter.NewLine(line)
		if err != nil {
			panic(err)
		}
		l.LineStyle.Width = vg.Points(2)
		l.LineStyle.Color = colors[i%len(colors)]
		p.Add(l)

		// Add R² value to legend
		p.Legend.Add(fmt.Sprintf("Cluster %d (R² = %.4f)", i+1, cluster.Model.R2), l)
	}

	plotPath := "_cache/plots/" + playerId + ".jpg"
	err = p.Save(6*vg.Inch, 6*vg.Inch, plotPath)
	if err != nil {
		return err
	}
	utils.OpenFile(plotPath)

	for _, cluster := range clusters {
		bts, err := buildJDConfig(cluster.Model)
		if err != nil {
			return err
		}
		jdPath := fmt.Sprintf("_cache/jd_configs/%s_%s.json", playerId, utils.RandomStr(4))
		_ = os.WriteFile(jdPath, *bts, 0666)
		slog.Info(fmt.Sprintf("Check \"%s\" for generated jd config", jdPath))
	}
	return nil
}

func buildJDPoint(playerId, fileName string) (*plotter.XY, error) {
	file, err := os.OpenFile("_replays/"+playerId+"/"+fileName, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	replay, err := bsor.Read(file)
	if err != nil {
		slog.Info(fmt.Sprintf("Error reading replay file: %s/%s", playerId, fileName))
		return nil, err
	}

	diff := storage.GetDiff(replay.Info.Hash, replay.Info.Difficulty)
	if diff == nil {
		return nil, errors.New(fmt.Sprintf("can't find leaderboard for %s and %s", replay.Info.Hash, replay.Info.Difficulty))
	}

	pair := plotter.XY{
		X: diff.Njs,
		Y: float64(replay.Info.JumpDistance),
	}

	_ = file.Close()

	return &pair, nil
}
func buildJDConfig(reg *regression.Regression) (*[]byte, error) {
	var configPairs []utils.JDPair
	var njs = utils.JDConfigLow

	for njs < utils.JDConfigHigh {

		predictedJD, err := utils.PredictDeg(reg, len(reg.GetCoeffs()), njs)
		if err != nil {
			return nil, err
		}

		configPairs = append(configPairs, utils.JDPair{
			NJS: njs,
			JD:  predictedJD,
		})

		njs += 0.5
	}
	bytes, err := json.MarshalIndent(configPairs, "", "   ")
	if err != nil {
		return nil, err
	}

	return &bytes, nil
}

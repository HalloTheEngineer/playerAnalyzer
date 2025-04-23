package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/motzel/go-bsor/bsor"
	"github.com/sajari/regression"
	"log"
	"log/slog"
	"math"
	"net/http"
	"os"
	"replayAnalyzer/analyser"
	"replayAnalyzer/models"
	"replayAnalyzer/storage"
	"strconv"
	"strings"
)

func handleFetchCmd(ctx context.Context, args []string) (err error) {
	var playerId string

	if len(args) > 0 {
		playerId = args[0]
	} else {
		playerId, err = models.GetInput("Enter player id: ")
		if err != nil {
			return err
		}
	}

	if !models.RequireNumbers(playerId) {
		return errors.New("valid player id is required")
	}
	queryStr := models.BeatleaderApiUrl + "/player/" + playerId + "/scores"

	params, err := models.GetInput("Provide beatleader-api query params (format: ?param=value&param2=value):\n" + queryStr)
	if err != nil {
		return err
	}
	queryStr = queryStr + params

	slog.Info("Collecting scores for request: " + queryStr)

	resp, err := http.Get(queryStr)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 && resp.StatusCode < 200 {
		return errors.New("beatleader-api returned status code " + strconv.Itoa(resp.StatusCode))
	}

	jsonBuf, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		return err
	}

	var cdnUrls []string
	var leaderboards []models.ALeaderboard

	for _, score := range jsonBuf.S("data").Children() {
		var l models.ALeaderboard

		err = json.Unmarshal(score.S("leaderboard").Bytes(), &l)
		if err != nil {
			return err
		}

		cdnUrls = append(cdnUrls, strings.Trim(score.S("replay").String(), "\""))
		leaderboards = append(leaderboards, l)
	}

	_ = os.Mkdir("_replays", 0755)

	for _, cdnUrl := range cdnUrls {
		err := storage.FetchReplay(playerId, cdnUrl)
		if err != nil {
			return err
		}
	}

	for _, leaderboard := range leaderboards {
		err := storage.SaveLeaderboard(&leaderboard)
		if err != nil {
			return err
		}
	}

	slog.Info(fmt.Sprintf("Fetched %d leaderboards", len(leaderboards)))
	slog.Info(fmt.Sprintf("Fetched %d replays of %s", len(cdnUrls), playerId))

	return nil
}

func handleJDGenCmd(ctx context.Context, args []string) (err error) {
	var playerId string

	if len(args) > 0 {
		playerId = args[0]
	} else {
		playerId, err = models.GetInput("Enter player id: ")
		if err != nil {
			return err
		}
	}

	slog.Info("Loading player's replays...")

	replays, err := storage.GetReplays(playerId)
	if err != nil {
		return err
	}

	slog.Info("Training jd prediction model...")

	var pairs []models.JDPair

	for _, fileName := range *replays {
		file, err := os.OpenFile("_replays/"+playerId+"/"+fileName, os.O_RDONLY, 0755)
		if err != nil {
			return err
		}

		replay, err := bsor.Read(file)
		if err != nil {
			return err
		}

		diff := storage.GetDiff(replay.Info.Hash, replay.Info.Difficulty)
		if diff == nil {
			return errors.New(fmt.Sprintf("can't find leaderboard for %s and %s", replay.Info.Hash, replay.Info.Difficulty))
		}

		pairs = append(pairs, models.JDPair{
			NJS: diff.Njs,
			JD:  float64(replay.Info.JumpDistance),
		})
		_ = file.Close()
	}

	bestR2 := -1.0
	bestDegree := 1
	bestFormula := ""
	var bestModel regression.Regression

	maxDegree := 5

	for degree := 1; degree <= maxDegree; degree++ {
		var r regression.Regression
		r.SetObserved("JD")

		for d := 1; d <= degree; d++ {
			r.SetVar(d-1, fmt.Sprintf("NJS^%d", d))
		}

		for _, dp := range pairs {
			var x []float64
			for d := 1; d <= degree; d++ {
				x = append(x, math.Pow(dp.NJS, float64(d)))
			}
			r.Train(regression.DataPoint(dp.JD, x))
		}

		err = r.Run()
		if err != nil {
			slog.Info(fmt.Sprintf("Error training model for degree %d: %v\n", degree, err))
			continue
		}

		slog.Info(fmt.Sprintf("Degree %d: R² = %.6f\n", degree, r.R2))

		if r.R2 > bestR2 {
			bestR2 = r.R2
			bestDegree = degree
			bestFormula = r.Formula
			bestModel = r
		}
	}

	slog.Info(fmt.Sprintf("Degree: %d\n", bestDegree))
	slog.Info(fmt.Sprintf("Formula: %s\n", bestFormula))
	slog.Info(fmt.Sprintf("R²: %.6f\n", bestR2))

	slog.Info(fmt.Sprintf("Training completed with %d data points.", len(pairs)))

	if len(pairs) < 50 {
		slog.Info("WARNING: Please note that the reliability of the model increases with more training data and a wider range of maps. Consider fetching more replays with different njs values.")
	}

	var configPairs []models.JDPair
	var njs = models.JDConfigLow

	for njs < models.JDConfigHigh {

		x := make([]float64, bestDegree)
		for d := 1; d <= bestDegree; d++ {
			x[d-1] = math.Pow(njs, float64(d))
		}

		predictedJD, err := bestModel.Predict(x)
		if err != nil {
			return err
		}

		configPairs = append(configPairs, models.JDPair{
			NJS: njs,
			JD:  predictedJD,
		})

		njs += 0.5
	}

	bytes, err := json.MarshalIndent(configPairs, "", "   ")
	if err != nil {
		return err
	}

	slog.Info(string(bytes))

	return nil
}

func ShowInfo() string {
	f, err := os.Open("path")
	if err != nil {
		log.Fatalf("Failed to open BSOR file: %v", err)
	}
	defer f.Close()

	r, err := bsor.Read(f)
	if err != nil {
		log.Fatalf("Failed to decode BSOR: %v", err)
	}

	result := analyser.AnalyzeControllers(r.Frames)
	output, err := json.MarshalIndent(struct {
		ReplayInfo      *bsor.Info            `json:"replay_info"`
		AnalysisResults models.AnalysisResult `json:"analysis_results"`
	}{&r.Info, result}, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}

	return string(output)
}

// Local Utils

func calculateRSq(r *regression.Regression, data []models.JDPair) (float64, error) {
	var ssTot, ssRes, meanJD float64
	for _, dp := range data {
		meanJD += dp.JD
	}
	meanJD /= float64(len(data))

	for _, dp := range data {
		pred, err := r.Predict([]float64{dp.NJS, math.Pow(dp.NJS, 2)})
		if err != nil {
			slog.Error("Failed to calculate RSq", err.Error())
			return 0, err
		}
		ssTot += math.Pow(dp.JD-meanJD, 2)
		ssRes += math.Pow(dp.JD-pred, 2)
	}

	return 1 - (ssRes / ssTot), nil
}

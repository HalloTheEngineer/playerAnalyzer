package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/motzel/go-bsor/bsor"
	"log"
	"log/slog"
	"os"
	"replayAnalyzer/analyser"
	"replayAnalyzer/logic"
	"replayAnalyzer/storage"
	"replayAnalyzer/utils"
)

func handleFetchCmd(ctx context.Context, args []string) (err error) {
	_ = os.Mkdir("_replays", 0755)

	var playerId string

	if len(args) > 0 {
		playerId = args[0]
	} else {
		playerId, err = utils.GetInput("Enter player id: ")
		if err != nil {
			return err
		}
	}

	if !utils.RequireNumbers(playerId) {
		return errors.New("valid player id is required")
	}

	var queries []string

	doCustom, err := utils.GetInput("Do you want to use custom url parameters for fetching? (Y|n) ")
	if err != nil {
		return err
	}
	if doCustom == "y" {
		params, err := utils.GetInput("Provide beatleader-api query params (format: ?param=value&param2=value):\n" + fmt.Sprintf(utils.BeatLeaderScoresUrl, playerId))
		if err != nil {
			return err
		}
		queries = append(queries, fmt.Sprintf(utils.BeatLeaderScoresUrl, playerId)+params)
	} else {
		for _, diff := range utils.BLDifficulties {
			queries = append(queries, fmt.Sprintf(utils.BeatLeaderScoresUrl, playerId)+fmt.Sprintf(utils.BeatLeaderDefaultScoreParams, diff))
		}
	}

	err = logic.FetchScores(playerId, queries)
	if err != nil {
		return err
	}

	return nil
}

// handleJDGenCmd Concept by HalloTheEngineer; logic implementation by Claude 3.7 Sonnet
func handleJDGenCmd(ctx context.Context, args []string) (err error) {
	_ = os.Mkdir("_cache/plots", 0777)
	_ = os.Mkdir("_cache/jd_configs", 0777)

	var playerId string

	if len(args) > 0 {
		playerId = args[0]
	} else {
		playerId, err = utils.GetInput("Enter player id: ")
		if err != nil {
			return err
		}
	}

	slog.Info("Fetching player info")

	player, err := utils.FetchToStruct[utils.SSPlayer](fmt.Sprintf("https://scoresaber.com/api/player/%s/basic", playerId))
	if err != nil {
		return err
	}

	err = logic.GenerateJDConfig(player)
	if err != nil {
		return err
	}

	return nil
}

func handleCleanupCmd(ctx context.Context, args []string) (err error) {

	var playerId string
	var removedReplays int

	if len(args) > 0 {
		playerId = args[0]
	} else {
		playerId, err = utils.GetInput("Enter player id: ")
		if err != nil {
			return err
		}
	}

	replays, err := storage.GetReplays(playerId)
	if err != nil {
		return err
	}

	for _, rep := range *replays {
		lstat, err := os.Lstat(fmt.Sprintf("_replays/%s/%s", playerId, rep))
		if err != nil {
			continue
		}
		if lstat.Size() < 200*1000 {
			_ = os.Remove(fmt.Sprintf("_replays/%s/%s", playerId, rep))
			removedReplays++
		}
	}

	slog.Info(fmt.Sprintf("Removed %d corrupted replays", removedReplays))

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
		ReplayInfo      *bsor.Info           `json:"replay_info"`
		AnalysisResults utils.AnalysisResult `json:"analysis_results"`
	}{&r.Info, result}, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}

	return string(output)
}

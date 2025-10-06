package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"playerAnalyzer/logic"
	"playerAnalyzer/models"
	"playerAnalyzer/utils"
)

// handleJDGenCmd Concept by HalloTheEngineer; logic implementation by Claude 3.7 Sonnet
func handleJDGenCmd(ctx context.Context, args []string) (err error) {
	_ = os.MkdirAll("_cache/plots", os.ModePerm)
	_ = os.MkdirAll("_cache/jd_configs", os.ModePerm)

	var playerId string
	var settings = models.Settings{
		Count:  100,
		Sort:   "top",
		Ranked: true,
	}

	playerId, err = utils.GetInput("Enter player id: ")
	if err != nil {
		return err
	}

	lCount, err := utils.GetInput("Enter score count: ")
	if err != nil {
		return err
	}
	settings.SetCount(lCount)

	lSort, err := utils.GetInput("Enter sort order (1=top, 2=recent): ")
	if err != nil {
		return err
	}
	settings.SetSort(lSort)

	lRanked, err := utils.GetInput("Enter ranked status (true,false): ")
	if err != nil {
		return err
	}
	settings.SetRanked(lRanked)

	slog.Info("Fetching player info")

	player, err := utils.FetchToStruct[utils.SSPlayer](fmt.Sprintf("https://scoresaber.com/api/player/%s/basic", playerId))
	if err != nil {
		return err
	}

	err = logic.GenerateJDConfig(player, settings)
	if err != nil {
		return err
	}

	return nil
}

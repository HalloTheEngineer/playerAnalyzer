package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"playerAnalyzer/logic"
	"playerAnalyzer/utils"
	"strconv"
)

// handleJDGenCmd Concept by HalloTheEngineer; logic implementation by Claude 3.7 Sonnet
func handleJDGenCmd(ctx context.Context, args []string) (err error) {
	_ = os.MkdirAll("_cache/plots", os.ModePerm)
	_ = os.MkdirAll("_cache/jd_configs", os.ModePerm)

	var playerId string
	var count int
	var sort string

	playerId, err = utils.GetInput("Enter player id: ")
	if err != nil {
		return err
	}
	lCount, err := utils.GetInput("Enter score count: ")
	if err != nil {
		return err
	}
	count, err = strconv.Atoi(lCount)
	if err != nil {
		return err
	}

	lSort, err := utils.GetInput("Enter sort order (1=top, 2=recent): ")
	if err != nil {
		return err
	}
	switch lSort {
	case "1":
		sort = "top"
	case "2":
		sort = "recent"
	default:
		sort = "top"
	}

	slog.Info("Fetching player info")

	player, err := utils.FetchToStruct[utils.SSPlayer](fmt.Sprintf("https://scoresaber.com/api/player/%s/basic", playerId))
	if err != nil {
		return err
	}

	err = logic.GenerateJDConfig(player, count, sort)
	if err != nil {
		return err
	}

	return nil
}

package storage

import (
	"fmt"
	"log/slog"
	"playerAnalyzer/models"
	"playerAnalyzer/utils"
)

const ssScoresUrl = "https://scoresaber.com/api/player/%s/scores?limit=%d&sort=%s&page=%d&withMetadata=true"
const blSpecScoreUrl = "https://api.beatleader.com/score/%s/%s/%s/Standard?leaderboardContext=general"
const blLeaderboardUrl = "https://api.beatleader.com/leaderboard/%s/%s/Standard"
const statsUrl = "https://cdn.scorestats.beatleader.com/%d.json"

// scoresaber scores > songHash + difficulty + gameMode > bl /leaderboard/hash/diff/mode > score > id > stats

func FetchStats(playerId string, settings models.Settings) ([]*utils.StatsResult, error) {
	var res []*utils.StatsResult

	ssScores, err := fetchAllScores(playerId, settings.Count, settings.Sort)
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Fetched %d scores of %s", len(ssScores.PlayerScores), playerId))

	for i, score := range ssScores.PlayerScores {
		if settings.Ranked && !score.Leaderboard.Ranked {
			continue
		}

		slog.Info(fmt.Sprintf("(%d) - %s", i+1, score.Leaderboard.SongName))
		// Fetching concrete BL play by criteria
		blScore, err := utils.FetchToStruct[utils.BLScore](fmt.Sprintf(blSpecScoreUrl, playerId, score.Leaderboard.SongHash, formatSSDiff(score.Leaderboard.Difficulty.Difficulty)))
		if err != nil {
			slog.Info("BL Score: " + err.Error())
			continue
		}

		blLead, err := utils.FetchToStruct[utils.BLLeaderboard](fmt.Sprintf(blLeaderboardUrl, score.Leaderboard.SongHash, formatSSDiff(score.Leaderboard.Difficulty.Difficulty)))
		if err != nil {
			slog.Info("BL Leaderboard: " + err.Error())
			continue
		}

		// Fetching corresponding stats of the play
		blStats, err := utils.FetchToStruct[utils.ScoreStats](fmt.Sprintf(statsUrl, blScore.Id))
		if err != nil {
			slog.Info("BL Stats: " + err.Error())
			continue
		}

		res = append(res, &utils.StatsResult{
			BLLead: blLead,
			Stats:  blStats,
		})
	}

	return res, nil
}

func fetchAllScores(playerId string, count int, sortOrder string) (*utils.SSScoreResponse, error) {
	const maxScoresPerPage = 100

	allScores := &utils.SSScoreResponse{
		PlayerScores: make([]utils.SSScore, 0, count),
	}

	page := 1
	remaining := count

	for remaining > 0 {
		limit := maxScoresPerPage
		if remaining < maxScoresPerPage {
			limit = remaining
		}

		pageScores, err := utils.FetchToStruct[utils.SSScoreResponse](fmt.Sprintf(ssScoresUrl, playerId, limit, sortOrder, page))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if len(pageScores.PlayerScores) == 0 {
			break
		}

		allScores.PlayerScores = append(allScores.PlayerScores, pageScores.PlayerScores...)

		allScores.Metadata = pageScores.Metadata

		remaining -= len(pageScores.PlayerScores)
		page++

		if len(pageScores.PlayerScores) < limit {
			break
		}

		slog.Debug(fmt.Sprintf("Fetched page %d: %d scores, %d remaining", page-1, len(pageScores.PlayerScores), remaining))
	}

	if len(allScores.PlayerScores) > count {
		allScores.PlayerScores = allScores.PlayerScores[:count]
	}

	return allScores, nil
}

func formatSSDiff(difficulty int) string {
	switch difficulty {
	case 1:
		return "Easy"
	case 3:
		return "Normal"
	case 5:
		return "Hard"
	case 7:
		return "Expert"
	case 9:
		return "ExpertPlus"
	}
	return ""
}

package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"replayAnalyzer/models"
	"strings"
)

func SaveLeaderboard(leaderboard *models.ALeaderboard) error {
	nBytes, err := json.MarshalIndent(leaderboard, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("_cache/leaderboards/"+leaderboard.Id+".json", nBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
func LoadLeaderboard(name string) (*models.ALeaderboard, error) {
	var leaderboard models.ALeaderboard

	bytes, err := os.ReadFile(filepath.Join("_cache/leaderboards", name))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &leaderboard)
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

func GetDiff(hash, diff string) *models.ADifficulty {
	for _, leaderboard := range LeaderboardCache {
		if strings.ToLower(leaderboard.Song.Hash) == strings.ToLower(hash) && strings.ToLower(leaderboard.Difficulty.DifficultyName) == strings.ToLower(diff) {
			return &leaderboard.Difficulty
		}
	}
	return nil
}

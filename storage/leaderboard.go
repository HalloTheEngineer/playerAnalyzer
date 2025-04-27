package storage

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"replayAnalyzer/utils"
	"strings"
)

func SaveLeaderboard(leaderboard *utils.ALeaderboard) error {
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
func LoadLeaderboard(name string) (*utils.ALeaderboard, error) {
	var leaderboard utils.ALeaderboard

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

func GetDiff(hash, diff string) *utils.ADifficulty {
	hash = strings.ToLower(hash)
	hash = strings.Split(hash, " ")[0]

	if strings.ContainsRune(hash, '_') {
		hash = hash[:strings.IndexRune(hash, '_')]
	}

	for _, leaderboard := range LeaderboardCache {
		if strings.ToLower(strings.Split(leaderboard.Song.Hash, " ")[0]) == hash && strings.ToLower(leaderboard.Difficulty.DifficultyName) == strings.ToLower(diff) {
			return &leaderboard.Difficulty
		}
	}
	slog.Info(hash)
	return nil
}

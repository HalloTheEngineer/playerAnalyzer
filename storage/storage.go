package storage

import (
	"os"
	"replayAnalyzer/models"
)

var ReplayCache []string
var LeaderboardCache []*models.ALeaderboard

func Load() error {

	_ = os.MkdirAll("_cache/leaderboards", 0755)
	_ = os.Mkdir("_replays", 0755)

	entries, _ := os.ReadDir("_cache/leaderboards")

	for _, entry := range entries {
		if !entry.IsDir() {

			l, err := LoadLeaderboard(entry.Name())
			if err != nil {
				return err
			}

			LeaderboardCache = append(LeaderboardCache, l)
		}
	}

	return nil
}

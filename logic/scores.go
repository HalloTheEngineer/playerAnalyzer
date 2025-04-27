package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"io"
	"log/slog"
	"net/http"
	"replayAnalyzer/storage"
	"replayAnalyzer/utils"
	"strconv"
	"strings"
)

func FetchScores(playerId string, queries []string) error {
	var leaderboardCount int
	var cdnUrlCount int

	for _, query := range queries {
		slog.Info("Collecting scores for request: " + query)

		resp, err := http.Get(query)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 200 {

			bts, _ := io.ReadAll(resp.Body)

			slog.Debug("Request returned status code: " + strconv.Itoa(resp.StatusCode) + "\n" + string(bts))
			return errors.New("beatleader-api returned status code " + strconv.Itoa(resp.StatusCode))
		}

		jsonBuf, err := gabs.ParseJSONBuffer(resp.Body)
		if err != nil {
			return err
		}

		var cdnUrls []string
		var leaderboards []utils.ALeaderboard

		for _, score := range jsonBuf.S("data").Children() {
			var l utils.ALeaderboard

			err = json.Unmarshal(score.S("leaderboard").Bytes(), &l)
			if err != nil {
				return err
			}

			cdnUrls = append(cdnUrls, strings.Trim(score.S("replay").String(), "\""))
			leaderboards = append(leaderboards, l)
		}

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

		leaderboardCount += len(leaderboards)
		cdnUrlCount = len(cdnUrls)
	}

	slog.Info(fmt.Sprintf("Fetched %d leaderboards", leaderboardCount))
	slog.Info(fmt.Sprintf("Fetched %d replays of %s", cdnUrlCount, playerId))

	return nil
}

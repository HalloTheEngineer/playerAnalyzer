package storage

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func FetchReplay(playerId, cdnUrl string) error {
	urlParts := strings.Split(cdnUrl, "/")

	if i, err := os.Lstat("_replays/" + playerId + "/" + urlParts[len(urlParts)-1]); err == nil && !i.IsDir() {
		return nil
	}

	slog.Info(fmt.Sprintf("Fetching replay (%s)", urlParts[len(urlParts)-1]))

	_ = os.Mkdir("_replays/"+playerId, 0755)

	replayResp, err := http.Get(cdnUrl)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(replayResp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile("_replays/"+playerId+"/"+urlParts[len(urlParts)-1], bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetReplays(playerId string) (*[]string, error) {
	entries, err := os.ReadDir("_replays/" + playerId)
	if err != nil {
		return nil, err
	}

	var replays []string

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".bsor") {
			replays = append(replays, entry.Name())
		}
	}
	return &replays, nil
}

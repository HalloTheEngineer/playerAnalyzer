package builder

import (
	"github.com/motzel/go-bsor/bsor"
	"log"
	"os"
	"path/filepath"
	"replayAnalyzer/models"
	"strconv"
)

func GenerateJDConfigFromFolder(folderPath string) models.JDConfig {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	mapNJS := make(map[float64][]float64)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".bsor" {
			continue
		}
		file, err := os.Open(filepath.Join(folderPath, entry.Name()))
		if err != nil {
			continue
		}

		r, err := bsor.Read(file)
		file.Close()

		if err != nil {
			continue
		}

		njs := float64(r.Info.NoteJumpSpeed)
		left := averageSwingDistance(r.Frames.LeftHand)
		right := averageSwingDistance(r.Frames.RightHand)
		avgJD := (left + right) / 2
		mapNJS[njs] = append(mapNJS[njs], avgJD)
	}

	jdmap := JDMap{Default: 18.0, Map: map[string]float64{}}
	for njs, distances := range mapNJS {
		sum := 0.0
		for _, d := range distances {
			sum += d
		}
		jdmap.Map[strconv.FormatFloat(njs, 'f', 1, 64)] = sum / float64(len(distances))
	}
	return jdmap
}

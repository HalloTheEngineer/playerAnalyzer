package utils

import (
	"fmt"
	"github.com/sajari/regression"
	"gonum.org/v1/plot/plotter"
	"time"
)

const (
	BeatLeaderApiUrl    = "https://api.beatleader.com"
	BeatLeaderScoresUrl = BeatLeaderApiUrl + "/player/%s/scores"

	BeatLeaderDefaultScoreParams = "?sortBy=date&order=desc&page=1&count=64&diff=%s&mode=Standard&requirements=none&scoreStatus=none&leaderboardContext=general&type=ranked&includeIO=false"

	JDConfigLow  = 8.0
	JDConfigHigh = 26.0
)

var (
	BLDifficulties = []string{"Easy", "Normal", "Hard", "Expert", "ExpertPlus"}
)

type (
	AnalysisResult struct {
		DominantHand         string  `json:"dominant_hand"`
		AvgSwingDistanceL    float64 `json:"avg_swing_distance_left"`
		AvgSwingDistanceR    float64 `json:"avg_swing_distance_right"`
		SwingIntensityDiff   float64 `json:"swing_intensity_diff"`
		EstimatedControllerL string  `json:"estimated_controller_left"`
		EstimatedControllerR string  `json:"estimated_controller_right"`
		EstimatedGripStyle   string  `json:"estimated_grip_style"`
	}

	JDConfig struct {
		PreferredValues []JDPair `json:"preferredValues"`
	}
	JDPair struct {
		NJS float64 `json:"njs"`
		JD  float64 `json:"jumpDistance"`
	}

	Point struct {
		X float64
		Y float64
	}
	Cluster struct {
		Points []plotter.XY
		Model  *regression.Regression
	}

	ALeaderboard struct {
		Id   string `json:"id"`
		Song struct {
			Id              string      `json:"id"`
			Hash            string      `json:"hash"`
			Name            string      `json:"name"`
			SubName         string      `json:"subName"`
			Author          string      `json:"author"`
			Mapper          string      `json:"mapper"`
			MapperId        int         `json:"mapperId"`
			CollaboratorIds interface{} `json:"collaboratorIds"`
			CoverImage      string      `json:"coverImage"`
			Bpm             float64     `json:"bpm"`
			Duration        int         `json:"duration"`
			FullCoverImage  string      `json:"fullCoverImage"`
			Explicity       int         `json:"explicity"`
		} `json:"song"`
		Difficulty ADifficulty `json:"difficulty"`
	}
	ADifficulty struct {
		Id             int     `json:"id"`
		Value          int     `json:"value"`
		Mode           int     `json:"mode"`
		DifficultyName string  `json:"difficultyName"`
		ModeName       string  `json:"modeName"`
		Status         int     `json:"status"`
		NominatedTime  int     `json:"nominatedTime"`
		QualifiedTime  int     `json:"qualifiedTime"`
		RankedTime     int     `json:"rankedTime"`
		SpeedTags      int     `json:"speedTags"`
		StyleTags      int     `json:"styleTags"`
		FeatureTags    int     `json:"featureTags"`
		Stars          float64 `json:"stars"`
		PredictedAcc   float64 `json:"predictedAcc"`
		PassRating     float64 `json:"passRating"`
		AccRating      float64 `json:"accRating"`
		TechRating     float64 `json:"techRating"`
		Type           int     `json:"type"`
		Njs            float64 `json:"njs"`
		Nps            float64 `json:"nps"`
		Notes          int     `json:"notes"`
		Bombs          int     `json:"bombs"`
		Walls          int     `json:"walls"`
		MaxScore       int     `json:"maxScore"`
		Duration       int     `json:"duration"`
		Requirements   int     `json:"requirements"`
	}

	SSPlayer struct {
		Id             string      `json:"id"`
		Name           string      `json:"name"`
		ProfilePicture string      `json:"profilePicture"`
		Bio            interface{} `json:"bio"`
		Country        string      `json:"country"`
		Pp             float64     `json:"pp"`
		Rank           int         `json:"rank"`
		CountryRank    int         `json:"countryRank"`
		Role           interface{} `json:"role"`
		Badges         interface{} `json:"badges"`
		Histories      string      `json:"histories"`
		Permissions    int         `json:"permissions"`
		Banned         bool        `json:"banned"`
		Inactive       bool        `json:"inactive"`
		ScoreStats     interface{} `json:"scoreStats"`
		FirstSeen      time.Time   `json:"firstSeen"`
	}
)

func (p JDPair) ToString() string {
	return fmt.Sprintf("NJS:%f  JD:%f", p.NJS, p.JD)
}

package models

const (
	BeatleaderApiUrl = "https://api.beatleader.com"

	JDConfigLow  = 5.0
	JDConfigHigh = 26.0
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
)

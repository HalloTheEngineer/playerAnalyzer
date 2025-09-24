package utils

import (
	"fmt"
	"time"

	"github.com/sajari/regression"
	"gonum.org/v1/plot/plotter"
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

	ScoreStats struct {
		HitTracker struct {
			MaxCombo     int     `json:"maxCombo"`
			MaxStreak    int     `json:"maxStreak"`
			LeftTiming   float64 `json:"leftTiming"`
			RightTiming  float64 `json:"rightTiming"`
			LeftMiss     int     `json:"leftMiss"`
			RightMiss    int     `json:"rightMiss"`
			LeftBadCuts  int     `json:"leftBadCuts"`
			RightBadCuts int     `json:"rightBadCuts"`
			LeftBombs    int     `json:"leftBombs"`
			RightBombs   int     `json:"rightBombs"`
		} `json:"hitTracker"`
		AccuracyTracker struct {
			AccRight            float64   `json:"accRight"`
			AccLeft             float64   `json:"accLeft"`
			LeftPreswing        float64   `json:"leftPreswing"`
			RightPreswing       float64   `json:"rightPreswing"`
			AveragePreswing     float64   `json:"averagePreswing"`
			LeftPostswing       float64   `json:"leftPostswing"`
			RightPostswing      float64   `json:"rightPostswing"`
			LeftTimeDependence  float64   `json:"leftTimeDependence"`
			RightTimeDependence float64   `json:"rightTimeDependence"`
			LeftAverageCut      []float64 `json:"leftAverageCut"`
			RightAverageCut     []float64 `json:"rightAverageCut"`
			GridAcc             []float64 `json:"gridAcc"`
			FcAcc               float64   `json:"fcAcc"`
		} `json:"accuracyTracker"`
		WinTracker struct {
			Won                 bool    `json:"won"`
			EndTime             float64 `json:"endTime"`
			FailTime            float64 `json:"failTime"`
			NbOfPause           int     `json:"nbOfPause"`
			TotalPauseDuration  float64 `json:"totalPauseDuration"`
			JumpDistance        float64 `json:"jumpDistance"`
			AverageHeight       float64 `json:"averageHeight"`
			AverageHeadPosition struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"averageHeadPosition"`
			TotalScore int `json:"totalScore"`
			MaxScore   int `json:"maxScore"`
		} `json:"winTracker"`
		ScoreGraphTracker struct {
			Graph []float64 `json:"graph"`
		} `json:"scoreGraphTracker"`
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

	BLScoreResponse struct {
		Metadata struct {
			ItemsPerPage int `json:"itemsPerPage"`
			Page         int `json:"page"`
			Total        int `json:"total"`
		} `json:"metadata"`
		Data []BLScore `json:"data"`
	}

	BLScore struct {
		Id            int     `json:"id"`
		BaseScore     int     `json:"baseScore"`
		ModifiedScore int     `json:"modifiedScore"`
		Accuracy      float64 `json:"accuracy"`
		PlayerId      string  `json:"playerId"`
		Pp            float64 `json:"pp"`
		BonusPp       float64 `json:"bonusPp"`
		PassPP        float64 `json:"passPP"`
		AccPP         float64 `json:"accPP"`
		TechPP        float64 `json:"techPP"`
		Qualification bool    `json:"qualification"`
		Weight        float64 `json:"weight"`
		Rank          int     `json:"rank"`
		CountryRank   int     `json:"countryRank"`
		Replay        string  `json:"replay"`
		Modifiers     string  `json:"modifiers"`
		ModifiedStars float64 `json:"modifiedStars"`
		BadCuts       int     `json:"badCuts"`
		MissedNotes   int     `json:"missedNotes"`
		BombCuts      int     `json:"bombCuts"`
		WallsHit      int     `json:"wallsHit"`
		Mistakes      int     `json:"mistakes"`
		Pauses        int     `json:"pauses"`
		FullCombo     bool    `json:"fullCombo"`
		MaxCombo      int     `json:"maxCombo"`
		FcAccuracy    float64 `json:"fcAccuracy"`
		FcPp          float64 `json:"fcPp"`
		Hmd           int     `json:"hmd"`
		Controller    int     `json:"controller"`
		AccRight      float64 `json:"accRight"`
		AccLeft       float64 `json:"accLeft"`
		Timeset       string  `json:"timeset"`
		Timepost      int     `json:"timepost"`
		Platform      string  `json:"platform"`
		Player        struct {
			Id                  string      `json:"id"`
			Name                string      `json:"name"`
			Platform            string      `json:"platform"`
			Avatar              string      `json:"avatar"`
			WebAvatar           string      `json:"webAvatar"`
			Country             string      `json:"country"`
			Alias               string      `json:"alias"`
			OldAlias            interface{} `json:"oldAlias"`
			Role                string      `json:"role"`
			MapperId            interface{} `json:"mapperId"`
			Mapper              interface{} `json:"mapper"`
			Pp                  float64     `json:"pp"`
			AccPp               float64     `json:"accPp"`
			TechPp              float64     `json:"techPp"`
			PassPp              float64     `json:"passPp"`
			AllContextsPp       float64     `json:"allContextsPp"`
			Rank                int         `json:"rank"`
			CountryRank         int         `json:"countryRank"`
			Level               int         `json:"level"`
			Experience          int         `json:"experience"`
			Prestige            int         `json:"prestige"`
			LastWeekPp          float64     `json:"lastWeekPp"`
			LastWeekRank        int         `json:"lastWeekRank"`
			LastWeekCountryRank int         `json:"lastWeekCountryRank"`
			Banned              bool        `json:"banned"`
			Bot                 bool        `json:"bot"`
			Temporary           bool        `json:"temporary"`
			Inactive            bool        `json:"inactive"`
			ExternalProfileUrl  string      `json:"externalProfileUrl"`
			RichBioTimeset      int         `json:"richBioTimeset"`
			CreatedAt           int         `json:"createdAt"`
			SpeedrunStart       int         `json:"speedrunStart"`
			ScoreStatsId        int         `json:"scoreStatsId"`
			ScoreStats          interface{} `json:"scoreStats"`
			ClanOrder           string      `json:"clanOrder"`
			Badges              interface{} `json:"badges"`
			PatreonFeatures     struct {
				Id              int    `json:"id"`
				Bio             string `json:"bio"`
				Message         string `json:"message"`
				LeftSaberColor  string `json:"leftSaberColor"`
				RightSaberColor string `json:"rightSaberColor"`
			} `json:"patreonFeatures"`
			ProfileSettings struct {
				Id                    int         `json:"id"`
				Bio                   interface{} `json:"bio"`
				Message               interface{} `json:"message"`
				EffectName            string      `json:"effectName"`
				ProfileAppearance     string      `json:"profileAppearance"`
				Hue                   int         `json:"hue"`
				Saturation            float64     `json:"saturation"`
				LeftSaberColor        interface{} `json:"leftSaberColor"`
				RightSaberColor       interface{} `json:"rightSaberColor"`
				ProfileCover          string      `json:"profileCover"`
				StarredFriends        string      `json:"starredFriends"`
				HorizontalRichBio     bool        `json:"horizontalRichBio"`
				RankedMapperSort      string      `json:"rankedMapperSort"`
				ShowBots              bool        `json:"showBots"`
				ShowAllRatings        bool        `json:"showAllRatings"`
				ShowExplicitCovers    bool        `json:"showExplicitCovers"`
				ShowStatsPublic       bool        `json:"showStatsPublic"`
				ShowStatsPublicPinned bool        `json:"showStatsPublicPinned"`
			} `json:"profileSettings"`
			Changes             interface{} `json:"changes"`
			EventsParticipating interface{} `json:"eventsParticipating"`
			Socials             interface{} `json:"socials"`
			Achievements        interface{} `json:"achievements"`
			EarthDayMap         interface{} `json:"earthDayMap"`
		} `json:"player"`
		ValidContexts           int           `json:"validContexts"`
		ValidForGeneral         bool          `json:"validForGeneral"`
		ContextExtensions       []interface{} `json:"contextExtensions"`
		LeaderboardId           string        `json:"leaderboardId"`
		Leaderboard             interface{}   `json:"leaderboard"`
		AuthorizedReplayWatched int           `json:"authorizedReplayWatched"`
		AnonimusReplayWatched   int           `json:"anonimusReplayWatched"`
		ReplayWatchedTotal      int           `json:"replayWatchedTotal"`
		ReplayOffsetsId         int           `json:"replayOffsetsId"`
		ReplayOffsets           interface{}   `json:"replayOffsets"`
		Country                 string        `json:"country"`
		MaxStreak               int           `json:"maxStreak"`
		PlayCount               int           `json:"playCount"`
		LastTryTime             int           `json:"lastTryTime"`
		LeftTiming              float64       `json:"leftTiming"`
		RightTiming             float64       `json:"rightTiming"`
		Priority                int           `json:"priority"`
		ScoreImprovementId      int           `json:"scoreImprovementId"`
		ScoreImprovement        interface{}   `json:"scoreImprovement"`
		Banned                  bool          `json:"banned"`
		Suspicious              bool          `json:"suspicious"`
		Bot                     bool          `json:"bot"`
		IgnoreForStats          bool          `json:"ignoreForStats"`
		Migrated                bool          `json:"migrated"`
		RankVoting              interface{}   `json:"rankVoting"`
		Metadata                interface{}   `json:"metadata"`
		Experience              float64       `json:"experience"`
		Status                  int           `json:"status"`
		ExternalStatuses        interface{}   `json:"externalStatuses"`
		SotwNominations         int           `json:"sotwNominations"`
	}

	BLLeaderboard struct {
		Id   string `json:"id"`
		Song struct {
			Id              string  `json:"id"`
			Hash            string  `json:"hash"`
			Name            string  `json:"name"`
			SubName         string  `json:"subName"`
			Author          string  `json:"author"`
			Mapper          string  `json:"mapper"`
			MapperId        int     `json:"mapperId"`
			CollaboratorIds *string `json:"collaboratorIds"`
			CoverImage      string  `json:"coverImage"`
			Bpm             float64 `json:"bpm"`
			Duration        int     `json:"duration"`
			FullCoverImage  string  `json:"fullCoverImage"`
			Explicity       int     `json:"explicity"`
		} `json:"song"`
		Difficulty struct {
			Id             int    `json:"id"`
			Value          int    `json:"value"`
			Mode           int    `json:"mode"`
			DifficultyName string `json:"difficultyName"`
			ModeName       string `json:"modeName"`
			Status         int    `json:"status"`
			ModifierValues struct {
				ModifierId int     `json:"modifierId"`
				Da         int     `json:"da"`
				Fs         float64 `json:"fs"`
				Sf         float64 `json:"sf"`
				Ss         float64 `json:"ss"`
				Gn         float64 `json:"gn"`
				Na         float64 `json:"na"`
				Nb         float64 `json:"nb"`
				Nf         float64 `json:"nf"`
				No         float64 `json:"no"`
				Pm         int     `json:"pm"`
				Sc         int     `json:"sc"`
				Sa         int     `json:"sa"`
				Op         float64 `json:"op"`
				Ez         float64 `json:"ez"`
				Hd         float64 `json:"hd"`
				Smc        float64 `json:"smc"`
				Ohp        int     `json:"ohp"`
			} `json:"modifierValues"`
			NominatedTime int     `json:"nominatedTime"`
			QualifiedTime int     `json:"qualifiedTime"`
			RankedTime    int     `json:"rankedTime"`
			SpeedTags     int     `json:"speedTags"`
			StyleTags     int     `json:"styleTags"`
			FeatureTags   int     `json:"featureTags"`
			Stars         float64 `json:"stars"`
			PredictedAcc  float64 `json:"predictedAcc"`
			PassRating    float64 `json:"passRating"`
			AccRating     float64 `json:"accRating"`
			TechRating    float64 `json:"techRating"`
			Type          int     `json:"type"`
			Njs           float64 `json:"njs"`
			Nps           float64 `json:"nps"`
			Notes         int     `json:"notes"`
			Bombs         int     `json:"bombs"`
			Walls         int     `json:"walls"`
			MaxScore      int     `json:"maxScore"`
			Duration      int     `json:"duration"`
			Requirements  int     `json:"requirements"`
		} `json:"difficulty"`
	}

	SSScoreResponse struct {
		PlayerScores []SSScore   `json:"playerScores"`
		Metadata     interface{} `json:"metadata"`
	}

	SSScore struct {
		Score struct {
			Id                    int         `json:"id"`
			LeaderboardPlayerInfo interface{} `json:"leaderboardPlayerInfo"`
			Rank                  int         `json:"rank"`
			BaseScore             int         `json:"baseScore"`
			ModifiedScore         int         `json:"modifiedScore"`
			Pp                    float64     `json:"pp"`
			Weight                float64     `json:"weight"`
			Modifiers             string      `json:"modifiers"`
			Multiplier            int         `json:"multiplier"`
			BadCuts               int         `json:"badCuts"`
			MissedNotes           int         `json:"missedNotes"`
			MaxCombo              int         `json:"maxCombo"`
			FullCombo             bool        `json:"fullCombo"`
			Hmd                   int         `json:"hmd"`
			TimeSet               time.Time   `json:"timeSet"`
			HasReplay             bool        `json:"hasReplay"`
			DeviceHmd             string      `json:"deviceHmd"`
			DeviceControllerLeft  string      `json:"deviceControllerLeft"`
			DeviceControllerRight string      `json:"deviceControllerRight"`
		} `json:"score"`
		Leaderboard struct {
			Id              int    `json:"id"`
			SongHash        string `json:"songHash"`
			SongName        string `json:"songName"`
			SongSubName     string `json:"songSubName"`
			SongAuthorName  string `json:"songAuthorName"`
			LevelAuthorName string `json:"levelAuthorName"`
			Difficulty      struct {
				LeaderboardId int    `json:"leaderboardId"`
				Difficulty    int    `json:"difficulty"`
				GameMode      string `json:"gameMode"`
				DifficultyRaw string `json:"difficultyRaw"`
			} `json:"difficulty"`
			MaxScore          int         `json:"maxScore"`
			CreatedDate       time.Time   `json:"createdDate"`
			RankedDate        time.Time   `json:"rankedDate"`
			QualifiedDate     *time.Time  `json:"qualifiedDate"`
			LovedDate         interface{} `json:"lovedDate"`
			Ranked            bool        `json:"ranked"`
			Qualified         bool        `json:"qualified"`
			Loved             bool        `json:"loved"`
			MaxPP             int         `json:"maxPP"`
			Stars             float64     `json:"stars"`
			Plays             int         `json:"plays"`
			DailyPlays        int         `json:"dailyPlays"`
			PositiveModifiers bool        `json:"positiveModifiers"`
			PlayerScore       interface{} `json:"playerScore"`
			CoverImage        string      `json:"coverImage"`
			Difficulties      interface{} `json:"difficulties"`
		} `json:"leaderboard"`
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

	StatsResult struct {
		BLLead *BLLeaderboard `json:"blLeaderboard"`
		Stats  *ScoreStats    `json:"stats"`
	}
)

func (p JDPair) ToString() string {
	return fmt.Sprintf("NJS:%f  JD:%f", p.NJS, p.JD)
}

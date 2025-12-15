package users_username_profile_get_handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type Request struct {
	Username string
}

type PersonalBest struct {
	Language string  `json:"language" example:"english"`
	Mode     string  `json:"mode" example:"words"`
	SubMode  string  `json:"submode" example:"25"`
	WPM      float64 `json:"wpm" example:"100"`
	Accuracy float64 `json:"accuracy" example:"300"`
	PlayedAt string  `json:"played_at" example:"2020-09-09T10:10:10Z"`
} //@name UsersUsernameProfileGetHandler.PersonalBest

type LanguageStats struct {
	Language       string `json:"language" example:"english"`
	TestsCompleted uint64 `json:"testsCompleted" example:"1"`
	TestsStarted   uint64 `json:"testsStarted" example:"1"`
	TimePlayedMs   int64  `json:"timePlayed" example:"1000"`
} //@name UsersUsernameProfileGetHandler.LanguageStats

type Profile struct {
	Username       string          `json:"username" example:"ffh"`
	JoinedAt       string          `json:"joinedAt" example:"2020-01-01T00:00:00Z"`
	CompletedTests uint64          `json:"completedTests" example:"1"`
	StartedTests   uint64          `json:"startedTests" example:"1"`
	TimePlayed     int64           `json:"timePlayed" example:"1000"`
	PersonalBests  []PersonalBest  `json:"personalBests"`
	LanguageStats  []LanguageStats `json:"languageStats"`
} //@name UsersUsernameProfileGetHandler.Profile

type ResponseBody struct {
	Profile Profile `json:"profile"`
} //@name UsersUsernameProfileGetHandler.ResponseBody

func newRequest(c *gin.Context) (*Request, error) {
	username := c.Param("username")
	if username == "" {
		return nil, errors.New("username is empty")
	}

	return &Request{
		Username: username,
	}, nil
}

func newProfile(profile *models.Profile) *Profile {
	personalBests := make([]PersonalBest, 0, len(profile.PersonalBests))
	for _, personalBest := range profile.PersonalBests {
		personalBests = append(personalBests, newPersonalBest(&personalBest))
	}

	languageStats := make([]LanguageStats, 0, len(profile.LanguageStats))
	for _, stat := range profile.LanguageStats {
		languageStats = append(languageStats, newLanguageStats(stat))
	}

	return &Profile{
		Username:       profile.Username,
		JoinedAt:       proto.MarshalTime(profile.JoinedAt),
		StartedTests:   profile.StartedTests,
		CompletedTests: profile.CompletedTests,
		TimePlayed:     profile.TimePlayed.Milliseconds(),
		PersonalBests:  personalBests,
		LanguageStats:  languageStats,
	}
}

func newLanguageStats(stats models.LanguageStats) LanguageStats {
	return LanguageStats{
		Language:       string(stats.Language),
		TestsCompleted: stats.TestsCompleted,
		TestsStarted:   stats.TestsStarted,
		TimePlayedMs:   stats.TimePlayed.Milliseconds(),
	}
}

func newPersonalBest(personalBest *models.PersonalBest) PersonalBest {
	return PersonalBest{
		Language: string(personalBest.Language),
		Mode:     string(personalBest.Mode),
		SubMode:  string(personalBest.SubMode),
		WPM:      float64(personalBest.WPM),
		Accuracy: float64(personalBest.Accuracy),
		PlayedAt: proto.MarshalTime(personalBest.PlayedAt),
	}
}

func newResponseBody(profile *models.Profile) *ResponseBody {
	return &ResponseBody{
		Profile: *newProfile(profile),
	}
}

package inbound

import (
	"context"
	"ludo-tournament/core/domain/models"
)

type LeagueService interface {
	CreateLeague(ctx context.Context, name string, organizerID string, settings models.LeagueSettings) (*models.League, error)
	GetLeague(ctx context.Context, id string) (*models.League, error)
	UpdateLeague(ctx context.Context, id string, settings models.LeagueSettings) error
	DeleteLeague(ctx context.Context, id string) error
	GenerateSchedule(ctx context.Context, leagueID string, playDates []string) error
	GeneratePairings(ctx context.Context, leagueID string, playDate string) ([]TablePairing, error)
	ReportLeagueMatch(ctx context.Context, matchID string, results []MatchResult, reportedBy string) error
	GetStandings(ctx context.Context, leagueID string) ([]PlayerStanding, error)
	AddTiebreaker(ctx context.Context, leagueID string, playerIDs []string) error
}

type TablePairing struct {
	MatchID     string   `json:"matchId"`
	PlayDate    string   `json:"playDate"`
	TableNumber int      `json:"tableNumber"`
	PlayerIDs   []string `json:"playerIds"`
}

type PlayerStanding struct {
	PlayerID    string  `json:"playerId"`
	DisplayName string  `json:"displayName"`
	GamesPlayed int     `json:"gamesPlayed"`
	TotalPoints float64 `json:"totalPoints"`
	Wins        int     `json:"wins"`
	Rank        int     `json:"rank"`
}

type MatchResult struct {
	PlayerID  string `json:"playerId"`
	Placement int    `json:"placement"` // 1, 2, 3, or 4
}

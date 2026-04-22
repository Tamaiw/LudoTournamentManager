package inbound

import (
	"context"
	"ludo-tournament/core/domain/models"
)

type TournamentService interface {
	CreateTournament(ctx context.Context, name string, organizerID string, settings models.TournamentSettings) (*models.Tournament, error)
	GetTournament(ctx context.Context, id string) (*models.Tournament, error)
	UpdateTournament(ctx context.Context, id string, settings models.TournamentSettings) error
	DeleteTournament(ctx context.Context, id string) error
	GenerateBracket(ctx context.Context, tournamentID string, playerIDs []string) error
	GetBracket(ctx context.Context, tournamentID string) (*models.KnockoutBracket, error)
	ReportMatch(ctx context.Context, matchID string, results []MatchResult, reportedBy string) error
	GetCurrentRoundPairings(ctx context.Context, tournamentID string) ([]GamePairing, error)
	CanEditGame(ctx context.Context, gameID string) (bool, error)
}

type MatchResult struct {
	PlayerID  string           `json:"playerId"`
	SeatColor models.SeatColor `json:"seatColor"`
	Placement int              `json:"placement"` // 1, 2, 3, or 4
}

type GamePairing struct {
	GameID      string             `json:"gameId"`
	Round       int                `json:"round"`
	TableNumber int                `json:"tableNumber"`
	PlayerIDs   []string           `json:"playerIds"`
	SeatColors  []models.SeatColor `json:"seatColors"`
	Status      models.MatchStatus `json:"status"`
}

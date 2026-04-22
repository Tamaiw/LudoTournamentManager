package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type MatchRepository interface {
	Create(ctx context.Context, match *models.Match) error
	GetByID(ctx context.Context, id string) (*models.Match, error)
	Update(ctx context.Context, match *models.Match) error
	ListByTournament(ctx context.Context, tournamentID string) ([]models.Match, error)
	ListByLeague(ctx context.Context, leagueID string) ([]models.Match, error)
	ListByRound(ctx context.Context, tournamentID string, round int) ([]models.Match, error)
	GetCompletedCountInRound(ctx context.Context, tournamentID string, round int) (int, error)
}

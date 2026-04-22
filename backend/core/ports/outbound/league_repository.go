package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type LeagueRepository interface {
	Create(ctx context.Context, league *models.League) error
	GetByID(ctx context.Context, id string) (*models.League, error)
	Update(ctx context.Context, league *models.League) error
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.League, error)
	ListByStatus(ctx context.Context, status models.LeagueStatus) ([]models.League, error)
}

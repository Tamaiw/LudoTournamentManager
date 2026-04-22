package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type TournamentRepository interface {
	Create(ctx context.Context, tournament *models.Tournament) error
	GetByID(ctx context.Context, id string) (*models.Tournament, error)
	Update(ctx context.Context, tournament *models.Tournament) error
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Tournament, error)
	ListByStatus(ctx context.Context, status models.TournamentStatus) ([]models.Tournament, error)
}

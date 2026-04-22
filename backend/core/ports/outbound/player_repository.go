package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type PlayerRepository interface {
	Create(ctx context.Context, player *models.Player) error
	GetByID(ctx context.Context, id string) (*models.Player, error)
	GetByUserID(ctx context.Context, userID string) (*models.Player, error)
	Update(ctx context.Context, player *models.Player) error
	SoftDelete(ctx context.Context, id string) error
}

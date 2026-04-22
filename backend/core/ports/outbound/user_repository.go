package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.User, error)
	UpdateLastActive(ctx context.Context, id string) error
}

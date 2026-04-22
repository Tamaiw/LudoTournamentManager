package persistence

import (
	"context"

	"gorm.io/gorm"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"
)

type GormPlayerRepository struct {
	db *gorm.DB
}

func NewGormPlayerRepository(db *gorm.DB) *GormPlayerRepository {
	return &GormPlayerRepository{db: db}
}

func (r *GormPlayerRepository) Create(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

func (r *GormPlayerRepository) GetByID(ctx context.Context, id string) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).First(&player, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &player, err
}

func (r *GormPlayerRepository) GetByUserID(ctx context.Context, userID string) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).First(&player, "user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &player, err
}

func (r *GormPlayerRepository) Update(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Save(player).Error
}

func (r *GormPlayerRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Player{}, "id = ?", id).Error
}

var _ outbound.PlayerRepository = (*GormPlayerRepository)(nil)

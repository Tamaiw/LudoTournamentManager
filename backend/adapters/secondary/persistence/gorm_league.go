package persistence

import (
	"context"

	"gorm.io/gorm"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"
)

type GormLeagueRepository struct {
	db *gorm.DB
}

func NewGormLeagueRepository(db *gorm.DB) *GormLeagueRepository {
	return &GormLeagueRepository{db: db}
}

func (r *GormLeagueRepository) Create(ctx context.Context, league *models.League) error {
	return r.db.WithContext(ctx).Create(league).Error
}

func (r *GormLeagueRepository) GetByID(ctx context.Context, id string) (*models.League, error) {
	var league models.League
	err := r.db.WithContext(ctx).First(&league, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &league, err
}

func (r *GormLeagueRepository) Update(ctx context.Context, league *models.League) error {
	return r.db.WithContext(ctx).Save(league).Error
}

func (r *GormLeagueRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.League{}, "id = ?", id).Error
}

func (r *GormLeagueRepository) List(ctx context.Context) ([]models.League, error) {
	var leagues []models.League
	err := r.db.WithContext(ctx).Find(&leagues).Error
	return leagues, err
}

func (r *GormLeagueRepository) ListByStatus(ctx context.Context, status models.LeagueStatus) ([]models.League, error) {
	var leagues []models.League
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&leagues).Error
	return leagues, err
}

var _ outbound.LeagueRepository = (*GormLeagueRepository)(nil)

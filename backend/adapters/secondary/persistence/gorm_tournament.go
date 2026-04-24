package persistence

import (
	"context"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"

	"gorm.io/gorm"
)

type GormTournamentRepository struct {
	db *gorm.DB
}

func NewGormTournamentRepository(db *gorm.DB) *GormTournamentRepository {
	return &GormTournamentRepository{db: db}
}

func (r *GormTournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
	return r.db.WithContext(ctx).Create(tournament).Error
}

func (r *GormTournamentRepository) GetByID(ctx context.Context, id string) (*models.Tournament, error) {
	var tournament models.Tournament
	err := r.db.WithContext(ctx).First(&tournament, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tournament, err
}

func (r *GormTournamentRepository) Update(ctx context.Context, tournament *models.Tournament) error {
	return r.db.WithContext(ctx).Save(tournament).Error
}

func (r *GormTournamentRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Tournament{}, "id = ?", id).Error
}

func (r *GormTournamentRepository) List(ctx context.Context) ([]models.Tournament, error) {
	var tournaments []models.Tournament
	err := r.db.WithContext(ctx).Find(&tournaments).Error
	return tournaments, err
}

func (r *GormTournamentRepository) ListByStatus(ctx context.Context, status models.TournamentStatus) ([]models.Tournament, error) {
	var tournaments []models.Tournament
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&tournaments).Error
	return tournaments, err
}

var _ outbound.TournamentRepository = (*GormTournamentRepository)(nil)

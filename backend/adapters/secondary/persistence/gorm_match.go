package persistence

import (
	"context"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"

	"gorm.io/gorm"
)

type GormMatchRepository struct {
	db *gorm.DB
}

func NewGormMatchRepository(db *gorm.DB) *GormMatchRepository {
	return &GormMatchRepository{db: db}
}

func (r *GormMatchRepository) Create(ctx context.Context, match *models.Match) error {
	return r.db.WithContext(ctx).Create(match).Error
}

func (r *GormMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	var match models.Match
	err := r.db.WithContext(ctx).First(&match, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &match, err
}

func (r *GormMatchRepository) Update(ctx context.Context, match *models.Match) error {
	return r.db.WithContext(ctx).Save(match).Error
}

func (r *GormMatchRepository) ListByTournament(ctx context.Context, tournamentID string) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.WithContext(ctx).Where("tournament_id = ?", tournamentID).Find(&matches).Error
	return matches, err
}

func (r *GormMatchRepository) ListByLeague(ctx context.Context, leagueID string) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.WithContext(ctx).Where("league_id = ?", leagueID).Find(&matches).Error
	return matches, err
}

func (r *GormMatchRepository) ListByRound(ctx context.Context, tournamentID string, round int) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.WithContext(ctx).Where("tournament_id = ? AND round = ?", tournamentID, round).Find(&matches).Error
	return matches, err
}

func (r *GormMatchRepository) GetCompletedCountInRound(ctx context.Context, tournamentID string, round int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Match{}).
		Where("tournament_id = ? AND round = ? AND status = ?", tournamentID, round, models.MatchStatusCompleted).
		Count(&count).Error
	return int(count), err
}

var _ outbound.MatchRepository = (*GormMatchRepository)(nil)

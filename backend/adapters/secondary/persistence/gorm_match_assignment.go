package persistence

import (
	"context"

	"gorm.io/gorm"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"
)

type GormMatchAssignmentRepository struct {
	db *gorm.DB
}

func NewGormMatchAssignmentRepository(db *gorm.DB) *GormMatchAssignmentRepository {
	return &GormMatchAssignmentRepository{db: db}
}

func (r *GormMatchAssignmentRepository) Create(ctx context.Context, assignment *models.MatchAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *GormMatchAssignmentRepository) GetByID(ctx context.Context, id string) (*models.MatchAssignment, error) {
	var assignment models.MatchAssignment
	err := r.db.WithContext(ctx).First(&assignment, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &assignment, err
}

func (r *GormMatchAssignmentRepository) Update(ctx context.Context, assignment *models.MatchAssignment) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}

func (r *GormMatchAssignmentRepository) GetByMatchAndPlayer(ctx context.Context, matchID, playerID string) (*models.MatchAssignment, error) {
	var assignment models.MatchAssignment
	err := r.db.WithContext(ctx).First(&assignment, "match_id = ? AND player_id = ?", matchID, playerID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &assignment, err
}

func (r *GormMatchAssignmentRepository) ListByMatch(ctx context.Context, matchID string) ([]models.MatchAssignment, error) {
	var assignments []models.MatchAssignment
	err := r.db.WithContext(ctx).Where("match_id = ?", matchID).Find(&assignments).Error
	return assignments, err
}

var _ outbound.MatchAssignmentRepository = (*GormMatchAssignmentRepository)(nil)

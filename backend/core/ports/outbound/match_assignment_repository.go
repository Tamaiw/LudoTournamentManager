package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type MatchAssignmentRepository interface {
	Create(ctx context.Context, assignment *models.MatchAssignment) error
	GetByID(ctx context.Context, id string) (*models.MatchAssignment, error)
	Update(ctx context.Context, assignment *models.MatchAssignment) error
	GetByMatchAndPlayer(ctx context.Context, matchID, playerID string) (*models.MatchAssignment, error)
	ListByMatch(ctx context.Context, matchID string) ([]models.MatchAssignment, error)
}

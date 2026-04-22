package outbound

import (
	"context"

	"ludo-tournament/core/domain/models"
)

type InvitationRepository interface {
	Create(ctx context.Context, invitation *models.Invitation) error
	GetByID(ctx context.Context, id string) (*models.Invitation, error)
	Update(ctx context.Context, invitation *models.Invitation) error
	ListByTournament(ctx context.Context, tournamentID string) ([]models.Invitation, error)
	ListByLeague(ctx context.Context, leagueID string) ([]models.Invitation, error)
	ListByInvitee(ctx context.Context, inviteeID string) ([]models.Invitation, error)
}

type UserInviteRepository interface {
	Create(ctx context.Context, invite *models.UserInvite) error
	GetByCode(ctx context.Context, code string) (*models.UserInvite, error)
	GetByEmail(ctx context.Context, email string) (*models.UserInvite, error)
	Update(ctx context.Context, invite *models.UserInvite) error
}

package outbound

import "ludo-tournament/core/domain/models"

type InvitationRepository interface {
	Create(invitation *models.Invitation) error
	GetByID(id string) (*models.Invitation, error)
	Update(invitation *models.Invitation) error
	ListByTournament(tournamentID string) ([]models.Invitation, error)
	ListByLeague(leagueID string) ([]models.Invitation, error)
	ListByInvitee(inviteeID string) ([]models.Invitation, error)
}

type UserInviteRepository interface {
	Create(invite *models.UserInvite) error
	GetByCode(code string) (*models.UserInvite, error)
	GetByEmail(email string) (*models.UserInvite, error)
	Update(invite *models.UserInvite) error
}

package persistence

import (
	"context"

	"gorm.io/gorm"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/outbound"
)

type GormInvitationRepository struct {
	db *gorm.DB
}

func NewGormInvitationRepository(db *gorm.DB) *GormInvitationRepository {
	return &GormInvitationRepository{db: db}
}

func (r *GormInvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *GormInvitationRepository) GetByID(ctx context.Context, id string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.WithContext(ctx).First(&invitation, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &invitation, err
}

func (r *GormInvitationRepository) Update(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *GormInvitationRepository) ListByTournament(ctx context.Context, tournamentID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.WithContext(ctx).Where("tournament_id = ?", tournamentID).Find(&invitations).Error
	return invitations, err
}

func (r *GormInvitationRepository) ListByLeague(ctx context.Context, leagueID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.WithContext(ctx).Where("league_id = ?", leagueID).Find(&invitations).Error
	return invitations, err
}

func (r *GormInvitationRepository) ListByInvitee(ctx context.Context, inviteeID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.WithContext(ctx).Where("invitee_id = ?", inviteeID).Find(&invitations).Error
	return invitations, err
}

var _ outbound.InvitationRepository = (*GormInvitationRepository)(nil)

type GormUserInviteRepository struct {
	db *gorm.DB
}

func NewGormUserInviteRepository(db *gorm.DB) *GormUserInviteRepository {
	return &GormUserInviteRepository{db: db}
}

func (r *GormUserInviteRepository) Create(ctx context.Context, invite *models.UserInvite) error {
	return r.db.WithContext(ctx).Create(invite).Error
}

func (r *GormUserInviteRepository) GetByCode(ctx context.Context, code string) (*models.UserInvite, error) {
	var invite models.UserInvite
	err := r.db.WithContext(ctx).First(&invite, "code = ?", code).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &invite, err
}

func (r *GormUserInviteRepository) GetByEmail(ctx context.Context, email string) (*models.UserInvite, error) {
	var invite models.UserInvite
	err := r.db.WithContext(ctx).First(&invite, "email = ?", email).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &invite, err
}

func (r *GormUserInviteRepository) Update(ctx context.Context, invite *models.UserInvite) error {
	return r.db.WithContext(ctx).Save(invite).Error
}

var _ outbound.UserInviteRepository = (*GormUserInviteRepository)(nil)

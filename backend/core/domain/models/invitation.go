package models

import "time"

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
)

type Invitation struct {
	ID           string           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TournamentID *string          `gorm:"type:uuid" json:"tournament_id,omitempty"`
	LeagueID     *string          `gorm:"type:uuid" json:"league_id,omitempty"`
	InviteeID    *string          `gorm:"type:uuid" json:"invitee_id,omitempty"`
	Status       InvitationStatus `gorm:"not null;default:pending" json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	ModifiedAt   time.Time        `json:"modified_at"`
	ModifiedBy   *string          `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt    *time.Time       `gorm:"index" json:"deleted_at,omitempty"`
}

type UserInvite struct {
	ID         string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email      string     `gorm:"uniqueIndex;not null" json:"email"`
	Code       string     `gorm:"uniqueIndex;not null" json:"code"`
	InvitedBy  *string    `gorm:"type:uuid" json:"invited_by,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	ModifiedBy *string    `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

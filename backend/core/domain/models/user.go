package models

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleGuest  Role = "guest"
)

type User struct {
	ID           string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	Role         Role       `gorm:"not null;default:guest" json:"role"`
	InvitedBy    *string    `gorm:"type:uuid" json:"invited_by,omitempty"`
	LastActive   *time.Time `json:"last_active,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ModifiedAt   time.Time  `json:"modified_at"`
	ModifiedBy   *string    `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

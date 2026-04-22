package models

import "time"

type KnockoutBracket struct {
	ID                string     `gorm:"type:uuid;primary_key" json:"id"`
	TournamentID      string     `gorm:"type:uuid;not null;uniqueIndex" json:"tournamentId"`
	Rounds            string     `gorm:"type:jsonb" json:"rounds"`            // JSON structure: {round_n: [{game_id, player_ids, status}]}
	AdvancementConfig string     `gorm:"type:jsonb" json:"advancementConfig"` // JSON structure per-round advancement rules
	CreatedAt         time.Time  `json:"createdAt"`
	ModifiedAt        time.Time  `json:"modifiedAt"`
	ModifiedBy        *string    `gorm:"type:uuid" json:"modifiedBy,omitempty"`
	DeletedAt         *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

package models

import "time"

type TournamentStatus string

const (
	TournamentStatusDraft    TournamentStatus = "draft"
	TournamentStatusLive     TournamentStatus = "live"
	TournamentStatusComplete TournamentStatus = "completed"
)

type AdvancementConfig struct {
	Round              string               `json:"round"`
	Games              int                  `json:"games"`
	AdvancementPerGame []AdvancementPerGame `json:"advancement_per_game"`
}

type AdvancementPerGame struct {
	GameIDs    []int `json:"game_ids"`
	Placements []int `json:"placements"` // e.g., [1, 2] means 1st and 2nd advance
}

type TournamentSettings struct {
	TablesCount     int                 `json:"tables_count"`
	Advancement     []AdvancementConfig `json:"advancement,omitempty"`
	DefaultReporter string              `json:"default_reporter"` // "lowest_advancing"
}

type Tournament struct {
	ID          string             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string             `gorm:"not null" json:"name"`
	Type        string             `gorm:"not null;default:knockout" json:"type"`
	OrganizerID string             `gorm:"type:uuid;not null" json:"organizer_id"`
	Status      TournamentStatus   `gorm:"not null;default:draft" json:"status"`
	Settings    TournamentSettings `gorm:"type:jsonb" json:"settings,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at"`
	ModifiedBy  *string            `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt   *time.Time         `gorm:"index" json:"deleted_at,omitempty"`
}

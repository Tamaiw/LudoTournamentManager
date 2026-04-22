package models

import "time"

type LeagueStatus string

const (
	LeagueStatusDraft    LeagueStatus = "draft"
	LeagueStatusLive     LeagueStatus = "live"
	LeagueStatusComplete LeagueStatus = "completed"
)

type ScoringRule struct {
	Placement int     `json:"placement"` // 1, 2, 3, or 4
	Points    float64 `json:"points"`
}

type LeagueSettings struct {
	ScoringRules   []ScoringRule `json:"scoring_rules"`
	GamesPerPlayer int           `json:"games_per_player"`
	TablesCount    int           `json:"tables_count"`
}

type League struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	OrganizerID string         `gorm:"type:uuid;not null" json:"organizer_id"`
	Status      LeagueStatus   `gorm:"not null;default:draft" json:"status"`
	Settings    LeagueSettings `gorm:"type:jsonb" json:"settings,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	ModifiedAt  time.Time      `json:"modified_at"`
	ModifiedBy  *string        `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt   *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
}

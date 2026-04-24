package models

import "time"

type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "pending"
	MatchStatusCompleted MatchStatus = "completed"
)

type SeatColor string

const (
	SeatYellow SeatColor = "yellow"
	SeatGreen  SeatColor = "green"
	SeatBlue   SeatColor = "blue"
	SeatRed    SeatColor = "red"
)

type Match struct {
	ID              string      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TournamentID    *string     `gorm:"type:uuid" json:"tournament_id,omitempty"`
	LeagueID        *string     `gorm:"type:uuid" json:"league_id,omitempty"`
	Round           int         `json:"round"`
	TableNumber     int         `json:"table_number"`
	Status          MatchStatus `gorm:"not null;default:pending" json:"status"`
	PlacementPoints []int       `gorm:"type:jsonb" json:"placement_points,omitempty"`
	CompletedAt     *time.Time  `json:"completed_at,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	ModifiedAt      time.Time   `json:"modified_at"`
	ModifiedBy      *string     `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt       *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

type MatchAssignment struct {
	ID           string     `gorm:"type:uuid;primary_key" json:"id"`
	MatchID      string     `gorm:"type:uuid;not null" json:"match_id"`
	PlayerID     string     `gorm:"type:uuid;not null" json:"player_id"`
	SeatColor    SeatColor  `json:"seat_color"`
	Result       *int       `json:"result,omitempty"` // 1, 2, 3, or 4
	SourceGameID *string    `gorm:"type:uuid" json:"source_game_id,omitempty"`
	ReportedBy   *string    `gorm:"type:uuid" json:"reported_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ModifiedAt   time.Time  `json:"modified_at"`
	ModifiedBy   *string    `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Pair represents two players who have been paired together
type Pair struct {
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
}

// LeagueMatchResult represents a player's result in a league match
type LeagueMatchResult struct {
	MatchID   string  `json:"matchId,omitempty"`
	PlayerID  string  `json:"playerId"`
	Placement int     `json:"placement"`
	Points    float64 `json:"points,omitempty"`
}

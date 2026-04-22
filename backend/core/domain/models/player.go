package models

import "time"

type TournamentHistoryEntry struct {
	TournamentID string    `json:"tournamentId"`
	RoundReached string    `json:"roundReached"`
	Date         time.Time `json:"date"`
}

type LeagueStatsEntry struct {
	LeagueID    string  `json:"leagueId"`
	GamesPlayed int     `json:"gamesPlayed"`
	TotalPoints float64 `json:"totalPoints"`
	Wins        int     `json:"wins"`
}

type Player struct {
	ID                string                   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            string                   `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	DisplayName       string                   `gorm:"not null" json:"display_name"`
	TournamentHistory []TournamentHistoryEntry `gorm:"type:jsonb" json:"tournament_history,omitempty"`
	LeagueStats       []LeagueStatsEntry       `gorm:"type:jsonb" json:"league_stats,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	ModifiedAt        time.Time                `json:"modified_at"`
	ModifiedBy        *string                  `gorm:"type:uuid" json:"modified_by,omitempty"`
	DeletedAt         *time.Time               `gorm:"index" json:"deleted_at,omitempty"`
}

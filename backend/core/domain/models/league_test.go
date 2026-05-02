package models

import (
	"encoding/json"
	"testing"
)

func TestLeagueStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   LeagueStatus
		expected string
	}{
		{"LeagueStatusDraft", LeagueStatusDraft, "draft"},
		{"LeagueStatusLive", LeagueStatusLive, "live"},
		{"LeagueStatusComplete", LeagueStatusComplete, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.status)
			}
		})
	}
}

func TestScoringRuleJSONSerialization(t *testing.T) {
	rule := ScoringRule{
		Placement: 1,
		Points:    3.0,
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("failed to marshal scoring rule: %v", err)
	}

	var unmarshaled ScoringRule
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal scoring rule: %v", err)
	}

	if unmarshaled.Placement != rule.Placement {
		t.Errorf("expected Placement %d, got %d", rule.Placement, unmarshaled.Placement)
	}
	if unmarshaled.Points != rule.Points {
		t.Errorf("expected Points %f, got %f", rule.Points, unmarshaled.Points)
	}
}

func TestLeagueSettingsJSONSerialization(t *testing.T) {
	settings := LeagueSettings{
		ScoringRules: []ScoringRule{
			{Placement: 1, Points: 3},
			{Placement: 2, Points: 2},
			{Placement: 3, Points: 1},
			{Placement: 4, Points: 0},
		},
		GamesPerPlayer: 5,
		TablesCount:    4,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal league settings: %v", err)
	}

	var unmarshaled LeagueSettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal league settings: %v", err)
	}

	if len(unmarshaled.ScoringRules) != len(settings.ScoringRules) {
		t.Errorf("expected %d scoring rules, got %d", len(settings.ScoringRules), len(unmarshaled.ScoringRules))
	}
	if unmarshaled.GamesPerPlayer != settings.GamesPerPlayer {
		t.Errorf("expected GamesPerPlayer %d, got %d", settings.GamesPerPlayer, unmarshaled.GamesPerPlayer)
	}
	if unmarshaled.TablesCount != settings.TablesCount {
		t.Errorf("expected TablesCount %d, got %d", settings.TablesCount, unmarshaled.TablesCount)
	}
}

func TestLeagueJSONSerialization(t *testing.T) {
	league := League{
		ID:          "league-uuid",
		Name:        "Test League",
		OrganizerID: "organizer-uuid",
		Status:      LeagueStatusLive,
		Settings: LeagueSettings{
			ScoringRules:   []ScoringRule{{Placement: 1, Points: 3}},
			GamesPerPlayer: 5,
			TablesCount:    4,
		},
	}

	data, err := json.Marshal(league)
	if err != nil {
		t.Fatalf("failed to marshal league: %v", err)
	}

	var unmarshaled League
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal league: %v", err)
	}

	if unmarshaled.ID != league.ID {
		t.Errorf("expected ID %q, got %q", league.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != league.Name {
		t.Errorf("expected Name %q, got %q", league.Name, unmarshaled.Name)
	}
	if unmarshaled.Status != league.Status {
		t.Errorf("expected Status %q, got %q", league.Status, unmarshaled.Status)
	}
}

func TestLeagueStatusTransitions(t *testing.T) {
	validStatuses := []LeagueStatus{LeagueStatusDraft, LeagueStatusLive, LeagueStatusComplete}

	for _, status := range validStatuses {
		if status == "" {
			t.Errorf("status should not be empty")
		}
	}
}
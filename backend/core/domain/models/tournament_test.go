package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTournamentStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   TournamentStatus
		expected string
	}{
		{"TournamentStatusDraft", TournamentStatusDraft, "draft"},
		{"TournamentStatusLive", TournamentStatusLive, "live"},
		{"TournamentStatusComplete", TournamentStatusComplete, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.status)
			}
		})
	}
}

func TestAdvancementPerGameJSON(t *testing.T) {
	apg := AdvancementPerGame{
		GameIDs:    []int{1, 2, 3, 4},
		Placements: []int{1, 2},
	}

	data, err := json.Marshal(apg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled AdvancementPerGame
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(unmarshaled.GameIDs) != len(apg.GameIDs) {
		t.Errorf("expected %d game IDs, got %d", len(apg.GameIDs), len(unmarshaled.GameIDs))
	}
	if len(unmarshaled.Placements) != len(apg.Placements) {
		t.Errorf("expected %d placements, got %d", len(apg.Placements), len(unmarshaled.Placements))
	}
}

func TestAdvancementConfigJSON(t *testing.T) {
	config := AdvancementConfig{
		Round:              "finals",
		Games:              4,
		AdvancementPerGame: []AdvancementPerGame{
			{GameIDs: []int{1, 2}, Placements: []int{1}},
			{GameIDs: []int{3, 4}, Placements: []int{1}},
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled AdvancementConfig
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.Round != config.Round {
		t.Errorf("expected round %q, got %q", config.Round, unmarshaled.Round)
	}
	if unmarshaled.Games != config.Games {
		t.Errorf("expected games %d, got %d", config.Games, unmarshaled.Games)
	}
}

func TestTournamentSettingsJSON(t *testing.T) {
	settings := TournamentSettings{
		TablesCount:     4,
		Advancement:    []AdvancementConfig{{Round: "round1", Games: 2}},
		DefaultReporter: "lowest_advancing",
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled TournamentSettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.TablesCount != settings.TablesCount {
		t.Errorf("expected TablesCount %d, got %d", settings.TablesCount, unmarshaled.TablesCount)
	}
	if unmarshaled.DefaultReporter != settings.DefaultReporter {
		t.Errorf("expected DefaultReporter %q, got %q", settings.DefaultReporter, unmarshaled.DefaultReporter)
	}
}

func TestTournamentJSONSerialization(t *testing.T) {
	tournament := Tournament{
		ID:          "tournament-uuid",
		Name:        "Test Tournament",
		Type:        "knockout",
		OrganizerID: "organizer-uuid",
		Status:      TournamentStatusDraft,
		Settings: TournamentSettings{
			TablesCount:     4,
			DefaultReporter: "lowest_advancing",
		},
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	data, err := json.Marshal(tournament)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Tournament
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.ID != tournament.ID {
		t.Errorf("expected ID %q, got %q", tournament.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != tournament.Name {
		t.Errorf("expected Name %q, got %q", tournament.Name, unmarshaled.Name)
	}
	if unmarshaled.Type != tournament.Type {
		t.Errorf("expected Type %q, got %q", tournament.Type, unmarshaled.Type)
	}
	if unmarshaled.Status != tournament.Status {
		t.Errorf("expected Status %q, got %q", tournament.Status, unmarshaled.Status)
	}
}
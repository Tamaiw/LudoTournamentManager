package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMatchStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   MatchStatus
		expected string
	}{
		{"MatchStatusPending", MatchStatusPending, "pending"},
		{"MatchStatusCompleted", MatchStatusCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.status)
			}
		})
	}
}

func TestSeatColorConstants(t *testing.T) {
	tests := []struct {
		name     string
		color    SeatColor
		expected string
	}{
		{"SeatYellow", SeatYellow, "yellow"},
		{"SeatGreen", SeatGreen, "green"},
		{"SeatBlue", SeatBlue, "blue"},
		{"SeatRed", SeatRed, "red"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.color) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.color)
			}
		})
	}
}

func TestMatchJSONSerialization(t *testing.T) {
	now := time.Now()
	match := Match{
		ID:           "match-uuid",
		TournamentID: strPtr("tournament-uuid"),
		LeagueID:     nil,
		Round:        1,
		TableNumber:  3,
		Status:       MatchStatusPending,
		CreatedAt:    now,
		ModifiedAt:   now,
	}

	data, err := json.Marshal(match)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Match
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.ID != match.ID {
		t.Errorf("expected ID %q, got %q", match.ID, unmarshaled.ID)
	}
	if unmarshaled.Round != match.Round {
		t.Errorf("expected Round %d, got %d", match.Round, unmarshaled.Round)
	}
	if unmarshaled.TableNumber != match.TableNumber {
		t.Errorf("expected TableNumber %d, got %d", match.TableNumber, unmarshaled.TableNumber)
	}
}

func TestMatchAssignmentJSONSerialization(t *testing.T) {
	result := 2
	assignment := MatchAssignment{
		ID:        "assignment-uuid",
		MatchID:   "match-uuid",
		PlayerID:  "player-uuid",
		SeatColor: SeatBlue,
		Result:    &result,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}

	data, err := json.Marshal(assignment)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled MatchAssignment
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.ID != assignment.ID {
		t.Errorf("expected ID %q, got %q", assignment.ID, unmarshaled.ID)
	}
	if unmarshaled.SeatColor != assignment.SeatColor {
		t.Errorf("expected SeatColor %q, got %q", assignment.SeatColor, unmarshaled.SeatColor)
	}
	if *unmarshaled.Result != *assignment.Result {
		t.Errorf("expected Result %d, got %d", *assignment.Result, *unmarshaled.Result)
	}
}

func TestPairJSONSerialization(t *testing.T) {
	pair := Pair{
		Player1: "player-1",
		Player2: "player-2",
	}

	data, err := json.Marshal(pair)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Pair
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.Player1 != pair.Player1 {
		t.Errorf("expected Player1 %q, got %q", pair.Player1, unmarshaled.Player1)
	}
	if unmarshaled.Player2 != pair.Player2 {
		t.Errorf("expected Player2 %q, got %q", pair.Player2, unmarshaled.Player2)
	}
}

func TestLeagueMatchResultJSONSerialization(t *testing.T) {
	result := LeagueMatchResult{
		MatchID:   "match-uuid",
		PlayerID:  "player-uuid",
		Placement: 1,
		Points:    3.0,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled LeagueMatchResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.PlayerID != result.PlayerID {
		t.Errorf("expected PlayerID %q, got %q", result.PlayerID, unmarshaled.PlayerID)
	}
	if unmarshaled.Placement != result.Placement {
		t.Errorf("expected Placement %d, got %d", result.Placement, unmarshaled.Placement)
	}
}

func TestMatchPlacementPointsJSON(t *testing.T) {
	points := []int{3, 2, 1, 0}
	match := Match{
		ID:              "match-uuid",
		PlacementPoints: points,
	}

	data, err := json.Marshal(match)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Match
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(unmarshaled.PlacementPoints) != len(points) {
		t.Errorf("expected %d placement points, got %d", len(points), len(unmarshaled.PlacementPoints))
	}
}

func TestAllSeatColors(t *testing.T) {
	colors := []SeatColor{SeatYellow, SeatGreen, SeatBlue, SeatRed}
	expected := []string{"yellow", "green", "blue", "red"}

	for i, color := range colors {
		if string(color) != expected[i] {
			t.Errorf("expected seat color %q, got %q", expected[i], color)
		}
	}
}

func strPtr(s string) *string {
	return &s
}
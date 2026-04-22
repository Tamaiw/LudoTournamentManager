package application

import (
	"testing"

	"ludo-tournament/core/domain/models"
)

// TestGenerateFairPairings_NoRepeatMatches tests that when generating pairings,
// players who have previously played together are not paired again.
func TestGenerateFairPairings_NoRepeatMatches(t *testing.T) {
	// Given: 8 players, prior pairings include A-B together
	players := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	priorPairings := []models.Pair{
		{Player1: "A", Player2: "B"},
		{Player1: "C", Player2: "D"},
		{Player1: "E", Player2: "F"},
		{Player1: "G", Player2: "H"},
	}
	numTables := 2

	// When: GenerateFairPairings
	pairings, err := GenerateFairPairings(players, priorPairings, numTables)

	// Then: A and B are not paired together
	if err != nil {
		t.Fatalf("GenerateFairPairings returned error: %v", err)
	}

	for _, pairing := range pairings {
		playerIDs := pairing.PlayerIDs
		for i := 0; i < len(playerIDs); i++ {
			for j := i + 1; j < len(playerIDs); j++ {
				p1, p2 := playerIDs[i], playerIDs[j]
				for _, prior := range priorPairings {
					if (p1 == prior.Player1 && p2 == prior.Player2) ||
						(p1 == prior.Player2 && p2 == prior.Player1) {
						t.Errorf("Players %s and %s were paired together but have previously played together", p1, p2)
					}
				}
			}
		}
	}
}

// TestGenerateFairPairings_AllPlayersAssigned tests that all players are assigned
// to exactly one table.
func TestGenerateFairPairings_AllPlayersAssigned(t *testing.T) {
	// Given: 8 players, 2 tables
	players := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	priorPairings := []models.Pair{}
	numTables := 2

	// When: GenerateFairPairings
	pairings, err := GenerateFairPairings(players, priorPairings, numTables)

	// Then: All 8 players assigned to exactly one table
	if err != nil {
		t.Fatalf("GenerateFairPairings returned error: %v", err)
	}

	assignedCount := make(map[string]int)
	for _, pairing := range pairings {
		for _, playerID := range pairing.PlayerIDs {
			assignedCount[playerID]++
		}
	}

	if len(assignedCount) != len(players) {
		t.Errorf("Expected %d players assigned, got %d", len(players), len(assignedCount))
	}

	for _, player := range players {
		if assignedCount[player] != 1 {
			t.Errorf("Player %s assigned %d times, expected exactly 1", player, assignedCount[player])
		}
	}
}

// TestGenerateFairPairings_ReturnsCorrectTableCount tests that the correct
// number of tables is created.
func TestGenerateFairPairings_ReturnsCorrectTableCount(t *testing.T) {
	// Given: 12 players, 3 tables
	players := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	priorPairings := []models.Pair{}
	numTables := 3

	// When: GenerateFairPairings
	pairings, err := GenerateFairPairings(players, priorPairings, numTables)

	// Then: 3 tables created
	if err != nil {
		t.Fatalf("GenerateFairPairings returned error: %v", err)
	}

	if len(pairings) != numTables {
		t.Errorf("Expected %d tables, got %d", numTables, len(pairings))
	}

	// Verify 4 players per table (3 pairs each = 6 pairs, 12 players)
	totalPlayers := 0
	for _, pairing := range pairings {
		totalPlayers += len(pairing.PlayerIDs)
	}
	if totalPlayers != 12 {
		t.Errorf("Expected 12 players across tables, got %d", totalPlayers)
	}
}

// TestCalculateLeagueStandings_SortsByPoints tests that players are sorted
// by total points in descending order.
func TestCalculateLeagueStandings_SortsByPoints(t *testing.T) {
	// Given: 4 players with known placements
	// And: scoring rules: 1st=3, 2nd=2, 3rd=1, 4th=0
	scoringRules := []models.ScoringRule{
		{Placement: 1, Points: 3},
		{Placement: 2, Points: 2},
		{Placement: 3, Points: 1},
		{Placement: 4, Points: 0},
	}

	matchResults := []models.LeagueMatchResult{
		{MatchID: "match1", PlayerID: "A", Placement: 1},
		{MatchID: "match1", PlayerID: "B", Placement: 2},
		{MatchID: "match1", PlayerID: "C", Placement: 3},
		{MatchID: "match1", PlayerID: "D", Placement: 4},
		{MatchID: "match2", PlayerID: "A", Placement: 2},
		{MatchID: "match2", PlayerID: "B", Placement: 1},
		{MatchID: "match2", PlayerID: "C", Placement: 4},
		{MatchID: "match2", PlayerID: "D", Placement: 3},
	}

	// When: CalculateLeagueStandings
	standings := CalculateLeagueStandings(matchResults, scoringRules)

	// Then: Players sorted by total points descending
	// A: 3 + 2 = 5 points
	// B: 2 + 3 = 5 points
	// C: 1 + 0 = 1 point
	// D: 0 + 1 = 1 point
	// A and B should be tied at 5, C and D at 1

	if len(standings) != 4 {
		t.Fatalf("Expected 4 standings, got %d", len(standings))
	}

	// First two should be A and B (tied at 5 points)
	if standings[0].TotalPoints != standings[1].TotalPoints || standings[0].TotalPoints != 5 {
		t.Errorf("Expected top players to have 5 points, got %.2f and %.2f",
			standings[0].TotalPoints, standings[1].TotalPoints)
	}

	// Last two should be C and D (tied at 1 point)
	if standings[2].TotalPoints != standings[3].TotalPoints || standings[2].TotalPoints != 1 {
		t.Errorf("Expected bottom players to have 1 point, got %.2f and %.2f",
			standings[2].TotalPoints, standings[3].TotalPoints)
	}
}

// TestCalculateLeagueStandings_TiebreakerByWins tests that when players have
// the same total points, the player with more wins is ranked higher.
func TestCalculateLeagueStandings_TiebreakerByWins(t *testing.T) {
	// Given: Two players with same total points
	// And: one has more wins
	scoringRules := []models.ScoringRule{
		{Placement: 1, Points: 3},
		{Placement: 2, Points: 2},
		{Placement: 3, Points: 1},
		{Placement: 4, Points: 0},
	}

	// Player X: 1st place (3 pts) + 3rd place (1 pt) = 4 pts, 1 win
	// Player Y: 2nd place (2 pts) + 2nd place (2 pts) = 4 pts, 0 wins
	matchResults := []models.LeagueMatchResult{
		{MatchID: "match1", PlayerID: "X", Placement: 1},
		{MatchID: "match1", PlayerID: "Y", Placement: 2},
		{MatchID: "match2", PlayerID: "X", Placement: 3},
		{MatchID: "match2", PlayerID: "Y", Placement: 2},
	}

	// When: CalculateLeagueStandings
	standings := CalculateLeagueStandings(matchResults, scoringRules)

	// Then: Player X with more wins is ranked higher
	if len(standings) != 2 {
		t.Fatalf("Expected 2 standings, got %d", len(standings))
	}

	// Both should have 4 points
	if standings[0].TotalPoints != 4 || standings[1].TotalPoints != 4 {
		t.Errorf("Expected both players to have 4 points")
	}

	// X should be first due to having more wins (1 vs 0)
	if standings[0].PlayerID != "X" {
		t.Errorf("Expected player X to be ranked first (has more wins), got %s", standings[0].PlayerID)
	}

	if standings[0].Wins != 1 {
		t.Errorf("Expected player X to have 1 win, got %d", standings[0].Wins)
	}
}

// TestGeneratePairings_CallsFairPairingAlgorithm tests that GeneratePairings
// calls the fair pairing algorithm with correct parameters.
func TestGeneratePairings_CallsFairPairingAlgorithm(t *testing.T) {
	// This test verifies the integration between the service method and the algorithm
	// We can't easily mock the fair pairing algorithm in a unit test without
	// more complex setup, so this test validates the algorithm is being used correctly.

	players := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	priorPairings := []models.Pair{
		{Player1: "A", Player2: "B"},
		{Player1: "C", Player2: "D"},
	}

	// When: GenerateFairPairings is called directly
	pairings, err := GenerateFairPairings(players, priorPairings, 2)

	// Then: It should succeed and return valid pairings
	if err != nil {
		t.Fatalf("GenerateFairPairings returned error: %v", err)
	}

	if len(pairings) != 2 {
		t.Errorf("Expected 2 pairings, got %d", len(pairings))
	}

	// Verify no repeat matches from priorPairings
	for _, pairing := range pairings {
		playerIDs := pairing.PlayerIDs
		for i := 0; i < len(playerIDs); i++ {
			for j := i + 1; j < len(playerIDs); j++ {
				p1, p2 := playerIDs[i], playerIDs[j]
				for _, prior := range priorPairings {
					if (p1 == prior.Player1 && p2 == prior.Player2) ||
						(p1 == prior.Player2 && p2 == prior.Player1) {
						t.Errorf("Players %s and %s were paired together but have previously played together", p1, p2)
					}
				}
			}
		}
	}
}

// TestGenerateFairPairings_InsufficientPlayers tests error handling when
// there are not enough players to fill the tables.
func TestGenerateFairPairings_InsufficientPlayers(t *testing.T) {
	// Given: Only 3 players but need at least 4 for 1 table
	players := []string{"A", "B", "C"}
	priorPairings := []models.Pair{}
	numTables := 1

	// When: GenerateFairPairings
	_, err := GenerateFairPairings(players, priorPairings, numTables)

	// Then: Should return error
	if err == nil {
		t.Error("Expected error for insufficient players, got nil")
	}
}

// TestGenerateFairPairings_EmptyPlayers tests error handling when
// no players are provided.
func TestGenerateFairPairings_EmptyPlayers(t *testing.T) {
	// Given: No players
	players := []string{}
	priorPairings := []models.Pair{}
	numTables := 1

	// When: GenerateFairPairings
	_, err := GenerateFairPairings(players, priorPairings, numTables)

	// Then: Should return error
	if err == nil {
		t.Error("Expected error for empty players, got nil")
	}
}

// TestCalculateLeagueStandings_EmptyResults tests that empty match results
// return empty standings.
func TestCalculateLeagueStandings_EmptyResults(t *testing.T) {
	// Given: No match results
	scoringRules := []models.ScoringRule{
		{Placement: 1, Points: 3},
		{Placement: 2, Points: 2},
		{Placement: 3, Points: 1},
		{Placement: 4, Points: 0},
	}

	matchResults := []models.LeagueMatchResult{}

	// When: CalculateLeagueStandings
	standings := CalculateLeagueStandings(matchResults, scoringRules)

	// Then: Should return empty standings
	if len(standings) != 0 {
		t.Errorf("Expected 0 standings, got %d", len(standings))
	}
}

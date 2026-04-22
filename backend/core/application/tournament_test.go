package application

import (
	"context"
	"testing"
	"time"

	"ludo-tournament/core/domain"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"
)

// Mock implementations

type mockTournamentRepository struct {
	tournaments map[string]*models.Tournament
}

func (m *mockTournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
	if m.tournaments == nil {
		m.tournaments = make(map[string]*models.Tournament)
	}
	tournament.CreatedAt = time.Now()
	tournament.ModifiedAt = time.Now()
	m.tournaments[tournament.ID] = tournament
	return nil
}

func (m *mockTournamentRepository) GetByID(ctx context.Context, id string) (*models.Tournament, error) {
	if t, ok := m.tournaments[id]; ok {
		return t, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockTournamentRepository) Update(ctx context.Context, tournament *models.Tournament) error {
	tournament.ModifiedAt = time.Now()
	m.tournaments[tournament.ID] = tournament
	return nil
}

func (m *mockTournamentRepository) SoftDelete(ctx context.Context, id string) error {
	delete(m.tournaments, id)
	return nil
}

func (m *mockTournamentRepository) List(ctx context.Context) ([]models.Tournament, error) {
	var result []models.Tournament
	for _, t := range m.tournaments {
		result = append(result, *t)
	}
	return result, nil
}

func (m *mockTournamentRepository) ListByStatus(ctx context.Context, status models.TournamentStatus) ([]models.Tournament, error) {
	var result []models.Tournament
	for _, t := range m.tournaments {
		if t.Status == status {
			result = append(result, *t)
		}
	}
	return result, nil
}

type mockMatchRepository struct {
	matches map[string]*models.Match
}

func (m *mockMatchRepository) Create(ctx context.Context, match *models.Match) error {
	if m.matches == nil {
		m.matches = make(map[string]*models.Match)
	}
	match.CreatedAt = time.Now()
	match.ModifiedAt = time.Now()
	m.matches[match.ID] = match
	return nil
}

func (m *mockMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	if mt, ok := m.matches[id]; ok {
		return mt, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockMatchRepository) Update(ctx context.Context, match *models.Match) error {
	match.ModifiedAt = time.Now()
	m.matches[match.ID] = match
	return nil
}

func (m *mockMatchRepository) ListByTournament(ctx context.Context, tournamentID string) ([]models.Match, error) {
	var result []models.Match
	for _, mt := range m.matches {
		if mt.TournamentID != nil && *mt.TournamentID == tournamentID {
			result = append(result, *mt)
		}
	}
	return result, nil
}

func (m *mockMatchRepository) ListByLeague(ctx context.Context, leagueID string) ([]models.Match, error) {
	var result []models.Match
	for _, mt := range m.matches {
		if mt.LeagueID != nil && *mt.LeagueID == leagueID {
			result = append(result, *mt)
		}
	}
	return result, nil
}

func (m *mockMatchRepository) ListByRound(ctx context.Context, tournamentID string, round int) ([]models.Match, error) {
	var result []models.Match
	for _, mt := range m.matches {
		if mt.TournamentID != nil && *mt.TournamentID == tournamentID && mt.Round == round {
			result = append(result, *mt)
		}
	}
	return result, nil
}

func (m *mockMatchRepository) GetCompletedCountInRound(ctx context.Context, tournamentID string, round int) (int, error) {
	count := 0
	for _, mt := range m.matches {
		if mt.TournamentID != nil && *mt.TournamentID == tournamentID && mt.Round == round && mt.Status == models.MatchStatusCompleted {
			count++
		}
	}
	return count, nil
}

type mockMatchAssignmentRepository struct {
	assignments map[string]*models.MatchAssignment
}

func (m *mockMatchAssignmentRepository) Create(ctx context.Context, assignment *models.MatchAssignment) error {
	if m.assignments == nil {
		m.assignments = make(map[string]*models.MatchAssignment)
	}
	assignment.CreatedAt = time.Now()
	assignment.ModifiedAt = time.Now()
	m.assignments[assignment.ID] = assignment
	return nil
}

func (m *mockMatchAssignmentRepository) GetByID(ctx context.Context, id string) (*models.MatchAssignment, error) {
	if a, ok := m.assignments[id]; ok {
		return a, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockMatchAssignmentRepository) Update(ctx context.Context, assignment *models.MatchAssignment) error {
	assignment.ModifiedAt = time.Now()
	m.assignments[assignment.ID] = assignment
	return nil
}

func (m *mockMatchAssignmentRepository) GetByMatchAndPlayer(ctx context.Context, matchID, playerID string) (*models.MatchAssignment, error) {
	for _, a := range m.assignments {
		if a.MatchID == matchID && a.PlayerID == playerID {
			return a, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockMatchAssignmentRepository) ListByMatch(ctx context.Context, matchID string) ([]models.MatchAssignment, error) {
	var result []models.MatchAssignment
	for _, a := range m.assignments {
		if a.MatchID == matchID {
			result = append(result, *a)
		}
	}
	return result, nil
}

// Tests

func TestGenerateBracket_SplitsPlayersIntoGamesOfFour(t *testing.T) {
	// Given: 8 players
	playerIDs := []string{"p1", "p2", "p3", "p4", "p5", "p6", "p7", "p8"}
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}

	// Create a tournament first
	tournament := &models.Tournament{
		ID:   "tournament-1",
		Name: "Test Tournament",
	}
	_ = tournamentRepo.Create(context.Background(), tournament)

	svc := NewTournamentService(tournamentRepo, matchRepo, &mockMatchAssignmentRepository{})

	// When: GenerateBracket
	err := svc.GenerateBracket(context.Background(), "tournament-1", playerIDs)

	// Then: 2 matches created with 4 players each
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	matches, _ := matchRepo.ListByTournament(context.Background(), "tournament-1")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	for _, match := range matches {
		if match.Round != 1 {
			t.Errorf("expected round 1, got %d", match.Round)
		}
	}
}

func TestGenerateBracket_RejectsInvalidPlayerCount(t *testing.T) {
	// Given: 5 players (not divisible by 4 without remainder)
	playerIDs := []string{"p1", "p2", "p3", "p4", "p5"}
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}

	// Create a tournament first
	tournament := &models.Tournament{
		ID:   "tournament-1",
		Name: "Test Tournament",
	}
	_ = tournamentRepo.Create(context.Background(), tournament)

	svc := NewTournamentService(tournamentRepo, matchRepo, &mockMatchAssignmentRepository{})

	// When: GenerateBracket with odd player count
	err := svc.GenerateBracket(context.Background(), "tournament-1", playerIDs)

	// Then: error returned
	if err == nil {
		t.Fatal("expected error for invalid player count")
	}
}

func TestValidateAdvancementConfig_ValidConfig(t *testing.T) {
	// Given: config where 6 players advance (3 games * 2 each)
	// And: next round has 6 spots (2 games * 3)
	config := []models.AdvancementConfig{
		{
			Round: "round_1",
			Games: 3,
			AdvancementPerGame: []models.AdvancementPerGame{
				{GameIDs: []int{1, 2, 3}, Placements: []int{1, 2}},
				{GameIDs: []int{4, 5, 6}, Placements: []int{1, 2}},
				{GameIDs: []int{7, 8, 9}, Placements: []int{1, 2}},
			},
		},
	}

	// When: ValidateAdvancementConfig
	err := ValidateAdvancementConfig(config, 6)

	// Then: no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateAdvancementConfig_InvalidConfig(t *testing.T) {
	// Given: config where 15 players advance
	// And: next round has only 12 spots (3 games * 4)
	config := []models.AdvancementConfig{
		{
			Round: "round_1",
			Games: 3,
			AdvancementPerGame: []models.AdvancementPerGame{
				{GameIDs: []int{1, 2, 3}, Placements: []int{1, 2, 3}},    // 3 advancing
				{GameIDs: []int{4, 5, 6}, Placements: []int{1, 2, 3}},    // 3 advancing
				{GameIDs: []int{7, 8, 9}, Placements: []int{1, 2, 3}},    // 3 advancing
				{GameIDs: []int{10, 11, 12}, Placements: []int{1, 2, 3}}, // 3 advancing
				{GameIDs: []int{13, 14, 15}, Placements: []int{1, 2, 3}}, // 3 advancing
			},
		},
	}

	// When: ValidateAdvancementConfig with 15 advancing but only 12 spots
	err := ValidateAdvancementConfig(config, 12)

	// Then: ErrInvalidAdvancement
	if err != domain.ErrInvalidAdvancement {
		t.Fatalf("expected ErrInvalidAdvancement, got %v", err)
	}
}

func TestCanEditGame_ReturnsTrueWhenNoDownstreamGames(t *testing.T) {
	// Given: Match in round 1, no matches in round 2 completed
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}

	// Create tournament
	tournament := &models.Tournament{
		ID:   "tournament-1",
		Name: "Test Tournament",
	}
	_ = tournamentRepo.Create(context.Background(), tournament)

	// Create match in round 1
	match1 := &models.Match{
		ID:           "match-1",
		TournamentID: strPtr("tournament-1"),
		Round:        1,
		TableNumber:  1,
		Status:       models.MatchStatusPending,
	}
	_ = matchRepo.Create(context.Background(), match1)

	svc := NewTournamentService(tournamentRepo, matchRepo, &mockMatchAssignmentRepository{})

	// When: CanEditGame for match-1
	canEdit, err := svc.CanEditGame(context.Background(), "match-1")

	// Then: returns true
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !canEdit {
		t.Error("expected CanEditGame to return true when no downstream games played")
	}
}

func TestCanEditGame_ReturnsFalseWhenDownstreamGamesPlayed(t *testing.T) {
	// Given: Match in round 1, a match in round 2 is completed
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}

	// Create tournament
	tournament := &models.Tournament{
		ID:   "tournament-1",
		Name: "Test Tournament",
	}
	_ = tournamentRepo.Create(context.Background(), tournament)

	// Create match in round 1
	match1 := &models.Match{
		ID:           "match-1",
		TournamentID: strPtr("tournament-1"),
		Round:        1,
		TableNumber:  1,
		Status:       models.MatchStatusPending,
	}
	_ = matchRepo.Create(context.Background(), match1)

	// Create match in round 2 that is completed
	completedTime := time.Now()
	match2 := &models.Match{
		ID:           "match-2",
		TournamentID: strPtr("tournament-1"),
		Round:        2,
		TableNumber:  1,
		Status:       models.MatchStatusCompleted,
		CompletedAt:  &completedTime,
	}
	_ = matchRepo.Create(context.Background(), match2)

	svc := NewTournamentService(tournamentRepo, matchRepo, &mockMatchAssignmentRepository{})

	// When: CanEditGame for match-1
	canEdit, err := svc.CanEditGame(context.Background(), "match-1")

	// Then: returns false
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canEdit {
		t.Error("expected CanEditGame to return false when downstream games played")
	}
}

func TestReportMatch_RecordsResults(t *testing.T) {
	// Given: A match with 4 players assigned
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}
	assignmentRepo := &mockMatchAssignmentRepository{assignments: make(map[string]*models.MatchAssignment)}

	// Create tournament
	tournament := &models.Tournament{
		ID:   "tournament-1",
		Name: "Test Tournament",
	}
	_ = tournamentRepo.Create(context.Background(), tournament)

	// Create match
	match := &models.Match{
		ID:           "match-1",
		TournamentID: strPtr("tournament-1"),
		Round:        1,
		TableNumber:  1,
		Status:       models.MatchStatusPending,
	}
	_ = matchRepo.Create(context.Background(), match)

	// Create match assignments for the 4 players
	seatColors := []models.SeatColor{models.SeatYellow, models.SeatGreen, models.SeatBlue, models.SeatRed}
	players := []string{"p1", "p2", "p3", "p4"}
	for i, playerID := range players {
		assignment := &models.MatchAssignment{
			ID:        playerID + "-assignment",
			MatchID:   "match-1",
			PlayerID:  playerID,
			SeatColor: seatColors[i],
		}
		_ = assignmentRepo.Create(context.Background(), assignment)
	}

	svc := NewTournamentService(tournamentRepo, matchRepo, assignmentRepo)

	// When: ReportMatch with results
	results := []inbound.MatchResult{
		{PlayerID: "p1", SeatColor: models.SeatYellow, Placement: 1},
		{PlayerID: "p2", SeatColor: models.SeatGreen, Placement: 2},
		{PlayerID: "p3", SeatColor: models.SeatBlue, Placement: 3},
		{PlayerID: "p4", SeatColor: models.SeatRed, Placement: 4},
	}
	err := svc.ReportMatch(context.Background(), "match-1", results, "reporter-1")

	// Then: MatchAssignment records updated with results
	// And: Match status set to completed
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedMatch, _ := matchRepo.GetByID(context.Background(), "match-1")
	if updatedMatch.Status != models.MatchStatusCompleted {
		t.Errorf("expected match status to be completed, got %v", updatedMatch.Status)
	}
	if updatedMatch.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}

	// Verify assignments were updated with results
	for _, result := range results {
		assignment, _ := assignmentRepo.GetByMatchAndPlayer(context.Background(), "match-1", result.PlayerID)
		if assignment.Result == nil {
			t.Errorf("expected assignment for player %s to have result set", result.PlayerID)
		} else if *assignment.Result != result.Placement {
			t.Errorf("expected placement %d for player %s, got %d", result.Placement, result.PlayerID, *assignment.Result)
		}
	}
}

func TestReportMatch_ValidatesResultsCount(t *testing.T) {
	// Given: A match with 4 players assigned
	tournamentRepo := &mockTournamentRepository{tournaments: make(map[string]*models.Tournament)}
	matchRepo := &mockMatchRepository{matches: make(map[string]*models.Match)}

	// Create match
	match := &models.Match{
		ID:           "match-1",
		TournamentID: strPtr("tournament-1"),
		Round:        1,
		TableNumber:  1,
		Status:       models.MatchStatusPending,
	}
	_ = matchRepo.Create(context.Background(), match)

	svc := NewTournamentService(tournamentRepo, matchRepo, &mockMatchAssignmentRepository{})

	// When: ReportMatch with only 2 results (invalid)
	results := []inbound.MatchResult{
		{PlayerID: "p1", SeatColor: models.SeatYellow, Placement: 1},
		{PlayerID: "p2", SeatColor: models.SeatGreen, Placement: 2},
	}

	// Then: error returned
	err := svc.ReportMatch(context.Background(), "match-1", results, "reporter-1")
	if err == nil {
		t.Fatal("expected error for invalid results count")
	}
}

// Helper
func strPtr(s string) *string {
	return &s
}

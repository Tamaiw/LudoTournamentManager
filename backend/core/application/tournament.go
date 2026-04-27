package application

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"ludo-tournament/core/domain"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"
	"ludo-tournament/core/ports/outbound"
)

// TournamentService implements inbound.TournamentService
type TournamentService struct {
	tournamentRepo outbound.TournamentRepository
	matchRepo      outbound.MatchRepository
	assignmentRepo outbound.MatchAssignmentRepository
	uuidCounter    int
}

// NewTournamentService creates a new TournamentService
func NewTournamentService(tournamentRepo outbound.TournamentRepository, matchRepo outbound.MatchRepository, assignmentRepo outbound.MatchAssignmentRepository) *TournamentService {
	return &TournamentService{
		tournamentRepo: tournamentRepo,
		matchRepo:      matchRepo,
		assignmentRepo: assignmentRepo,
		uuidCounter:    0,
	}
}

// Compile time check that TournamentService implements inbound.TournamentService
var _ inbound.TournamentService = (*TournamentService)(nil)

// CreateTournament creates a new tournament
func (s *TournamentService) CreateTournament(ctx context.Context, name string, organizerID string, settings models.TournamentSettings) (*models.Tournament, error) {
	tournament := &models.Tournament{
		ID:          s.generateUUID(),
		Name:        name,
		Type:        "knockout",
		OrganizerID: organizerID,
		Status:      models.TournamentStatusDraft,
		Settings:    settings,
	}

	if err := s.tournamentRepo.Create(ctx, tournament); err != nil {
		return nil, err
	}

	return tournament, nil
}

// GetTournament retrieves a tournament by ID
func (s *TournamentService) GetTournament(ctx context.Context, id string) (*models.Tournament, error) {
	return s.tournamentRepo.GetByID(ctx, id)
}

// UpdateTournament updates tournament settings
func (s *TournamentService) UpdateTournament(ctx context.Context, id string, settings models.TournamentSettings) error {
	tournament, err := s.tournamentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if tournament.Status != models.TournamentStatusDraft {
		return domain.ErrTournamentActive
	}

	tournament.Settings = settings
	return s.tournamentRepo.Update(ctx, tournament)
}

// DeleteTournament soft deletes a tournament
func (s *TournamentService) DeleteTournament(ctx context.Context, id string) error {
	return s.tournamentRepo.SoftDelete(ctx, id)
}

// GenerateBracket creates the initial bracket for a tournament
func (s *TournamentService) GenerateBracket(ctx context.Context, tournamentID string, playerIDs []string) error {
	// Validate player count - must be divisible by 4 for standard knockout
	if len(playerIDs) == 0 || len(playerIDs)%4 != 0 {
		return domain.ErrInvalidInput
	}

	// Shuffle players using Fisher-Yates
	shuffled := make([]string, len(playerIDs))
	copy(shuffled, playerIDs)
	shuffleStrings(shuffled)

	// Get tournament
	tournament, err := s.tournamentRepo.GetByID(ctx, tournamentID)
	if err != nil {
		return err
	}

	// Calculate number of games (groups of 4)
	numGames := len(playerIDs) / 4

	seatColors := []models.SeatColor{models.SeatYellow, models.SeatGreen, models.SeatBlue, models.SeatRed}

	// Create matches for round 1 - each match gets 4 players
	matches := make([]*models.Match, numGames)
	for i := 0; i < numGames; i++ {
		match := &models.Match{
			ID:           s.generateUUID(),
			TournamentID: &tournamentID,
			Round:        1,
			TableNumber:  i + 1,
			Status:       models.MatchStatusPending,
		}
		if err := s.matchRepo.Create(ctx, match); err != nil {
			return err
		}
		matches[i] = match
	}

	// Assign players to matches round-robin
	playerIndex := 0
	for _, match := range matches {
		for seatNum := 0; seatNum < 4 && playerIndex < len(shuffled); seatNum++ {
			assignment := &models.MatchAssignment{
				ID:        s.generateUUID(),
				MatchID:   match.ID,
				PlayerID:  shuffled[playerIndex],
				SeatColor: seatColors[seatNum],
			}
			if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
				return err
			}
			playerIndex++
		}
	}

	// Update tournament status to live
	tournament.Status = models.TournamentStatusLive
	return s.tournamentRepo.Update(ctx, tournament)
}

// GetBracket returns the knockout bracket for a tournament
func (s *TournamentService) GetBracket(ctx context.Context, tournamentID string) (*models.KnockoutBracket, error) {
	matches, err := s.matchRepo.ListByTournament(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	roundsJSON, err := buildRoundsJSON(ctx, matches, s.assignmentRepo)
	if err != nil {
		return nil, err
	}

	bracket := &models.KnockoutBracket{
		ID:           s.generateUUID(),
		TournamentID: tournamentID,
		Rounds:       roundsJSON,
	}

	return bracket, nil
}

// ReportMatch records the result of a completed game
func (s *TournamentService) ReportMatch(ctx context.Context, matchID string, results []inbound.MatchResult, reportedBy string) error {
	// Validate results has exactly 4 entries (or fewer if forfeits allowed)
	// For now, require exactly 4 results
	if len(results) != 4 {
		return domain.ErrInvalidInput
	}

	// Get the match
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	if match.Status == models.MatchStatusCompleted {
		return domain.ErrGameAlreadyPlayed
	}

	// Update match assignments with results
	for _, result := range results {
		placement := result.Placement // Create a copy for pointer stability
		assignment, err := s.assignmentRepo.GetByMatchAndPlayer(ctx, matchID, result.PlayerID)
		if err != nil {
			return err
		}
		assignment.Result = &placement
		assignment.SeatColor = result.SeatColor
		assignment.ReportedBy = &reportedBy
		if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
			return err
		}
	}

	// Update match status to completed
	now := time.Now()
	match.Status = models.MatchStatusCompleted
	match.CompletedAt = &now

	return s.matchRepo.Update(ctx, match)
}

// GetCurrentRoundPairings returns the pairings for the current active round
func (s *TournamentService) GetCurrentRoundPairings(ctx context.Context, tournamentID string) ([]inbound.GamePairing, error) {
	tournament, err := s.tournamentRepo.GetByID(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	if tournament.Status != models.TournamentStatusLive {
		return nil, domain.ErrTournamentActive
	}

	matches, err := s.matchRepo.ListByTournament(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	var pairings []inbound.GamePairing
	for _, match := range matches {
		if match.Status == models.MatchStatusPending {
			pairing := inbound.GamePairing{
				GameID:      match.ID,
				Round:       match.Round,
				TableNumber: match.TableNumber,
				Status:      match.Status,
			}
			pairings = append(pairings, pairing)
		}
	}

	return pairings, nil
}

// CanEditGame checks if a game can still be edited
func (s *TournamentService) CanEditGame(ctx context.Context, gameID string) (bool, error) {
	// Get the match
	match, err := s.matchRepo.GetByID(ctx, gameID)
	if err != nil {
		return false, err
	}

	// Check if any match in higher rounds is completed
	// Round N can be edited as long as no games in round N+1 or beyond are completed
	currentRound := match.Round

	// Get all matches for this tournament
	matches, err := s.matchRepo.ListByTournament(ctx, *match.TournamentID)
	if err != nil {
		return false, err
	}

	for _, m := range matches {
		if m.Round > currentRound && m.Status == models.MatchStatusCompleted {
			return false, nil // Cannot edit - downstream games already played
		}
	}

	return true, nil
}

// ValidateAdvancementConfig validates that total advancing players from round N
// equals available spots in round N+1
func ValidateAdvancementConfig(config []models.AdvancementConfig, nextRoundSpots int) error {
	// Count total advancing players
	totalAdvancing := 0
	for _, adv := range config {
		for _, apg := range adv.AdvancementPerGame {
			totalAdvancing += len(apg.Placements)
		}
	}

	// Compare total advancing to available spots in next round
	if totalAdvancing != nextRoundSpots {
		return domain.ErrInvalidAdvancement
	}

	return nil
}

// shuffleStrings performs Fisher-Yates shuffle
func shuffleStrings(arr []string) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// buildRoundsJSON builds the JSON string for bracket rounds from matches
func buildRoundsJSON(ctx context.Context, matches []models.Match, assignmentRepo outbound.MatchAssignmentRepository) (string, error) {
	// Group matches by round
	roundsMap := make(map[int][]models.Match)
	for _, match := range matches {
		roundsMap[match.Round] = append(roundsMap[match.Round], match)
	}

	// Build the structure
	type matchEntry struct {
		GameID      string             `json:"game_id"`
		PlayerIDs   []string           `json:"player_ids"`
		Status      models.MatchStatus `json:"status"`
		TableNumber int                `json:"table_number"`
	}

	type roundEntry struct {
		Matches []matchEntry `json:"matches"`
	}

	result := make(map[string]roundEntry)
	for round, roundMatches := range roundsMap {
		entries := make([]matchEntry, len(roundMatches))
		for i, match := range roundMatches {
			// Get player assignments for this match
			assignments, err := assignmentRepo.ListByMatch(ctx, match.ID)
			if err != nil {
				return "", err
			}
			playerIDs := make([]string, len(assignments))
			for j, a := range assignments {
				playerIDs[j] = a.PlayerID
			}
			entries[i] = matchEntry{
				GameID:      match.ID,
				PlayerIDs:   playerIDs,
				Status:      match.Status,
				TableNumber: match.TableNumber,
			}
		}
		result[fmt.Sprintf("round_%d", round)] = roundEntry{Matches: entries}
	}

	// Marshal to JSON
	roundsJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(roundsJSON), nil
}

// generateUUID is a placeholder - in real implementation use uuid package
func (s *TournamentService) generateUUID() string {
	s.uuidCounter++
	return fmt.Sprintf("uuid-%d-%d", time.Now().UnixNano(), s.uuidCounter)
}

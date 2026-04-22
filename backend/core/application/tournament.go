package application

import (
	"context"
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
	playerRepo     outbound.PlayerRepository
	uuidCounter    int
}

// NewTournamentService creates a new TournamentService
func NewTournamentService(tournamentRepo outbound.TournamentRepository, matchRepo outbound.MatchRepository, playerRepo outbound.PlayerRepository) *TournamentService {
	return &TournamentService{
		tournamentRepo: tournamentRepo,
		matchRepo:      matchRepo,
		playerRepo:     playerRepo,
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

	// Create matches for round 1 - each match gets 4 players
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

	bracket := &models.KnockoutBracket{
		ID:           s.generateUUID(),
		TournamentID: tournamentID,
	}

	// Build rounds structure from matches
	_ = matches // Currently returning empty structure

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

	// Update match status to completed
	now := time.Now()
	match.Status = models.MatchStatusCompleted
	match.CompletedAt = &now

	// Update match assignments with results
	for _, result := range results {
		_ = result // In real implementation, would update MatchAssignment records
		_ = reportedBy
	}

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
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// generateUUID is a placeholder - in real implementation use uuid package
func (s *TournamentService) generateUUID() string {
	s.uuidCounter++
	return fmt.Sprintf("uuid-%d-%d", time.Now().UnixNano(), s.uuidCounter)
}

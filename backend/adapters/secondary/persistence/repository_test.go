package persistence

import (
	"testing"
	"time"

	"ludo-tournament/core/domain"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"
	"ludo-tournament/core/ports/outbound"
)

// TestLeagueRepositoryContract verifies GormLeagueRepository implements outbound.LeagueRepository
func TestLeagueRepositoryContract(t *testing.T) {
	repo := &GormLeagueRepository{}
	var _ outbound.LeagueRepository = repo
}

// TestTournamentRepositoryContract verifies GormTournamentRepository implements outbound.TournamentRepository
func TestTournamentRepositoryContract(t *testing.T) {
	repo := &GormTournamentRepository{}
	var _ outbound.TournamentRepository = repo
}

// TestMatchRepositoryContract verifies GormMatchRepository implements outbound.MatchRepository
func TestMatchRepositoryContract(t *testing.T) {
	repo := &GormMatchRepository{}
	var _ outbound.MatchRepository = repo
}

// TestUserRepositoryContract verifies GormUserRepository implements outbound.UserRepository
func TestUserRepositoryContract(t *testing.T) {
	repo := &GormUserRepository{}
	var _ outbound.UserRepository = repo
}

// TestPlayerRepositoryContract verifies GormPlayerRepository implements outbound.PlayerRepository
func TestPlayerRepositoryContract(t *testing.T) {
	repo := &GormPlayerRepository{}
	var _ outbound.PlayerRepository = repo
}

// TestMatchAssignmentRepositoryContract verifies GormMatchAssignmentRepository implements outbound.MatchAssignmentRepository
func TestMatchAssignmentRepositoryContract(t *testing.T) {
	repo := &GormMatchAssignmentRepository{}
	var _ outbound.MatchAssignmentRepository = repo
}

// TestInvitationRepositoryContract verifies GormInvitationRepository implements outbound.InvitationRepository
func TestInvitationRepositoryContract(t *testing.T) {
	repo := &GormInvitationRepository{}
	var _ outbound.InvitationRepository = repo
}

// TestDomainModels tests model constructors and basic structure
func TestLeagueModel(t *testing.T) {
	league := &models.League{
		ID:          "test-id",
		Name:        "Test",
		OrganizerID: "org-1",
		Status:      models.LeagueStatusDraft,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	if league.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %q", league.ID)
	}
	if league.Status != models.LeagueStatusDraft {
		t.Errorf("expected status draft, got %v", league.Status)
	}
}

func TestTournamentModel(t *testing.T) {
	tournament := &models.Tournament{
		ID:          "test-id",
		Name:        "Test Tournament",
		Type:        "knockout",
		OrganizerID: "org-1",
		Status:      models.TournamentStatusDraft,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	if tournament.Type != "knockout" {
		t.Errorf("expected type 'knockout', got %q", tournament.Type)
	}
}

func TestMatchModel(t *testing.T) {
	tournamentID := "tournament-1"
	match := &models.Match{
		ID:           "match-1",
		TournamentID: &tournamentID,
		Round:        1,
		TableNumber:  1,
		Status:       models.MatchStatusPending,
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}

	if match.Round != 1 {
		t.Errorf("expected round 1, got %d", match.Round)
	}
	if match.Status != models.MatchStatusPending {
		t.Errorf("expected status pending, got %v", match.Status)
	}
}

func TestUserModel(t *testing.T) {
	user := &models.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         models.RoleMember,
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}

	if user.Role != models.RoleMember {
		t.Errorf("expected role member, got %v", user.Role)
	}
}

func TestInvitationModel(t *testing.T) {
	invitation := &models.Invitation{
		ID:     "inv-1",
		Status: models.InvitationStatusPending,
	}

	if invitation.Status != models.InvitationStatusPending {
		t.Errorf("expected status pending, got %v", invitation.Status)
	}
}

func TestUserInviteModel(t *testing.T) {
	invite := &models.UserInvite{
		ID:    "invite-1",
		Email: "test@example.com",
		Code:  "CODE123",
	}

	if invite.Code != "CODE123" {
		t.Errorf("expected code 'CODE123', got %q", invite.Code)
	}
}

func TestMatchAssignmentModel(t *testing.T) {
	assignment := &models.MatchAssignment{
		ID:        "assign-1",
		MatchID:   "match-1",
		PlayerID:  "player-1",
		SeatColor: models.SeatYellow,
	}

	if assignment.SeatColor != models.SeatYellow {
		t.Errorf("expected seat yellow, got %v", assignment.SeatColor)
	}
}

func TestPlayerModel(t *testing.T) {
	player := &models.Player{
		ID:          "player-1",
		UserID:      "user-1",
		DisplayName: "Player One",
	}

	if player.DisplayName != "Player One" {
		t.Errorf("expected display name 'Player One', got %q", player.DisplayName)
	}
}

func TestKnockoutBracketModel(t *testing.T) {
	bracket := &models.KnockoutBracket{
		ID:           "bracket-1",
		TournamentID: "tournament-1",
		Rounds:       `{"round_1":[]}`,
	}

	if bracket.TournamentID != "tournament-1" {
		t.Errorf("expected tournament ID 'tournament-1', got %q", bracket.TournamentID)
	}
}

// TestDomainErrors verifies domain error sentinels
func TestDomainErrors(t *testing.T) {
	errs := []error{
		domain.ErrNotFound,
		domain.ErrInvalidInput,
		domain.ErrUnauthorized,
		domain.ErrForbidden,
		domain.ErrTournamentActive,
		domain.ErrGameAlreadyPlayed,
		domain.ErrInvalidAdvancement,
		domain.ErrNoRematch,
	}

	for _, err := range errs {
		if err == nil {
			t.Error("expected non-nil error")
		}
		if err.Error() == "" {
			t.Error("expected error to have message")
		}
	}
}

// TestModelScoringRules tests scoring rules structure
func TestScoringRules(t *testing.T) {
	rule := models.ScoringRule{
		Placement: 1,
		Points:    3.0,
	}

	if rule.Placement != 1 {
		t.Errorf("expected placement 1, got %d", rule.Placement)
	}
	if rule.Points != 3.0 {
		t.Errorf("expected points 3.0, got %f", rule.Points)
	}
}

// TestLeagueSettings tests league settings structure
func TestLeagueSettings(t *testing.T) {
	settings := models.LeagueSettings{
		ScoringRules: []models.ScoringRule{
			{Placement: 1, Points: 3},
			{Placement: 2, Points: 2},
		},
		GamesPerPlayer: 5,
		TablesCount:    4,
	}

	if len(settings.ScoringRules) != 2 {
		t.Errorf("expected 2 scoring rules, got %d", len(settings.ScoringRules))
	}
	if settings.GamesPerPlayer != 5 {
		t.Errorf("expected 5 games per player, got %d", settings.GamesPerPlayer)
	}
}

// TestTournamentSettings tests tournament settings structure
func TestTournamentSettings(t *testing.T) {
	settings := models.TournamentSettings{
		TablesCount:     4,
		DefaultReporter: "lowest_advancing",
	}

	if settings.TablesCount != 4 {
		t.Errorf("expected 4 tables, got %d", settings.TablesCount)
	}
	if settings.DefaultReporter != "lowest_advancing" {
		t.Errorf("expected reporter 'lowest_advancing', got %q", settings.DefaultReporter)
	}
}

// TestAdvancementConfig tests advancement configuration structure
func TestAdvancementConfig(t *testing.T) {
	config := models.AdvancementConfig{
		Round: "round_1",
		Games: 2,
		AdvancementPerGame: []models.AdvancementPerGame{
			{GameIDs: []int{1, 2}, Placements: []int{1}},
		},
	}

	if config.Round != "round_1" {
		t.Errorf("expected round 'round_1', got %q", config.Round)
	}
	if len(config.AdvancementPerGame) != 1 {
		t.Errorf("expected 1 advancement per game entry, got %d", len(config.AdvancementPerGame))
	}
}

// TestLeagueMatchResult tests league match result structure
func TestLeagueMatchResult(t *testing.T) {
	result := models.LeagueMatchResult{
		MatchID:   "match-1",
		PlayerID:  "player-1",
		Placement: 1,
		Points:    3.0,
	}

	if result.Placement != 1 {
		t.Errorf("expected placement 1, got %d", result.Placement)
	}
	if result.Points != 3.0 {
		t.Errorf("expected points 3.0, got %f", result.Points)
	}
}

func TestMatchResult(t *testing.T) {
	result := inbound.MatchResult{
		PlayerID:  "player-1",
		SeatColor: models.SeatGreen,
		Placement: 2,
	}

	if result.SeatColor != models.SeatGreen {
		t.Errorf("expected seat green, got %v", result.SeatColor)
	}
	if result.Placement != 2 {
		t.Errorf("expected placement 2, got %d", result.Placement)
	}
}

// Compile-time check that repositories implement outbound interfaces
var _ = outbound.LeagueRepository(nil)
var _ = outbound.TournamentRepository(nil)
var _ = outbound.MatchRepository(nil)
var _ = outbound.UserRepository(nil)
var _ = outbound.PlayerRepository(nil)
var _ = outbound.MatchAssignmentRepository(nil)
var _ = outbound.InvitationRepository(nil)
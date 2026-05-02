package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

type mockLeagueService struct {
	createFn       func(ctx context.Context, name, organizerID string, settings models.LeagueSettings) (*models.League, error)
	getFn          func(ctx context.Context, id string) (*models.League, error)
	updateFn       func(ctx context.Context, id string, settings models.LeagueSettings) error
	deleteFn       func(ctx context.Context, id string) error
	generateSchFn  func(ctx context.Context, leagueID string, playDates []string) error
	generatePairFn func(ctx context.Context, leagueID, playDate string) ([]inbound.TablePairing, error)
	reportMatchFn  func(ctx context.Context, matchID string, results []inbound.LeagueMatchResultInput, reportedBy string) error
	getStandingsFn func(ctx context.Context, leagueID string) ([]inbound.PlayerStanding, error)
	addTiebreakFn  func(ctx context.Context, leagueID string, playerIDs []string) error
}

func (m *mockLeagueService) CreateLeague(ctx context.Context, name, organizerID string, settings models.LeagueSettings) (*models.League, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name, organizerID, settings)
	}
	return &models.League{ID: "league-1", Name: name, OrganizerID: organizerID, Settings: settings, Status: models.LeagueStatusDraft, CreatedAt: time.Now(), ModifiedAt: time.Now()}, nil
}

func (m *mockLeagueService) GetLeague(ctx context.Context, id string) (*models.League, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return &models.League{ID: id, Name: "Test League", Status: models.LeagueStatusDraft}, nil
}

func (m *mockLeagueService) UpdateLeague(ctx context.Context, id string, settings models.LeagueSettings) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, settings)
	}
	return nil
}

func (m *mockLeagueService) DeleteLeague(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockLeagueService) GenerateSchedule(ctx context.Context, leagueID string, playDates []string) error {
	if m.generateSchFn != nil {
		return m.generateSchFn(ctx, leagueID, playDates)
	}
	return nil
}

func (m *mockLeagueService) GeneratePairings(ctx context.Context, leagueID, playDate string) ([]inbound.TablePairing, error) {
	if m.generatePairFn != nil {
		return m.generatePairFn(ctx, leagueID, playDate)
	}
	return []inbound.TablePairing{}, nil
}

func (m *mockLeagueService) ReportLeagueMatch(ctx context.Context, matchID string, results []inbound.LeagueMatchResultInput, reportedBy string) error {
	if m.reportMatchFn != nil {
		return m.reportMatchFn(ctx, matchID, results, reportedBy)
	}
	return nil
}

func (m *mockLeagueService) GetStandings(ctx context.Context, leagueID string) ([]inbound.PlayerStanding, error) {
	if m.getStandingsFn != nil {
		return m.getStandingsFn(ctx, leagueID)
	}
	return []inbound.PlayerStanding{}, nil
}

func (m *mockLeagueService) AddTiebreaker(ctx context.Context, leagueID string, playerIDs []string) error {
	if m.addTiebreakFn != nil {
		return m.addTiebreakFn(ctx, leagueID, playerIDs)
	}
	return nil
}

func TestCreateLeagueHandler_Success(t *testing.T) {
	svc := &mockLeagueService{
		createFn: func(ctx context.Context, name, organizerID string, settings models.LeagueSettings) (*models.League, error) {
			return &models.League{
				ID:          "new-league-id",
				Name:        name,
				OrganizerID: organizerID,
				Settings:    settings,
				Status:      models.LeagueStatusDraft,
				CreatedAt:   time.Now(),
				ModifiedAt:  time.Now(),
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":"My League","organizerId":"org-123","settings":{"tablesCount":4}}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/leagues", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateLeagueHandler(svc)
	handler(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var league models.League
	if err := json.Unmarshal(w.Body.Bytes(), &league); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if league.Name != "My League" {
		t.Errorf("expected name 'My League', got %q", league.Name)
	}
	if league.ID != "new-league-id" {
		t.Errorf("expected ID 'new-league-id', got %q", league.ID)
	}
}

func TestCreateLeagueHandler_InvalidJSON(t *testing.T) {
	svc := &mockLeagueService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/leagues", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateLeagueHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateLeagueHandler_MissingName(t *testing.T) {
	svc := &mockLeagueService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"organizerId":"org-123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/leagues", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateLeagueHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListLeaguesHandler_NotImplemented(t *testing.T) {
	svc := &mockLeagueService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/leagues", nil)

	handler := ListLeaguesHandler(svc)
	handler(c)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}
}

func TestGetLeagueHandler_Success(t *testing.T) {
	svc := &mockLeagueService{
		getFn: func(ctx context.Context, id string) (*models.League, error) {
			return &models.League{
				ID:     id,
				Name:   "Test League",
				Status: models.LeagueStatusLive,
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "league-123"}}

	c.Request, _ = http.NewRequest(http.MethodGet, "/leagues/league-123", nil)

	handler := GetLeagueHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteLeagueHandler_Success(t *testing.T) {
	svc := &mockLeagueService{
		deleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "league-123"}}

	c.Request, _ = http.NewRequest(http.MethodDelete, "/leagues/league-123", nil)

	handler := DeleteLeagueHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetLeagueStandingsHandler_Success(t *testing.T) {
	svc := &mockLeagueService{
		getStandingsFn: func(ctx context.Context, leagueID string) ([]inbound.PlayerStanding, error) {
			return []inbound.PlayerStanding{
				{PlayerID: "p1", DisplayName: "Player 1", GamesPlayed: 5, TotalPoints: 15, Wins: 3, Rank: 1},
				{PlayerID: "p2", DisplayName: "Player 2", GamesPlayed: 5, TotalPoints: 12, Wins: 2, Rank: 2},
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "league-123"}}

	c.Request, _ = http.NewRequest(http.MethodGet, "/leagues/league-123/standings", nil)

	handler := GetLeagueStandingsHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenerateLeaguePairingsHandler_Success(t *testing.T) {
	svc := &mockLeagueService{
		generatePairFn: func(ctx context.Context, leagueID, playDate string) ([]inbound.TablePairing, error) {
			return []inbound.TablePairing{
				{MatchID: "match-1", PlayDate: playDate, TableNumber: 1, PlayerIDs: []string{"p1", "p2", "p3", "p4"}},
				{MatchID: "match-2", PlayDate: playDate, TableNumber: 2, PlayerIDs: []string{"p5", "p6", "p7", "p8"}},
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "league-123"}}

	body := `{"playDate":"2024-01-15"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/leagues/league-123/pairings/generate", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := GenerateLeaguePairingsHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
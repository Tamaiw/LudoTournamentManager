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

type mockTournamentService struct {
	createFn           func(ctx context.Context, name, organizerID string, settings models.TournamentSettings) (*models.Tournament, error)
	getFn              func(ctx context.Context, id string) (*models.Tournament, error)
	updateFn           func(ctx context.Context, id string, settings models.TournamentSettings) error
	deleteFn           func(ctx context.Context, id string) error
	generateBracketFn  func(ctx context.Context, tournamentID string, playerIDs []string) error
	getBracketFn       func(ctx context.Context, tournamentID string) (*models.KnockoutBracket, error)
	reportMatchFn      func(ctx context.Context, matchID string, results []inbound.MatchResult, reportedBy string) error
	getCurrentPairingsFn func(ctx context.Context, tournamentID string) ([]inbound.GamePairing, error)
	canEditGameFn      func(ctx context.Context, gameID string) (bool, error)
}

func (m *mockTournamentService) CreateTournament(ctx context.Context, name, organizerID string, settings models.TournamentSettings) (*models.Tournament, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name, organizerID, settings)
	}
	return &models.Tournament{ID: "tournament-1", Name: name, OrganizerID: organizerID, Settings: settings, Status: models.TournamentStatusDraft, CreatedAt: time.Now(), ModifiedAt: time.Now()}, nil
}

func (m *mockTournamentService) GetTournament(ctx context.Context, id string) (*models.Tournament, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return &models.Tournament{ID: id, Name: "Test Tournament", Status: models.TournamentStatusDraft}, nil
}

func (m *mockTournamentService) UpdateTournament(ctx context.Context, id string, settings models.TournamentSettings) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, settings)
	}
	return nil
}

func (m *mockTournamentService) DeleteTournament(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTournamentService) GenerateBracket(ctx context.Context, tournamentID string, playerIDs []string) error {
	if m.generateBracketFn != nil {
		return m.generateBracketFn(ctx, tournamentID, playerIDs)
	}
	return nil
}

func (m *mockTournamentService) GetBracket(ctx context.Context, tournamentID string) (*models.KnockoutBracket, error) {
	if m.getBracketFn != nil {
		return m.getBracketFn(ctx, tournamentID)
	}
	return &models.KnockoutBracket{ID: "bracket-1", TournamentID: tournamentID}, nil
}

func (m *mockTournamentService) ReportMatch(ctx context.Context, matchID string, results []inbound.MatchResult, reportedBy string) error {
	if m.reportMatchFn != nil {
		return m.reportMatchFn(ctx, matchID, results, reportedBy)
	}
	return nil
}

func (m *mockTournamentService) GetCurrentRoundPairings(ctx context.Context, tournamentID string) ([]inbound.GamePairing, error) {
	if m.getCurrentPairingsFn != nil {
		return m.getCurrentPairingsFn(ctx, tournamentID)
	}
	return []inbound.GamePairing{}, nil
}

func (m *mockTournamentService) CanEditGame(ctx context.Context, gameID string) (bool, error) {
	if m.canEditGameFn != nil {
		return m.canEditGameFn(ctx, gameID)
	}
	return true, nil
}

func TestCreateTournamentHandler_Success(t *testing.T) {
	svc := &mockTournamentService{
		createFn: func(ctx context.Context, name, organizerID string, settings models.TournamentSettings) (*models.Tournament, error) {
			return &models.Tournament{
				ID:          "new-tournament-id",
				Name:        name,
				OrganizerID: organizerID,
				Settings:    settings,
				Status:      models.TournamentStatusDraft,
				CreatedAt:   time.Now(),
				ModifiedAt:  time.Now(),
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":"My Tournament","organizerId":"org-123","settings":{"tablesCount":4}}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tournaments", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var tournament models.Tournament
	if err := json.Unmarshal(w.Body.Bytes(), &tournament); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if tournament.Name != "My Tournament" {
		t.Errorf("expected name 'My Tournament', got %q", tournament.Name)
	}
	if tournament.ID != "new-tournament-id" {
		t.Errorf("expected ID 'new-tournament-id', got %q", tournament.ID)
	}
}

func TestCreateTournamentHandler_InvalidJSON(t *testing.T) {
	svc := &mockTournamentService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/tournaments", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTournamentHandler_MissingName(t *testing.T) {
	svc := &mockTournamentService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"organizerId":"org-123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tournaments", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := CreateTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListTournamentsHandler_NotImplemented(t *testing.T) {
	svc := &mockTournamentService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/tournaments", nil)

	handler := ListTournamentsHandler(svc)
	handler(c)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}
}

func TestGetTournamentHandler_Success(t *testing.T) {
	svc := &mockTournamentService{
		getFn: func(ctx context.Context, id string) (*models.Tournament, error) {
			return &models.Tournament{
				ID:     id,
				Name:   "Test Tournament",
				Status: models.TournamentStatusDraft,
			}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "tournament-123"}}

	c.Request, _ = http.NewRequest(http.MethodGet, "/tournaments/tournament-123", nil)

	handler := GetTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateTournamentHandler_Success(t *testing.T) {
	svc := &mockTournamentService{
		updateFn: func(ctx context.Context, id string, settings models.TournamentSettings) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "tournament-123"}}

	body := `{"tablesCount":8,"defaultReporter":"lowest_advancing"}`
	c.Request, _ = http.NewRequest(http.MethodPatch, "/tournaments/tournament-123", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := UpdateTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteTournamentHandler_Success(t *testing.T) {
	svc := &mockTournamentService{
		deleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "tournament-123"}}

	c.Request, _ = http.NewRequest(http.MethodDelete, "/tournaments/tournament-123", nil)

	handler := DeleteTournamentHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReportMatchHandler_Success(t *testing.T) {
	svc := &mockTournamentService{
		reportMatchFn: func(ctx context.Context, matchID string, results []inbound.MatchResult, reportedBy string) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "tournament-123"}}

	body := `{"matchId":"match-1","results":[{"playerId":"p1","seatColor":"yellow","placement":1},{"playerId":"p2","seatColor":"green","placement":2},{"playerId":"p3","seatColor":"blue","placement":3},{"playerId":"p4","seatColor":"red","placement":4}],"reportedBy":"admin-1"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tournaments/tournament-123/matches", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ReportMatchHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReportMatchHandler_InvalidResults(t *testing.T) {
	svc := &mockTournamentService{
		reportMatchFn: func(ctx context.Context, matchID string, results []inbound.MatchResult, reportedBy string) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "tournament-123"}}

	// Only 2 results instead of 4
	body := `{"matchId":"match-1","results":[{"playerId":"p1","placement":1},{"playerId":"p2","placement":2}]}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tournaments/tournament-123/matches", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ReportMatchHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}